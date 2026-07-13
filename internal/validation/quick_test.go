package validation

import (
	"context"
	"testing"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

type fakeProber struct{ tcp, tls, handshake int }

func (p *fakeProber) TCP(context.Context, string, time.Duration) (time.Duration, error) {
	p.tcp++
	return time.Millisecond, nil
}
func (p *fakeProber) TLS(context.Context, TLSRequest) (time.Duration, error) {
	p.tls++
	return time.Millisecond, nil
}
func (p *fakeProber) ProxyHandshake(context.Context, *node.Node, time.Duration) (time.Duration, error) {
	p.handshake++
	return time.Millisecond, nil
}

type fakeResolver struct{}

func (fakeResolver) LookupHost(context.Context, string) ([]string, error) {
	return []string{"192.0.2.1"}, nil
}

func TestQuickCheckerDoesNotTCPPrefilterQUIC(t *testing.T) {
	prober := &fakeProber{}
	checker := QuickChecker{Prober: prober, Resolver: fakeResolver{}, Timeout: time.Second}
	result := checker.Check(context.Background(), &node.Node{Protocol: node.ProtoHysteria2, Server: "example.com", Port: 443, Password: "pw", Security: "tls"})
	if !result.Passed || prober.tcp != 0 || prober.tls != 0 {
		t.Fatalf("result=%+v prober=%+v", result, prober)
	}
}
