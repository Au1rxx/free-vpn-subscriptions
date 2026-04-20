// Package node defines the unified proxy node representation used across the
// aggregator pipeline. Upstream subscriptions arrive in many shapes (URI
// strings, base64 blobs, Clash YAML) and are normalized into []Node before
// probing, deduplication, and output.
package node

import "fmt"

// Protocol identifiers.
const (
	ProtoVLESS      = "vless"
	ProtoVMess      = "vmess"
	ProtoTrojan     = "trojan"
	ProtoSS         = "shadowsocks"
	ProtoHysteria2  = "hysteria2"
)

// Node is the normalized representation of a proxy endpoint. Fields are a
// superset of what each protocol needs; unused fields are simply empty.
type Node struct {
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
	Server   string `json:"server"`
	Port     int    `json:"port"`

	// Credentials — only one of UUID / Password / (Cipher+Password) is set
	// depending on protocol.
	UUID     string `json:"uuid,omitempty"`
	Password string `json:"password,omitempty"`
	Cipher   string `json:"cipher,omitempty"`
	AlterID  int    `json:"alter_id,omitempty"`

	// Transport.
	Network     string `json:"network,omitempty"`      // tcp|ws|grpc|quic
	Security    string `json:"security,omitempty"`     // none|tls|reality
	SNI         string `json:"sni,omitempty"`
	ALPN        string `json:"alpn,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
	PublicKey   string `json:"public_key,omitempty"`   // reality pbk
	ShortID     string `json:"short_id,omitempty"`     // reality sid
	SpiderX     string `json:"spider_x,omitempty"`     // reality spx
	Flow        string `json:"flow,omitempty"`         // vless flow
	Path        string `json:"path,omitempty"`         // ws path
	Host        string `json:"host,omitempty"`         // ws host header
	ServiceName string `json:"service_name,omitempty"` // grpc
	Insecure    bool   `json:"insecure,omitempty"`

	// Runtime fields populated by probe/aggregate.
	// LatencyMS holds the primary ranking latency: TCP RTT after probe.TCP,
	// then overwritten by HTTP-over-proxy median after verify.Run.
	// TCPLatencyMS preserves the raw TCP RTT for display purposes.
	LatencyMS    int    `json:"latency_ms,omitempty"`
	TCPLatencyMS int    `json:"tcp_latency_ms,omitempty"`
	Country      string `json:"country,omitempty"`
	SourceName   string `json:"source_name,omitempty"`
}

// Key returns a stable deduplication key: protocol + server + port.
// Two different credentials on the same endpoint are treated as one node —
// this is a deliberate trade-off (free sharers rotate creds, endpoint is the
// real resource).
func (n *Node) Key() string {
	return fmt.Sprintf("%s|%s|%d", n.Protocol, n.Server, n.Port)
}

// Valid checks the minimum required fields are present.
func (n *Node) Valid() bool {
	if n.Server == "" || n.Port <= 0 || n.Port > 65535 {
		return false
	}
	switch n.Protocol {
	case ProtoVLESS, ProtoVMess:
		return n.UUID != ""
	case ProtoTrojan, ProtoHysteria2:
		return n.Password != ""
	case ProtoSS:
		return n.Password != "" && n.Cipher != ""
	}
	return false
}
