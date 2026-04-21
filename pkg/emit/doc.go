// Package emit renders []*node.Node into subscription formats consumed by
// end-user clients.
//
// Stable entry points:
//
//   - Clash       — Clash YAML (full coverage of VLESS/VMess/Trojan/SS/Hysteria2)
//   - Singbox     — sing-box JSON (full coverage)
//   - V2RayBase64 — v2rayN-style base64-wrapped URI list (full coverage)
//   - Surge       — Surge conf (partial: SS/Trojan/VMess; VLESS & Hy2 skipped)
//   - QuantumultX — QuanX server_local (partial: SS/Trojan/VMess)
//   - Loon        — Loon conf (partial: SS/Trojan/VMess)
//
// Stability
//
// This package is public but pre-1.0. The Clash/Singbox/V2RayBase64 outputs
// include an opinionated selector/url-test group; the Surge/Loon outputs add
// a matching [Proxy Group] section. Callers that need a bare proxies list
// should post-process the rendered text.
//
// Partial-coverage emitters (Surge/QuanX/Loon) silently drop nodes whose
// protocol has no mapping — inspect the returned string or count proxies in
// the output if you need to detect loss.
package emit
