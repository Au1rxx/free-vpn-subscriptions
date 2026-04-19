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
	Hook2:       "No signup. No payment. No installation of binaries. Refreshed hourly from public sources with every node tested.",
	KeywordLine: "Free VPN subscriptions · free proxy list · free v2ray / clash / sing-box · VLESS / Reality / VMess / Trojan / Shadowsocks / Hysteria2 · hourly refreshed · TCP + TLS probed · by country",

	WhyHeading: "## 💡 Why This Project?",
	WhyBody:    "Every \"free VPN\" list on GitHub is either stale, full of dead nodes, or asks you to install a sketchy binary. This repo **only publishes nodes that passed a live TCP handshake AND a TLS handshake minutes ago**, from curated public sources, sorted by latency. You get 3 portable subscription files — drop them into Clash, sing-box, or v2rayN and go.",

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

	FAQHeading: "## ❓ FAQ",
	FAQ1Q:      "Is this actually free?",
	FAQ1A:      "Yes. Nodes are operated by third-party volunteers who publish their own free subscriptions. We don't run any servers ourselves — we just test, rank, and repackage what's already public.",
	FAQ2Q:      "How fresh is the data?",
	FAQ2A:      "A GitHub Action runs every hour: pulls all upstream sources, TCP + TLS probes every node, drops anything dead, sorts by latency, and commits new output files. Check the `Last updated` timestamp above.",
	FAQ3Q:      "Can I trust these nodes?",
	FAQ3A:      "Free nodes see all your traffic. **Never use them for banking, login, or anything sensitive.** Fine for bypassing geo-blocks on public content. Use your own VPS / paid provider for real privacy.",
	FAQ4Q:      "Why do some nodes fail even though they're listed?",
	FAQ4A:      "We verify TCP reachability and TLS handshakes, but a node can still have an expired quota, wrong routing, or GFW poisoning. Try a few; the selector group gives you fallbacks.",

	ContributingHeading: "## 🤝 Contributing",
	ContributingBody:    "Know a reliable public subscription source we should add? Open an issue with the URL and format.",

	DisclaimerHeading: "## ⚠️ Disclaimer",
	DisclaimerBody:    "This repository aggregates **publicly shared** proxy configurations from third-party volunteers. We do not operate any servers, do not warrant availability or security, and are not responsible for how you use them. Intended for educational and personal connectivity use. Comply with all applicable laws in your jurisdiction.",

	StarHistoryHeading: "## ⭐ Star History",
	FinalCTA:           "If this project helped you, give it a ⭐ — every star makes it easier for others to find.",
}
