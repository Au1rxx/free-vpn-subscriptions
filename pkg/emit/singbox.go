package emit

import (
	"encoding/json"
	"fmt"

	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

// Singbox emits a sing-box-compatible outbounds JSON.
func Singbox(nodes []*node.Node) (string, error) {
	outbounds := []map[string]any{}
	tags := []string{}

	for i, n := range nodes {
		ob := singboxOutbound(n, i)
		if ob == nil {
			continue
		}
		outbounds = append(outbounds, ob)
		tags = append(tags, ob["tag"].(string))
	}

	// selector + urltest groups
	outbounds = append(outbounds, map[string]any{
		"type":      "urltest",
		"tag":       "auto",
		"outbounds": tags,
		"url":       "https://www.gstatic.com/generate_204",
		"interval":  "5m",
	})
	outbounds = append(outbounds, map[string]any{
		"type":      "selector",
		"tag":       "select",
		"outbounds": append([]string{"auto"}, tags...),
	})
	outbounds = append(outbounds, map[string]any{"type": "direct", "tag": "direct"})

	cfg := map[string]any{
		"outbounds": outbounds,
		"route": map[string]any{
			"final": "select",
		},
	}
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func singboxOutbound(n *node.Node, idx int) map[string]any {
	tag := fmt.Sprintf("%02d-%s-%s", idx+1, n.Protocol, safe(n.Name))
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
		if n.Security == "reality" || n.Security == "tls" {
			tls := map[string]any{
				"enabled":     true,
				"server_name": n.SNI,
				"insecure":    n.Insecure,
			}
			if n.Security == "reality" {
				tls["reality"] = map[string]any{
					"enabled":    true,
					"public_key": n.PublicKey,
					"short_id":   n.ShortID,
				}
				tls["utls"] = map[string]any{
					"enabled":     true,
					"fingerprint": or(n.Fingerprint, "chrome"),
				}
			}
			ob["tls"] = tls
		}

	case node.ProtoVMess:
		ob["type"] = "vmess"
		ob["uuid"] = n.UUID
		ob["alter_id"] = n.AlterID
		ob["security"] = or(n.Cipher, "auto")
		if n.Security == "tls" {
			ob["tls"] = map[string]any{
				"enabled":     true,
				"server_name": n.SNI,
				"insecure":    n.Insecure,
			}
		}

	case node.ProtoTrojan:
		ob["type"] = "trojan"
		ob["password"] = n.Password
		ob["tls"] = map[string]any{
			"enabled":     true,
			"server_name": n.SNI,
			"insecure":    n.Insecure,
		}

	case node.ProtoSS:
		ob["type"] = "shadowsocks"
		ob["method"] = n.Cipher
		ob["password"] = n.Password

	case node.ProtoHysteria2:
		ob["type"] = "hysteria2"
		ob["password"] = n.Password
		ob["tls"] = map[string]any{
			"enabled":     true,
			"server_name": n.SNI,
			"insecure":    n.Insecure,
		}

	default:
		return nil
	}
	return ob
}
