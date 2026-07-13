package validation

import "github.com/Au1rxx/free-vpn-subscriptions/pkg/node"

type TransportKind string

const (
	TransportTCP            TransportKind = "tcp"
	TransportProxyHandshake TransportKind = "proxy_handshake"
	TransportQUIC           TransportKind = "quic"
	TransportWireGuard      TransportKind = "wireguard"
)

type Route struct {
	Kind         TransportKind
	PrecheckTLS  bool
	NeedsSingBox bool
}

func RouteFor(n *node.Node) Route {
	route := Route{Kind: TransportTCP, NeedsSingBox: true, PrecheckTLS: n.Security == "tls" || n.Security == "reality"}
	switch n.Protocol {
	case node.ProtoHysteria, node.ProtoHysteria2, node.ProtoTUIC:
		route.Kind = TransportQUIC
	case node.ProtoWireGuard:
		route.Kind = TransportWireGuard
	case node.ProtoSOCKS4, node.ProtoSOCKS5, node.ProtoHTTP, node.ProtoHTTPS:
		route.Kind = TransportProxyHandshake
	}
	return route
}
