package probe

import (
	"context"
	"crypto/tls"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/node"
)

// TLS filters TCP-alive nodes further by attempting a TLS ClientHello against
// nodes that advertise TLS/Reality. Nodes that do not use TLS (e.g. plain
// VMess/SS) are passed through unchanged. Failed handshakes drop the node.
//
// Free proxies often present self-signed certs so InsecureSkipVerify is true —
// we only care that the peer speaks TLS at all, not about trust anchors.
// Cancelling ctx aborts queued handshakes.
func TLS(ctx context.Context, nodes []*node.Node, timeout time.Duration, concurrency int) []*node.Node {
	if concurrency <= 0 {
		concurrency = 50
	}
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	keep := make([]*node.Node, len(nodes))
	for i, n := range nodes {
		if !needsTLSProbe(n) {
			keep[i] = n
			continue
		}
		if ctx.Err() != nil {
			break
		}
		wg.Add(1)
		sem <- struct{}{}
		go func(i int, n *node.Node) {
			defer wg.Done()
			defer func() { <-sem }()
			if tlsHandshake(ctx, n, timeout) {
				keep[i] = n
			}
		}(i, n)
	}
	wg.Wait()

	out := make([]*node.Node, 0, len(nodes))
	for _, n := range keep {
		if n != nil {
			out = append(out, n)
		}
	}
	return out
}

// needsTLSProbe reports whether the node's protocol+security implies TLS on
// the wire. Reality is deliberately excluded: its ClientHello is spoofed to
// look like a legitimate target (e.g. microsoft.com) and a real TLS handshake
// against the proxy server doesn't tell us anything useful.
func needsTLSProbe(n *node.Node) bool {
	switch n.Protocol {
	case node.ProtoTrojan, node.ProtoHysteria2:
		return true
	case node.ProtoVLESS, node.ProtoVMess:
		return n.Security == "tls"
	}
	return false
}

func tlsHandshake(ctx context.Context, n *node.Node, timeout time.Duration) bool {
	addr := net.JoinHostPort(n.Server, strconv.Itoa(n.Port))
	dialer := &net.Dialer{Timeout: timeout}
	sni := n.SNI
	if sni == "" {
		sni = n.Server
	}
	// tls.Dialer honors ctx cancellation during the underlying TCP dial;
	// the handshake itself is bounded by the dialer timeout above.
	tlsDialer := &tls.Dialer{
		NetDialer: dialer,
		Config: &tls.Config{
			ServerName:         sni,
			InsecureSkipVerify: true,
		},
	}
	conn, err := tlsDialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
