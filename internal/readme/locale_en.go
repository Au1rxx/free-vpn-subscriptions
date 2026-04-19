package readme

var EN = Locale{
	Code:        "en",
	DisplayName: "English",
	FileName:    "README.md",
	LangAttr:    "en",

	BadgeNodes:   "nodes",
	BadgeAlive:   "alive",
	BadgeMedian:  "median--rtt",
	BadgeUpdated: "updated",

	Hook1:       "**The easiest way to get a working free VPN — copy a subscription link, paste it into your client, connect.**",
	Hook2:       "No signup. No payment. No installation of binaries. Refreshed hourly from public sources — every node is TCP + TLS probed before publishing.",
	KeywordLine: "Free VPN subscriptions · free proxy list · free v2ray / clash / sing-box · VLESS / Reality / VMess / Trojan / Shadowsocks / Hysteria2 · hourly refreshed · TCP + TLS probed · by country",

	WhyHeading: "## 💡 Why This Project?",
	WhyBody:    "Every \"free VPN\" list on GitHub is either stale, full of dead nodes, or asks you to install a sketchy binary. This repo **only publishes nodes that passed a live TCP handshake AND a TLS handshake minutes ago**, from curated public sources, sorted by latency. You get 3 portable subscription files — drop them into Clash, sing-box, or v2rayN and go.",

	VerificationHeading: "## 🔬 How we verify nodes actually work",
	VerificationBody: `**Honest answer first: we cannot *guarantee* a node will pass your traffic.** No aggregator can, without running real traffic through it. Here is exactly what we verify, what we cannot, and where the real guarantee comes from.

### ✅ What we verify at aggregation time (before publishing)

1. **TCP reachability** — we open a TCP connection to every ` + "`server:port`" + `. Dead hosts, bad DNS, and blocked ports get dropped. Drops roughly 40 % of raw entries.
2. **TLS handshake** — for every TLS / Reality / WS-TLS node we complete the full handshake. Expired certs, SNI mismatches, and broken Reality short-ids get dropped. Drops another ~10 %.
3. **Latency sort** — survivors are ranked by RTT and the top N are kept.

Typical numbers from a recent run: **17 sources → ~4,800 raw → ~2,900 TCP-alive → ~2,600 TLS-OK → top 200 published**.

### ❌ What we cannot verify

- Proxy protocol auth. A wrong UUID / password is only rejected *after* the TLS handshake by the upstream server.
- Actual HTTP-through-proxy success.
- Bandwidth or throughput.
- Geolocation beyond what GeoIP tells us about the exit IP.

### 🛡️ Runtime verification — the real guarantee

The ` + "`clash.yaml`" + ` we publish ships with a ` + "`url-test`" + ` proxy group that probes **real HTTP through each node** every 5 minutes:

` + "```yaml" + `
proxy-groups:
  - name: AUTO
    type: url-test
    url: http://www.gstatic.com/generate_204
    interval: 300
` + "```" + `

Your client ranks the node list by *actual* HTTP-through-proxy latency and auto-picks the fastest working node. sing-box and v2ray have equivalent mechanisms. If a selected node dies, the client switches to the next without intervention.

### 🧮 Expected outcome

Of the top 200 published each run, a typical client will find 30-50 that serve HTTP cleanly at any given moment. Rotate if one gets slow — the URL-test group makes that one click.`,

	SubscribeHeading:   "## 🚀 One-Click Subscribe",
	SubscribeIntro:     "Copy the URL that matches your client and paste it into the subscription import field:",
	SubscribeColClient: "Client",
	SubscribeColFormat: "Format",
	SubscribeColURL:    "Subscribe URL",

	ClientsHeading: "## 🧩 Supported Clients",
	ClientsWindows: "**Windows**: v2rayN, Clash Verge, Hiddify, NekoRay",
	ClientsMacOS:   "**macOS**: ClashX Pro, Clash Verge, sing-box, Hiddify",
	ClientsIOS:     "**iOS**: Shadowrocket, Stash, Loon, sing-box, Hiddify",
	ClientsAndroid: "**Android**: v2rayNG, NekoBox, Clash Meta for Android, Hiddify, sing-box",
	ClientsLinux:   "**Linux**: mihomo (Clash.Meta), sing-box, v2ray-core",

	StatsHeading:     "## 📊 Live Stats",
	StatsNodes:       "**Nodes selected**",
	StatsAlive:       "**Alive across all sources**",
	StatsFastest:     "**Fastest node RTT**",
	StatsMedian:      "**Median RTT**",
	StatsUpdated:     "**Last updated (UTC)**",
	ProtocolMixLabel: "**Protocol mix:**",
	SourcesLabel:     "**Sources used this run:**",

	ByCountryHeading: "## 🌍 By Country",
	ByCountryIntro:   "Want nodes in a specific region only? Use one of these targeted subscription URLs:",
	ByCountryColCC:   "Country",
	ByCountryColN:    "Nodes",

	GuidesHeading:     "## 📖 Step-by-step Guides",
	GuidesIntro:       "New to VPN clients? Pick your platform and follow the tutorial:",
	GuideLocaleSuffix: "",

	FAQHeading: "## ❓ FAQ",
	FAQ1Q:      "Is this actually free?",
	FAQ1A:      "Yes. Nodes are operated by third-party volunteers who publish their own free subscriptions. We don't run any servers ourselves — we just test, rank, and repackage what's already public.",
	FAQ2Q:      "How fresh is the data?",
	FAQ2A:      "Every hour (with a small random delay to avoid hammering upstream on the `:00` mark): pulls all upstream sources, TCP + TLS probes every node, drops anything dead, sorts by latency, publishes new output files. See the `Last updated` badge above.",
	FAQ3Q:      "Can I trust these nodes?",
	FAQ3A:      "Free nodes see all your traffic. **Never use them for banking, login, or anything sensitive.** Fine for bypassing geo-blocks on public content. Use your own VPS / paid provider for real privacy.",
	FAQ4Q:      "Why do some nodes fail even though they're listed?",
	FAQ4A:      "We verify TCP reachability and TLS handshakes only; a node can still have an expired quota, wrong routing, or GFW poisoning. The published `clash.yaml` pairs every node with a `url-test` proxy group (`http://www.gstatic.com/generate_204`, 300 s interval) — your client auto-picks the fastest node that actually serves HTTP. If one dies, pick the next.",

	ContributingHeading: "## 🤝 Contributing",
	ContributingBody:    "Know a reliable public subscription source we should add? Open an issue with the URL and format.",

	DisclaimerHeading: "## ⚠️ Disclaimer",
	DisclaimerBody:    "This repository aggregates **publicly shared** proxy configurations from third-party volunteers. We do not operate any servers, do not warrant availability or security, and are not responsible for how you use them. Intended for educational and personal connectivity use. Comply with all applicable laws in your jurisdiction.",

	StarHistoryHeading: "## ⭐ Star History",
	FinalCTA:           "If this project helped you, give it a ⭐ — every star makes it easier for others to find.",
}
