package pages

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/aggregate"
	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

func testInput() Input {
	return Input{
		Title:   "Free VPN Subscriptions",
		RepoURL: "https://github.com/example/repo",
		SiteURL: "https://example.github.io/repo",
		Summary: aggregate.Summary{
			TotalAlive:      900,
			TotalSelected:   120,
			ByCountry:       map[string]int{"US": 40, "JP": 20, "XX": 99, "SG": 1},
			ByProtocol:      map[string]int{"vless": 80, "trojan": 40},
			MedianLatencyMS: 210,
			GeneratedAtUnix: time.Date(2026, 7, 25, 13, 30, 0, 0, time.UTC).Unix(),
		},
		Selected:      []*node.Node{{Protocol: node.ProtoVLESS, Server: "one.example", Port: 443, Country: "US"}},
		MinPerCountry: 3,
	}
}

func TestGenerateEmitsIndexingMetaTags(t *testing.T) {
	dir := t.TempDir()
	if err := Generate(testInput(), dir); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(filepath.Join(dir, "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`<meta name="robots" content="index,follow,max-image-preview:large,max-snippet:-1">`,
		`<meta name="theme-color" content="#0f172a">`,
		`<meta property="og:site_name" content="Free VPN Subscriptions">`,
		`<meta property="og:image:alt" content="Free VPN Subscriptions">`,
		`<meta name="twitter:image:alt" content="Free VPN Subscriptions">`,
		`<link rel="canonical" href="https://example.github.io/repo/">`,
	} {
		if !strings.Contains(string(body), want) {
			t.Errorf("index.html missing %s", want)
		}
	}
}

// Country and guide pages are separate templates, so they can drift from the
// index template's meta block.
func TestGenerateEmitsIndexingMetaOnEveryTemplate(t *testing.T) {
	dir := t.TempDir()
	if err := Generate(testInput(), dir); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"us.html", "index.zh.html", filepath.Join("guides", guides[0].Slug+".html")} {
		body, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		for _, want := range []string{`name="robots"`, `property="og:site_name"`, `name="twitter:image:alt"`} {
			if !strings.Contains(string(body), want) {
				t.Errorf("%s missing %s", name, want)
			}
		}
	}
}

func TestSitemapUsesFullDatetimeAndAlternates(t *testing.T) {
	dir := t.TempDir()
	if err := Generate(testInput(), dir); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(filepath.Join(dir, "sitemap.xml"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(body)
	if !strings.Contains(text, "<lastmod>2026-07-25T13:30:00Z</lastmod>") {
		t.Errorf("sitemap lastmod is not the full generation datetime:\n%s", firstLines(text, 8))
	}
	if !strings.Contains(text, `<xhtml:link rel="alternate" hreflang="zh-Hans"`) {
		t.Error("sitemap missing hreflang alternates")
	}
	// XX is a placeholder and SG is below MinPerCountry, so neither is a page.
	for _, absent := range []string{"/xx.html", "/sg.html"} {
		if strings.Contains(text, absent) {
			t.Errorf("sitemap should not list %s", absent)
		}
	}
}

func firstLines(s string, n int) string {
	parts := strings.SplitN(s, "\n", n+1)
	if len(parts) > n {
		parts = parts[:n]
	}
	return strings.Join(parts, "\n")
}
