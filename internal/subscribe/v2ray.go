package subscribe

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

// V2RayBase64 emits a newline-separated list of proxy URIs, base64-encoded as
// a whole — the standard format consumed by v2rayN / v2rayNG / Shadowrocket.
func V2RayBase64(nodes []*node.Node) string {
	var lines []string
	for _, n := range nodes {
		uri := toURI(n)
		if uri != "" {
			lines = append(lines, uri)
		}
	}
	joined := strings.Join(lines, "\n")
	return base64.StdEncoding.EncodeToString([]byte(joined))
}

// toURI re-emits a node into its canonical URI scheme.
func toURI(n *node.Node) string {
	switch n.Protocol {
	case node.ProtoVLESS:
		return vlessURI(n)
	case node.ProtoVMess:
		return vmessURI(n)
	case node.ProtoTrojan:
		return trojanURI(n)
	case node.ProtoSS:
		return ssURI(n)
	case node.ProtoHysteria2:
		return hy2URI(n)
	}
	return ""
}

func vlessURI(n *node.Node) string {
	q := url.Values{}
	setIf(q, "type", n.Network)
	setIf(q, "security", n.Security)
	setIf(q, "sni", n.SNI)
	setIf(q, "alpn", n.ALPN)
	setIf(q, "fp", n.Fingerprint)
	setIf(q, "pbk", n.PublicKey)
	setIf(q, "sid", n.ShortID)
	setIf(q, "spx", n.SpiderX)
	setIf(q, "flow", n.Flow)
	setIf(q, "path", n.Path)
	setIf(q, "host", n.Host)
	setIf(q, "serviceName", n.ServiceName)
	return fmt.Sprintf("vless://%s@%s:%d?%s#%s",
		n.UUID, n.Server, n.Port, q.Encode(), url.QueryEscape(fallbackName(n)))
}

func vmessURI(n *node.Node) string {
	obj := map[string]any{
		"v":    "2",
		"ps":   fallbackName(n),
		"add":  n.Server,
		"port": strconv.Itoa(n.Port),
		"id":   n.UUID,
		"aid":  n.AlterID,
		"net":  n.Network,
		"type": "none",
		"host": n.Host,
		"path": n.Path,
		"tls":  tlsForVmess(n),
		"sni":  n.SNI,
	}
	b, _ := json.Marshal(obj)
	return "vmess://" + base64.StdEncoding.EncodeToString(b)
}

func trojanURI(n *node.Node) string {
	q := url.Values{}
	setIf(q, "sni", n.SNI)
	setIf(q, "type", n.Network)
	if n.Insecure {
		q.Set("allowInsecure", "1")
	}
	return fmt.Sprintf("trojan://%s@%s:%d?%s#%s",
		url.QueryEscape(n.Password), n.Server, n.Port, q.Encode(), url.QueryEscape(fallbackName(n)))
}

func ssURI(n *node.Node) string {
	creds := base64.RawURLEncoding.EncodeToString([]byte(n.Cipher + ":" + n.Password))
	return fmt.Sprintf("ss://%s@%s:%d#%s",
		creds, n.Server, n.Port, url.QueryEscape(fallbackName(n)))
}

func hy2URI(n *node.Node) string {
	q := url.Values{}
	setIf(q, "sni", n.SNI)
	if n.Insecure {
		q.Set("insecure", "1")
	}
	return fmt.Sprintf("hy2://%s@%s:%d?%s#%s",
		url.QueryEscape(n.Password), n.Server, n.Port, q.Encode(), url.QueryEscape(fallbackName(n)))
}

func setIf(q url.Values, key, val string) {
	if val != "" {
		q.Set(key, val)
	}
}

func tlsForVmess(n *node.Node) string {
	if n.Security == "tls" || n.Security == "reality" {
		return n.Security
	}
	return ""
}

func fallbackName(n *node.Node) string {
	if n.Name != "" {
		return n.Name
	}
	return fmt.Sprintf("%s-%s-%d", n.Protocol, n.Server, n.Port)
}
