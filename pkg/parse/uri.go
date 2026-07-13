package parse

import (
	"fmt"

	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

// URIList parses a newline-separated list of proxy URIs. Empty lines and
// lines starting with '#' are treated as comments. Per-line parse failures
// are swallowed silently; only Valid() nodes are returned.
func URIList(body string) []*node.Node {
	return Parse([]byte(body), FormatURIList).Nodes
}

// Base64List decodes a whole-blob base64 body and then treats it as a URI
// list. Returns an error only if the base64 decode fails.
func Base64List(body []byte) ([]*node.Node, error) {
	result := Parse(body, FormatBase64)
	if len(result.Errors) > 0 && len(result.Nodes) == 0 {
		return nil, fmt.Errorf("%s", result.Errors[0].Message)
	}
	return result.Nodes, nil
}
