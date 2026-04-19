package readme

var ZH = Locale{
	Code:        "zh",
	DisplayName: "简体中文",
	FileName:    "README_CN.md",
	LangAttr:    "zh-Hans",

	BadgeNodes:   "节点",
	BadgeAlive:   "存活",
	BadgeMedian:  "中位延迟",
	BadgeUpdated: "更新",

	Hook1:       "**获取可用免费 VPN 的最简单方式 —— 复制订阅链接,粘贴到客户端,连上。**",
	Hook2:       "无需注册。无需付费。无需安装任何二进制。每小时从公共源自动抓取,逐个节点测试。",
	KeywordLine: "免费 VPN 订阅 · 免费机场 · 免费梯子 · 免费科学上网 · Clash 订阅 · v2ray 订阅 · sing-box 订阅 · VLESS Reality · VMess · Trojan · Shadowsocks · Hysteria2 · 每小时刷新 · TCP+TLS 探测 · 按国家分类",

	WhyHeading: "## 💡 为什么用这个项目?",
	WhyBody:    "GitHub 上几乎所有的\"免费 VPN\"列表都有三个问题:数据过期、全是死节点、或者要你装来路不明的二进制。本仓库**只发布几分钟前通过 TCP 握手并通过 TLS 握手的节点**,来源于筛选过的公共订阅,按延迟排序。直接给你 3 种通用订阅文件 —— 粘到 Clash / sing-box / v2rayN 即用。",

	SubscribeHeading:   "## 🚀 一键订阅",
	SubscribeIntro:     "复制对应客户端的 URL,粘贴到订阅导入框:",
	SubscribeColClient: "客户端",
	SubscribeColFormat: "格式",
	SubscribeColURL:    "订阅链接",

	ClientsHeading: "## 🧩 支持的客户端",
	ClientsWindows: "**Windows**:v2rayN、Clash Verge、Hiddify、NekoRay",
	ClientsMacOS:   "**macOS**:ClashX Pro、Clash Verge、sing-box、Hiddify",
	ClientsIOS:     "**iOS**:Shadowrocket、Stash、Loon、sing-box、Hiddify",
	ClientsAndroid: "**Android**:v2rayNG、NekoBox、Clash Meta for Android、Hiddify、sing-box",
	ClientsLinux:   "**Linux**:mihomo (Clash.Meta)、sing-box、v2ray-core",

	StatsHeading:     "## 📊 实时统计",
	StatsNodes:       "**精选节点数**",
	StatsAlive:       "**全源存活总数**",
	StatsFastest:     "**最快延迟**",
	StatsMedian:      "**中位延迟**",
	StatsUpdated:     "**最后更新 (UTC)**",
	ProtocolMixLabel: "**协议分布:**",
	SourcesLabel:     "**本次使用的源:**",

	ByCountryHeading: "## 🌍 按国家订阅",
	ByCountryIntro:   "只想要特定地区的节点?选一个针对性订阅链接:",
	ByCountryColCC:   "国家/地区",
	ByCountryColN:    "节点数",

	GuidesHeading:     "## 📖 客户端图文教程",
	GuidesIntro:       "新手不知道怎么配?按平台选一篇跟着做:",
	GuideLocaleSuffix: ".zh",

	FAQHeading: "## ❓ 常见问题",
	FAQ1Q:      "真的完全免费吗?",
	FAQ1A:      "是的。所有节点由第三方志愿者自己运营并公开免费订阅。本仓库不运营任何服务器,只是做测试、排名、重新打包公开内容。",
	FAQ2Q:      "数据多新?",
	FAQ2A:      "GitHub Actions 每小时运行一次:拉取所有上游源 → TCP+TLS 探测每个节点 → 丢弃死节点 → 按延迟排序 → 提交新的输出文件。见顶部徽章上的更新时间。",
	FAQ3Q:      "这些节点可以信任吗?",
	FAQ3A:      "免费节点能看到你所有流量。**绝不要用来登录银行、邮箱等敏感账号。**用来突破地区限制访问公开内容没问题。真正需要隐私请自建 VPS 或付费服务。",
	FAQ4Q:      "列表里的节点为什么有的连不上?",
	FAQ4A:      "我们验证 TCP 可达和 TLS 握手,但节点仍可能配额用完、路由被污染、证书到期。多试几个,selector 组自带 fallback。",

	ContributingHeading: "## 🤝 贡献",
	ContributingBody:    "知道稳定的公共订阅源可以加入?提 issue 给我们 URL 和格式。",

	DisclaimerHeading: "## ⚠️ 免责声明",
	DisclaimerBody:    "本仓库聚合第三方志愿者**公开分享**的代理配置。我们不运营任何服务器,不保证可用性或安全性,不为使用行为负责。仅供学习和个人连接使用。请遵守所在司法管辖区的法律。",

	StarHistoryHeading: "## ⭐ Star 历史",
	FinalCTA:           "如果这个项目帮到你,点个 ⭐ —— 每一颗 star 都能帮更多人发现它。",
}
