// Package readme renders public README.md files (one per locale) from
// aggregation summary data. Each README is designed for SEO, scan-ability,
// and star conversion — modelled after the free-llm-api-keys structure.
package readme

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/internal/aggregate"
	"github.com/Au1rxx/free-vpn-subscriptions/internal/node"
)

type Input struct {
	Title          string
	RepoURL        string
	Nodes          []*node.Node
	Summary        aggregate.Summary
	CountryEnabled bool
	MinPerCountry  int
}

// Generate returns a complete README markdown document for the given locale.
func Generate(in Input, loc Locale) string {
	updated := time.Unix(in.Summary.GeneratedAtUnix, 0).UTC().Format("2006-01-02 15:04 UTC")

	var b strings.Builder

	// Title
	fmt.Fprintf(&b, "# %s\n\n", in.Title)

	// Language switcher (every locale sees the same bar so navigation is
	// symmetric — users can always get back to their preferred language).
	b.WriteString(renderLangSwitcher(loc))
	b.WriteString("\n\n")

	// Hero image — SVG workflow diagram is the fallback when no real
	// screenshot exists yet; the alt text carries SEO weight.
	fmt.Fprintf(&b, "<p align=\"center\"><img src=\"%s/raw/main/assets/workflow.svg\" alt=\"%s — aggregation workflow\" width=\"780\"></p>\n\n",
		in.RepoURL, in.Title)

	// Badges
	fmt.Fprintf(&b, "![%s](https://img.shields.io/badge/%s-%d-brightgreen) ",
		loc.BadgeNodes, loc.BadgeNodes, in.Summary.TotalSelected)
	fmt.Fprintf(&b, "![%s](https://img.shields.io/badge/%s-%d-blue) ",
		loc.BadgeAlive, loc.BadgeAlive, in.Summary.TotalAlive)
	fmt.Fprintf(&b, "![%s](https://img.shields.io/badge/%s-%dms-orange) ",
		loc.BadgeMedian, loc.BadgeMedian, in.Summary.MedianLatencyMS)
	fmt.Fprintf(&b, "![%s](https://img.shields.io/badge/%s-%s-informational)\n\n",
		loc.BadgeUpdated, loc.BadgeUpdated, strings.ReplaceAll(updated, " ", "_"))

	// Hook
	fmt.Fprintf(&b, "> %s  \n", loc.Hook1)
	fmt.Fprintf(&b, "> %s\n\n", loc.Hook2)
	fmt.Fprintf(&b, "> %s\n\n", loc.KeywordLine)

	// Why
	fmt.Fprintf(&b, "%s\n\n%s\n\n", loc.WhyHeading, loc.WhyBody)
	fmt.Fprintf(&b, "> 📖 How the fetch → probe → rank pipeline works: [ARCHITECTURE.md](./ARCHITECTURE.md)\n\n")

	// Subscribe
	fmt.Fprintf(&b, "%s\n\n%s\n\n", loc.SubscribeHeading, loc.SubscribeIntro)
	fmt.Fprintf(&b, "| %s | %s | %s |\n|---|---|---|\n",
		loc.SubscribeColClient, loc.SubscribeColFormat, loc.SubscribeColURL)
	fmt.Fprintf(&b, "| Clash / Clash Verge / ClashX | `clash.yaml` | `%s/raw/main/output/clash.yaml` |\n", in.RepoURL)
	fmt.Fprintf(&b, "| sing-box | `singbox.json` | `%s/raw/main/output/singbox.json` |\n", in.RepoURL)
	fmt.Fprintf(&b, "| v2rayN / v2rayNG / Shadowrocket / NekoBox | `v2ray-base64` | `%s/raw/main/output/v2ray-base64.txt` |\n\n", in.RepoURL)

	// Per-country
	if in.CountryEnabled && in.MinPerCountry > 0 && len(in.Summary.ByCountry) > 0 {
		renderByCountry(&b, in, loc)
	}

	// Guides (client tutorials) — one entry per supported client, pointing at
	// the locale's HTML guide page (or English fallback).
	renderGuides(&b, in, loc)

	// Clients
	fmt.Fprintf(&b, "%s\n\n- %s\n- %s\n- %s\n- %s\n- %s\n\n",
		loc.ClientsHeading,
		loc.ClientsWindows, loc.ClientsMacOS,
		loc.ClientsIOS, loc.ClientsAndroid, loc.ClientsLinux)

	// Stats
	fmt.Fprintf(&b, "%s\n\n", loc.StatsHeading)
	fmt.Fprintf(&b, "- %s: %d\n", loc.StatsNodes, in.Summary.TotalSelected)
	fmt.Fprintf(&b, "- %s: %d\n", loc.StatsAlive, in.Summary.TotalAlive)
	fmt.Fprintf(&b, "- %s: %d ms\n", loc.StatsFastest, in.Summary.MinLatencyMS)
	fmt.Fprintf(&b, "- %s: %d ms\n", loc.StatsMedian, in.Summary.MedianLatencyMS)
	fmt.Fprintf(&b, "- %s: %s\n\n", loc.StatsUpdated, updated)

	// Protocol breakdown
	if len(in.Summary.ByProtocol) > 0 {
		b.WriteString(loc.ProtocolMixLabel + " ")
		keys := sortedKeys(in.Summary.ByProtocol)
		parts := make([]string, 0, len(keys))
		for _, k := range keys {
			parts = append(parts, fmt.Sprintf("%s × %d", k, in.Summary.ByProtocol[k]))
		}
		b.WriteString(strings.Join(parts, " · "))
		b.WriteString("\n\n")
	}

	// Source breakdown
	if len(in.Summary.BySource) > 0 {
		b.WriteString(loc.SourcesLabel + " ")
		keys := sortedKeys(in.Summary.BySource)
		parts := make([]string, 0, len(keys))
		for _, k := range keys {
			parts = append(parts, fmt.Sprintf("`%s` × %d", k, in.Summary.BySource[k]))
		}
		b.WriteString(strings.Join(parts, " · "))
		b.WriteString("\n\n")
	}

	// FAQ
	fmt.Fprintf(&b, "%s\n\n", loc.FAQHeading)
	fmt.Fprintf(&b, "<details><summary>%s</summary>\n\n%s\n\n</details>\n\n", loc.FAQ1Q, loc.FAQ1A)
	fmt.Fprintf(&b, "<details><summary>%s</summary>\n\n%s\n\n</details>\n\n", loc.FAQ2Q, loc.FAQ2A)
	fmt.Fprintf(&b, "<details><summary>%s</summary>\n\n%s\n\n</details>\n\n", loc.FAQ3Q, loc.FAQ3A)
	fmt.Fprintf(&b, "<details><summary>%s</summary>\n\n%s\n\n</details>\n\n", loc.FAQ4Q, loc.FAQ4A)

	// Contributing
	fmt.Fprintf(&b, "%s\n\n%s\n\n", loc.ContributingHeading, loc.ContributingBody)

	// Disclaimer
	fmt.Fprintf(&b, "%s\n\n%s\n\n", loc.DisclaimerHeading, loc.DisclaimerBody)

	// Star History
	fmt.Fprintf(&b, "%s\n\n", loc.StarHistoryHeading)
	fmt.Fprintf(&b, "[![Star History Chart](https://api.star-history.com/svg?repos=%s&type=Date)](https://www.star-history.com/#%s&Date)\n\n",
		repoSlug(in.RepoURL), repoSlug(in.RepoURL))

	b.WriteString("---\n\n")
	b.WriteString(loc.FinalCTA + "\n")

	return b.String()
}

