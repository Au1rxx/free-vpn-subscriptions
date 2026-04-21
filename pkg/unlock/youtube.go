package unlock

import (
	"context"
	"io"
	"net/http"
	"strings"
)

const youtubePremiumURL = "https://www.youtube.com/premium"

// CheckYouTubePremium fetches the YouTube Premium landing page and
// classifies the response by body text:
//
//   - "countryCode":"XX" in the embedded config → Unlocked + region
//   - "Premium is not available in your country" → Blocked
//   - Otherwise → Unknown
func CheckYouTubePremium(ctx context.Context, client *http.Client) Result {
	r := Result{Target: "youtube-premium"}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, youtubePremiumURL, nil)
	if err != nil {
		r.Status = Unknown
		r.Detail = err.Error()
		return r
	}
	req.Header.Set("User-Agent", defaultUA)
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	resp, err := client.Do(req)
	if err != nil {
		r.Status = Unknown
		r.Detail = err.Error()
		return r
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 256*1024))
	bodyStr := string(body)

	if strings.Contains(bodyStr, "Premium is not available in your country") {
		r.Status = Blocked
		r.Detail = "landing page advertises geo-block"
		return r
	}
	if code := extractYTCountry(bodyStr); code != "" {
		r.Status = Unlocked
		r.Region = code
		r.Detail = "countryCode=" + code
		return r
	}
	r.Status = Unknown
	r.Detail = "no country signal"
	return r
}

func extractYTCountry(body string) string {
	// YouTube embeds an INNERTUBE_CONTEXT JSON blob; the country code
	// lives in `"gl":"XX"` (Google-locale) or `"countryCode":"XX"`.
	for _, key := range []string{`"gl":"`, `"countryCode":"`} {
		idx := strings.Index(body, key)
		if idx < 0 {
			continue
		}
		rest := body[idx+len(key):]
		end := strings.IndexByte(rest, '"')
		if end == 2 {
			return strings.ToUpper(rest[:end])
		}
	}
	return ""
}
