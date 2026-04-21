// Package sources fetches upstream subscription feeds over HTTP and hands
// the raw bodies off to pkg/parse for format decoding.
package sources

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/config"
	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
	"github.com/Au1rxx/free-vpn-subscriptions/pkg/parse"
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
		nodes = parse.URIList(string(body))
	case "base64":
		nodes, err = parse.Base64List(body)
		if err != nil {
			return nil, fmt.Errorf("source %q: base64: %w", src.Name, err)
		}
	case "clash":
		nodes, err = parse.Clash(body)
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
