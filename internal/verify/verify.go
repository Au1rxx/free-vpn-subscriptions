// Package verify runs real HTTP-over-proxy probes through candidate nodes
// using sing-box as a batched subprocess.
//
// Unlike internal/probe (TCP + TLS handshake, cannot see proxy auth or real
// HTTP success) this package forwards a small request through each proxy
// and checks that an end-to-end HTTP 204 actually comes back. Nodes that
// reach here with http_ok==true have demonstrably functioning credentials,
// routing, and exit networking at probe time.
//
// Design: start one sing-box subprocess per batch of N nodes, each node
// getting its own mixed inbound on 127.0.0.1:base+i routed to its own
// outbound. Then do 2 probe rounds spaced apart so we can distinguish
// transient passes from actually-stable nodes.
package verify

import (
	"context"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

// Config controls the verification stage. Populated from YAML.
type Config struct {
	Enabled         bool     // run the stage at all
	CandidatePool   int      // how many top-RTT nodes to send through HTTP probing
	BatchSize       int      // number of outbounds per sing-box subprocess
	BasePort        int      // listen port for node i in a batch: base + i
	Concurrency     int      // parallel HTTP probes within a batch
	TimeoutMS       int      // per-request deadline
	Rounds          int      // number of probe rounds (typically 2)
	RoundGapMS      int      // sleep between rounds
	Targets         []string // URLs to GET; all must 2xx/204 for a round to pass
	SingBoxBin      string   // path to sing-box executable
	StartupTimeout  time.Duration
}

// Result is verification output for one node.
type Result struct {
	Node        *node.Node
	Passed      bool
	HTTPMedian  int // milliseconds, median across passing requests
	SuccessRate int // percent of (rounds × targets) requests that succeeded
}

// Run verifies nodes and returns the subset confirmed to forward HTTP traffic.
// `candidates` are nodes that already passed TCP+TLS, sorted ascending by RTT.
func Run(ctx context.Context, candidates []*node.Node, cfg Config) []*node.Node {
	if !cfg.Enabled || len(candidates) == 0 {
		return candidates
	}
	cfg = applyDefaults(cfg)

	pool := candidates
	if cfg.CandidatePool > 0 && len(pool) > cfg.CandidatePool {
		pool = pool[:cfg.CandidatePool]
	}

	// Pre-filter: drop any node sing-box itself rejects (corrupt cipher,
	// missing uuid, unsupported flow, etc.). Saves us from batch-abort
	// failure modes where one bad outbound kills the whole subprocess.
	pool = prefilterValid(ctx, pool, cfg)
	fmt.Fprintf(os.Stderr, "  [verify] %d / %d candidates accepted by sing-box config check\n", len(pool), min(cfg.CandidatePool, len(candidates)))

	results := make([]Result, len(pool))
	for i, n := range pool {
		results[i] = Result{Node: n}
	}

	for round := 1; round <= cfg.Rounds; round++ {
		if ctx.Err() != nil {
			break
		}
		fmt.Fprintf(os.Stderr, "  [verify] round %d/%d over %d candidates\n", round, cfg.Rounds, len(pool))
		runRound(ctx, pool, results, cfg)
		if round < cfg.Rounds && cfg.RoundGapMS > 0 {
			select {
			case <-ctx.Done():
			case <-time.After(time.Duration(cfg.RoundGapMS) * time.Millisecond):
			}
		}
	}

	return finalize(results, cfg)
}

// prefilterValid runs `sing-box check` against each candidate's outbound
// in parallel and keeps only the ones sing-box accepts.
func prefilterValid(ctx context.Context, nodes []*node.Node, cfg Config) []*node.Node {
	out := make([]*node.Node, len(nodes))
	sem := make(chan struct{}, 24)
	var wg sync.WaitGroup
	for i, n := range nodes {
		if ctx.Err() != nil {
			break
		}
		wg.Add(1)
		sem <- struct{}{}
		go func(i int, n *node.Node) {
			defer wg.Done()
			defer func() { <-sem }()
			if validOutbound(cfg.SingBoxBin, n) {
				out[i] = n
			}
		}(i, n)
	}
	wg.Wait()
	kept := make([]*node.Node, 0, len(nodes))
	for _, n := range out {
		if n != nil {
			kept = append(kept, n)
		}
	}
	return kept
}

// runRound mutates `results` in-place, running one round of probing over all
// candidates, batched.
func runRound(ctx context.Context, pool []*node.Node, results []Result, cfg Config) {
	total := len(pool)
	for start := 0; start < total; start += cfg.BatchSize {
		if ctx.Err() != nil {
			return
		}
		end := start + cfg.BatchSize
		if end > total {
			end = total
		}
		batch := pool[start:end]
		offsets := make([]int, len(batch))
		for i := range offsets {
			offsets[i] = start + i
		}
		probeBatch(ctx, batch, offsets, results, cfg)
	}
}

// probeBatch starts a sing-box subprocess for the given slice of nodes and
// collects HTTP probe results. Failure to start sing-box marks the whole
// batch as untested (they inherit TCP+TLS-only standing).
func probeBatch(ctx context.Context, batch []*node.Node, offsets []int, results []Result, cfg Config) {
	bc, err := buildBatchConfig(batch, cfg.BasePort)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [verify] batch config build failed: %v\n", err)
		return
	}
	sb, err := startSingbox(ctx, bc, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  [verify] sing-box start failed: %v\n", err)
		return
	}
	defer sb.stop()

	sem := make(chan struct{}, cfg.Concurrency)
	var wg sync.WaitGroup
	for i := range batch {
		if ctx.Err() != nil {
			break
		}
		wg.Add(1)
		sem <- struct{}{}
		go func(i int) {
			defer wg.Done()
			defer func() { <-sem }()
			port := cfg.BasePort + i
			outcome := probeViaSocks(ctx, port, cfg)
			r := &results[offsets[i]]
			r.SuccessRate += outcome.successes
			if outcome.passed {
				r.Passed = true
			}
			// Track the best (lowest) median across rounds. Nodes that only
			// partially pass still get a latency so we can rank them.
			if outcome.medianMS > 0 && (r.HTTPMedian == 0 || outcome.medianMS < r.HTTPMedian) {
				r.HTTPMedian = outcome.medianMS
			}
		}(i)
	}
	wg.Wait()
}

