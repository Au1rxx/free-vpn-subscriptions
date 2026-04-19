package pages

// pageL10n holds every user-facing string used by the index and per-country
// page templates. Field values are looked up by locale code.
type pageL10n struct {
	// <html lang=...>
	LangAttr string
	// Displayed in language switcher for this locale (shown in this locale).
	NativeName string

	// <title> / heading
	IndexTitleTpl        string // %d selected
	IndexDescriptionTpl  string // %d selected
	IndexKeywords        string
	IndexHeading         string
	IndexSubTagline      string // short pitch beneath hero

	BadgeNodes        string
	BadgeAlive        string
	BadgeMedianRTT    string
	BadgeUpdated      string
	BadgeUpdatedUnit  string // "ms" suffix if needed — leave empty normally

	OneClickHeading  string
	OneClickIntro    string
	ColClash         string // label for Clash row
	ColSing          string
	ColV2ray         string

	ByCountryHeading string
	ByCountryIntro   string
	NodesSuffix      string // "nodes" / "节点" rendered after count

	GuidesHeading string
	GuidesIntro   string

	ClientsHeading string
	ClientsWindows string
	ClientsMacOS   string
	ClientsIOS     string
	ClientsAndroid string
	ClientsLinux   string

	FAQHeading string
	FAQ1Q      string
	FAQ1A      string
	FAQ2Q      string
	FAQ2A      string
	FAQ3Q      string
	FAQ3A      string
	FAQ4Q      string
	FAQ4A      string

	StarButton       string
	FooterLicense    string // "MIT licensed. Open source on GitHub."
	FooterDisclaimer string

	// Country page strings
	CountryTitleTpl        string // %s name, %d count
	CountryDescriptionTpl  string
	CountryKeywordsTpl     string // %s name (6 insertions)
	CountryHeadingTpl      string // %s flag, %s name
	CountrySubTpl          string // %d count, %s name
	CountryBreadcrumb      string // "← All countries"
	CountryOtherHeading    string // "Other countries"
	CountrySubscribeHeadingTpl string // "Subscribe to %s nodes only"
	CountryStarHistoryNote string

	// Language switcher preamble
	LanguageLabel string // "Language:" / "语言:"
}

