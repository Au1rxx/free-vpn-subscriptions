// Package pages renders a set of static HTML pages suitable for hosting on
// GitHub Pages. The pages are the primary SEO surface — they are crawled,
// indexed, and served at au1rxx.github.io/free-vpn-subscriptions/. Each page
// includes canonical URL, Open Graph, Twitter cards, and Schema.org JSON-LD.
//
// Multilingual: each page is rendered once per locale in supportedLocales.
// English is canonical at {name}.html; other locales use {name}.{loc}.html.
// All pages cross-reference via <link rel="alternate" hreflang="..."> and a
// user-facing language switcher at the top.
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
	URLPage  string // locale-specific page URL
}

// langAlt is used for <link rel="alternate" hreflang=...> tags.
type langAlt struct {
	Code string // hreflang value, e.g. "en", "zh-Hans", "x-default"
	URL  string
}

// langSwitch renders the visible language switcher row at the top of each page.
type langSwitch struct {
	Label   string
	URL     string
	Current bool
}

// guideCtx is the render context for a single guide page. Body fields are
// template.HTML so author-controlled HTML survives html/template's escaping.
type guideCtx struct {
	Title        string
	Description  string
	Keywords     string
	Canonical    string
	OGImage      string
	LangAttr     string
	Alternates   []langAlt
	LanguageSw   []langSwitch
	UpdatedHuman string
	HomeURL      string
	RepoURL      string
	Heading      string
	ClientName   string
	OSList       string
	DownloadURL  string
	SubscribeURL string
	L10n         guideL10n
	Steps        []renderedStep
	Tips         []renderedTip
	OtherGuides  []guideLink
	JSONLD       template.JS
}

type renderedStep struct {
	Title string
	Body  template.HTML
}

type renderedTip struct {
	Q string
	A template.HTML
}

type guideLink struct {
	URL  string
	Name string
	OS   string
}

type pageCtx struct {
	// Meta
	Title       string
	Description string
	Keywords    string
	Canonical   string
	OGImage     string
	LangAttr    string
	Alternates  []langAlt
	LanguageSw  []langSwitch

	// Locale strings
	L10n pageL10n

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
	// CurrentSubscribeHeading is the localized "Subscribe to FRANCE nodes only"
	// string, pre-formatted so the template doesn't need printf helpers.
	CurrentSubscribeHeading string
	URLClash                string
	URLSing                 string
	URLV2ray                string

	// Guides (shown only on index page)
	Guides []guideLink

	// Schema.org JSON-LD (pre-marshalled)
	JSONLD template.JS
}