// renderLangSwitcher builds a Markdown link row for every supported locale,
// bolding the current language.
func renderLangSwitcher(current Locale) string {
	var parts []string
	for _, loc := range Locales() {
		label := loc.DisplayName
		if loc.Code == current.Code {
			parts = append(parts, fmt.Sprintf("**%s**", label))
		} else {
			parts = append(parts, fmt.Sprintf("[%s](./%s)", label, loc.FileName))
		}
	}
	return strings.Join(parts, " · ")
}

// renderByCountry appends a "By Country" section listing available
// per-country subscription files, sorted by node count desc.
func renderByCountry(b *strings.Builder, in Input, loc Locale) {
	type row struct {
		cc    string
		count int
	}
	var rows []row
	for cc, n := range in.Summary.ByCountry {
		if cc == "" || cc == "XX" {
			continue
		}
		if n < in.MinPerCountry {
			continue
		}
		rows = append(rows, row{cc, n})
	}
	if len(rows) == 0 {
		return
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].count != rows[j].count {
			return rows[i].count > rows[j].count
		}
		return rows[i].cc < rows[j].cc
	})

	fmt.Fprintf(b, "%s\n\n%s\n\n", loc.ByCountryHeading, loc.ByCountryIntro)
	fmt.Fprintf(b, "| %s | %s | Clash | sing-box | v2ray |\n", loc.ByCountryColCC, loc.ByCountryColN)
	b.WriteString("|---|---|---|---|---|\n")
	for _, r := range rows {
		flag := countryFlag(r.cc)
		label := countryLabel(r.cc)
		base := in.RepoURL + "/raw/main/output/by-country"
		fmt.Fprintf(b, "| %s %s (`%s`) | %d | [clash-%s.yaml](%s/clash-%s.yaml) | [singbox-%s.json](%s/singbox-%s.json) | [v2ray-base64-%s.txt](%s/v2ray-base64-%s.txt) |\n",
			flag, label, r.cc, r.count,
			r.cc, base, r.cc,
			r.cc, base, r.cc,
			r.cc, base, r.cc)
	}
	b.WriteString("\n")
}

