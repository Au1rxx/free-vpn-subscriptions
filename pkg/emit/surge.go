package emit

import (
	"fmt"
	"strings"

	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

// Surge emits a Surge-compatible conf snippet (a [Proxy] section plus a
// [Proxy Group] block). Coverage is partial: Shadowsocks, Trojan, and VMess
// (with optional ws/tls) are supported; VLESS and Hysteria2 are skipped
// because Surge does not natively handle them.
func Surge(nodes []*node.Node) (string, error) {
	var lines []string
	var names []string
	lines = append(lines, "[Proxy]")
	for i, n := range nodes {
		line, name := surgeLine(n, i)
		if line == "" {
			continue
		}
		lines = append(lines, line)
		names = append(names, name)
	}
	lines = append(lines, "", "[Proxy Group]")
	if len(names) > 0 {
		lines = append(lines, "auto = url-test, "+strings.Join(names, ", ")+", url = https://www.gstatic.com/generate_204, interval = 300")
		lines = append(lines, "select = select, auto, "+strings.Join(names, ", "))
	}
	return strings.Join(lines, "\n") + "\n", nil
}

func surgeLine(n *node.Node, idx int) (string, string) {
	name := fmt.Sprintf("%02d-%s-%s", idx+1, n.Protocol, safe(n.Name))
	switch n.Protocol {
	case node.ProtoSS:
		return fmt.Sprintf("%s = ss, %s, %d, encrypt-method=%s, password=%s",
			name, n.Server, n.Port, n.Cipher, n.Password), name
	case node.ProtoTrojan:
		parts := []string{
			fmt.Sprintf("%s = trojan, %s, %d", name, n.Server, n.Port),
			"password=" + n.Password,
		}
		if n.SNI != "" {
			parts = append(parts, "sni="+n.SNI)
		}
		parts = append(parts, fmt.Sprintf("skip-cert-verify=%t", n.Insecure))
		return strings.Join(parts, ", "), name
	case node.ProtoVMess:
		parts := []string{
			fmt.Sprintf("%s = vmess, %s, %d", name, n.Server, n.Port),
			"username=" + n.UUID,
		}
		if n.Network == "ws" {
			parts = append(parts, "ws=true")
			if n.Path != "" {
				parts = append(parts, "ws-path="+n.Path)
			}
			if n.Host != "" {
				parts = append(parts, "ws-headers=Host:"+n.Host)
			}
		}
		if n.Security == "tls" {
			parts = append(parts, "tls=true")
			if n.SNI != "" {
				parts = append(parts, "sni="+n.SNI)
			}
			parts = append(parts, fmt.Sprintf("skip-cert-verify=%t", n.Insecure))
		}
		return strings.Join(parts, ", "), name
	}
	return "", ""
}