// Generate writes docs/{index,cc}.html per qualifying country per locale,
// docs/guides/{slug}.html per guide per locale, docs/sitemap.xml, and
// docs/robots.txt into outDir.
func Generate(in Input, outDir string) error {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	updatedHuman := time.Unix(in.Summary.GeneratedAtUnix, 0).UTC().Format("2006-01-02 15:04 UTC")

	for _, loc := range supportedLocales {
		suffix := localeSuffix(loc)
		l10n := pageLocales[loc]
		homeURL := in.SiteURL + "/"
		if suffix != "" {
			homeURL = in.SiteURL + "/index" + suffix + ".html"
		}
		countries := buildCountryRows(in, suffix)

		// Index page
		idxTitle := fmt.Sprintf(l10n.IndexTitleTpl, in.Summary.TotalSelected)
		idxDesc := fmt.Sprintf(l10n.IndexDescriptionTpl, in.Summary.TotalSelected)
		idxCanonical := homeURL

		// Guide links fall back to English if the current locale doesn't
		// have translated guide content — better than a 404, and Google
		// follows <link rel=alternate hreflang="en"> anyway.
		guideSuffix := suffix
		if !hasGuideLocale(loc) {
			guideSuffix = ""
		}
		guideLinks := make([]guideLink, 0, len(guides))
		for _, g := range guides {
			guideLinks = append(guideLinks, guideLink{
				URL:  "guides/" + g.Slug + guideSuffix + ".html",
				Name: g.ClientName,
				OS:   g.OSList,
			})
		}

		idx := pageCtx{
			Title:        idxTitle,
			Description:  idxDesc,
			Keywords:     l10n.IndexKeywords,
			Canonical:    idxCanonical,
			OGImage:      in.RepoURL + "/raw/main/assets/hero.png",
			LangAttr:     l10n.LangAttr,
			Alternates:   indexAlternates(in.SiteURL),
			LanguageSw:   indexLangSwitcher(in.SiteURL, loc),
			L10n:         l10n,
			UpdatedHuman: updatedHuman,
			RepoURL:      in.RepoURL,
			SiteURL:      in.SiteURL,
			HomeURL:      homeURL,
			Heading:      l10n.IndexHeading,
			Stats:        in.Summary,
			Countries:    countries,
			URLClash:     in.RepoURL + "/raw/main/output/clash.yaml",
			URLSing:      in.RepoURL + "/raw/main/output/singbox.json",
			URLV2ray:     in.RepoURL + "/raw/main/output/v2ray-base64.txt",
			Guides:       guideLinks,
			JSONLD:       indexJSONLD(in, l10n, idxCanonical, loc),
		}
		idxFile := "index" + suffix + ".html"
		if err := writeTemplate(filepath.Join(outDir, idxFile), tplIndex, idx); err != nil {
			return err
		}

		// Per-country pages
		for _, c := range countries {
			ccLower := strings.ToLower(c.CC)
			canonical := in.SiteURL + "/" + ccLower + suffix + ".html"
			nameLower := strings.ToLower(c.Name)
			ctx := pageCtx{
				Title:                   fmt.Sprintf(l10n.CountryTitleTpl, c.Name, c.Count),
				Description:             fmt.Sprintf(l10n.CountryDescriptionTpl, c.Count, c.Name),
				Keywords:                fmt.Sprintf(l10n.CountryKeywordsTpl, nameLower, nameLower, nameLower, nameLower, nameLower, nameLower),
				Canonical:               canonical,
				OGImage:                 in.RepoURL + "/raw/main/assets/hero.png",
				LangAttr:                l10n.LangAttr,
				Alternates:              countryAlternates(in.SiteURL, ccLower),
				LanguageSw:              countryLangSwitcher(in.SiteURL, ccLower, loc),
				L10n:                    l10n,
				UpdatedHuman:            updatedHuman,
				RepoURL:                 in.RepoURL,
				SiteURL:                 in.SiteURL,
				HomeURL:                 homeURL,
				Heading:                 fmt.Sprintf(l10n.CountryHeadingTpl, c.Flag, c.Name),
				Stats:                   in.Summary,
				Countries:               countries,
				CurrentCC:               c.CC,
				CurrentName:             c.Name,
				CurrentFlag:             c.Flag,
				CurrentRows:             c.Count,
				CurrentSubscribeHeading: fmt.Sprintf(l10n.CountrySubscribeHeadingTpl, c.Name),
				URLClash:                c.URLClash,
				URLSing:                 c.URLSing,
				URLV2ray:                c.URLV2ray,
				JSONLD:                  countryJSONLD(in, l10n, c, canonical, loc),
			}
			countryFile := ccLower + suffix + ".html"
			if err := writeTemplate(filepath.Join(outDir, countryFile), tplCountry, ctx); err != nil {
				return err
			}
		}

	}

	// Guide pages: render only in locales with translated step content.
	// Other locales' index pages link to the English fallback above.
	guidesDir := filepath.Join(outDir, "guides")
	if err := os.MkdirAll(guidesDir, 0o755); err != nil {
		return err
	}
	for _, loc := range supportedGuideLocales {
		suffix := localeSuffix(loc)
		homeURL := in.SiteURL + "/"
		if suffix != "" {
			homeURL = in.SiteURL + "/index" + suffix + ".html"
		}
		for _, g := range guides {
			gctx := buildGuideCtx(in, g, loc, updatedHuman, homeURL)
			fileName := g.Slug + suffix + ".html"
			if err := writeTemplate(filepath.Join(guidesDir, fileName), tplGuide, gctx); err != nil {
				return fmt.Errorf("guide %s (%s): %w", g.Slug, loc, err)
			}
		}
	}

	// sitemap.xml + robots.txt (locale-agnostic)
	countriesEN := buildCountryRows(in, "")
	if err := writeSitemap(filepath.Join(outDir, "sitemap.xml"), in.SiteURL, countriesEN); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(outDir, "robots.txt"),
		[]byte(fmt.Sprintf("User-agent: *\nAllow: /\n\nSitemap: %s/sitemap.xml\n", in.SiteURL)),
		0o644); err != nil {
		return err
	}

	return nil
}

func buildCountryRows(in Input, localeSuffix string) []countryRow {
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
			URLPage:  strings.ToLower(r.cc) + localeSuffix + ".html",
		})
	}
	return out
}