// guideEntry is the small subset of pages.guideSpec that readme rendering
// needs. Kept local so the readme package doesn't depend on pages.
type guideEntry struct {
	Slug       string
	ClientName string
	OSList     string
}

// guideEntries mirrors the list in internal/pages/guides.go — if that list
// changes, add matching entries here. Kept hand-synced rather than importing
// pages (which would create a cycle).
var guideEntries = []guideEntry{
	{Slug: "clash-verge", ClientName: "Clash Verge", OSList: "Windows / macOS / Linux"},
	{Slug: "v2rayng", ClientName: "v2rayNG", OSList: "Android"},
	{Slug: "shadowrocket", ClientName: "Shadowrocket", OSList: "iOS / iPadOS"},
	{Slug: "sing-box", ClientName: "sing-box", OSList: "Windows / macOS / Linux / iOS / Android"},
}

func renderGuides(b *strings.Builder, in Input, loc Locale) {
	fmt.Fprintf(b, "%s\n\n%s\n\n", loc.GuidesHeading, loc.GuidesIntro)
	base := "https://au1rxx.github.io/free-vpn-subscriptions/guides"
	for _, g := range guideEntries {
		url := fmt.Sprintf("%s/%s%s.html", base, g.Slug, loc.GuideLocaleSuffix)
		fmt.Fprintf(b, "- [**%s**](%s) · %s\n", g.ClientName, url, g.OSList)
	}
	b.WriteString("\n")
}

// countryFlag returns a regional-indicator emoji for a 2-letter country code.
func countryFlag(cc string) string {
	if len(cc) != 2 {
		return ""
	}
	c := strings.ToUpper(cc)
	r1 := rune(c[0]) - 'A' + 0x1F1E6
	r2 := rune(c[1]) - 'A' + 0x1F1E6
	return string([]rune{r1, r2})
}

func countryLabel(cc string) string {
	if name, ok := countryNames[strings.ToUpper(cc)]; ok {
		return name
	}
	return cc
}

var countryNames = map[string]string{
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
	"CR": "Costa Rica", "PA": "Panama", "VG": "British Virgin Islands",
	"KY": "Cayman Islands", "PR": "Puerto Rico", "CO": "Colombia",
	"PE": "Peru", "EG": "Egypt", "NG": "Nigeria", "KE": "Kenya",
	"PK": "Pakistan", "BD": "Bangladesh", "LK": "Sri Lanka", "NP": "Nepal",
	"MM": "Myanmar", "KH": "Cambodia", "LA": "Laos", "MO": "Macau",
}

func sortedKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func repoSlug(repoURL string) string {
	s := strings.TrimPrefix(repoURL, "https://github.com/")
	s = strings.TrimSuffix(s, "/")
	return s
}
