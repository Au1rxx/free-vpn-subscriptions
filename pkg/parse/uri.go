package parse

import (
	"strings"

	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

// URIList parses a newline-separated list of proxy URIs. Empty lines and
// lines starting with '#' are treated as comments. Per-line parse failures
// are swallowed silently; only Valid() nodes are returned.
func URIList(body string) []*node.Node {
	var out []*node.Node
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		n, err := node.ParseURI(line)
		if err != nil {
			continue
		}
		if n.Valid() {
			out = append(out, n)
		}
	}
	return out
}

// Base64List decodes a whole-blob base64 body and then treats it as a URI
// list. Returns an error only if the base64 decode fails.
func Base64List(body []byte) ([]*node.Node, error) {
	decoded, err := node.B64Decode(strings.TrimSpace(string(body)))
	if err != nil {
		return nil, err
	}
	return URIList(string(decoded)), nil
}
