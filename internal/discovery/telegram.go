package discovery

import (
	"context"
	"net/http"
	"regexp"
)

var absoluteURLPattern = regexp.MustCompile(`https?://[^\s"'<>]+`)

// TelegramDiscoverer uses only public preview pages; API credentials are not
// required. Mentioned channels and subscription links become candidates.
type TelegramDiscoverer struct{ Client *http.Client }

func (d TelegramDiscoverer) Discover(ctx context.Context, seed Seed) ([]Candidate, error) {
	body, _, err := fetch(ctx, d.Client, seed.URL, "")
	if err != nil {
		return nil, err
	}
	var candidates []Candidate
	for _, rawURL := range absoluteURLPattern.FindAllString(string(body), -1) {
		candidates = append(candidates, Candidate{URL: rawURL, Kind: "telegram-public", ParentURL: seed.URL, Evidence: "public-preview", Depth: seed.Depth + 1})
	}
	return deduplicate(candidates), nil
}
