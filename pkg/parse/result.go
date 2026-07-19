package parse

import "github.com/Au1rxx/free-vpn-subscriptions/pkg/node"

// Format identifies the subscription representation that was actually used.
type Format string

const (
	FormatAuto    Format = "auto"
	FormatURIList Format = "uri-list"
	FormatBase64  Format = "base64"
	FormatClash   Format = "clash"
	FormatSingBox Format = "singbox"
	FormatXray    Format = "xray"
)

// EntryError is a bounded, credential-safe description of one rejected item.
type EntryError struct {
	Line       int    `json:"line"`
	Code       string `json:"code"`
	Scheme     string `json:"scheme,omitempty"`
	SampleHash string `json:"sample_hash"`
	Message    string `json:"message"`
}

// Result permits partial success and carries nested source discoveries.
type Result struct {
	Format         Format
	Nodes          []*node.Node
	Errors         []EntryError
	DiscoveredURLs []string
}
