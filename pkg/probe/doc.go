// Package probe measures the reachability of proxy endpoints with
// lightweight handshake-only probes: concurrent TCP dials and, for nodes
// that advertise TLS, a TLS ClientHello.
//
// Scope vs. pkg/verify (internal to the aggregator): probe only confirms
// that the endpoint speaks the expected protocol at the socket/TLS layer.
// It cannot detect expired credentials, misrouted exits, or GFW-poisoned
// responses. For those, run an HTTP-over-proxy verify as a second stage.
//
// Stability
//
// This package is public but pre-1.0. TCP and TLS are the stable entry
// points. Internal helpers may change without notice. InsecureSkipVerify
// is hardcoded true in the TLS prober because free proxies routinely
// present self-signed certs and we only care that the peer speaks TLS
// at all.
package probe
