package discovery

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// GitHubDiscoverer searches GitHub/Gist code with bounded pagination.
type GitHubDiscoverer struct {
	BaseURL  string
	Client   *http.Client
	MaxPages int
}

func (d GitHubDiscoverer) Discover(ctx context.Context, seed Seed) ([]Candidate, error) {
	base := strings.TrimRight(d.BaseURL, "/")
	if base == "" {
		base = "https://api.github.com"
	}
	pages := d.MaxPages
	if pages <= 0 || pages > 20 {
		pages = 5
	}
	token := readOptionalCredential(seed.TokenFile)
	var candidates []Candidate
	for page := 1; page <= pages; page++ {
		endpoint := base + "/search/code?q=" + url.QueryEscape(seed.Query) + "&per_page=100&page=" + strconv.Itoa(page)
		body, headers, err := fetch(ctx, d.Client, endpoint, token)
		if err != nil {
			return candidates, err
		}
		var response struct {
			Items []struct {
				Name        string `json:"name"`
				HTMLURL     string `json:"html_url"`
				DownloadURL string `json:"download_url"`
				Repository  struct {
					FullName string `json:"full_name"`
				} `json:"repository"`
			} `json:"items"`
		}
		if err := json.Unmarshal(body, &response); err != nil {
			return candidates, &Error{Code: "invalid_json", Err: err}
		}
		for _, item := range response.Items {
			candidateURL := item.DownloadURL
			if candidateURL == "" {
				candidateURL = githubRawURL(item.HTMLURL)
			}
			candidates = append(candidates, Candidate{URL: candidateURL, Kind: "github-code", Name: item.Name, ParentURL: seed.URL, Evidence: item.Repository.FullName, Depth: seed.Depth + 1})
		}
		if !strings.Contains(headers.Get("Link"), `rel="next"`) {
			break
		}
	}
	return deduplicate(candidates), nil
}

func githubRawURL(value string) string {
	parsed, err := url.Parse(value)
	if err != nil || parsed.Host != "github.com" {
		return value
	}
	parts := strings.Split(strings.TrimPrefix(parsed.Path, "/"), "/")
	if len(parts) >= 5 && parts[2] == "blob" {
		return "https://raw.githubusercontent.com/" + parts[0] + "/" + parts[1] + "/" + strings.Join(parts[3:], "/")
	}
	return value
}
