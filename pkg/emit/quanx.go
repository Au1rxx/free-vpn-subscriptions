package emit

import (
	"fmt"
	"strings"

	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

// QuantumultX emits a QuantumultX-compatible [server_local] section.
// Coverage is partial: Shadowsocks, Trojan, and VMess are supported;
// VLESS and Hysteria2 are skipped because QuanX has no native mapping.
func QuantumultX(nodes []*node.Node) (string, error) {
	var lines []string
	lines = append(lines, "[server_local]")
	for i, n := range nodes {
		line := quanxLine(n, i)
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n") + "\n", nil
}

func quanxLine(n *node.Node, idx int) string {
	tag := fmt.Sprintf("%02d-%s-%s", idx+1, n.Protocol, safe(n.Name))
	switch n.Protocol {
	case node.ProtoSS:
		parts := []string{
			fmt.Sprintf("shadowsocks=%s:%d", n.Server, n.Port),
			"method=" + n.Cipher,
			"password=" + n.Password,
			"tag=" + tag,
		}
		return strings.Join(parts, ", ")
	case node.ProtoTrojan:
		parts := []string{
			fmt.Sprintf("trojan=%s:%d", n.Server, n.Port),
			"password=" + n.Password,
			"over-tls=true",
		}
		if n.SNI != "" {
			parts = append(parts, "tls-host="+n.SNI)
		}
		if n.Insecure {
			parts = append(parts, "tls-verification=false")
		}
		parts = append(parts, "tag="+tag)
		return strings.Join(parts, ", ")
	case node.ProtoVMess:
		parts := []string{
			fmt.Sprintf("vmess=%s:%d", n.Server, n.Port),
			"method=none",
			"password=" + n.UUID,
		}
		obfs := ""
		switch {
		case n.Network == "ws" && n.Security == "tls":
			obfs = "wss"
		case n.Network == "ws":
			obfs = "ws"
		case n.Security == "tls":
			obfs = "over-tls"
		}
		if obfs != "" {
			parts = append(parts, "obfs="+obfs)
		}
		if n.Network == "ws" {
			if n.Host != "" {
				parts = append(parts, "obfs-host="+n.Host)
			}
			if n.Path != "" {
				parts = append(parts, "obfs-uri="+n.Path)
			}
		}
		if n.Security == "tls" {
			if n.SNI != "" {
				parts = append(parts, "tls-host="+n.SNI)
			}
			if n.Insecure {
				parts = append(parts, "tls-verification=false")
			}
		}
		parts = append(parts, "tag="+tag)
		return strings.Join(parts, ", ")
	}
	return ""
}
