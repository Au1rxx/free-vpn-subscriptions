package discovery

import (
	"context"
	"html"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var (
	hrefPattern = regexp.MustCompile(`(?i)(?:href|src)=["']([^"']+)["']`)
	locPattern  = regexp.MustCompile(`(?is)<loc>\s*([^<]+?)\s*</loc>`)
)

// WebDiscoverer extracts candidates from HTML, RSS/Atom and Sitemap XML.
type WebDiscoverer struct{ Client *http.Client }

func (d WebDiscoverer) Discover(ctx context.Context, seed Seed) ([]Candidate, error) {
	body, _, err := fetch(ctx, d.Client, seed.URL, "")
	if err != nil {
		return nil, err
	}
	base, _ := url.Parse(seed.URL)
	var values []string
	for _, match := range locPattern.FindAllSubmatch(body, -1) {
		values = append(values, string(match[1]))
	}
	for _, match := range hrefPattern.FindAllSubmatch(body, -1) {
		values = append(values, string(match[1]))
	}
	var candidates []Candidate
	for _, value := range values {
		value = html.UnescapeString(strings.TrimSpace(value))
		parsed, err := url.Parse(value)
		if err != nil {
			continue
		}
		if !parsed.IsAbs() && base != nil {
			parsed = base.ResolveReference(parsed)
		}
		candidates = append(candidates, Candidate{URL: parsed.String(), Kind: "web", ParentURL: seed.URL, Evidence: "page-link", Depth: seed.Depth + 1})
	}
	return deduplicate(candidates), nil
}
