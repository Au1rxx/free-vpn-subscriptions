package subscribe

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/node"
)

// Clash emits a Clash-compatible YAML profile (proxies + one URL-test group).
func Clash(nodes []*node.Node) (string, error) {
	proxies := make([]map[string]any, 0, len(nodes))
	names := make([]string, 0, len(nodes))
	for i, n := range nodes {
		p := clashProxy(n, i)
		if p == nil {
			continue
		}
		proxies = append(proxies, p)
		names = append(names, p["name"].(string))
	}

	cfg := map[string]any{
		"proxies": proxies,
		"proxy-groups": []map[string]any{
			{
				"name":     "auto",
				"type":     "url-test",
				"proxies":  names,
				"url":      "https://www.gstatic.com/generate_204",
				"interval": 300,
			},
			{
				"name":    "select",
				"type":    "select",
				"proxies": append([]string{"auto"}, names...),
			},
		},
		"rules": []string{"MATCH,select"},
	}
	b, err := yaml.Marshal(cfg)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func clashProxy(n *node.Node, idx int) map[string]any {
	name := fmt.Sprintf("%02d-%s-%s", idx+1, n.Protocol, safe(n.Name))
	p := map[string]any{
		"name":   name,
		"server": n.Server,
		"port":   n.Port,
	}
	switch n.Protocol {
	case node.ProtoVLESS:
		p["type"] = "vless"
		p["uuid"] = n.UUID
		p["network"] = or(n.Network, "tcp")
		if n.Flow != "" {
			p["flow"] = n.Flow
		}
		p["udp"] = true
		if n.Security == "reality" {
			p["tls"] = true
			p["servername"] = n.SNI
			p["client-fingerprint"] = or(n.Fingerprint, "chrome")
			p["reality-opts"] = map[string]any{
				"public-key": n.PublicKey,
				"short-id":   n.ShortID,
			}
		} else if n.Security == "tls" {
			p["tls"] = true
			p["servername"] = n.SNI
			p["skip-cert-verify"] = n.Insecure
		}

	case node.ProtoVMess:
		p["type"] = "vmess"
		p["uuid"] = n.UUID
		p["alterId"] = n.AlterID
		p["cipher"] = or(n.Cipher, "auto")
		p["network"] = or(n.Network, "tcp")
		p["udp"] = true
		if n.Security == "tls" {
			p["tls"] = true
			p["servername"] = n.SNI
			p["skip-cert-verify"] = n.Insecure
		}

	case node.ProtoTrojan:
		p["type"] = "trojan"
		p["password"] = n.Password
		p["sni"] = n.SNI
		p["skip-cert-verify"] = n.Insecure
		p["udp"] = true

	case node.ProtoSS:
		p["type"] = "ss"
		p["cipher"] = n.Cipher
		p["password"] = n.Password
		p["udp"] = true

	case node.ProtoHysteria2:
		p["type"] = "hysteria2"
		p["password"] = n.Password
		p["sni"] = n.SNI
		p["skip-cert-verify"] = n.Insecure

	default:
		return nil
	}
	return p
}

func or(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func safe(s string) string {
	// Clash doesn't love commas/newlines in names — strip aggressively.
	out := []rune{}
	for _, r := range s {
		if r == ',' || r == '\n' || r == '\r' || r == '#' {
			continue
		}
		out = append(out, r)
	}
	if len(out) > 40 {
		out = out[:40]
	}
	if len(out) == 0 {
		return "node"
	}
	return string(out)
}
