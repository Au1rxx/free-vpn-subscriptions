package discovery

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// GitLabDiscoverer searches public GitLab blobs with bounded pagination.
type GitLabDiscoverer struct {
	BaseURL  string
	Client   *http.Client
	MaxPages int
}

func (d GitLabDiscoverer) Discover(ctx context.Context, seed Seed) ([]Candidate, error) {
	base := strings.TrimRight(d.BaseURL, "/")
	if base == "" {
		base = "https://gitlab.com"
	}
	pages := d.MaxPages
	if pages <= 0 || pages > 20 {
		pages = 5
	}
	token := readOptionalCredential(seed.TokenFile)
	var candidates []Candidate
	for page := 1; page <= pages; page++ {
		endpoint := base + "/api/v4/search?scope=blobs&search=" + url.QueryEscape(seed.Query) + "&per_page=100&page=" + strconv.Itoa(page)
		body, headers, err := fetch(ctx, d.Client, endpoint, token)
		if err != nil {
			return candidates, err
		}
		var results []struct {
			Path      string `json:"path"`
			Ref       string `json:"ref"`
			ProjectID int    `json:"project_id"`
		}
		if err := json.Unmarshal(body, &results); err != nil {
			return candidates, &Error{Code: "invalid_json", Err: err}
		}
		for _, item := range results {
			raw := base + "/api/v4/projects/" + strconv.Itoa(item.ProjectID) + "/repository/files/" + url.PathEscape(item.Path) + "/raw?ref=" + url.QueryEscape(item.Ref)
			candidates = append(candidates, Candidate{URL: raw, Kind: "gitlab-code", Name: item.Path, ParentURL: seed.URL, Evidence: "blob-search", Depth: seed.Depth + 1})
		}
		if headers.Get("X-Next-Page") == "" {
			break
		}
	}
	return deduplicate(candidates), nil
}
