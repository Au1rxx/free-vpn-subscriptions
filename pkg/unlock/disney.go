package unlock

import (
	"context"
	"io"
	"net/http"
	"strings"
)

const disneyHome = "https://www.disneyplus.com/"

// CheckDisney requests the Disney+ homepage and inspects the response
// body for region code hints. Disney+ serves a "not available" landing
// page in blocked regions and a localized portal otherwise.
//
// Region extraction reads the <html lang="xx-YY"> attribute as a best
// effort; it is often absent and Region stays empty. This is fine —
// the probe's primary job is the unlock verdict, not regional routing.
func CheckDisney(ctx context.Context, client *http.Client) Result {
	r := Result{Target: "disney"}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, disneyHome, nil)
	if err != nil {
		r.Status = Unknown
		r.Detail = err.Error()
		return r
	}
	req.Header.Set("User-Agent", defaultUA)
	resp, err := client.Do(req)
	if err != nil {
		r.Status = Unknown
		r.Detail = err.Error()
		return r
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
	bodyStr := string(body)

	lower := strings.ToLower(bodyStr)
	switch {
	case strings.Contains(lower, "not available in your region"),
		strings.Contains(lower, "unsupported_location"),
		strings.Contains(lower, "region is not supported"):
		r.Status = Blocked
		r.Detail = "landing page advertises geo-block"
	case resp.StatusCode == http.StatusOK:
		r.Status = Unlocked
		r.Region = extractHTMLLang(bodyStr)
		r.Detail = "homepage reachable"
	default:
		r.Status = Unknown
		r.Detail = "http " + resp.Status
	}
	return r
}

func extractHTMLLang(body string) string {
	// Very tolerant: find `lang="xx-YY"` or `lang='xx-YY'` and return YY.
	for _, marker := range []string{`lang="`, `lang='`} {
		idx := strings.Index(body, marker)
		if idx < 0 {
			continue
		}
		rest := body[idx+len(marker):]
		end := strings.IndexAny(rest, `"'`)
		if end < 0 {
			continue
		}
		tag := rest[:end]
		parts := strings.Split(tag, "-")
		if len(parts) == 2 && len(parts[1]) == 2 {
			return strings.ToUpper(parts[1])
		}
	}
	return ""
}
