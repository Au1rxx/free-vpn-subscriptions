package unlock

import (
	"context"
	"net/http"
	"strings"
)

// Netflix title IDs used as canaries. These are well-known in the
// open-source unlock-detection community and historically stable:
//
//   - 81280792 Squid Game (Netflix original — available wherever
//     Netflix itself is available)
//   - 70143836 Friends (licensed content — region-gated)
const (
	netflixOriginal = "https://www.netflix.com/title/81280792"
	netflixLicensed = "https://www.netflix.com/title/70143836"
)

// CheckNetflix hits two title URLs: a Netflix-original and a
// region-licensed title. The combination of responses tells us Blocked
// vs Partial (originals-only) vs Unlocked.
//
// This mirrors the well-known open-source detection pattern. It does
// NOT attempt to extract the user's Netflix region — that requires
// logging in or scraping more pages and is out of scope for a probe.
func CheckNetflix(ctx context.Context, client *http.Client) Result {
	r := Result{Target: "netflix"}
	origStatus, origErr := netflixHit(ctx, client, netflixOriginal)
	licStatus, licErr := netflixHit(ctx, client, netflixLicensed)

	if origErr != nil && licErr != nil {
		r.Status = Unknown
		r.Detail = "network error: " + origErr.Error()
		return r
	}

	switch {
	case origStatus == http.StatusNotFound && licStatus == http.StatusNotFound:
		// Netflix returns 404 for blocked regions on both titles.
		r.Status = Blocked
		r.Detail = "both titles 404 (geo-block)"
	case licStatus == http.StatusOK:
		r.Status = Unlocked
		r.Detail = "licensed title reachable"
	case origStatus == http.StatusOK && licStatus == http.StatusNotFound:
		r.Status = Partial
		r.Detail = "originals only"
	default:
		r.Status = Unknown
		r.Detail = "unexpected response"
	}
	return r
}

func netflixHit(ctx context.Context, client *http.Client, url string) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("User-Agent", defaultUA)
	// Netflix sometimes rewrites based on Accept-Language; leave it
	// blank so CDN edge defaults to IP-geo.
	req.Header.Set("Accept-Language", "")
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	_ = drain(resp)
	// Some mirrors return 3xx to a region error page; count those as
	// "not OK" by collapsing to 404-equivalent.
	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		if strings.Contains(resp.Header.Get("Location"), "NotAvailable") {
			return http.StatusNotFound, nil
		}
	}
	return resp.StatusCode, nil
}
