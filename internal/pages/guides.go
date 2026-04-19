package pages

// guideSpec describes one client tutorial. Rendered once per supported
// locale via tplGuide.
type guideSpec struct {
	Slug        string // filename without .html, used in URL
	ClientName  string // e.g. "Clash Verge"
	OSList      string // e.g. "Windows / macOS / Linux"
	DownloadURL string
	URLField    string // which subscription URL to import: "clash" | "singbox" | "v2ray"
	L10n        map[string]guideL10n
}

// guideL10n is the translated content for one guide in one language.
type guideL10n struct {
	Title       string
	Description string
	Keywords    string

	// UI chrome shown on the page itself
	StepsHeading       string
	TipsHeading        string
	SubscribeHeading   string
	SubscribeIntro     string
	SubscribeLabel     string // e.g. "Clash Verge URL" / "Clash Verge 订阅"
	OtherGuidesHeading string
	HomeLinkText       string
	DownloadLabelTpl   string // Printf template with one %s for ClientName
	UpdatedLabel       string // "updated" / "更新"

	Steps []guideStep
	Tips  []qaItem
}

type guideStep struct {
	Title string
	Body  string // HTML allowed
}

type qaItem struct {
	Q string
	A string // HTML allowed
}

// supportedLocales lists the codes we render guide pages for. Order matters
// for hreflang alternates.
var supportedLocales = []string{"en", "zh"}

// localeSuffix returns the filename suffix for a locale: "" for English,
// ".zh" for Chinese. This keeps /guides/clash-verge.html as canonical English
// while /guides/clash-verge.zh.html serves Chinese.
func localeSuffix(loc string) string {
	if loc == "en" {
		return ""
	}
	return "." + loc
}

// localeLangAttr maps a locale code to an HTML lang attribute value.
func localeLangAttr(loc string) string {
	switch loc {
	case "zh":
		return "zh-Hans"
	default:
		return "en"
	}
}

