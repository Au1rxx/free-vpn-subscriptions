// Package probe measures the reachability and latency of proxy endpoints.
// It performs concurrent TCP handshakes to server:port; successful nodes have
// their LatencyMS populated.
package probe

import (
	"context"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/node"
)

// TCP probes every node in parallel (bounded by concurrency) and returns only
// those that completed a TCP handshake within the timeout. Nodes' LatencyMS
// field is populated. Cancelling ctx aborts pending dials; in-flight dials
// still respect their own per-node timeout.
func TCP(ctx context.Context, nodes []*node.Node, timeout time.Duration, concurrency int) []*node.Node {
	if concurrency <= 0 {
		concurrency = 50
	}
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	results := make([]*node.Node, len(nodes))
	for i, n := range nodes {
		if ctx.Err() != nil {
			break
		}
		wg.Add(1)
		sem <- struct{}{}
		go func(i int, n *node.Node) {
			defer wg.Done()
			defer func() { <-sem }()
			latency, ok := dial(ctx, n.Server, n.Port, timeout)
			if !ok {
				return
			}
			ms := int(latency / time.Millisecond)
			n.LatencyMS = ms
			n.TCPLatencyMS = ms
			results[i] = n
		}(i, n)
	}
	wg.Wait()

	alive := make([]*node.Node, 0, len(nodes))
	for _, n := range results {
		if n != nil {
			alive = append(alive, n)
		}
	}
	return alive
}

func dial(ctx context.Context, host string, port int, timeout time.Duration) (time.Duration, bool) {
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	dialer := &net.Dialer{Timeout: timeout}
	start := time.Now()
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return 0, false
	}
	_ = conn.Close()
	return time.Since(start), true
}
