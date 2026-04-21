// Package parse converts raw subscription bodies (URI lists, base64 blobs,
// Clash YAML) into normalized []*node.Node. It is the format-recognition
// layer: callers are responsible for fetching bytes over the network and for
// picking the right entry point based on the source's declared format.
//
// Stability
//
// This package is public but pre-1.0. Entry points (URIList, Base64List,
// Clash) are the stable surface. Per-entry parse errors are swallowed: the
// returned slice contains only successfully parsed, Valid() nodes. Hard
// failures (invalid base64, invalid YAML) surface as error returns.
package parse