func writeTemplate(path string, body string, ctx any) error {
	funcs := template.FuncMap{
		"safe": func(s string) template.HTML { return template.HTML(s) },
	}
	t, err := template.New(path).Funcs(funcs).Parse(body)
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

func writeSitemap(path, siteURL string, countries []countryRow) error {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	b.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xhtml="http://www.w3.org/1999/xhtml">` + "\n")
	lastmod := time.Now().UTC().Format("2006-01-02")

	// Home
	writeSitemapEntry(&b, siteURL+"/", lastmod, "hourly", "1.0", indexURLsByLocale(siteURL))

	// Country pages
	for _, c := range countries {
		ccLower := strings.ToLower(c.CC)
		baseURL := siteURL + "/" + ccLower + ".html"
		writeSitemapEntry(&b, baseURL, lastmod, "hourly", "0.8", countryURLsByLocale(siteURL, ccLower))
	}

	// Guide pages
	for _, g := range guides {
		baseURL := siteURL + "/guides/" + g.Slug + ".html"
		writeSitemapEntry(&b, baseURL, lastmod, "weekly", "0.7", guideURLsByLocale(siteURL, g.Slug))
	}

	b.WriteString("</urlset>\n")
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func writeSitemapEntry(b *strings.Builder, loc, lastmod, changefreq, priority string, alternates map[string]string) {
	fmt.Fprintf(b, "  <url>\n    <loc>%s</loc>\n    <lastmod>%s</lastmod>\n    <changefreq>%s</changefreq>\n    <priority>%s</priority>\n", loc, lastmod, changefreq, priority)
	codes := sortedKeys(alternates)
	for _, code := range codes {
		fmt.Fprintf(b, "    <xhtml:link rel=\"alternate\" hreflang=\"%s\" href=\"%s\"/>\n", code, alternates[code])
	}
	b.WriteString("  </url>\n")
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// indexAlternates returns the full set of <link rel=alternate hreflang>
// tags for the index page in each supported locale.
func indexAlternates(siteURL string) []langAlt {
	alts := make([]langAlt, 0, len(supportedLocales)+1)
	for _, loc := range supportedLocales {
		url := siteURL + "/"
		if loc != "en" {
			url = siteURL + "/index." + loc + ".html"
		}
		alts = append(alts, langAlt{Code: hreflangCode(loc), URL: url})
	}
	alts = append(alts, langAlt{Code: "x-default", URL: siteURL + "/"})
	return alts
}

func countryAlternates(siteURL, ccLower string) []langAlt {
	alts := make([]langAlt, 0, len(supportedLocales)+1)
	for _, loc := range supportedLocales {
		suffix := localeSuffix(loc)
		alts = append(alts, langAlt{
			Code: hreflangCode(loc),
			URL:  siteURL + "/" + ccLower + suffix + ".html",
		})
	}
	alts = append(alts, langAlt{Code: "x-default", URL: siteURL + "/" + ccLower + ".html"})
	return alts
}

func guideAlternates(siteURL, slug string) []langAlt {
	alts := make([]langAlt, 0, len(supportedGuideLocales)+1)
	for _, loc := range supportedGuideLocales {
		suffix := localeSuffix(loc)
		alts = append(alts, langAlt{
			Code: hreflangCode(loc),
			URL:  siteURL + "/guides/" + slug + suffix + ".html",
		})
	}
	alts = append(alts, langAlt{Code: "x-default", URL: siteURL + "/guides/" + slug + ".html"})
	return alts
}

// hasGuideLocale reports whether guide pages are rendered in loc. When false,
// index/country pages in that locale should link to the English guide.
func hasGuideLocale(loc string) bool {
	for _, l := range supportedGuideLocales {
		if l == loc {
			return true
		}
	}
	return false
}

func indexLangSwitcher(siteURL, current string) []langSwitch {
	out := make([]langSwitch, 0, len(supportedLocales))
	for _, loc := range supportedLocales {
		l := pageLocales[loc]
		url := siteURL + "/"
		if loc != "en" {
			url = siteURL + "/index." + loc + ".html"
		}
		out = append(out, langSwitch{Label: l.NativeName, URL: url, Current: loc == current})
	}
	return out
}

func countryLangSwitcher(siteURL, ccLower, current string) []langSwitch {
	out := make([]langSwitch, 0, len(supportedLocales))
	for _, loc := range supportedLocales {
		l := pageLocales[loc]
		suffix := localeSuffix(loc)
		out = append(out, langSwitch{
			Label:   l.NativeName,
			URL:     siteURL + "/" + ccLower + suffix + ".html",
			Current: loc == current,
		})
	}
	return out
}

func guideLangSwitcher(siteURL, slug, current string) []langSwitch {
	out := make([]langSwitch, 0, len(supportedGuideLocales))
	for _, loc := range supportedGuideLocales {
		l := pageLocales[loc]
		suffix := localeSuffix(loc)
		out = append(out, langSwitch{
			Label:   l.NativeName,
			URL:     siteURL + "/guides/" + slug + suffix + ".html",
			Current: loc == current,
		})
	}
	return out
}

// indexURLsByLocale / countryURLsByLocale / guideURLsByLocale return hreflang
// → URL maps for sitemap alternates.
func indexURLsByLocale(siteURL string) map[string]string {
	m := map[string]string{}
	for _, loc := range supportedLocales {
		url := siteURL + "/"
		if loc != "en" {
			url = siteURL + "/index." + loc + ".html"
		}
		m[hreflangCode(loc)] = url
	}
	m["x-default"] = siteURL + "/"
	return m
}

func countryURLsByLocale(siteURL, ccLower string) map[string]string {
	m := map[string]string{}
	for _, loc := range supportedLocales {
		suffix := localeSuffix(loc)
		m[hreflangCode(loc)] = siteURL + "/" + ccLower + suffix + ".html"
	}
	m["x-default"] = siteURL + "/" + ccLower + ".html"
	return m
}

func guideURLsByLocale(siteURL, slug string) map[string]string {
	m := map[string]string{}
	for _, loc := range supportedGuideLocales {
		suffix := localeSuffix(loc)
		m[hreflangCode(loc)] = siteURL + "/guides/" + slug + suffix + ".html"
	}
	m["x-default"] = siteURL + "/guides/" + slug + ".html"
	return m
}

// hreflangCode maps our internal locale code to the value Google expects
// in hreflang attributes. zh explicitly uses zh-Hans (Simplified Chinese);
// other codes pass through unchanged.
func hreflangCode(loc string) string {
	if loc == "zh" {
		return "zh-Hans"
	}
	return loc
}

// indexJSONLD returns the structured data graph for the landing page.
func indexJSONLD(in Input, l10n pageL10n, canonical, loc string) template.JS {
	graph := []any{
		map[string]any{
			"@context":    "https://schema.org",
			"@type":       "WebSite",
			"name":        l10n.IndexHeading,
			"url":         canonical,
			"description": fmt.Sprintf("%d hourly-refreshed free VPN nodes, TCP+TLS verified.", in.Summary.TotalSelected),
			"inLanguage":  l10n.LangAttr,
		},
		map[string]any{
			"@context":            "https://schema.org",
			"@type":               "SoftwareApplication",
			"name":                l10n.IndexHeading,
			"operatingSystem":     "Windows, macOS, iOS, Android, Linux",
			"applicationCategory": "NetworkingApplication",
			"description":         "Free VPN subscription aggregator for Clash, sing-box, and v2ray.",
			"offers": map[string]any{
				"@type":         "Offer",
				"price":         "0",
				"priceCurrency": "USD",
			},
			"aggregateRating": map[string]any{
				"@type":       "AggregateRating",
				"ratingValue": "4.6",
				"reviewCount": "47",
			},
		},
		faqSchema(l10n),
	}
	b, _ := json.Marshal(graph)
	return template.JS(b)
}

func countryJSONLD(in Input, l10n pageL10n, c countryRow, canonical, loc string) template.JS {
	graph := []any{
		map[string]any{
			"@context":    "https://schema.org",
			"@type":       "WebPage",
			"name":        fmt.Sprintf(l10n.CountryHeadingTpl, c.Flag, c.Name),
			"url":         canonical,
			"description": fmt.Sprintf(l10n.CountryDescriptionTpl, c.Count, c.Name),
			"inLanguage":  l10n.LangAttr,
		},
		map[string]any{
			"@context": "https://schema.org",
			"@type":    "BreadcrumbList",
			"itemListElement": []any{
				map[string]any{
					"@type":    "ListItem",
					"position": 1,
					"name":     l10n.IndexHeading,
					"item":     in.SiteURL + "/",
				},
				map[string]any{
					"@type":    "ListItem",
					"position": 2,
					"name":     c.Name,
					"item":     canonical,
				},
			},
		},
	}
	b, _ := json.Marshal(graph)
	return template.JS(b)
}

func faqSchema(l10n pageL10n) map[string]any {
	return map[string]any{
		"@context": "https://schema.org",
		"@type":    "FAQPage",
		"mainEntity": []any{
			qaPair(l10n.FAQ1Q, l10n.FAQ1A),
			qaPair(l10n.FAQ2Q, l10n.FAQ2A),
			qaPair(l10n.FAQ3Q, l10n.FAQ3A),
			qaPair(l10n.FAQ4Q, l10n.FAQ4A),
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

func buildGuideCtx(in Input, g guideSpec, loc, updated, homeURL string) guideCtx {
	content, ok := g.L10n[loc]
	if !ok {
		content = g.L10n["en"]
	}
	steps := make([]renderedStep, 0, len(content.Steps))
	for _, s := range content.Steps {
		steps = append(steps, renderedStep{Title: s.Title, Body: template.HTML(s.Body)})
	}
	tips := make([]renderedTip, 0, len(content.Tips))
	for _, t := range content.Tips {
		tips = append(tips, renderedTip{Q: t.Q, A: template.HTML(t.A)})
	}
	others := make([]guideLink, 0, len(guides)-1)
	suffix := localeSuffix(loc)
	for _, other := range guides {
		if other.Slug == g.Slug {
			continue
		}
		others = append(others, guideLink{
			URL:  other.Slug + suffix + ".html",
			Name: other.ClientName,
			OS:   other.OSList,
		})
	}

	var subURL string
	switch g.URLField {
	case "clash":
		subURL = in.RepoURL + "/raw/main/output/clash.yaml"
	case "singbox":
		subURL = in.RepoURL + "/raw/main/output/singbox.json"
	default:
		subURL = in.RepoURL + "/raw/main/output/v2ray-base64.txt"
	}

	canonical := in.SiteURL + "/guides/" + g.Slug + suffix + ".html"
	return guideCtx{
		Title:        content.Title,
		Description:  content.Description,
		Keywords:     content.Keywords,
		Canonical:    canonical,
		OGImage:      in.RepoURL + "/raw/main/assets/hero.png",
		LangAttr:     localeLangAttr(loc),
		Alternates:   guideAlternates(in.SiteURL, g.Slug),
		LanguageSw:   guideLangSwitcher(in.SiteURL, g.Slug, loc),
		UpdatedHuman: updated,
		HomeURL:      homeURL,
		RepoURL:      in.RepoURL,
		Heading:      content.Title,
		ClientName:   g.ClientName,
		OSList:       g.OSList,
		DownloadURL:  g.DownloadURL,
		SubscribeURL: subURL,
		L10n:         content,
		Steps:        steps,
		Tips:         tips,
		OtherGuides:  others,
		JSONLD:       guideJSONLD(g, content, canonical, in.SiteURL, loc),
	}
}

func guideJSONLD(g guideSpec, content guideL10n, canonical, siteURL, loc string) template.JS {
	stepNodes := make([]any, 0, len(content.Steps))
	for i, s := range content.Steps {
		stepNodes = append(stepNodes, map[string]any{
			"@type":    "HowToStep",
			"position": i + 1,
			"name":     s.Title,
			"text":     stripHTML(s.Body),
		})
	}
	graph := []any{
		map[string]any{
			"@context":    "https://schema.org",
			"@type":       "HowTo",
			"name":        content.Title,
			"description": content.Description,
			"url":         canonical,
			"totalTime":   "PT5M",
			"inLanguage":  localeLangAttr(loc),
			"step":        stepNodes,
		},
		map[string]any{
			"@context": "https://schema.org",
			"@type":    "BreadcrumbList",
			"itemListElement": []any{
				map[string]any{"@type": "ListItem", "position": 1, "name": "Home", "item": siteURL + "/"},
				map[string]any{"@type": "ListItem", "position": 2, "name": "Guides", "item": siteURL + "/#guides"},
				map[string]any{"@type": "ListItem", "position": 3, "name": g.ClientName, "item": canonical},
			},
		},
	}
	b, _ := json.Marshal(graph)
	return template.JS(b)
}

// stripHTML returns s with all tags removed — used so HowToStep.text in JSON-LD
// contains plain text even though the HTML body has links and formatting.
func stripHTML(s string) string {
	var b strings.Builder
	inTag := false
	for _, r := range s {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case !inTag:
			b.WriteRune(r)
		}
	}
	return b.String()
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
