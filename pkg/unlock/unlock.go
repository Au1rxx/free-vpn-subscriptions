package unlock

import (
	"context"
	"net/http"
	"time"
)

// Status describes the coarse outcome of a single target probe.
type Status int

const (
	// Unknown means the probe could not reach a verdict (network error,
	// unexpected response). Treat as "not unlocked" for ranking purposes
	// but preserve the Detail for debugging.
	Unknown Status = iota
	// Blocked means the service explicitly denied access from the
	// client's apparent origin (403, geo-block landing page, or the
	// service's documented "not available in your country" signal).
	Blocked
	// Partial means some content tiers are reachable but region-limited
	// catalogs are not — classic Netflix "originals-only" state.
	Partial
	// Unlocked means the service is fully reachable. Region may or may
	// not be populated depending on whether the probe can extract it.
	Unlocked
)

// String returns a short lowercase token ("unknown" / "blocked" /
// "partial" / "unlocked") suitable for tables and JSON.
func (s Status) String() string {
	switch s {
	case Blocked:
		return "blocked"
	case Partial:
		return "partial"
	case Unlocked:
		return "unlocked"
	}
	return "unknown"
}

// Result carries one target's probe outcome.
type Result struct {
	Target string `json:"target"`
	Status Status `json:"-"`
	// StatusText mirrors Status.String() for JSON consumers.
	StatusText string `json:"status"`
	// Region is an ISO 3166-1 alpha-2 country code when the probe can
	// determine it (e.g. from a CDN trace endpoint or a locale hint).
	// Empty when unknown.
	Region string `json:"region,omitempty"`
	// Detail is free-form diagnostic text — HTTP status code, body
	// fragment, or error message. Not for end-user display.
	Detail string `json:"detail,omitempty"`
}

// CheckFunc runs a single target's probe against the given client and
// returns its verdict. Implementations must respect ctx for cancellation.
type CheckFunc func(ctx context.Context, client *http.Client) Result

// Target is one streaming/service probe.
type Target struct {
	// Name is a short stable identifier ("netflix", "disney",
	// "youtube-premium", "chatgpt") used in CLI flags and JSON output.
	Name string
	// Check runs the probe. It must populate result.Target itself with
	// the Target.Name — Run will fix up StatusText after the call.
	Check CheckFunc
}

// All returns the default set of targets shipped with this package.
// The order is meaningful for table output: most-requested first.
func All() []Target {
	return []Target{
		{Name: "netflix", Check: CheckNetflix},
		{Name: "disney", Check: CheckDisney},
		{Name: "youtube-premium", Check: CheckYouTubePremium},
		{Name: "chatgpt", Check: CheckChatGPT},
	}
}

// Run probes each target sequentially against the same client and
// returns results in the input order. Each target gets its own derived
// context capped at perTarget; pass 0 to skip the cap and let the
// parent context govern.
//
// Sequential (not parallel) by design: a single proxied HTTP client is
// the common case, and parallel probes would just fight over the same
// socket. Callers who have parallel clients can fan out themselves.
func Run(ctx context.Context, client *http.Client, targets []Target, perTarget time.Duration) []Result {
	results := make([]Result, 0, len(targets))
	for _, t := range targets {
		tctx := ctx
		var cancel context.CancelFunc
		if perTarget > 0 {
			tctx, cancel = context.WithTimeout(ctx, perTarget)
		}
		r := t.Check(tctx, client)
		if cancel != nil {
			cancel()
		}
		if r.Target == "" {
			r.Target = t.Name
		}
		r.StatusText = r.Status.String()
		results = append(results, r)
	}
	return results
}
