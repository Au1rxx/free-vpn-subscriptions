package pages

// guideSpec describes one client tutorial page. Rendered via tplGuide.
type guideSpec struct {
	Slug        string // filename without .html, used in URL
	Title       string // <title> and h1
	Description string // meta description
	Keywords    string
	ClientName  string // e.g. "Clash Verge"
	OSList      string // e.g. "Windows, macOS, Linux"
	DownloadURL string
	URLField    string // which subscription URL to import: "clash" | "singbox" | "v2ray"
	Steps       []guideStep
	Tips        []qaItem
}

type guideStep struct {
	Title string
	Body  string // HTML allowed
}

type qaItem struct {
	Q string
	A string // HTML allowed
}

var guides = []guideSpec{
	{
		Slug:        "clash-verge",
		Title:       "How to use a free VPN with Clash Verge (Windows / macOS / Linux)",
		Description: "Step-by-step guide to import a free VPN subscription into Clash Verge on Windows, macOS, or Linux. TCP+TLS verified nodes, updated hourly.",
		Keywords:    "clash verge, free vpn clash, clash subscription url, clash windows, clash macos, mihomo",
		ClientName:  "Clash Verge",
		OSList:      "Windows, macOS, Linux",
		DownloadURL: "https://github.com/clash-verge-rev/clash-verge-rev/releases/latest",
		URLField:    "clash",
		Steps: []guideStep{
			{
				Title: "Install Clash Verge",
				Body:  `Download the latest release from <a href="https://github.com/clash-verge-rev/clash-verge-rev/releases/latest" target="_blank" rel="noopener">Clash Verge Rev</a>. Pick the installer matching your OS (<code>.msi</code> for Windows, <code>.dmg</code> for macOS, <code>.deb</code>/<code>.rpm</code>/<code>.AppImage</code> for Linux). Run it and launch the app.`,
			},
			{
				Title: "Open the Subscriptions panel",
				Body:  `In the left sidebar, click <strong>Profiles</strong> (订阅). You'll see an empty list the first time. Click the <strong>+</strong> or <strong>Import</strong> button at the top.`,
			},
			{
				Title: "Paste the subscription URL",
				Body:  `Copy the Clash URL from the <a href="../index.html">homepage</a> and paste it into the URL field. Give it any name you like. Click <strong>Save</strong>. Clash Verge will download the profile in a few seconds.`,
			},
			{
				Title: "Enable the profile",
				Body:  `Click the profile card so it shows as active (a blue border or check mark). In the <strong>Proxies</strong> tab you should now see a selector group with all the free nodes.`,
			},
			{
				Title: "Turn on system proxy",
				Body:  `Go back to the dashboard (首页). Toggle <strong>System Proxy</strong> or <strong>TUN Mode</strong> ON. System proxy works for browsers; TUN mode routes all traffic including CLI tools and games — TUN mode requires admin/sudo the first time.`,
			},
			{
				Title: "Test it",
				Body:  `Open a browser and go to <code>ipinfo.io</code>. The country should match the node you picked. If it still shows your home country, double-check the proxy toggle and that your browser isn't bypassing system proxy (Firefox has its own proxy settings).`,
			},
		},
		Tips: []qaItem{
			{Q: "Which node should I pick?", A: `Start with the selector group's auto-select (URL-Test) — it picks the fastest one by latency. If a node feels slow, switch to a different country manually.`},
			{Q: "How do I update the subscription?", A: `Right-click the profile card → <strong>Update</strong>. Or enable auto-update in the profile settings (we recommend every 1 hour since the upstream is refreshed hourly).`},
			{Q: "Can I use this with browsers only (no TUN)?", A: `Yes. Turn on <strong>System Proxy</strong> only. Most browsers respect it. Chrome/Edge use system proxy by default; Firefox needs "Use system proxy settings" in its Network settings.`},
			{Q: "Clash Verge says \"profile update failed\"", A: `Usually a transient GitHub raw content outage. Retry in a few minutes. If it persists, check that you can open the subscription URL directly in a browser.`},
		},
	},
	{
		Slug:        "v2rayng",
		Title:       "How to use a free VPN with v2rayNG (Android)",
		Description: "Import a free VPN subscription into v2rayNG on Android. Works on any Android 5.0+ phone, no root required.",
		Keywords:    "v2rayng, android free vpn, v2ray android, vmess android, vless android, free vpn apk",
		ClientName:  "v2rayNG",
		OSList:      "Android 5.0+",
		DownloadURL: "https://github.com/2dust/v2rayNG/releases/latest",
		URLField:    "v2ray",
		Steps: []guideStep{
			{
				Title: "Install v2rayNG",
				Body:  `Download the latest <code>.apk</code> from <a href="https://github.com/2dust/v2rayNG/releases/latest" target="_blank" rel="noopener">v2rayNG releases</a>. Install it (Android may ask you to allow "Install unknown apps" for your browser). v2rayNG is also on Google Play in some regions.`,
			},
			{
				Title: "Copy the subscription URL",
				Body:  `On the <a href="../index.html">homepage</a>, long-press the <strong>v2rayN / v2rayNG</strong> URL and select <strong>Copy link</strong>. This is the v2ray-base64 URL.`,
			},
			{
				Title: "Add it as a subscription group",
				Body:  `Open v2rayNG. Tap the <strong>≡</strong> menu → <strong>Subscription group setting</strong> → <strong>+</strong>. Give it a name (e.g. "free nodes"), paste the URL into <strong>URL</strong>, leave other fields default, tap <strong>✓</strong>.`,
			},
			{
				Title: "Update subscriptions",
				Body:  `Back on the main screen, tap the <strong>≡</strong> menu → <strong>Update subscriptions (no proxy)</strong>. After a few seconds you'll see all the free nodes listed.`,
			},
			{
				Title: "Test one and connect",
				Body:  `Tap any server row, then tap the big <strong>V</strong> button at the bottom right to start the VPN. Android will ask for a VPN permission prompt once — tap OK. The V turns green when connected.`,
			},
			{
				Title: "Verify the country",
				Body:  `Open a browser and visit <code>ipinfo.io</code>. The shown country should match the node. If it still shows your home IP, tap V again to disconnect then reconnect — the VPN sometimes takes a moment on first use.`,
			},
		},
		Tips: []qaItem{
			{Q: "Should I use per-app VPN?", A: `For casual browsing, no — default (global) works fine. If you want to keep banking or local apps off the VPN, enable <strong>Settings → Per-app proxy</strong> and exclude them.`},
			{Q: "How do I find the fastest server?", A: `Tap <strong>≡</strong> → <strong>Real ping test</strong>. v2rayNG will measure each server. Sort by ping, pick the lowest.`},
			{Q: "One server stopped working", A: `Tap it and pick a different one from the same country. Free nodes rotate — we refresh the list every hour.`},
			{Q: "Can I use v2rayNG on Android TV?", A: `Yes — install via a sideloader or adb. Navigation is a bit awkward with a remote; consider Clash Meta for Android TV instead if you have issues.`},
		},
	},
	{
		Slug:        "shadowrocket",
		Title:       "How to use a free VPN with Shadowrocket (iOS / iPadOS)",
		Description: "Import a free VPN subscription into Shadowrocket on iPhone or iPad. Requires a non-China Apple ID.",
		Keywords:    "shadowrocket, ios free vpn, iphone free vpn, ipad vpn, shadowrocket subscription",
		ClientName:  "Shadowrocket",
		OSList:      "iOS 15+, iPadOS 15+",
		DownloadURL: "https://apps.apple.com/app/shadowrocket/id932747118",
		URLField:    "v2ray",
		Steps: []guideStep{
			{
				Title: "Install Shadowrocket",
				Body:  `Shadowrocket is a paid app (~$2.99) on the App Store. It's <strong>not available on the China App Store</strong> — sign into an Apple ID from a supported region (US, JP, HK, etc.). If you can't, use the free <a href="https://apps.apple.com/app/loon/id1373567447" target="_blank" rel="noopener">Loon</a> alternative or <a href="https://apps.apple.com/app/sing-box/id6451272673" target="_blank" rel="noopener">sing-box</a> (free).`,
			},
			{
				Title: "Copy the v2ray subscription URL",
				Body:  `On the <a href="../index.html">homepage</a>, tap and hold the <strong>v2rayN / v2rayNG / Shadowrocket</strong> URL, pick <strong>Copy</strong>. This is the v2ray base64 format that Shadowrocket understands.`,
			},
			{
				Title: "Import via Subscribe",
				Body:  `Open Shadowrocket → tap <strong>+</strong> (top right) → <strong>Type</strong> → <strong>Subscribe</strong>. Paste the URL into <strong>URL</strong>, name it anything, tap <strong>Save</strong>.`,
			},
			{
				Title: "Update the subscription",
				Body:  `Back on the main list, pull down to refresh, or tap the subscription row → <strong>Update</strong>. You should see dozens of nodes appear.`,
			},
			{
				Title: "Choose connection mode",
				Body:  `At the bottom tap <strong>Config</strong>. <strong>Global Proxy</strong> routes everything; <strong>Proxy (Rule-based)</strong> only proxies what the ruleset says (default is fine for most people). Toggle the main switch on the home screen to connect. iOS will request VPN permission the first time.`,
			},
			{
				Title: "Verify",
				Body:  `Safari → <code>ipinfo.io</code>. Country should match your chosen node. If not, try another node — free nodes sometimes have broken routing on iOS's stricter TLS stack.`,
			},
		},
		Tips: []qaItem{
			{Q: "The App Store says Shadowrocket isn't available in my region", A: `You need an Apple ID from the US, Japan, Hong Kong, or similar. Create a new Apple ID with a non-China region and no payment method — follow Apple's guide. Alternatively, use sing-box (free) or Loon.`},
			{Q: "Auto-update the subscription?", A: `Tap the subscription row → enable <strong>Auto Update</strong>. Set interval to <strong>1 hour</strong> since we refresh every hour.`},
			{Q: "How do I pick the fastest node?", A: `Tap any server → <strong>Latency Test</strong> (the lightning icon). Or long-press the subscription → <strong>Latency Test All</strong>.`},
			{Q: "Nothing loads after I connect", A: `A known iOS quirk — toggle the VPN switch off and on once. If still broken, try a different node. VLESS+Reality and Hysteria2 sometimes need ruleset adjustments on iOS.`},
		},
	},
	{
		Slug:        "sing-box",
		Title:       "How to use a free VPN with sing-box (cross-platform, free)",
		Description: "Run a free VPN subscription through sing-box on Windows, macOS, Linux, iOS, or Android. Completely free.",
		Keywords:    "sing-box, sing box free vpn, sing-box subscription, sing-box ios, sing-box android, free vpn cross platform",
		ClientName:  "sing-box",
		OSList:      "Windows, macOS, Linux, iOS, Android",
		DownloadURL: "https://github.com/SagerNet/sing-box/releases/latest",
		URLField:    "singbox",
		Steps: []guideStep{
			{
				Title: "Pick a client",
				Body:  `sing-box ships in two forms: a CLI binary (Windows/macOS/Linux) and native apps for iOS/Android. For desktop GUIs that wrap sing-box, try <a href="https://github.com/hiddify/hiddify-app/releases" target="_blank" rel="noopener">Hiddify Next</a>. On iOS grab the free <a href="https://apps.apple.com/app/sing-box/id6451272673" target="_blank" rel="noopener">sing-box app</a>; on Android grab it from <a href="https://github.com/SagerNet/sing-box/releases/latest" target="_blank" rel="noopener">GitHub releases</a>.`,
			},
			{
				Title: "Copy the sing-box subscription URL",
				Body:  `On the <a href="../index.html">homepage</a>, copy the <strong>sing-box</strong> URL (ends in <code>singbox.json</code>).`,
			},
			{
				Title: "Import (GUI route)",
				Body:  `In Hiddify Next or the sing-box app: <strong>Profiles</strong> → <strong>Add Profile</strong> → paste the URL. Enable the profile, toggle the VPN switch on the home screen. Done.`,
			},
			{
				Title: "Import (CLI route)",
				Body:  `Save the JSON locally: <code>curl -o config.json &lt;subscription URL&gt;</code>. Run: <code>sing-box run -c config.json</code>. By default this starts a SOCKS5 + HTTP proxy on <code>127.0.0.1:2080</code>/<code>2081</code>. Point your browser or system proxy there.`,
			},
			{
				Title: "Enable TUN (optional, CLI route)",
				Body:  `For full-system routing without per-app config, edit <code>config.json</code> and add a TUN inbound (see <a href="https://sing-box.sagernet.org/configuration/inbound/tun/" target="_blank" rel="noopener">docs</a>). Run <code>sing-box</code> as root/admin. All TCP+UDP traffic will route through the selected node.`,
			},
			{
				Title: "Verify",
				Body:  `<code>curl --proxy socks5://127.0.0.1:2080 https://ipinfo.io</code> (CLI) or open a browser (GUI). Country should match your picked node.`,
			},
		},
		Tips: []qaItem{
			{Q: "Why use sing-box over Clash?", A: `sing-box has first-class support for modern protocols (VLESS+Reality, Hysteria2, TUIC) and is lighter. Clash Meta supports them too but sing-box's implementation tends to be closer to upstream specs.`},
			{Q: "Can I auto-update the subscription?", A: `In Hiddify Next / sing-box app: enable auto-update in profile settings, set to 1 hour. For CLI: add a cron/systemd timer: <code>0 * * * * curl -o /etc/sing-box/config.json &lt;URL&gt; && systemctl reload sing-box</code>.`},
			{Q: "macOS menu bar client?", A: `Try <a href="https://github.com/hiddify/hiddify-app/releases" target="_blank" rel="noopener">Hiddify Next</a> or <a href="https://github.com/yichengchen/ClashX/releases" target="_blank" rel="noopener">ClashX Meta</a>. sing-box itself doesn't ship a native macOS GUI yet.`},
			{Q: "Android has no system proxy — how does sing-box route traffic?", A: `The sing-box Android app uses VPNService, same as v2rayNG and Clash. All traffic goes through it when the toggle is on.`},
		},
	},
}
