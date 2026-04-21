package probe

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

// TestTCP_Loopback spins a local TCP listener and asserts the prober sees it
// as alive with a reasonable LatencyMS. Failure cases share a non-routable
// IP so the dial errors out fast (reserved TEST-NET-1).
func TestTCP_Loopback(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer l.Close()
	_, portStr, _ := net.SplitHostPort(l.Addr().String())
	port, _ := strconv.Atoi(portStr)

	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				return
			}
			_ = c.Close()
		}
	}()

	nodes := []*node.Node{
		{Name: "live", Protocol: node.ProtoSS, Server: "127.0.0.1", Port: port, Cipher: "aes-256-gcm", Password: "x"},
		{Name: "dead", Protocol: node.ProtoSS, Server: "192.0.2.1", Port: 1, Cipher: "aes-256-gcm", Password: "x"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	alive := TCP(ctx, nodes, 500*time.Millisecond, 4)
	if len(alive) != 1 {
		t.Fatalf("got %d alive nodes, want 1", len(alive))
	}
	if alive[0].Name != "live" {
		t.Errorf("alive node = %q, want live", alive[0].Name)
	}
	if alive[0].LatencyMS < 0 {
		t.Errorf("LatencyMS = %d, want >= 0", alive[0].LatencyMS)
	}
}

func TestNeedsTLSProbe(t *testing.T) {
	cases := []struct {
		n    *node.Node
		want bool
	}{
		{&node.Node{Protocol: node.ProtoTrojan}, true},
		{&node.Node{Protocol: node.ProtoHysteria2}, true},
		{&node.Node{Protocol: node.ProtoVLESS, Security: "tls"}, true},
		{&node.Node{Protocol: node.ProtoVLESS, Security: "reality"}, false},
		{&node.Node{Protocol: node.ProtoVMess, Security: "tls"}, true},
		{&node.Node{Protocol: node.ProtoVMess, Security: ""}, false},
		{&node.Node{Protocol: node.ProtoSS}, false},
	}
	for _, c := range cases {
		if got := needsTLSProbe(c.n); got != c.want {
			t.Errorf("needsTLSProbe(%s/%s) = %v, want %v", c.n.Protocol, c.n.Security, got, c.want)
		}
	}
}
