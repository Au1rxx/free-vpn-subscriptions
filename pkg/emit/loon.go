package emit

import (
	"fmt"
	"strings"

	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

// Loon emits a Loon-compatible [Proxy] section plus a [Proxy Group] block.
// Coverage is partial: Shadowsocks, Trojan, and VMess are supported;
// VLESS and Hysteria2 are skipped for this MVP.
func Loon(nodes []*node.Node) (string, error) {
	var lines []string
	var names []string
	lines = append(lines, "[Proxy]")
	for i, n := range nodes {
		line, name := loonLine(n, i)
		if line == "" {
			continue
		}
		lines = append(lines, line)
		names = append(names, name)
	}
	lines = append(lines, "", "[Proxy Group]")
	if len(names) > 0 {
		lines = append(lines, "auto = url-test,"+strings.Join(names, ",")+",url=https://www.gstatic.com/generate_204,interval=300")
		lines = append(lines, "select = select,auto,"+strings.Join(names, ","))
	}
	return strings.Join(lines, "\n") + "\n", nil
}

func loonLine(n *node.Node, idx int) (string, string) {
	name := fmt.Sprintf("%02d-%s-%s", idx+1, n.Protocol, safe(n.Name))
	switch n.Protocol {
	case node.ProtoSS:
		return fmt.Sprintf("%s = Shadowsocks,%s,%d,%s,%q",
			name, n.Server, n.Port, n.Cipher, n.Password), name
	case node.ProtoTrojan:
		parts := []string{
			fmt.Sprintf("%s = trojan,%s,%d,%q", name, n.Server, n.Port, n.Password),
		}
		if n.SNI != "" {
			parts = append(parts, "tls-name="+n.SNI)
		}
		parts = append(parts, fmt.Sprintf("skip-cert-verify=%t", n.Insecure))
		return strings.Join(parts, ","), name
	case node.ProtoVMess:
		cipher := n.Cipher
		if cipher == "" {
			cipher = "auto"
		}
		parts := []string{
			fmt.Sprintf("%s = vmess,%s,%d,%s,%q", name, n.Server, n.Port, cipher, n.UUID),
		}
		if n.Network == "ws" {
			parts = append(parts, "transport:ws")
			if n.Path != "" {
				parts = append(parts, "path="+n.Path)
			}
			if n.Host != "" {
				parts = append(parts, "host="+n.Host)
			}
		}
		if n.Security == "tls" {
			parts = append(parts, "over-tls=true")
			if n.SNI != "" {
				parts = append(parts, "tls-name="+n.SNI)
			}
			parts = append(parts, fmt.Sprintf("skip-cert-verify=%t", n.Insecure))
		}
		return strings.Join(parts, ","), name
	}
	return "", ""
}
