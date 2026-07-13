package emit

import (
	"encoding/json"
	"fmt"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/verify"
	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

// Singbox emits a sing-box-compatible outbounds JSON.
func Singbox(nodes []*node.Node) (string, error) {
	outbounds := []map[string]any{}
	endpoints := []map[string]any{}
	tags := []string{}

	for i, n := range nodes {
		ob := singboxOutbound(n, i)
		if ob == nil {
			continue
		}
		if verify.IsEndpoint(n) {
			endpoints = append(endpoints, ob)
		} else {
			outbounds = append(outbounds, ob)
		}
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
	if len(endpoints) > 0 {
		cfg["endpoints"] = endpoints
	}
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func singboxOutbound(n *node.Node, idx int) map[string]any {
	tag := fmt.Sprintf("%02d-%s-%s", idx+1, n.Protocol, safe(n.Name))
	return SingboxOutbound(n, tag)
}

// SingboxOutbound converts a single node into a sing-box outbound map
// with the caller-supplied tag, mirroring the per-node logic behind
// Singbox() but without wrapping it in a full config. Sibling tools
// that want to build their own inbound/route (e.g. a single-proxy
// launcher for ad-hoc probes) should use this and compose the rest.
//
// Returns nil for protocols without a sing-box mapping.
func SingboxOutbound(n *node.Node, tag string) map[string]any {
	outbound, err := verify.BuildOutbound(n, tag)
	if err != nil {
		return nil
	}
	return outbound
}
