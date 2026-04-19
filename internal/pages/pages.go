// Package pages renders a set of static HTML pages suitable for hosting on
// GitHub Pages. The pages are the primary SEO surface — they are crawled,
// indexed, and served at au1rxx.github.io/free-vpn-subscriptions/. Each page
// includes canonical URL, Open Graph, Twitter cards, and Schema.org JSON-LD.
package pages

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/aggregate"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/node"
)

// Input matches readme.Input in shape — intentionally duplicated so packages
// stay independent and this package can evolve its own fields later.
type Input struct {
	Title         string
	RepoURL       string
	SiteURL       string // canonical origin, e.g. https://au1rxx.github.io/free-vpn-subscriptions
	Summary       aggregate.Summary
	Selected      []*node.Node
	MinPerCountry int
}

type countryRow struct {
	CC       string
	Name     string
	Flag     string
	Count    int
	URLClash string
	URLSing  string
	URLV2ray string
	URLPage  string
}

type pageCtx struct {
	// Meta
	Title        string
	Description  string
	Keywords     string
	Canonical    string
	OGImage      string
	LangAttr     string
	UpdatedHuman string
	RepoURL      string
	SiteURL      string
	HomeURL      string

	// Body
	Heading     string
	Stats       aggregate.Summary
	Countries   []countryRow
	CurrentCC   string
	CurrentName string
	CurrentFlag string
	CurrentRows int
	URLClash    string
	URLSing     string
	URLV2ray    string

	// Schema.org JSON-LD (pre-marshalled)
	JSONLD template.JS
}

// Generate writes docs/index.html, docs/{cc}.html per qualifying country,
// docs/sitemap.xml, and docs/robots.txt into outDir.
func Generate(in Input, outDir string) error {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	countries := buildCountryRows(in)
	updatedHuman := time.Unix(in.Summary.GeneratedAtUnix, 0).UTC().Format("2006-01-02 15:04 UTC")
	homeURL := in.SiteURL + "/"

	// Index page
	idx := pageCtx{
		Title:        "Free VPN Subscriptions · hourly refreshed · Clash / sing-box / v2ray",
		Description:  fmt.Sprintf("%d TCP+TLS-verified nodes from public sources, refreshed hourly. One-click Clash, sing-box, and v2ray subscription URLs. By-country filters for US, HK, JP, UK, and more.", in.Summary.TotalSelected),
		Keywords:     "free vpn, vpn subscription, clash, sing-box, v2ray, vless, reality, trojan, shadowsocks, hysteria2, proxy list, free proxy",
		Canonical:    homeURL,
		OGImage:      in.RepoURL + "/raw/main/assets/workflow.svg",
		LangAttr:     "en",
		UpdatedHuman: updatedHuman,
		RepoURL:      in.RepoURL,
		SiteURL:      in.SiteURL,
		HomeURL:      homeURL,
		Heading:      "Free VPN Subscriptions",
		Stats:        in.Summary,
		Countries:    countries,
		URLClash:     in.RepoURL + "/raw/main/output/clash.yaml",
		URLSing:      in.RepoURL + "/raw/main/output/singbox.json",
		URLV2ray:     in.RepoURL + "/raw/main/output/v2ray-base64.txt",
		JSONLD:       indexJSONLD(in, countries, updatedHuman),
	}
	if err := writeTemplate(filepath.Join(outDir, "index.html"), tplIndex, idx); err != nil {
		return err
	}

	// Per-country pages
	for _, c := range countries {
		ccLower := strings.ToLower(c.CC)
		ctx := pageCtx{
			Title:        fmt.Sprintf("Free %s VPN Subscription · %d nodes · Clash / sing-box / v2ray", c.Name, c.Count),
			Description:  fmt.Sprintf("%d TCP+TLS-verified free VPN nodes in %s, refreshed hourly. Copy a Clash, sing-box, or v2ray subscription URL and paste it into your client.", c.Count, c.Name),
			Keywords:     fmt.Sprintf("free %s vpn, %s vpn subscription, %s clash, %s v2ray, %s proxy, %s free vpn", strings.ToLower(c.Name), strings.ToLower(c.Name), strings.ToLower(c.Name), strings.ToLower(c.Name), strings.ToLower(c.Name), strings.ToLower(c.Name)),
			Canonical:    in.SiteURL + "/" + ccLower + ".html",
			OGImage:      in.RepoURL + "/raw/main/assets/workflow.svg",
			LangAttr:     "en",
			UpdatedHuman: updatedHuman,
			RepoURL:      in.RepoURL,
			SiteURL:      in.SiteURL,
			HomeURL:      homeURL,
			Heading:      fmt.Sprintf("Free %s %s VPN Subscription", c.Flag, c.Name),
			Stats:        in.Summary,
			Countries:    countries,
			CurrentCC:    c.CC,
			CurrentName:  c.Name,
			CurrentFlag:  c.Flag,
			CurrentRows:  c.Count,
			URLClash:     c.URLClash,
			URLSing:      c.URLSing,
			URLV2ray:     c.URLV2ray,
			JSONLD:       countryJSONLD(in, c, updatedHuman),
		}
		if err := writeTemplate(filepath.Join(outDir, ccLower+".html"), tplCountry, ctx); err != nil {
			return err
		}
	}

	// sitemap.xml
	if err := writeSitemap(filepath.Join(outDir, "sitemap.xml"), in.SiteURL, countries, updatedHuman); err != nil {
		return err
	}

	// robots.txt
	if err := os.WriteFile(filepath.Join(outDir, "robots.txt"),
		[]byte(fmt.Sprintf("User-agent: *\nAllow: /\n\nSitemap: %s/sitemap.xml\n", in.SiteURL)),
		0o644); err != nil {
		return err
	}

	return nil
}

