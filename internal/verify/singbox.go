package verify

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/node"
)

type batchConfig struct {
	path      string
	basePort  int
	batchSize int
}

// buildBatchConfig writes a sing-box config where each node gets its own
// mixed (SOCKS + HTTP) inbound routed to its own outbound. Returns the temp
// path; caller is responsible for nothing (startSingbox cleans up on stop).
func buildBatchConfig(batch []*node.Node, basePort int) (*batchConfig, error) {
	inbounds := make([]map[string]any, 0, len(batch))
	outbounds := make([]map[string]any, 0, len(batch)+1)
	rules := make([]map[string]any, 0, len(batch))

	for i, n := range batch {
		outTag := fmt.Sprintf("out-%d", i)
		inTag := fmt.Sprintf("in-%d", i)
		ob := buildOutbound(n, outTag)
		if ob == nil {
			// Can't forward this protocol — still reserve the inbound/rule
			// so port offsets line up with the batch index. Route it to
			// direct; the probe will fail against gstatic from our own IP
			// only if we have no egress, which we do, so mark these as
			// effectively undetected: we'll log and skip.
			ob = map[string]any{"type": "direct", "tag": outTag}
		}
		outbounds = append(outbounds, ob)
		inbounds = append(inbounds, map[string]any{
			"type":        "mixed",
			"tag":         inTag,
			"listen":      "127.0.0.1",
			"listen_port": basePort + i,
		})
		rules = append(rules, map[string]any{
			"inbound":  inTag,
			"outbound": outTag,
		})
	}
	outbounds = append(outbounds, map[string]any{"type": "direct", "tag": "direct"})

	cfg := map[string]any{
		"log":       map[string]any{"disabled": true},
		"inbounds":  inbounds,
		"outbounds": outbounds,
		"route":     map[string]any{"rules": rules, "final": "direct"},
	}
	raw, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	f, err := os.CreateTemp("", "fvs-verify-*.json")
	if err != nil {
		return nil, err
	}
	if _, err := f.Write(raw); err != nil {
		_ = f.Close()
		_ = os.Remove(f.Name())
		return nil, err
	}
	if err := f.Close(); err != nil {
		return nil, err
	}
	return &batchConfig{path: f.Name(), basePort: basePort, batchSize: len(batch)}, nil
}

// validOutbound single-node-checks whether sing-box accepts the outbound
// JSON for n. Used to pre-filter the candidate pool so one malformed node
// (e.g. corrupt SS cipher) doesn't abort an entire batch.
func validOutbound(bin string, n *node.Node) bool {
	ob := buildOutbound(n, "out-0")
	if ob == nil {
		return false
	}
	cfg := map[string]any{
		"log": map[string]any{"disabled": true},
		"inbounds": []map[string]any{
			{"type": "mixed", "tag": "in-0", "listen": "127.0.0.1", "listen_port": 0},
		},
		"outbounds": []map[string]any{
			ob,
			{"type": "direct", "tag": "direct"},
		},
		"route": map[string]any{
			"rules": []map[string]any{{"inbound": "in-0", "outbound": "out-0"}},
			"final": "direct",
		},
	}
	raw, err := json.Marshal(cfg)
	if err != nil {
		return false
	}
	f, err := os.CreateTemp("", "fvs-check-*.json")
	if err != nil {
		return false
	}
	defer os.Remove(f.Name())
	if _, err := f.Write(raw); err != nil {
		_ = f.Close()
		return false
	}
	_ = f.Close()
	return exec.Command(bin, "check", "-c", f.Name()).Run() == nil
}

type singboxProc struct {
	cmd     *exec.Cmd
	cfgPath string
	cancel  context.CancelFunc
}

// startSingbox launches sing-box with the generated batch config and waits
// until at least the first inbound port is accepting connections. If the
// process exits early (bad config), the returned error includes stderr.
func startSingbox(parent context.Context, bc *batchConfig, cfg Config) (*singboxProc, error) {
	ctx, cancel := context.WithCancel(parent)

	// Pre-flight: validate the config before committing to a full boot.
	// `sing-box check` catches most schema/field mistakes in ~20ms; faster
	// than waiting for a silent startup failure.
	checkCmd := exec.Command(cfg.SingBoxBin, "check", "-c", bc.path)
	if out, err := checkCmd.CombinedOutput(); err != nil {
		cancel()
		_ = os.Remove(bc.path)
		return nil, fmt.Errorf("config check: %w (%s)", err, truncate(string(out), 400))
	}

	cmd := exec.CommandContext(ctx, cfg.SingBoxBin, "run", "-c", bc.path)
	// Separate process group so we can SIGKILL the whole subtree on stop.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Start(); err != nil {
		cancel()
		_ = os.Remove(bc.path)
		return nil, err
	}

	// Wait for the first inbound to accept TCP. Most inbounds come up in
	// the same scheduler tick; probing port 0 is a good liveness signal.
	deadline := time.Now().Add(cfg.StartupTimeout)
	ready := false
	for time.Now().Before(deadline) {
		if tryDial(bc.basePort, 200*time.Millisecond) {
			ready = true
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	if !ready {
		cancel()
		_ = cmd.Process.Kill()
		_, _ = cmd.Process.Wait()
		_ = os.Remove(bc.path)
		return nil, fmt.Errorf("sing-box did not open port %d within %s", bc.basePort, cfg.StartupTimeout)
	}

	return &singboxProc{cmd: cmd, cfgPath: bc.path, cancel: cancel}, nil
}

func (s *singboxProc) stop() {
	if s == nil {
		return
	}
	if s.cmd != nil && s.cmd.Process != nil {
		// Kill the whole process group to cover any child goroutines/pids.
		_ = syscall.Kill(-s.cmd.Process.Pid, syscall.SIGKILL)
		_, _ = s.cmd.Process.Wait()
	}
	if s.cancel != nil {
		s.cancel()
	}
	if s.cfgPath != "" {
		_ = os.Remove(s.cfgPath)
	}
}

func tryDial(port int, timeout time.Duration) bool {
	c, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(port)), timeout)
	if err != nil {
		return false
	}
	_ = c.Close()
	return true
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

