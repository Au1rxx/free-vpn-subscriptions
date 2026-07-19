package discovery

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

// GiteeDiscoverer searches Gitee's public code API.
type GiteeDiscoverer struct {
	BaseURL string
	Client  *http.Client
}

func (d GiteeDiscoverer) Discover(ctx context.Context, seed Seed) ([]Candidate, error) {
	base := strings.TrimRight(d.BaseURL, "/")
	if base == "" {
		base = "https://gitee.com"
	}
	endpoint := base + "/api/v5/search/code?q=" + url.QueryEscape(seed.Query) + "&per_page=100&page=1"
	body, _, err := fetch(ctx, d.Client, endpoint, readOptionalCredential(seed.TokenFile))
	if err != nil {
		return nil, err
	}
	var items []struct {
		Name    string `json:"name"`
		HTMLURL string `json:"html_url"`
	}
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, &Error{Code: "invalid_json", Err: err}
	}
	var candidates []Candidate
	for _, item := range items {
		candidates = append(candidates, Candidate{URL: item.HTMLURL, Kind: "gitee-code", Name: item.Name, ParentURL: seed.URL, Evidence: "code-search", Depth: seed.Depth + 1})
	}
	return deduplicate(candidates), nil
}
