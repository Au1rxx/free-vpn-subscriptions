package discovery

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGitHubDiscoveryPaginatesAndReportsRateLimit(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.URL.Query().Get("q") == "limited" {
			w.Header().Set("Retry-After", "30")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		page := r.URL.Query().Get("page")
		if page == "1" {
			w.Header().Set("Link", fmt.Sprintf(`<%s/search/code?page=2>; rel="next"`, serverURL(r)))
			_, _ = fmt.Fprint(w, `{"items":[{"name":"a.txt","html_url":"https://github.com/o/r/blob/a.txt","repository":{"full_name":"o/r"}}]}`)
			return
		}
		_, _ = fmt.Fprint(w, `{"items":[{"name":"b.txt","html_url":"https://github.com/o/r/blob/b.txt","repository":{"full_name":"o/r"}}]}`)
	}))
	defer server.Close()
	discoverer := GitHubDiscoverer{BaseURL: server.URL, Client: server.Client(), MaxPages: 3}
	candidates, err := discoverer.Discover(context.Background(), Seed{Query: "vless://"})
	if err != nil || len(candidates) != 2 || requests != 2 {
		t.Fatalf("candidates=%+v requests=%d err=%v", candidates, requests, err)
	}
	_, err = discoverer.Discover(context.Background(), Seed{Query: "limited"})
	if typed, ok := err.(*Error); !ok || typed.Code != "rate_limited" || typed.RetryAfterSeconds != 30 {
		t.Fatalf("unexpected rate limit: %#v", err)
	}
}

func TestTelegramAndSitemapDiscovery(t *testing.T) {
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "sitemap") {
			w.Header().Set("Content-Type", "application/xml")
			_, _ = fmt.Fprintf(w, `<urlset><url><loc>%s/sub/a.txt</loc></url><url><loc>%s/sub/b.yaml</loc></url></urlset>`, server.URL, server.URL)
			return
		}
		_, _ = fmt.Fprint(w, `<a href="https://t.me/another_channel">mention</a><a href="https://example.com/sub.txt">subscription</a>`)
	}))
	defer server.Close()
	telegram := TelegramDiscoverer{Client: server.Client()}
	candidates, err := telegram.Discover(context.Background(), Seed{URL: server.URL + "/channel"})
	if err != nil || len(candidates) != 2 {
		t.Fatalf("telegram candidates=%+v err=%v", candidates, err)
	}
	web := WebDiscoverer{Client: server.Client()}
	candidates, err = web.Discover(context.Background(), Seed{URL: server.URL + "/sitemap.xml"})
	if err != nil || len(candidates) != 2 {
		t.Fatalf("sitemap candidates=%+v err=%v", candidates, err)
	}
}

func serverURL(r *http.Request) string {
	return "http://" + r.Host
}
