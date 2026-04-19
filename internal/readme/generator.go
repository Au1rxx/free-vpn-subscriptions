// Package readme renders the public README.md from aggregation summary data.
// The output is designed for SEO, scan-ability, and star conversion — modelled
// after the free-llm-api-keys repo's structure.
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
	Title   string
	RepoURL string
	Nodes   []*node.Node
	Summary aggregate.Summary
}

// Generate returns the complete README markdown.
func Generate(in Input) string {
	updated := time.Unix(in.Summary.GeneratedAtUnix, 0).UTC().Format("2006-01-02 15:04 UTC")

	var b strings.Builder

	fmt.Fprintf(&b, "# %s\n\n", in.Title)

	// Badges
	fmt.Fprintf(&b, "![Nodes](https://img.shields.io/badge/nodes-%d-brightgreen) ", in.Summary.TotalSelected)
	fmt.Fprintf(&b, "![Alive](https://img.shields.io/badge/alive-%d-blue) ", in.Summary.TotalAlive)
	fmt.Fprintf(&b, "![Median RTT](https://img.shields.io/badge/median--rtt-%dms-orange) ", in.Summary.MedianLatencyMS)
	fmt.Fprintf(&b, "![Updated](https://img.shields.io/badge/updated-%s-informational)\n\n", strings.ReplaceAll(updated, " ", "_"))

	fmt.Fprintf(&b, "> **The easiest way to get a working free VPN — copy a subscription link, paste it into your client, connect.**  \n")
	fmt.Fprintf(&b, "> No signup. No payment. No installation of binaries. Refreshed hourly from public sources with every node tested.\n\n")
	fmt.Fprintf(&b, "> 免费 VPN 订阅 · 免费梯子 · 免费科学上网 · free proxy · v2ray/clash/sing-box · VLESS / Reality / VMess / Trojan / Shadowsocks / Hysteria2\n\n")

	// Why this project
	b.WriteString("## 💡 Why This Project?\n\n")
	b.WriteString("Every \"free VPN\" list on GitHub is either stale, full of dead nodes, or asks you to install a sketchy binary. This repo **only publishes nodes that passed a live TCP health check minutes ago**, from curated public sources, sorted by latency. You get 3 portable subscription files — drop them into Clash, sing-box, or v2rayN and go.\n\n")

	// One-click subscribe
	b.WriteString("## 🚀 One-Click Subscribe\n\n")
	b.WriteString("Copy the URL that matches your client and paste it into the subscription import field:\n\n")
	b.WriteString("| Client | Format | Subscribe URL |\n")
	b.WriteString("|---|---|---|\n")
	fmt.Fprintf(&b, "| Clash / Clash Verge / ClashX | `clash.yaml` | `%s/raw/main/output/clash.yaml` |\n", in.RepoURL)
	fmt.Fprintf(&b, "| sing-box | `singbox.json` | `%s/raw/main/output/singbox.json` |\n", in.RepoURL)
	fmt.Fprintf(&b, "| v2rayN / v2rayNG / Shadowrocket / NekoBox | `v2ray-base64` | `%s/raw/main/output/v2ray-base64.txt` |\n\n", in.RepoURL)

	// Client compatibility
	b.WriteString("## 🧩 Supported Clients\n\n")
	b.WriteString("- **Windows**: v2rayN, Clash Verge, Hiddify, NekoRay\n")
	b.WriteString("- **macOS**: ClashX Pro, Clash Verge, sing-box, Hiddify\n")
	b.WriteString("- **iOS**: Shadowrocket, Stash, Loon, sing-box, Hiddify\n")
	b.WriteString("- **Android**: v2rayNG, NekoBox, Clash Meta for Android, Hiddify, sing-box\n")
	b.WriteString("- **Linux**: mihomo (Clash.Meta), sing-box, v2ray-core\n\n")

	// Stats
	b.WriteString("## 📊 Live Stats\n\n")
	fmt.Fprintf(&b, "- **Nodes selected**: %d\n", in.Summary.TotalSelected)
	fmt.Fprintf(&b, "- **Alive across all sources**: %d\n", in.Summary.TotalAlive)
	fmt.Fprintf(&b, "- **Fastest node RTT**: %d ms\n", in.Summary.MinLatencyMS)
	fmt.Fprintf(&b, "- **Median RTT**: %d ms\n", in.Summary.MedianLatencyMS)
	fmt.Fprintf(&b, "- **Last updated (UTC)**: %s\n\n", updated)

	// Protocol breakdown
	if len(in.Summary.ByProtocol) > 0 {
		b.WriteString("**Protocol mix:** ")
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
		b.WriteString("**Sources used this run:** ")
		keys := sortedKeys(in.Summary.BySource)
		parts := make([]string, 0, len(keys))
		for _, k := range keys {
			parts = append(parts, fmt.Sprintf("`%s` × %d", k, in.Summary.BySource[k]))
		}
		b.WriteString(strings.Join(parts, " · "))
		b.WriteString("\n\n")
	}

	// FAQ
	b.WriteString("## ❓ FAQ\n\n")
	b.WriteString("<details><summary>Is this actually free?</summary>\n\n")
	b.WriteString("Yes. Nodes are operated by third-party volunteers who publish their own free subscriptions. We don't run any servers ourselves — we just test, rank, and repackage what's already public.\n\n</details>\n\n")
	b.WriteString("<details><summary>How fresh is the data?</summary>\n\n")
	b.WriteString("A GitHub Action runs every hour: pulls all upstream sources, TCP-probes every node, drops anything dead, sorts by latency, and commits new output files. Check the `Last updated` timestamp above.\n\n</details>\n\n")
	b.WriteString("<details><summary>Can I trust these nodes?</summary>\n\n")
	b.WriteString("Free nodes see all your traffic. **Never use them for banking, login, or anything sensitive.** Fine for bypassing geo-blocks on public content. Use your own VPS / paid provider for real privacy.\n\n</details>\n\n")
	b.WriteString("<details><summary>Why do some nodes fail even though they're listed?</summary>\n\n")
	b.WriteString("We only do TCP reachability checks. A node that handshakes may still have an expired cert, full bandwidth quota, or GFW-poisoned routes. Try a few; that's why the selector group gives you fallbacks.\n\n</details>\n\n")

	// Contributing
	b.WriteString("## 🤝 Contributing\n\n")
	b.WriteString("Know a reliable public subscription source we should add? Open an issue with the URL and format.\n\n")

	// Disclaimer
	b.WriteString("## ⚠️ Disclaimer\n\n")
	b.WriteString("This repository aggregates **publicly shared** proxy configurations from third-party volunteers. We do not operate any servers, do not warrant availability or security, and are not responsible for how you use them. Intended for educational and personal connectivity use. Comply with all applicable laws in your jurisdiction.\n\n")

	b.WriteString("## ⭐ Star History\n\n")
	fmt.Fprintf(&b, "[![Star History Chart](https://api.star-history.com/svg?repos=%s&type=Date)](https://www.star-history.com/#%s&Date)\n\n",
		repoSlug(in.RepoURL), repoSlug(in.RepoURL))

	b.WriteString("---\n\n")
	b.WriteString("If this project helped you, give it a ⭐ — every star makes it easier for others to find.\n")

	return b.String()
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
	// extract "owner/repo" from https://github.com/owner/repo
	s := strings.TrimPrefix(repoURL, "https://github.com/")
	s = strings.TrimSuffix(s, "/")
	return s
}