func buildCountryRows(in Input) []countryRow {
	type row struct {
		cc    string
		count int
	}
	var rows []row
	for cc, n := range in.Summary.ByCountry {
		if cc == "" || cc == "XX" || n < in.MinPerCountry {
			continue
		}
		rows = append(rows, row{cc, n})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].count != rows[j].count {
			return rows[i].count > rows[j].count
		}
		return rows[i].cc < rows[j].cc
	})

	out := make([]countryRow, 0, len(rows))
	base := in.RepoURL + "/raw/main/output/by-country"
	for _, r := range rows {
		out = append(out, countryRow{
			CC:       r.cc,
			Name:     countryName(r.cc),
			Flag:     countryFlag(r.cc),
			Count:    r.count,
			URLClash: fmt.Sprintf("%s/clash-%s.yaml", base, r.cc),
			URLSing:  fmt.Sprintf("%s/singbox-%s.json", base, r.cc),
			URLV2ray: fmt.Sprintf("%s/v2ray-base64-%s.txt", base, r.cc),
			URLPage:  strings.ToLower(r.cc) + ".html",
		})
	}
	return out
}

func writeTemplate(path string, body string, ctx pageCtx) error {
	t, err := template.New(path).Parse(body)
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return t.Execute(f, ctx)
}

func writeSitemap(path, siteURL string, countries []countryRow, updatedHuman string) error {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	b.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")
	lastmod := time.Now().UTC().Format("2006-01-02")
	fmt.Fprintf(&b, "  <url><loc>%s/</loc><lastmod>%s</lastmod><changefreq>hourly</changefreq><priority>1.0</priority></url>\n", siteURL, lastmod)
	for _, c := range countries {
		fmt.Fprintf(&b, "  <url><loc>%s/%s.html</loc><lastmod>%s</lastmod><changefreq>hourly</changefreq><priority>0.8</priority></url>\n",
			siteURL, strings.ToLower(c.CC), lastmod)
	}
	b.WriteString("</urlset>\n")
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

// indexJSONLD returns the structured data graph for the landing page.
func indexJSONLD(in Input, countries []countryRow, updated string) template.JS {
	graph := []any{
		map[string]any{
			"@context":    "https://schema.org",
			"@type":       "WebSite",
			"name":        "Free VPN Subscriptions",
			"url":         in.SiteURL + "/",
			"description": fmt.Sprintf("%d hourly-refreshed free VPN nodes, TCP+TLS verified, published as Clash/sing-box/v2ray subscription URLs.", in.Summary.TotalSelected),
			"inLanguage":  "en",
		},
		map[string]any{
			"@context":    "https://schema.org",
			"@type":       "SoftwareApplication",
			"name":        "Free VPN Subscriptions",
			"operatingSystem": "Windows, macOS, iOS, Android, Linux",
			"applicationCategory": "NetworkingApplication",
			"description": "Free VPN subscription aggregator for Clash, sing-box, and v2ray. 150+ nodes tested hourly.",
			"offers": map[string]any{
				"@type":         "Offer",
				"price":         "0",
				"priceCurrency": "USD",
			},
		},
		faqSchema(),
	}
	b, _ := json.Marshal(graph)
	return template.JS(b)
}

func countryJSONLD(in Input, c countryRow, updated string) template.JS {
	graph := []any{
		map[string]any{
			"@context":    "https://schema.org",
			"@type":       "WebPage",
			"name":        fmt.Sprintf("Free %s VPN Subscription", c.Name),
			"url":         in.SiteURL + "/" + strings.ToLower(c.CC) + ".html",
			"description": fmt.Sprintf("%d free VPN nodes in %s, refreshed hourly and TCP+TLS verified.", c.Count, c.Name),
			"inLanguage":  "en",
		},
		map[string]any{
			"@context": "https://schema.org",
			"@type":    "BreadcrumbList",
			"itemListElement": []any{
				map[string]any{
					"@type":    "ListItem",
					"position": 1,
					"name":     "Home",
					"item":     in.SiteURL + "/",
				},
				map[string]any{
					"@type":    "ListItem",
					"position": 2,
					"name":     c.Name,
					"item":     in.SiteURL + "/" + strings.ToLower(c.CC) + ".html",
				},
			},
		},
	}
	b, _ := json.Marshal(graph)
	return template.JS(b)
}

func faqSchema() map[string]any {
	return map[string]any{
		"@context": "https://schema.org",
		"@type":    "FAQPage",
		"mainEntity": []any{
			qaPair("Is this actually free?",
				"Yes. Nodes are operated by third-party volunteers who publish their own free subscriptions. We don't run any servers — we just test, rank, and repackage what's already public."),
			qaPair("How fresh is the data?",
				"A GitHub Action runs every hour: pulls all upstream sources, TCP+TLS probes every node, drops anything dead, sorts by latency, and commits new output files."),
			qaPair("Can I trust these nodes?",
				"Free nodes see all your traffic. Never use them for banking, login, or anything sensitive. Fine for bypassing geo-blocks on public content."),
			qaPair("Why do some nodes fail?",
				"We verify TCP reachability and TLS handshakes, but a node may still have expired quota, bad routing, or revoked certs. Try a few; the selector group gives you fallbacks."),
		},
	}
}

func qaPair(q, a string) map[string]any {
	return map[string]any{
		"@type": "Question",
		"name":  q,
		"acceptedAnswer": map[string]any{
			"@type": "Answer",
			"text":  a,
		},
	}
}

func countryFlag(cc string) string {
	if len(cc) != 2 {
		return ""
	}
	c := strings.ToUpper(cc)
	r1 := rune(c[0]) - 'A' + 0x1F1E6
	r2 := rune(c[1]) - 'A' + 0x1F1E6
	return string([]rune{r1, r2})
}

// countryName mirrors readme.countryNames; duplicated to keep packages loose.
var names = map[string]string{
	"US": "United States", "HK": "Hong Kong", "JP": "Japan", "KR": "Korea",
	"SG": "Singapore", "TW": "Taiwan", "GB": "United Kingdom", "DE": "Germany",
	"FR": "France", "NL": "Netherlands", "CA": "Canada", "AU": "Australia",
	"RU": "Russia", "IN": "India", "TR": "Turkey", "BR": "Brazil",
	"IT": "Italy", "ES": "Spain", "SE": "Sweden", "CH": "Switzerland",
	"PL": "Poland", "AT": "Austria", "BE": "Belgium", "DK": "Denmark",
	"FI": "Finland", "NO": "Norway", "IE": "Ireland", "IL": "Israel",
	"AE": "UAE", "SA": "Saudi Arabia", "VN": "Vietnam", "TH": "Thailand",
	"MY": "Malaysia", "ID": "Indonesia", "PH": "Philippines", "MX": "Mexico",
	"AR": "Argentina", "CL": "Chile", "ZA": "South Africa", "UA": "Ukraine",
	"CZ": "Czechia", "HU": "Hungary", "RO": "Romania", "LU": "Luxembourg",
	"IS": "Iceland", "NZ": "New Zealand", "PT": "Portugal", "GR": "Greece",
	"EE": "Estonia", "LT": "Lithuania", "LV": "Latvia", "BG": "Bulgaria",
	"SK": "Slovakia", "SI": "Slovenia", "HR": "Croatia", "RS": "Serbia",
	"MD": "Moldova", "BY": "Belarus", "KZ": "Kazakhstan", "CN": "China",
	"BZ": "Belize", "IM": "Isle of Man", "CY": "Cyprus", "MT": "Malta",
	"MA": "Morocco", "SC": "Seychelles", "PA": "Panama", "PE": "Peru",
	"EC": "Ecuador", "CO": "Colombia", "MK": "North Macedonia", "PY": "Paraguay",
	"BO": "Bolivia",
}

func countryName(cc string) string {
	if v, ok := names[strings.ToUpper(cc)]; ok {
		return v
	}
	return cc
}
