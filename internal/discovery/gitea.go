package discovery

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

// GiteaDiscoverer supports Codeberg and other public Gitea installations.
// A seed URL may point at a richer instance-specific search API.
type GiteaDiscoverer struct {
	BaseURL string
	Client  *http.Client
}

func (d GiteaDiscoverer) Discover(ctx context.Context, seed Seed) ([]Candidate, error) {
	endpoint := seed.URL
	if endpoint == "" {
		base := strings.TrimRight(d.BaseURL, "/")
		if base == "" {
			base = "https://codeberg.org"
		}
		endpoint = base + "/api/v1/repos/search?q=" + url.QueryEscape(seed.Query) + "&limit=50"
	}
	return APIDiscoverer{Client: d.Client}.Discover(ctx, Seed{URL: endpoint, TokenFile: seed.TokenFile, Depth: seed.Depth})
}
