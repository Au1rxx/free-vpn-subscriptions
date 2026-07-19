package validation

import (
	"testing"

	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

func TestRouteForAllProtocolFamilies(t *testing.T) {
	tests := []struct {
		protocol string
		kind     TransportKind
	}{
		{node.ProtoVLESS, TransportTCP}, {node.ProtoVMess, TransportTCP},
		{node.ProtoTrojan, TransportTCP}, {node.ProtoSS, TransportTCP}, {node.ProtoSSR, TransportTCP},
		{node.ProtoHysteria, TransportQUIC}, {node.ProtoHysteria2, TransportQUIC}, {node.ProtoTUIC, TransportQUIC},
		{node.ProtoWireGuard, TransportWireGuard}, {node.ProtoSOCKS4, TransportProxyHandshake},
		{node.ProtoSOCKS5, TransportProxyHandshake}, {node.ProtoHTTP, TransportProxyHandshake}, {node.ProtoHTTPS, TransportProxyHandshake},
	}
	for _, tt := range tests {
		route := RouteFor(&node.Node{Protocol: tt.protocol})
		if route.Kind != tt.kind || !route.NeedsSingBox {
			t.Fatalf("protocol=%s route=%+v", tt.protocol, route)
		}
	}
	if !RouteFor(&node.Node{Protocol: node.ProtoVLESS, Security: "tls"}).PrecheckTLS {
		t.Fatal("declared TLS was not routed to precheck")
	}
	if RouteFor(&node.Node{Protocol: node.ProtoVLESS, Security: "none"}).PrecheckTLS {
		t.Fatal("plain config was routed to TLS precheck")
	}
}
