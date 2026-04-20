package verify

import (
	"github.com/Au1rxx/free-vpn-subscriptions/internal/node"
)

// Outbound renders a node into a sing-box outbound JSON object.
// Returns nil for protocols we cannot translate (keeps the batch clean).
// Tag is assigned by the caller (e.g. "out-0").
//
// Deliberately a superset of internal/subscribe.Singbox — that function is
// shaped for end-user clients and skips WS/gRPC/Reality details we need
// here to actually forward traffic through the probe pipeline.
func buildOutbound(n *node.Node, tag string) map[string]any {
	ob := map[string]any{
		"tag":         tag,
		"server":      n.Server,
		"server_port": n.Port,
	}
	switch n.Protocol {
	case node.ProtoVLESS:
		ob["type"] = "vless"
		ob["uuid"] = n.UUID
		if n.Flow != "" {
			ob["flow"] = n.Flow
		}
		if tls := buildTLS(n); tls != nil {
			ob["tls"] = tls
		}
		if tr := buildTransport(n); tr != nil {
			ob["transport"] = tr
		}
	case node.ProtoVMess:
		ob["type"] = "vmess"
		ob["uuid"] = n.UUID
		ob["alter_id"] = n.AlterID
		ob["security"] = orDefault(n.Cipher, "auto")
		if n.Security == "tls" {
			if tls := buildTLS(n); tls != nil {
				ob["tls"] = tls
			}
		}
		if tr := buildTransport(n); tr != nil {
			ob["transport"] = tr
		}
	case node.ProtoTrojan:
		ob["type"] = "trojan"
		ob["password"] = n.Password
		if tls := buildTLS(n); tls != nil {
			ob["tls"] = tls
		} else {
			ob["tls"] = map[string]any{"enabled": true, "server_name": orDefault(n.SNI, n.Server), "insecure": n.Insecure}
		}
		if tr := buildTransport(n); tr != nil {
			ob["transport"] = tr
		}
	case node.ProtoSS:
		ob["type"] = "shadowsocks"
		ob["method"] = n.Cipher
		ob["password"] = n.Password
	case node.ProtoHysteria2:
		ob["type"] = "hysteria2"
		ob["password"] = n.Password
		if tls := buildTLS(n); tls != nil {
			ob["tls"] = tls
		} else {
			ob["tls"] = map[string]any{"enabled": true, "server_name": orDefault(n.SNI, n.Server), "insecure": n.Insecure}
		}
	default:
		return nil
	}
	return ob
}

func buildTLS(n *node.Node) map[string]any {
	if n.Security != "tls" && n.Security != "reality" {
		return nil
	}
	sni := orDefault(n.SNI, n.Server)
	tls := map[string]any{
		"enabled":     true,
		"server_name": sni,
		"insecure":    n.Insecure,
	}
	if n.ALPN != "" {
		tls["alpn"] = []string{n.ALPN}
	}
	if n.Fingerprint != "" {
		tls["utls"] = map[string]any{"enabled": true, "fingerprint": n.Fingerprint}
	}
	if n.Security == "reality" {
		reality := map[string]any{"enabled": true}
		if n.PublicKey != "" {
			reality["public_key"] = n.PublicKey
		}
		if n.ShortID != "" {
			reality["short_id"] = n.ShortID
		}
		tls["reality"] = reality
		// Reality strictly requires utls; default to chrome when upstream omits.
		if _, ok := tls["utls"]; !ok {
			tls["utls"] = map[string]any{"enabled": true, "fingerprint": "chrome"}
		}
	}
	return tls
}

func buildTransport(n *node.Node) map[string]any {
	switch n.Network {
	case "ws":
		t := map[string]any{"type": "ws"}
		if n.Path != "" {
			t["path"] = n.Path
		}
		if n.Host != "" {
			t["headers"] = map[string]any{"Host": n.Host}
		}
		return t
	case "grpc":
		if n.ServiceName == "" {
			return nil
		}
		return map[string]any{"type": "grpc", "service_name": n.ServiceName}
	}
	return nil
}

func orDefault(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}
