package unlock

import (
	"context"
	"io"
	"net/http"
	"strings"
)

const (
	openaiTraceURL      = "https://chat.openai.com/cdn-cgi/trace"
	openaiComplianceURL = "https://api.openai.com/compliance/cookie_requirements"
)

// CheckChatGPT combines two probes:
//
//  1. Cloudflare trace (chat.openai.com/cdn-cgi/trace) — always
//     reachable from CF edges; `loc=XX` gives the PoP's country.
//  2. OpenAI compliance cookie endpoint — returns a body containing
//     "unsupported_country" from regions where ChatGPT is not served.
//
// A loc value alone is not sufficient: CF edges in allowed countries
// can still be blocked at the application layer. The compliance check
// is the authoritative unlock signal; loc provides region info.
func CheckChatGPT(ctx context.Context, client *http.Client) Result {
	r := Result{Target: "chatgpt"}

	// 1. trace — cheap, also gives us the region.
	if loc := openaiTrace(ctx, client); loc != "" {
		r.Region = loc
	}

	// 2. compliance — authoritative block signal.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, openaiComplianceURL, nil)
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
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 16*1024))
	bodyStr := strings.ToLower(string(body))

	switch {
	case strings.Contains(bodyStr, "unsupported_country"):
		r.Status = Blocked
		r.Detail = "api compliance: unsupported_country"
	case resp.StatusCode == http.StatusOK:
		// The endpoint returns a JSON body in supported regions; the
		// exact shape has shifted over time ("cookie_requirements" vs
		// "cookie_consent_required"). A 200 without the
		// unsupported_country marker is the reliable "unlocked" signal.
		r.Status = Unlocked
		r.Detail = "api compliance ok"
	default:
		r.Status = Unknown
		r.Detail = "http " + resp.Status
	}
	return r
}

func openaiTrace(ctx context.Context, client *http.Client) string {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, openaiTraceURL, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("User-Agent", defaultUA)
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4*1024))
	for _, line := range strings.Split(string(body), "\n") {
		if strings.HasPrefix(line, "loc=") {
			return strings.ToUpper(strings.TrimSpace(line[4:]))
		}
	}
	return ""
}
