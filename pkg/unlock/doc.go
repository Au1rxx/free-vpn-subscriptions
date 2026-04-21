// Package unlock detects whether an HTTP client (typically one whose
// transport dials through a proxy) can reach region-locked streaming
// services such as Netflix, Disney+, YouTube Premium, and ChatGPT.
//
// The package is deliberately transport-agnostic: callers pass an
// *http.Client they have already configured (e.g. with a SOCKS5 dialer
// backed by sing-box) and the Check functions only issue GETs against
// the target's probe URLs.
//
// Stability
//
// Public but pre-1.0. The stable surface is:
//
//   - Status constants (Unknown / Blocked / Partial / Unlocked)
//   - Result struct
//   - Target struct + CheckFunc type
//   - All() — default target set
//   - Run(ctx, client, targets) — sequential probe runner
//
// Individual per-target Check functions (Netflix, Disney, YouTube,
// ChatGPT) are NOT part of the stable surface — they may tighten their
// heuristics between minor releases as upstream anti-probe defenses
// evolve. Callers that need a fixed heuristic should pin a specific
// version.
package unlock
