package pages

import (
	"math"
	"os"
	"path/filepath"
	"strconv"
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

func TestGenerateLiveReportForEveryLocale(t *testing.T) {
	in := testInput()
	now := time.Unix(in.Summary.GeneratedAtUnix, 0).UTC()
	in.Summary.TotalVerified = 260
	in.Summary.ByCountry = map[string]int{
		"US": 40, "JP": 20, "DE": 15, "FR": 12, "NL": 10,
		"SG": 9, "GB": 8, "CA": 7, "AU": 6, "XX": 99,
	}
	in.History = []HistoryPoint{
		{GeneratedAt: now.Add(-30 * 24 * time.Hour), Selected: 80, Verified: 200, MedianLatencyMS: 300, Countries: 7},
		{GeneratedAt: now.Add(-7 * 24 * time.Hour), Selected: 90, Verified: 220, MedianLatencyMS: 260, Countries: 8},
		{GeneratedAt: now.Add(-24 * time.Hour), Selected: 100, Verified: 240, MedianLatencyMS: 230, Countries: 8},
		{GeneratedAt: now, Selected: 120, Verified: 260, MedianLatencyMS: 210, Countries: 9},
	}

	dir := t.TempDir()
	if err := Generate(in, dir); err != nil {
		t.Fatal(err)
	}
	for _, locale := range supportedLocales {
		l10n := pageLocales[locale]
		name := "index" + localeSuffix(locale) + ".html"
		body, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		text := string(body)
		for _, want := range []string{
			`id="live-report"`,
			l10n.LiveHeading,
			l10n.ProtocolHeading,
			l10n.TopCountriesHeading,
			l10n.TrendHeading,
			l10n.VerificationHeading,
			l10n.VerificationText,
			l10n.LimitationsHeading,
			l10n.LimitationsText,
			l10n.SelectionHeading,
			l10n.SelectionText,
			`data-window="24h"`,
			`data-window="7d"`,
			`data-window="30d"`,
			`aria-label="30-day selected node trend"`,
			"VLESS",
			"66.7%",
		} {
			if want == "" || !strings.Contains(text, want) {
				t.Errorf("%s missing live report content %q", name, want)
			}
		}
		if strings.Index(text, in.RepoURL+"/raw/main/output/clash.yaml") > strings.Index(text, `id="live-report"`) {
			t.Errorf("%s moved the subscription card below the live report", name)
		}

		topSection := sectionByID(t, text, "top-countries")
		if got := strings.Count(topSection, `class="metric-row"`); got != 8 {
			t.Errorf("%s top country rows = %d, want 8", name, got)
		}
	}
}

func TestGenerateShowsAccumulatingStateWhenTrendHistoryIsInsufficient(t *testing.T) {
	in := testInput()
	now := time.Unix(in.Summary.GeneratedAtUnix, 0).UTC()
	in.History = []HistoryPoint{{
		GeneratedAt: now, Selected: in.Summary.TotalSelected,
		Verified: in.Summary.TotalVerified, MedianLatencyMS: in.Summary.MedianLatencyMS,
	}}
	dir := t.TempDir()
	if err := Generate(in, dir); err != nil {
		t.Fatal(err)
	}
	for _, locale := range supportedLocales {
		body, err := os.ReadFile(filepath.Join(dir, "index"+localeSuffix(locale)+".html"))
		if err != nil {
			t.Fatal(err)
		}
		if got := strings.Count(string(body), pageLocales[locale].TrendAccumulating); got != 3 {
			t.Errorf("%s accumulating labels = %d, want 3", locale, got)
		}
	}
}

func TestCountryPageShowsCountryProtocolComposition(t *testing.T) {
	in := testInput()
	in.Summary.ByCountry = map[string]int{"US": 4}
	in.Selected = []*node.Node{
		{Protocol: node.ProtoVLESS, Country: "US"},
		{Protocol: node.ProtoVLESS, Country: "US"},
		{Protocol: node.ProtoVLESS, Country: "US"},
		{Protocol: node.ProtoTrojan, Country: "US"},
	}
	dir := t.TempDir()
	if err := Generate(in, dir); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(filepath.Join(dir, "us.html"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(body)
	for _, want := range []string{
		`id="country-protocols"`,
		pageLocales["en"].CountryProtocolHeading,
		"VLESS", "75.0%", "TROJAN", "25.0%",
	} {
		if want == "" || !strings.Contains(text, want) {
			t.Errorf("country page missing %q", want)
		}
	}
}

func TestMetricRowsPercentagesAddToExactly100(t *testing.T) {
	rows := buildMetricRows(map[string]int{"vless": 1, "trojan": 1, "shadowsocks": 1})
	var total float64
	for _, row := range rows {
		value, err := strconv.ParseFloat(strings.TrimSuffix(row.Percent, "%"), 64)
		if err != nil {
			t.Fatal(err)
		}
		total += value
	}
	if math.Abs(total-100) > 0.001 {
		t.Fatalf("percent total = %.1f, want 100.0", total)
	}
	if rows[0].Name != "SHADOWSOCKS" || rows[0].Count != 1 {
		t.Fatalf("equal-count rows not sorted by name: %+v", rows)
	}
}

func TestRenderTrendSVGIsDeterministicAccessibleAndScriptFree(t *testing.T) {
	now := time.Date(2026, 7, 26, 12, 0, 0, 0, time.UTC)
	points := []HistoryPoint{
		{GeneratedAt: now.Add(-48 * time.Hour), Selected: 100},
		{GeneratedAt: now.Add(-24 * time.Hour), Selected: 80},
		{GeneratedAt: now, Selected: 120},
	}
	first := string(renderTrendSVG(points))
	second := string(renderTrendSVG(points))
	if first != second {
		t.Fatal("trend SVG is not deterministic")
	}
	for _, want := range []string{
		"<svg", "<title>30-day selected node trend</title>",
		"<desc>Selected nodes over the retained hourly history.</desc>",
		`aria-label="30-day selected node trend"`, "<polyline",
	} {
		if !strings.Contains(first, want) {
			t.Errorf("trend SVG missing %q", want)
		}
	}
	if strings.Contains(strings.ToLower(first), "<script") {
		t.Fatal("trend SVG contains script")
	}
}

func sectionByID(t *testing.T, body, id string) string {
	t.Helper()
	start := strings.Index(body, `id="`+id+`"`)
	if start < 0 {
		t.Fatalf("section %s not found", id)
	}
	end := strings.Index(body[start:], "</section>")
	if end < 0 {
		t.Fatalf("section %s has no closing tag", id)
	}
	return body[start : start+end]
}

func firstLines(s string, n int) string {
	parts := strings.SplitN(s, "\n", n+1)
	if len(parts) > n {
		parts = parts[:n]
	}
	return strings.Join(parts, "\n")
}
