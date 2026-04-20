// Package sources fetches and parses upstream subscription feeds into []Node.
package sources

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/config"
	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

// Fetch returns nodes collected from a single source. Parse errors on
// individual entries are swallowed; only hard failures (HTTP error, decode
// error on the whole blob) propagate. The passed context is honored for
// dial + read, so a workflow-wide deadline cancels in-flight requests.
func Fetch(ctx context.Context, src config.Source, timeout time.Duration) ([]*node.Node, error) {
	if !src.Enabled {
		return nil, nil
	}

	client := &http.Client{Timeout: timeout}
	req, err := http.NewRequestWithContext(ctx, "GET", src.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("source %q: build request: %w", src.Name, err)
	}
	req.Header.Set("User-Agent", "free-vpn-subscriptions/1.0 (+https://github.com/Au1rxx/free-vpn-subscriptions)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("source %q: http: %w", src.Name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("source %q: http %d", src.Name, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("source %q: read: %w", src.Name, err)
	}

	var nodes []*node.Node
	switch src.Format {
	case "uri-list":
		nodes = parseURIList(string(body))
	case "base64":
		decoded, err := node.B64Decode(strings.TrimSpace(string(body)))
		if err != nil {
			return nil, fmt.Errorf("source %q: base64: %w", src.Name, err)
		}
		nodes = parseURIList(string(decoded))
	case "clash":
		nodes, err = parseClash(body)
		if err != nil {
			return nil, fmt.Errorf("source %q: clash: %w", src.Name, err)
		}
	default:
		return nil, fmt.Errorf("source %q: unknown format %q", src.Name, src.Format)
	}

	for _, n := range nodes {
		n.SourceName = src.Name
	}
	return nodes, nil
}

// parseURIList walks a newline-separated list of proxy URIs.
func parseURIList(s string) []*node.Node {
	var out []*node.Node
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		n, err := node.ParseURI(line)
		if err != nil {
			continue
		}
		if n.Valid() {
			out = append(out, n)
		}
	}
	return out
}

// clashConfig is the subset of a Clash YAML we need.
type clashConfig struct {
	Proxies []map[string]any `yaml:"proxies"`
}

// parseClash extracts proxies from a Clash-style YAML into normalized Nodes.
func parseClash(body []byte) ([]*node.Node, error) {
	var cc clashConfig
	if err := yaml.Unmarshal(body, &cc); err != nil {
		return nil, err
	}
	var out []*node.Node
	for _, p := range cc.Proxies {
		n := clashProxyToNode(p)
		if n != nil && n.Valid() {
			out = append(out, n)
		}
	}
	return out, nil
}
