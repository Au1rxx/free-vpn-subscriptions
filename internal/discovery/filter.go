package discovery

import (
	"net/url"
	"path"
	"strings"
)

// LikelySubscriptionURL rejects navigation and static-asset links before they
// enter the fetch scheduler. Explicit configured seeds are not subject to this
// heuristic; it only governs automatically discovered candidates.
func LikelySubscriptionURL(raw string) bool {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return false
	}
	host := strings.ToLower(parsed.Hostname())
	lowerPath := strings.ToLower(parsed.Path)
	extension := strings.ToLower(path.Ext(lowerPath))
	switch extension {
	case ".css", ".js", ".png", ".jpg", ".jpeg", ".gif", ".svg", ".ico", ".woff", ".woff2", ".map":
		return false
	case ".txt", ".yaml", ".yml", ".json", ".conf", ".list":
		return true
	}
	if host == "github.com" || strings.HasSuffix(host, "githubassets.com") {
		return false
	}
	evidence := host + lowerPath + "?" + strings.ToLower(parsed.RawQuery)
	for _, marker := range []string{
		"subscription", "/sub", "proxy", "proxies", "config", "v2ray", "vmess", "vless",
		"trojan", "clash", "singbox", "shadow", "hysteria", "tuic", "wireguard", "/api/", "feed", "nodes", "mixed",
	} {
		if strings.Contains(evidence, marker) {
			return true
		}
	}
	return false
}
