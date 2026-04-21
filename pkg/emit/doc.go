// Package emit renders []*node.Node into the three canonical subscription
// formats consumed by end-user clients: Clash YAML, sing-box JSON, and the
// v2rayN-style base64-wrapped URI list.
//
// Stability
//
// This package is public but pre-1.0. Entry points (Clash, Singbox,
// V2RayBase64) are the stable surface. The output includes an opinionated
// selector/url-test group configuration; callers that need a bare proxies
// list should post-process the YAML/JSON themselves for now.
package emit