var guides = []guideSpec{
	{
		Slug:        "clash-verge",
		ClientName:  "Clash Verge",
		OSList:      "Windows, macOS, Linux",
		DownloadURL: "https://github.com/clash-verge-rev/clash-verge-rev/releases/latest",
		URLField:    "clash",
		L10n: map[string]guideL10n{
			"en": {
				Title:              "How to use a free VPN with Clash Verge (Windows / macOS / Linux)",
				Description:        "Step-by-step guide to import a free VPN subscription into Clash Verge on Windows, macOS, or Linux. TCP+TLS verified nodes, updated hourly.",
				Keywords:           "clash verge, free vpn clash, clash subscription url, clash windows, clash macos, mihomo",
				StepsHeading:       "🪜 Steps",
				TipsHeading:        "💡 Tips & Troubleshooting",
				SubscribeHeading:   "📋 Subscription URL",
				SubscribeIntro:     "Copy this URL and paste it into Clash Verge's subscription import field:",
				SubscribeLabel:     "Clash Verge URL",
				OtherGuidesHeading: "📚 Other guides",
				HomeLinkText:       "← Home",
				DownloadLabelTpl:   "⬇ Download %s",
				UpdatedLabel:       "updated",
				Steps: []guideStep{
					{Title: "Install Clash Verge", Body: `Download the latest release from <a href="https://github.com/clash-verge-rev/clash-verge-rev/releases/latest" target="_blank" rel="noopener">Clash Verge Rev</a>. Pick the installer matching your OS (<code>.msi</code> for Windows, <code>.dmg</code> for macOS, <code>.deb</code>/<code>.rpm</code>/<code>.AppImage</code> for Linux). Run it and launch the app.`},
					{Title: "Open the Subscriptions panel", Body: `In the left sidebar, click <strong>Profiles</strong>. You'll see an empty list the first time. Click the <strong>+</strong> or <strong>Import</strong> button at the top.`},
					{Title: "Paste the subscription URL", Body: `Copy the Clash URL from the <a href="../index.html">homepage</a> and paste it into the URL field. Give it any name you like. Click <strong>Save</strong>. Clash Verge will download the profile in a few seconds.`},
					{Title: "Enable the profile", Body: `Click the profile card so it shows as active. In the <strong>Proxies</strong> tab you should now see a selector group with all the free nodes.`},
					{Title: "Turn on system proxy", Body: `Back on the dashboard, toggle <strong>System Proxy</strong> or <strong>TUN Mode</strong> ON. System proxy works for browsers; TUN mode routes all traffic including CLI tools and games — TUN mode requires admin/sudo the first time.`},
					{Title: "Test it", Body: `Open a browser and go to <code>ipinfo.io</code>. The country should match the node you picked. If it still shows your home country, double-check the proxy toggle and that your browser isn't bypassing system proxy (Firefox has its own proxy settings).`},
				},
				Tips: []qaItem{
					{Q: "Which node should I pick?", A: `Start with the selector group's auto-select (URL-Test) — it picks the fastest by latency. If a node feels slow, switch to a different country manually.`},
					{Q: "How do I update the subscription?", A: `Right-click the profile card → <strong>Update</strong>. Or enable auto-update in the profile settings (we recommend every 1 hour since upstream is refreshed hourly).`},
					{Q: "Can I use this with browsers only (no TUN)?", A: `Yes. Turn on <strong>System Proxy</strong> only. Most browsers respect it. Chrome/Edge use system proxy by default; Firefox needs "Use system proxy settings" in its Network settings.`},
					{Q: "Clash Verge says \"profile update failed\"", A: `Usually a transient GitHub raw content outage. Retry in a few minutes. If it persists, check that you can open the subscription URL directly in a browser.`},
				},
			},
			"zh": {
				Title:              "Clash Verge 导入免费 VPN 订阅教程 (Windows / macOS / Linux)",
				Description:        "一步步教你把免费 VPN 订阅导入 Clash Verge。所有节点每小时 TCP+TLS 实测过滤,Windows / macOS / Linux 通用。",
				Keywords:           "clash verge 教程, clash 订阅, 免费 clash 订阅, clash 免费节点, clash windows, clash macos, mihomo, 免费机场 clash",
				StepsHeading:       "🪜 操作步骤",
				TipsHeading:        "💡 常见问题排错",
				SubscribeHeading:   "📋 订阅链接",
				SubscribeIntro:     "复制下面这条 URL,粘贴到 Clash Verge 的订阅导入框:",
				SubscribeLabel:     "Clash Verge 订阅链接",
				OtherGuidesHeading: "📚 其他客户端教程",
				HomeLinkText:       "← 返回首页",
				DownloadLabelTpl:   "⬇ 下载 %s",
				UpdatedLabel:       "更新",
				Steps: []guideStep{
					{Title: "安装 Clash Verge", Body: `去 <a href="https://github.com/clash-verge-rev/clash-verge-rev/releases/latest" target="_blank" rel="noopener">Clash Verge Rev</a> 下载最新版本。Windows 选 <code>.msi</code>,macOS 选 <code>.dmg</code>,Linux 选 <code>.deb</code>/<code>.rpm</code>/<code>.AppImage</code>。双击安装后打开。`},
					{Title: "打开订阅面板", Body: `点左侧栏的 <strong>订阅</strong> (Profiles)。第一次打开是空的。点顶部的 <strong>+</strong> 或 <strong>导入</strong>。`},
					{Title: "粘贴订阅链接", Body: `从<a href="../index.html">首页</a>复制 Clash 的订阅 URL,粘贴到 URL 输入框。名称随便填。点 <strong>保存</strong>,几秒后配置下载完成。`},
					{Title: "启用订阅", Body: `点一下刚才的订阅卡片,让它变成激活状态(边框变蓝)。切到 <strong>代理</strong> (Proxies) 标签,能看到选择器组里列出了所有免费节点。`},
					{Title: "打开系统代理", Body: `回到<strong>主页</strong>,打开 <strong>系统代理</strong> 或 <strong>TUN 模式</strong> 开关。系统代理只对浏览器生效;TUN 模式接管所有流量(包括终端/游戏),首次启用需要管理员权限。`},
					{Title: "验证", Body: `浏览器打开 <code>ipinfo.io</code>,显示的国家应该和你选的节点一致。如果还是国内 IP,检查系统代理开关,以及浏览器有没有绕过系统代理(Firefox 要手动勾 "使用系统代理设置")。`},
				},
				Tips: []qaItem{
					{Q: "选哪个节点最好?", A: `先用选择器组里的 URL-Test (自动选最快),它会按延迟挑最快那个。觉得慢就手动换一个国家。`},
					{Q: "怎么更新订阅?", A: `右键订阅卡片 → <strong>更新</strong>。或在订阅设置里开自动更新,建议间隔 1 小时(上游每小时刷新)。`},
					{Q: "只用浏览器不想开 TUN 可以吗?", A: `可以。只打开 <strong>系统代理</strong> 就行,大部分浏览器都会走。Chrome/Edge 默认走系统代理;Firefox 要在网络设置里勾 "使用系统代理"。`},
					{Q: `提示 "profile update failed" 订阅更新失败`, A: `通常是 GitHub raw 临时抽风,等几分钟再试。如果一直失败,把订阅 URL 粘到浏览器看能不能打开,确认网络可达。`},
				},
			},
		},
	},
	{
		Slug:        "v2rayng",
		ClientName:  "v2rayNG",
		OSList:      "Android 5.0+",
		DownloadURL: "https://github.com/2dust/v2rayNG/releases/latest",
		URLField:    "v2ray",
		L10n: map[string]guideL10n{
			"en": {
				Title:              "How to use a free VPN with v2rayNG (Android)",
				Description:        "Import a free VPN subscription into v2rayNG on Android. Works on any Android 5.0+ phone, no root required.",
				Keywords:           "v2rayng, android free vpn, v2ray android, vmess android, vless android, free vpn apk",
				StepsHeading:       "🪜 Steps",
				TipsHeading:        "💡 Tips & Troubleshooting",
				SubscribeHeading:   "📋 Subscription URL",
				SubscribeIntro:     "Copy this URL and paste it into v2rayNG's subscription group:",
				SubscribeLabel:     "v2rayNG URL",
				OtherGuidesHeading: "📚 Other guides",
				HomeLinkText:       "← Home",
				DownloadLabelTpl:   "⬇ Download %s",
				UpdatedLabel:       "updated",
				Steps: []guideStep{
					{Title: "Install v2rayNG", Body: `Download the latest <code>.apk</code> from <a href="https://github.com/2dust/v2rayNG/releases/latest" target="_blank" rel="noopener">v2rayNG releases</a>. Install it (Android may ask you to allow "Install unknown apps" for your browser). v2rayNG is also on Google Play in some regions.`},
					{Title: "Copy the subscription URL", Body: `On the <a href="../index.html">homepage</a>, long-press the <strong>v2rayN / v2rayNG</strong> URL and select <strong>Copy link</strong>. This is the v2ray-base64 URL.`},
					{Title: "Add it as a subscription group", Body: `Open v2rayNG. Tap <strong>≡</strong> → <strong>Subscription group setting</strong> → <strong>+</strong>. Name it, paste the URL into <strong>URL</strong>, leave other fields default, tap <strong>✓</strong>.`},
					{Title: "Update subscriptions", Body: `Back on the main screen, tap <strong>≡</strong> → <strong>Update subscriptions (no proxy)</strong>. After a few seconds you'll see the free nodes listed.`},
					{Title: "Test one and connect", Body: `Tap any server row, then tap the big <strong>V</strong> button at the bottom right to start the VPN. Android will ask for VPN permission once — tap OK. The V turns green when connected.`},
					{Title: "Verify the country", Body: `Open a browser and visit <code>ipinfo.io</code>. The shown country should match the node. If it still shows your home IP, tap V again to disconnect then reconnect.`},
				},
				Tips: []qaItem{
					{Q: "Should I use per-app VPN?", A: `For casual browsing, no — default (global) works fine. If you want to keep banking or local apps off the VPN, enable <strong>Settings → Per-app proxy</strong> and exclude them.`},
					{Q: "How do I find the fastest server?", A: `Tap <strong>≡</strong> → <strong>Real ping test</strong>. v2rayNG will measure each server. Sort by ping, pick the lowest.`},
					{Q: "One server stopped working", A: `Tap it and pick a different one from the same country. Free nodes rotate — we refresh the list every hour.`},
					{Q: "Can I use v2rayNG on Android TV?", A: `Yes — install via a sideloader or adb. Navigation is a bit awkward with a remote; consider Clash Meta for Android TV instead if you have issues.`},
				},
			},
			"zh": {
				Title:              "v2rayNG 安卓免费 VPN 订阅配置教程",
				Description:        "一步步教你把免费 VPN 订阅导入安卓 v2rayNG。无需 root,支持 Android 5.0 以上所有机型。",
				Keywords:           "v2rayng 教程, v2rayng 安卓, 安卓 免费 vpn, v2ray 安卓, vmess 安卓, vless 安卓, 安卓 免费机场, v2rayng apk",
				StepsHeading:       "🪜 操作步骤",
				TipsHeading:        "💡 常见问题排错",
				SubscribeHeading:   "📋 订阅链接",
				SubscribeIntro:     "复制下面这条 URL,粘贴到 v2rayNG 的订阅组设置:",
				SubscribeLabel:     "v2rayNG 订阅链接",
				OtherGuidesHeading: "📚 其他客户端教程",
				HomeLinkText:       "← 返回首页",
				DownloadLabelTpl:   "⬇ 下载 %s",
				UpdatedLabel:       "更新",
				Steps: []guideStep{
					{Title: "安装 v2rayNG", Body: `去 <a href="https://github.com/2dust/v2rayNG/releases/latest" target="_blank" rel="noopener">v2rayNG Releases</a> 下载最新 <code>.apk</code>。安装时安卓可能提示"允许安装未知来源应用",给浏览器授权即可。部分地区的 Google Play 也能装。`},
					{Title: "复制订阅链接", Body: `在<a href="../index.html">首页</a>,长按 <strong>v2rayN / v2rayNG</strong> 那一行的 URL,选 <strong>复制链接</strong>。这是 v2ray-base64 格式。`},
					{Title: "添加为订阅组", Body: `打开 v2rayNG,点 <strong>≡</strong> → <strong>订阅设置</strong> → <strong>+</strong>。名字随便填,把 URL 粘进 <strong>URL</strong>,其他默认,点 <strong>✓</strong>。`},
					{Title: "更新订阅", Body: `回主界面,点 <strong>≡</strong> → <strong>更新订阅 (不使用代理)</strong>。几秒后能看到节点列表。`},
					{Title: "测试并连接", Body: `点任意一个服务器,然后点右下角大 <strong>V</strong> 按钮启动 VPN。安卓会弹 VPN 权限提示,点"允许"。V 变绿就是连上了。`},
					{Title: "验证国家", Body: `浏览器打开 <code>ipinfo.io</code>,显示的国家应该和你选的节点一致。如果还是国内 IP,点 V 断开重连一下。`},
				},
				Tips: []qaItem{
					{Q: "要不要开分应用代理?", A: `普通上网默认全局就行。想让银行 App、本地应用不走 VPN,在 <strong>设置 → 分应用代理</strong> 把它们排除。`},
					{Q: "怎么找最快的节点?", A: `点 <strong>≡</strong> → <strong>真连接延迟测试</strong>,v2rayNG 会测每个节点。按延迟排序,选最低那个。`},
					{Q: "某个节点突然连不上", A: `换同一国家的另一个节点。免费节点会轮换,我们每小时刷新一次列表。`},
					{Q: "Android TV 能用 v2rayNG 吗?", A: `能,但要用 sideload 或 adb 安装。遥控器操作有点别扭,Android TV 建议用 Clash Meta for Android TV。`},
				},
			},
		},
	},
	{
		Slug:        "shadowrocket",
		ClientName:  "Shadowrocket",
		OSList:      "iOS 15+, iPadOS 15+",
		DownloadURL: "https://apps.apple.com/app/shadowrocket/id932747118",
		URLField:    "v2ray",
		L10n: map[string]guideL10n{
			"en": {
				Title:              "How to use a free VPN with Shadowrocket (iOS / iPadOS)",
				Description:        "Import a free VPN subscription into Shadowrocket on iPhone or iPad. Requires a non-China Apple ID.",
				Keywords:           "shadowrocket, ios free vpn, iphone free vpn, ipad vpn, shadowrocket subscription",
				StepsHeading:       "🪜 Steps",
				TipsHeading:        "💡 Tips & Troubleshooting",
				SubscribeHeading:   "📋 Subscription URL",
				SubscribeIntro:     "Copy this URL and paste it into Shadowrocket's subscribe field:",
				SubscribeLabel:     "Shadowrocket URL",
				OtherGuidesHeading: "📚 Other guides",
				HomeLinkText:       "← Home",
				DownloadLabelTpl:   "⬇ Download %s",
				UpdatedLabel:       "updated",
				Steps: []guideStep{
					{Title: "Install Shadowrocket", Body: `Shadowrocket is a paid app (~$2.99) on the App Store. It's <strong>not available on the China App Store</strong> — sign into an Apple ID from a supported region (US, JP, HK, etc.). If you can't, use the free <a href="https://apps.apple.com/app/loon/id1373567447" target="_blank" rel="noopener">Loon</a> alternative or <a href="https://apps.apple.com/app/sing-box/id6451272673" target="_blank" rel="noopener">sing-box</a> (free).`},
					{Title: "Copy the v2ray subscription URL", Body: `On the <a href="../index.html">homepage</a>, tap and hold the <strong>v2rayN / v2rayNG / Shadowrocket</strong> URL, pick <strong>Copy</strong>. This is the v2ray base64 format that Shadowrocket understands.`},
					{Title: "Import via Subscribe", Body: `Open Shadowrocket → tap <strong>+</strong> (top right) → <strong>Type</strong> → <strong>Subscribe</strong>. Paste the URL into <strong>URL</strong>, name it anything, tap <strong>Save</strong>.`},
					{Title: "Update the subscription", Body: `Back on the main list, pull down to refresh, or tap the subscription row → <strong>Update</strong>. You should see dozens of nodes appear.`},
					{Title: "Choose connection mode", Body: `At the bottom tap <strong>Config</strong>. <strong>Global Proxy</strong> routes everything; <strong>Proxy (Rule-based)</strong> only proxies what the ruleset says. Toggle the main switch on the home screen to connect. iOS will request VPN permission the first time.`},
					{Title: "Verify", Body: `Safari → <code>ipinfo.io</code>. Country should match your chosen node. If not, try another node — free nodes sometimes have broken routing on iOS's stricter TLS stack.`},
				},
				Tips: []qaItem{
					{Q: "The App Store says Shadowrocket isn't available in my region", A: `You need an Apple ID from the US, Japan, Hong Kong, or similar. Create a new Apple ID with a non-China region and no payment method. Alternatively, use sing-box (free) or Loon.`},
					{Q: "Auto-update the subscription?", A: `Tap the subscription row → enable <strong>Auto Update</strong>. Set interval to <strong>1 hour</strong> since we refresh every hour.`},
					{Q: "How do I pick the fastest node?", A: `Tap any server → <strong>Latency Test</strong> (the lightning icon). Or long-press the subscription → <strong>Latency Test All</strong>.`},
					{Q: "Nothing loads after I connect", A: `Toggle the VPN switch off and on once. If still broken, try a different node. VLESS+Reality and Hysteria2 sometimes need ruleset adjustments on iOS.`},
				},
			},
			"zh": {
				Title:              "Shadowrocket iPhone 免费 VPN 订阅配置教程",
				Description:        "一步步教你把免费 VPN 订阅导入 iOS Shadowrocket (小火箭)。需要非中国区 Apple ID。",
				Keywords:           "shadowrocket 教程, 小火箭, iphone 免费 vpn, ios 免费机场, shadowrocket 订阅, ipad 翻墙, 小火箭 订阅",
				StepsHeading:       "🪜 操作步骤",
				TipsHeading:        "💡 常见问题排错",
				SubscribeHeading:   "📋 订阅链接",
				SubscribeIntro:     "复制下面这条 URL,粘贴到 Shadowrocket 的订阅导入框:",
				SubscribeLabel:     "Shadowrocket 订阅链接",
				OtherGuidesHeading: "📚 其他客户端教程",
				HomeLinkText:       "← 返回首页",
				DownloadLabelTpl:   "⬇ 下载 %s",
				UpdatedLabel:       "更新",
				Steps: []guideStep{
					{Title: "安装 Shadowrocket", Body: `Shadowrocket (小火箭) 是 App Store 付费应用,约 $2.99。<strong>中国区 App Store 没有</strong>,得用非中国区 Apple ID 登录(美区、日区、港区都行)。没非中国区账号的,改用免费 <a href="https://apps.apple.com/app/loon/id1373567447" target="_blank" rel="noopener">Loon</a> 或 <a href="https://apps.apple.com/app/sing-box/id6451272673" target="_blank" rel="noopener">sing-box</a> (免费)。`},
					{Title: "复制 v2ray 订阅链接", Body: `在<a href="../index.html">首页</a>,长按 <strong>v2rayN / v2rayNG / Shadowrocket</strong> 那行 URL,选 <strong>复制</strong>。这是 v2ray base64 格式,Shadowrocket 认得。`},
					{Title: "通过 Subscribe 导入", Body: `打开 Shadowrocket → 点右上 <strong>+</strong> → <strong>Type</strong> → <strong>Subscribe</strong>。把 URL 粘进 <strong>URL</strong>,名字随便,点 <strong>Save</strong>。`},
					{Title: "更新订阅", Body: `回主界面,下拉刷新,或点订阅那一行 → <strong>Update</strong>。能看到几十个节点。`},
					{Title: "选连接模式", Body: `底部点 <strong>Config</strong>。<strong>Global Proxy</strong> 是全局;<strong>Proxy (Rule-based)</strong> 是按规则走(推荐)。回主屏幕拨开关连接。第一次 iOS 会弹 VPN 权限,点允许。`},
					{Title: "验证", Body: `Safari 打开 <code>ipinfo.io</code>,国家应该和所选节点一致。如果没变,换一个节点——iOS 的 TLS 栈比较严,部分免费节点在 iOS 上会抽风。`},
				},
				Tips: []qaItem{
					{Q: "App Store 显示 Shadowrocket 在该地区不可用", A: `需要美区/日区/港区 Apple ID。自己注册一个非中国区账号,不用绑定付款方式。实在搞不定就用免费的 sing-box 或 Loon。`},
					{Q: "订阅怎么自动更新?", A: `点订阅那一行 → 打开 <strong>Auto Update</strong>,间隔设 <strong>1 小时</strong> (和我们上游刷新频率一致)。`},
					{Q: "怎么挑最快的节点?", A: `点任意服务器 → <strong>Latency Test</strong> (闪电图标)。或长按订阅 → <strong>Latency Test All</strong> 批量测。`},
					{Q: "连上了但打不开网页", A: `VPN 开关关了再开一次。还是不行就换节点。VLESS+Reality 和 Hysteria2 在 iOS 上偶尔需要调分流规则。`},
				},
			},
		},
	},
	{
		Slug:        "sing-box",
		ClientName:  "sing-box",
		OSList:      "Windows, macOS, Linux, iOS, Android",
		DownloadURL: "https://github.com/SagerNet/sing-box/releases/latest",
		URLField:    "singbox",
		L10n: map[string]guideL10n{
			"en": {
				Title:              "How to use a free VPN with sing-box (cross-platform, free)",
				Description:        "Run a free VPN subscription through sing-box on Windows, macOS, Linux, iOS, or Android. Completely free.",
				Keywords:           "sing-box, sing box free vpn, sing-box subscription, sing-box ios, sing-box android, free vpn cross platform",
				StepsHeading:       "🪜 Steps",
				TipsHeading:        "💡 Tips & Troubleshooting",
				SubscribeHeading:   "📋 Subscription URL",
				SubscribeIntro:     "Copy this URL and paste it into sing-box's profile import:",
				SubscribeLabel:     "sing-box URL",
				OtherGuidesHeading: "📚 Other guides",
				HomeLinkText:       "← Home",
				DownloadLabelTpl:   "⬇ Download %s",
				UpdatedLabel:       "updated",
				Steps: []guideStep{
					{Title: "Pick a client", Body: `sing-box ships as a CLI (Windows/macOS/Linux) and native apps for iOS/Android. For desktop GUIs try <a href="https://github.com/hiddify/hiddify-app/releases" target="_blank" rel="noopener">Hiddify Next</a>. On iOS grab the free <a href="https://apps.apple.com/app/sing-box/id6451272673" target="_blank" rel="noopener">sing-box app</a>; on Android grab it from <a href="https://github.com/SagerNet/sing-box/releases/latest" target="_blank" rel="noopener">GitHub releases</a>.`},
					{Title: "Copy the sing-box subscription URL", Body: `On the <a href="../index.html">homepage</a>, copy the <strong>sing-box</strong> URL (ends in <code>singbox.json</code>).`},
					{Title: "Import (GUI route)", Body: `In Hiddify Next or the sing-box app: <strong>Profiles</strong> → <strong>Add Profile</strong> → paste the URL. Enable the profile, toggle the VPN switch.`},
					{Title: "Import (CLI route)", Body: `Save the JSON locally: <code>curl -o config.json &lt;subscription URL&gt;</code>. Run: <code>sing-box run -c config.json</code>. By default this starts a SOCKS5 + HTTP proxy on <code>127.0.0.1:2080</code>/<code>2081</code>.`},
					{Title: "Enable TUN (optional, CLI route)", Body: `For full-system routing, edit <code>config.json</code> and add a TUN inbound (see <a href="https://sing-box.sagernet.org/configuration/inbound/tun/" target="_blank" rel="noopener">docs</a>). Run <code>sing-box</code> as root/admin. All TCP+UDP routes through the selected node.`},
					{Title: "Verify", Body: `<code>curl --proxy socks5://127.0.0.1:2080 https://ipinfo.io</code> (CLI) or open a browser (GUI). Country should match your picked node.`},
				},
				Tips: []qaItem{
					{Q: "Why use sing-box over Clash?", A: `sing-box has first-class support for modern protocols (VLESS+Reality, Hysteria2, TUIC) and is lighter. Clash Meta supports them too but sing-box's implementation tends to be closer to upstream specs.`},
					{Q: "Can I auto-update the subscription?", A: `In GUI: enable auto-update in profile settings, set 1 hour. For CLI: add a cron entry — <code>0 * * * * curl -o /etc/sing-box/config.json &lt;URL&gt; && systemctl reload sing-box</code>.`},
					{Q: "macOS menu bar client?", A: `Try <a href="https://github.com/hiddify/hiddify-app/releases" target="_blank" rel="noopener">Hiddify Next</a> or <a href="https://github.com/yichengchen/ClashX/releases" target="_blank" rel="noopener">ClashX Meta</a>.`},
					{Q: "Android has no system proxy — how does sing-box route traffic?", A: `The sing-box Android app uses VPNService, same as v2rayNG. All traffic goes through it when the toggle is on.`},
				},
			},
			"zh": {
				Title:              "sing-box 跨平台免费 VPN 订阅配置教程",
				Description:        "一步步教你用 sing-box 跑免费 VPN 订阅。Windows、macOS、Linux、iOS、Android 全平台可用,完全免费。",
				Keywords:           "sing-box 教程, sing-box 订阅, sing-box ios, sing-box android, sing-box macos, 免费 vpn 跨平台, hiddify, reality, hysteria2",
				StepsHeading:       "🪜 操作步骤",
				TipsHeading:        "💡 常见问题排错",
				SubscribeHeading:   "📋 订阅链接",
				SubscribeIntro:     "复制下面这条 URL,粘贴到 sing-box 的订阅导入:",
				SubscribeLabel:     "sing-box 订阅链接",
				OtherGuidesHeading: "📚 其他客户端教程",
				HomeLinkText:       "← 返回首页",
				DownloadLabelTpl:   "⬇ 下载 %s",
				UpdatedLabel:       "更新",
				Steps: []guideStep{
					{Title: "挑一种客户端", Body: `sing-box 有两种形态:命令行二进制 (Windows/macOS/Linux) 和移动端原生 App (iOS/Android)。桌面想要图形界面,装 <a href="https://github.com/hiddify/hiddify-app/releases" target="_blank" rel="noopener">Hiddify Next</a>。iOS 装免费的 <a href="https://apps.apple.com/app/sing-box/id6451272673" target="_blank" rel="noopener">sing-box</a>;安卓去 <a href="https://github.com/SagerNet/sing-box/releases/latest" target="_blank" rel="noopener">GitHub Releases</a> 下 APK。`},
					{Title: "复制 sing-box 订阅链接", Body: `在<a href="../index.html">首页</a>复制 <strong>sing-box</strong> 那一行的 URL (结尾是 <code>singbox.json</code>)。`},
					{Title: "导入 (图形界面路线)", Body: `打开 Hiddify Next 或 sing-box App:<strong>Profiles</strong> → <strong>Add Profile</strong> → 粘 URL。启用 profile,拨 VPN 开关。`},
					{Title: "导入 (命令行路线)", Body: `把 JSON 保存到本地:<code>curl -o config.json &lt;订阅 URL&gt;</code>。运行:<code>sing-box run -c config.json</code>。默认在 <code>127.0.0.1:2080</code>/<code>2081</code> 起 SOCKS5 + HTTP 代理。浏览器或系统代理指到这里。`},
					{Title: "开启 TUN (可选,CLI 路线)", Body: `想全局接管所有流量,改 <code>config.json</code> 加一个 TUN inbound (参考<a href="https://sing-box.sagernet.org/configuration/inbound/tun/" target="_blank" rel="noopener">官方文档</a>)。<code>sing-box</code> 以 root/管理员运行。所有 TCP+UDP 都走所选节点。`},
					{Title: "验证", Body: `<code>curl --proxy socks5://127.0.0.1:2080 https://ipinfo.io</code> (CLI) 或直接打开浏览器 (GUI),国家应该匹配所选节点。`},
				},
				Tips: []qaItem{
					{Q: "sing-box 和 Clash 选哪个?", A: `sing-box 对新协议 (VLESS+Reality、Hysteria2、TUIC) 原生支持更好,体量更轻。Clash Meta 也支持这些协议,但 sing-box 实现更贴近官方规范。`},
					{Q: "能自动更新订阅吗?", A: `GUI 里开自动更新,间隔设 1 小时。CLI 加个 cron:<code>0 * * * * curl -o /etc/sing-box/config.json &lt;URL&gt; && systemctl reload sing-box</code>。`},
					{Q: "macOS 有菜单栏版吗?", A: `用 <a href="https://github.com/hiddify/hiddify-app/releases" target="_blank" rel="noopener">Hiddify Next</a> 或 <a href="https://github.com/yichengchen/ClashX/releases" target="_blank" rel="noopener">ClashX Meta</a>。sing-box 自己暂时没发 macOS 原生 GUI。`},
					{Q: "安卓没系统代理,sing-box 怎么接管流量?", A: `sing-box 安卓 App 用 VPNService,和 v2rayNG 一样。开关打开就全局生效。`},
				},
			},
		},
	},
}
