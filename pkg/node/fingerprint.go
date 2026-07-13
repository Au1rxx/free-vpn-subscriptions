package node

import (
	"crypto/sha256"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
)

// CanonicalConfig contains every setting that changes connection semantics.
// Display labels, source attribution and runtime measurements are excluded.
type CanonicalConfig struct {
	Protocol    string            `json:"protocol"`
	Server      string            `json:"server"`
	Port        int               `json:"port"`
	UUID        string            `json:"uuid,omitempty"`
	Username    string            `json:"username,omitempty"`
	Password    string            `json:"password,omitempty"`
	Cipher      string            `json:"cipher,omitempty"`
	AlterID     int               `json:"alter_id,omitempty"`
	Network     string            `json:"network"`
	Security    string            `json:"security"`
	SNI         string            `json:"sni,omitempty"`
	ALPN        string            `json:"alpn,omitempty"`
	Fingerprint string            `json:"fingerprint,omitempty"`
	PublicKey   string            `json:"public_key,omitempty"`
	ShortID     string            `json:"short_id,omitempty"`
	SpiderX     string            `json:"spider_x,omitempty"`
	Flow        string            `json:"flow,omitempty"`
	Path        string            `json:"path,omitempty"`
	Host        string            `json:"host,omitempty"`
	ServiceName string            `json:"service_name,omitempty"`
	Insecure    bool              `json:"insecure"`
	Extra       map[string]string `json:"extra,omitempty"`
}

// CanonicalJSON serializes a normalized, stable representation suitable for
// persistent identity and cross-run deduplication.
func (n *Node) CanonicalJSON() ([]byte, error) {
	config := CanonicalConfig{
		Protocol: strings.ToLower(strings.TrimSpace(n.Protocol)),
		Server:   normalizeHost(n.Server), Port: n.Port,
		UUID: n.UUID, Username: n.Username, Password: n.Password,
		Cipher: strings.ToLower(strings.TrimSpace(n.Cipher)), AlterID: n.AlterID,
		Network:  strings.ToLower(strings.TrimSpace(n.Network)),
		Security: strings.ToLower(strings.TrimSpace(n.Security)),
		SNI:      normalizeHost(n.SNI), ALPN: normalizeALPN(n.ALPN),
		Fingerprint: strings.ToLower(strings.TrimSpace(n.Fingerprint)),
		PublicKey:   n.PublicKey, ShortID: n.ShortID, SpiderX: n.SpiderX,
		Flow: n.Flow, Path: n.Path, Host: normalizeHost(n.Host),
		ServiceName: n.ServiceName, Insecure: n.Insecure,
	}
	if config.Network == "" {
		config.Network = "tcp"
	}
	if config.Security == "" {
		config.Security = "none"
	}
	if len(n.Extra) > 0 {
		config.Extra = make(map[string]string, len(n.Extra))
		for key, value := range n.Extra {
			config.Extra[strings.ToLower(strings.TrimSpace(key))] = strings.TrimSpace(value)
		}
	}
	return json.Marshal(config)
}

// EndpointFingerprint identifies a protocol endpoint independently of its
// credentials and transport settings.
func (n *Node) EndpointFingerprint() [32]byte {
	identity := strings.ToLower(strings.TrimSpace(n.Protocol)) + "\x00" + normalizeHost(n.Server) + "\x00" + strconv.Itoa(n.Port)
	return sha256.Sum256([]byte(identity))
}

// ConfigFingerprint identifies one complete usable configuration.
func (n *Node) ConfigFingerprint() [32]byte {
	body, err := n.CanonicalJSON()
	if err != nil {
		return [32]byte{}
	}
	return sha256.Sum256(body)
}

func normalizeHost(value string) string {
	return strings.TrimSuffix(strings.ToLower(strings.TrimSpace(value)), ".")
}

func normalizeALPN(value string) string {
	if strings.TrimSpace(value) == "" {
		return ""
	}
	items := strings.Split(value, ",")
	for i := range items {
		items[i] = strings.ToLower(strings.TrimSpace(items[i]))
	}
	sort.Strings(items)
	return strings.Join(items, ",")
}