var pageLocales = map[string]pageL10n{
	"en": {
		LangAttr:   "en",
		NativeName: "English",

		IndexTitleTpl:       "Free VPN Subscriptions · %d nodes · hourly · Clash / sing-box / v2ray",
		IndexDescriptionTpl: "%d TCP+TLS-verified free VPN nodes from public sources. Hourly refresh. Copy a Clash, sing-box or v2ray URL and paste into your client.",
		IndexKeywords:       "free vpn, vpn subscription, clash, sing-box, v2ray, vless, reality, trojan, shadowsocks, hysteria2, proxy list, free proxy",
		IndexHeading:        "Free VPN Subscriptions",
		IndexSubTagline:     "The easiest way to get a working free VPN — copy a subscription link, paste it into your client, connect.",

		BadgeNodes:     "nodes",
		BadgeAlive:     "alive",
		BadgeMedianRTT: "median RTT",
		BadgeUpdated:   "updated",

		OneClickHeading: "🚀 One-Click Subscribe",
		OneClickIntro:   "Copy the URL that matches your client and paste it into the subscription import field:",
		ColClash:        "Clash / Clash Verge / ClashX",
		ColSing:         "sing-box",
		ColV2ray:        "v2rayN / v2rayNG / Shadowrocket / NekoBox",

		ByCountryHeading: "🌍 By Country",
		ByCountryIntro:   "Want nodes in a specific region only? Choose a targeted subscription:",
		NodesSuffix:      "nodes",

		GuidesHeading: "📖 Step-by-step Guides",
		GuidesIntro:   "New to VPN clients? Pick your platform:",

		ClientsHeading: "🧩 Supported Clients",
		ClientsWindows: "<strong>Windows</strong>: v2rayN, Clash Verge, Hiddify, NekoRay",
		ClientsMacOS:   "<strong>macOS</strong>: ClashX Pro, Clash Verge, sing-box, Hiddify",
		ClientsIOS:     "<strong>iOS</strong>: Shadowrocket, Stash, Loon, sing-box, Hiddify",
		ClientsAndroid: "<strong>Android</strong>: v2rayNG, NekoBox, Clash Meta for Android, Hiddify, sing-box",
		ClientsLinux:   "<strong>Linux</strong>: mihomo (Clash.Meta), sing-box, v2ray-core",

		FAQHeading: "❓ FAQ",
		FAQ1Q:      "Is this actually free?",
		FAQ1A:      "Yes. Nodes are operated by third-party volunteers. We don't run any servers — we just test, rank, and repackage what's already public.",
		FAQ2Q:      "How fresh is the data?",
		FAQ2A:      "A GitHub Action runs every hour: pulls all upstream sources, TCP+TLS probes every node, drops anything dead, sorts by latency, commits new output files.",
		FAQ3Q:      "Can I trust these nodes?",
		FAQ3A:      "Free nodes see all your traffic. Never use them for banking, login, or anything sensitive. Fine for bypassing geo-blocks on public content. Use your own VPS or a paid provider for real privacy.",
		FAQ4Q:      "Why do some nodes fail?",
		FAQ4A:      "We verify TCP reachability and TLS handshakes, but a node may still have an expired quota, bad routing, or revoked certs. Try a few; the selector group gives you fallbacks.",

		StarButton:       "⭐ Star on GitHub",
		FooterLicense:    "Open source on GitHub. MIT licensed.",
		FooterDisclaimer: "This project aggregates publicly-shared proxy configurations. We do not operate any servers and make no guarantees. Comply with the laws of your jurisdiction.",

		CountryTitleTpl:            "Free %s VPN Subscription · %d nodes · Clash / sing-box / v2ray",
		CountryDescriptionTpl:      "%d TCP+TLS-verified free VPN nodes in %s, refreshed hourly. Copy a Clash, sing-box, or v2ray subscription URL.",
		CountryKeywordsTpl:         "free %s vpn, %s vpn subscription, %s clash, %s v2ray, %s proxy, %s free vpn",
		CountryHeadingTpl:          "Free %s %s VPN Subscription",
		CountrySubTpl:              "%d free VPN nodes in %s, TCP+TLS verified, refreshed hourly.",
		CountryBreadcrumb:          "← All countries",
		CountryOtherHeading:        "🌍 Other countries",
		CountrySubscribeHeadingTpl: "🚀 Subscribe to %s nodes only",

		LanguageLabel: "Language:",
	},
	"zh": {
		LangAttr:   "zh-Hans",
		NativeName: "简体中文",

		IndexTitleTpl:       "免费 VPN 订阅 · %d 节点 · 每小时刷新 · Clash / sing-box / v2ray",
		IndexDescriptionTpl: "来自公共源的 %d 个 TCP+TLS 实测节点,每小时自动刷新。复制订阅链接粘贴到客户端即可,支持 Clash、sing-box、v2ray。",
		IndexKeywords:       "免费 VPN 订阅,免费机场,免费梯子,clash 订阅,v2ray 订阅,sing-box 订阅,VLESS,Reality,Trojan,Shadowsocks,Hysteria2,每小时刷新,免费节点",
		IndexHeading:        "免费 VPN 订阅",
		IndexSubTagline:     "获取可用免费 VPN 的最简单方式——复制订阅链接,粘贴到客户端,连上就用。",

		BadgeNodes:     "节点",
		BadgeAlive:     "存活",
		BadgeMedianRTT: "中位延迟",
		BadgeUpdated:   "更新",

		OneClickHeading: "🚀 一键订阅",
		OneClickIntro:   "复制对应客户端的 URL,粘贴到订阅导入框:",
		ColClash:        "Clash / Clash Verge / ClashX",
		ColSing:         "sing-box",
		ColV2ray:        "v2rayN / v2rayNG / Shadowrocket / NekoBox",

		ByCountryHeading: "🌍 按国家订阅",
		ByCountryIntro:   "只想要特定地区的节点?选一个针对性订阅链接:",
		NodesSuffix:      "节点",

		GuidesHeading: "📖 客户端图文教程",
		GuidesIntro:   "新手?按平台挑一篇跟着做:",

		ClientsHeading: "🧩 支持的客户端",
		ClientsWindows: "<strong>Windows</strong>:v2rayN、Clash Verge、Hiddify、NekoRay",
		ClientsMacOS:   "<strong>macOS</strong>:ClashX Pro、Clash Verge、sing-box、Hiddify",
		ClientsIOS:     "<strong>iOS</strong>:Shadowrocket、Stash、Loon、sing-box、Hiddify",
		ClientsAndroid: "<strong>Android</strong>:v2rayNG、NekoBox、Clash Meta for Android、Hiddify、sing-box",
		ClientsLinux:   "<strong>Linux</strong>:mihomo (Clash.Meta)、sing-box、v2ray-core",

		FAQHeading: "❓ 常见问题",
		FAQ1Q:      "真的完全免费吗?",
		FAQ1A:      "是的。所有节点由第三方志愿者运营并公开免费订阅。本项目不运营任何服务器,只做测试、排名和重新打包。",
		FAQ2Q:      "数据多新?",
		FAQ2A:      "GitHub Actions 每小时运行一次:拉取所有上游源,对每个节点做 TCP+TLS 探测,丢弃死节点,按延迟排序,提交新文件。",
		FAQ3Q:      "这些节点可以信任吗?",
		FAQ3A:      "免费节点能看到你所有流量。绝不要用来登录银行、邮箱等敏感账号。用来突破地区限制访问公开内容没问题。真正需要隐私请自建 VPS 或付费服务。",
		FAQ4Q:      "列表里有些节点连不上?",
		FAQ4A:      "我们验证 TCP 可达和 TLS 握手,但节点仍可能配额用完、路由被污染或证书到期。多试几个,selector 组自带 fallback。",

		StarButton:       "⭐ GitHub 上点 Star",
		FooterLicense:    "开源于 GitHub,MIT 许可证。",
		FooterDisclaimer: "本项目聚合公开分享的代理配置,不运营任何服务器,不作任何保证。请遵守所在司法管辖区的法律。",

		CountryTitleTpl:            "%s 免费 VPN 订阅 · %d 节点 · Clash / sing-box / v2ray",
		CountryDescriptionTpl:      "%d 个 %s 免费 VPN 节点,每小时 TCP+TLS 实测过滤。复制 Clash、sing-box 或 v2ray 订阅链接即可使用。",
		CountryKeywordsTpl:         "%s 免费 vpn,%s 免费机场,%s clash 订阅,%s v2ray,%s 免费节点,%s 梯子",
		CountryHeadingTpl:          "%s %s 免费 VPN 订阅",
		CountrySubTpl:              "%d 个 %s 节点,TCP+TLS 实测过滤,每小时刷新。",
		CountryBreadcrumb:          "← 所有国家",
		CountryOtherHeading:        "🌍 其他国家",
		CountrySubscribeHeadingTpl: "🚀 仅订阅 %s 节点",

		LanguageLabel: "语言:",
	},
}
