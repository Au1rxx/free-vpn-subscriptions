package sources

import (
	"fmt"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/node"
)

// clashProxyToNode maps a single Clash proxy map into a normalized Node.
// Returns nil if the proxy type is unsupported.
func clashProxyToNode(p map[string]any) *node.Node {
	ptype, _ := p["type"].(string)
	name, _ := p["name"].(string)
	server, _ := p["server"].(string)
	port := clashInt(p["port"])
	if server == "" || port == 0 {
		return nil
	}

	n := &node.Node{
		Name:   name,
		Server: server,
		Port:   port,
	}

	switch ptype {
	case "vless":
		n.Protocol = node.ProtoVLESS
		n.UUID, _ = p["uuid"].(string)
		n.Flow, _ = p["flow"].(string)
		n.Network = clashStr(p["network"], "tcp")
		if tls, _ := p["tls"].(bool); tls {
			n.Security = "tls"
		}
		if ro, ok := p["reality-opts"].(map[string]any); ok {
			n.Security = "reality"
			n.PublicKey, _ = ro["public-key"].(string)
			n.ShortID, _ = ro["short-id"].(string)
		}
		n.SNI = clashStr(p["servername"], clashStr(p["sni"], ""))
		n.Fingerprint, _ = p["client-fingerprint"].(string)
		n.Insecure, _ = p["skip-cert-verify"].(bool)

	case "vmess":
		n.Protocol = node.ProtoVMess
		n.UUID, _ = p["uuid"].(string)
		n.AlterID = clashInt(p["alterId"])
		n.Cipher, _ = p["cipher"].(string)
		n.Network = clashStr(p["network"], "tcp")
		if tls, _ := p["tls"].(bool); tls {
			n.Security = "tls"
		}
		n.SNI = clashStr(p["servername"], "")
		n.Insecure, _ = p["skip-cert-verify"].(bool)

	case "trojan":
		n.Protocol = node.ProtoTrojan
		n.Password, _ = p["password"].(string)
		n.Network = clashStr(p["network"], "tcp")
		n.Security = "tls"
		n.SNI, _ = p["sni"].(string)
		n.Insecure, _ = p["skip-cert-verify"].(bool)

	case "ss", "shadowsocks":
		n.Protocol = node.ProtoSS
		n.Cipher, _ = p["cipher"].(string)
		n.Password, _ = p["password"].(string)

	case "hysteria2", "hy2":
		n.Protocol = node.ProtoHysteria2
		n.Password = clashStr(p["password"], clashStr(p["auth"], ""))
		n.SNI, _ = p["sni"].(string)
		n.Security = "tls"
		n.Insecure, _ = p["skip-cert-verify"].(bool)

	default:
		return nil
	}
	return n
}

func clashStr(v any, fallback string) string {
	if s, ok := v.(string); ok && s != "" {
		return s
	}
	return fallback
}

func clashInt(v any) int {
	switch x := v.(type) {
	case int:
		return x
	case int64:
		return int(x)
	case float64:
		return int(x)
	case string:
		var n int
		_, _ = fmt.Sscanf(x, "%d", &n)
		return n
	}
	return 0
}