// finalize selects results whose success rate across all probe attempts
// (rounds × targets) meets the 50 % threshold, and sets Node.LatencyMS to
// the real HTTP median so downstream ranking uses probe-real numbers.
//
// 50 % is lenient on purpose: free nodes often block exactly one of our
// two targets (Cloudflare captive portal is sometimes IP-banned even when
// gstatic works). A node that passes *any* half of the 2R×2T = 4 checks
// is still useful — the client-side url-test weeds out remaining noise.
func finalize(results []Result, cfg Config) []*node.Node {
	required := cfg.Rounds * len(cfg.Targets)
	threshold := required / 2
	if threshold < 1 {
		threshold = 1
	}
	kept := make([]*node.Node, 0, len(results))
	for _, r := range results {
		if r.SuccessRate < threshold {
			continue
		}
		if r.HTTPMedian > 0 {
			r.Node.LatencyMS = r.HTTPMedian
		}
		kept = append(kept, r.Node)
	}
	sort.Slice(kept, func(i, j int) bool {
		return kept[i].LatencyMS < kept[j].LatencyMS
	})
	return kept
}

func applyDefaults(c Config) Config {
	if c.CandidatePool == 0 {
		c.CandidatePool = 600
	}
	if c.BatchSize == 0 {
		c.BatchSize = 40
	}
	if c.BasePort == 0 {
		c.BasePort = 20000
	}
	if c.Concurrency == 0 {
		c.Concurrency = 20
	}
	if c.TimeoutMS == 0 {
		c.TimeoutMS = 6000
	}
	if c.Rounds == 0 {
		c.Rounds = 2
	}
	if c.RoundGapMS == 0 {
		c.RoundGapMS = 45000
	}
	if len(c.Targets) == 0 {
		c.Targets = []string{
			"http://www.gstatic.com/generate_204",
			"https://www.cloudflare.com/cdn-cgi/trace",
		}
	}
	if c.SingBoxBin == "" {
		c.SingBoxBin = "sing-box"
	}
	if c.StartupTimeout == 0 {
		c.StartupTimeout = 10 * time.Second
	}
	return c
}
