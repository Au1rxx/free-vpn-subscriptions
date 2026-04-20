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
	Hook2:       "无需注册。无需付费。无需安装任何二进制。每小时从公共源自动抓取 —— 每个发布的节点都在几分钟前通过 sing-box 真实转发过 HTTP 流量。",
	KeywordLine: "免费 VPN 订阅 · 免费机场 · 免费梯子 · 免费科学上网 · Clash 订阅 · v2ray 订阅 · sing-box 订阅 · VLESS Reality · VMess · Trojan · Shadowsocks · Hysteria2 · 每小时刷新 · HTTP 实测验证 · 按国家分类",

	WhyHeading: "## 💡 为什么用这个项目?",
	WhyBody:    "GitHub 上几乎所有的\"免费 VPN\"列表都有三个问题:数据过期、全是死节点、或者要你装来路不明的二进制。本仓库比任何其他列表都更进一步 —— **我们不只是检查节点端口能不能通,而是用 sing-box 把真实 HTTP 请求经代理打出去、确认能收到 204 才发布**,全部在几分钟内完成。直接给你 3 种通用订阅文件 —— 粘到 Clash / sing-box / v2rayN 即用。",

	VerificationHeading: "## 🔬 我们如何验证节点可用",
	VerificationBody: `多数免费 VPN 列表停在"端口能连上"就发布了。我们不这样。下面是节点必须跑通的完整验证管线。

### ✅ 聚合阶段(发布前)的验证

1. **TCP 可达性** —— 对每个 ` + "`server:port`" + ` 发起一次 TCP 连接。服务器宕机、DNS 错误、端口被封全部被丢。大约过滤掉 40% 的原始条目。
2. **TLS 握手** —— 对所有 TLS / Reality / WS-TLS 节点完整跑一遍 TLS 握手。证书过期、SNI 不匹配、Reality short-id 失效都会被丢。再过滤掉约 10%。
3. **sing-box 配置校验** —— 每个幸存节点都被翻译成真实的 sing-box outbound,过 ` + "`sing-box check`" + ` 验证。密码错、UUID 畸形、flow 选项不支持 —— 全部在占用探测资源之前剔除。
4. **HTTP 过代理实测(这是关键的一步)** —— 我们把最快的约 900 个候选节点分批塞进 sing-box 子进程,每个节点分配一个本地 SOCKS5 入口,然后通过它真正发 HTTP 和 HTTPS 请求:
   - ` + "`http://www.gstatic.com/generate_204`" + `(期望 204)
   - ` + "`https://www.cloudflare.com/cdn-cgi/trace`" + `(期望 200)

   请求会完整走一遍代理协议(VLESS / VMess / Trojan / Shadowsocks / Hysteria2),所以能过这一关的节点,就证明它的认证、路由、TLS 内握手、出口网络全部是工作的。
5. **两轮,间隔 45 秒** —— 过了一次但 45 秒后就死的节点会被筛掉。只有在(轮数 × 目标数)中成功率 ≥ 50% 的节点才会留下。
6. **按真实延迟中位数排序** —— 幸存节点按 HTTP 过代理的真实往返时间中位数(不是原始 TCP RTT)排序,取前 N 发布。

最近一次运行的典型数字:**17 个源 → ~4,800 原始 → ~2,900 TCP 存活 → ~2,600 TLS OK → ~840 配置有效 → ~280 HTTP 实测通过 → 发布前 150**。发布出去的 150 个节点,每一个都在过去十分钟内真正转发过流量。

### ❌ 我们仍然验证不了什么

- **带宽 / 吞吐** —— 我们测延迟不测 Mbps。50ms 的节点看视频可能仍然慢。
- **精确地理位置** —— GeoIP 能告诉你出口 IP 是哪国,但城市或 ISP 级别的判断不可靠。
- **特定地区的封锁** —— 我们的探测机器能通的节点,不代表你的网络也能通(ISP 层封锁、captive portal 等)。
- **发布之后是否还活着** —— 十分钟前它是活的,之后可能挂了。

### 🛡️ 运行时兜底 —— 对付上面最后一条

我们发布的 ` + "`clash.yaml`" + ` 自带 ` + "`url-test`" + ` 组,**在你本地**每 5 分钟重新跑一次 HTTP 实测:

` + "```yaml" + `
proxy-groups:
  - name: AUTO
    type: url-test
    url: http://www.gstatic.com/generate_204
    interval: 300
` + "```" + `

客户端按 *你的网络* 下实时的 HTTP 过代理延迟排序,自动挑最快可用节点。sing-box 和 v2ray 有等价机制。如果聚合到下一次运行中途某节点挂了,客户端会自动切下一个,不需要你手动处理。

### 🧮 实际效果

每次发布的约 150 个节点里,客户端通常能找到 **80-120 个在你的网络下能稳定过 HTTP**,比只做 TCP 探测的列表命中率高 2-3 倍。其中一个挂了 url-test 组会透明地切换。`,

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
	FAQ2A:      "每小时刷新一次(带一点随机延迟避免整点集中打上游):拉取所有源 → TCP → TLS → sing-box 配置校验 → HTTP 过代理探测(两轮,间隔 45 秒)→ 按真实 HTTP 延迟排序 → 发布新的输出文件。完整管线约 10 分钟。见顶部徽章上的更新时间。",
	FAQ3Q:      "这些节点可以信任吗?",
	FAQ3A:      "免费节点能看到你所有流量。**绝不要用来登录银行、邮箱等敏感账号。**用来突破地区限制访问公开内容没问题。真正需要隐私请自建 VPS 或付费服务。",
	FAQ4Q:      "列表里的节点为什么有的连不上?",
	FAQ4A:      "即使过了我们的 HTTP 实测,节点在聚合后也可能死掉:配额用尽、上游吊销了 key、你的 ISP 封锁了出口 IP,或者运营者干脆下架了。发布的 `clash.yaml` 自带 `url-test` 组(每 300 秒对 `http://www.gstatic.com/generate_204` 打一次),客户端会在 *你的网络下* 自动选真正能过 HTTP 的最快节点。挂了就换下一个。通常 150 个里随时有 80-120 个能用。",

	ContributingHeading: "## 🤝 贡献",
	ContributingBody:    "知道稳定的公共订阅源可以加入?提 issue 给我们 URL 和格式。",

	DisclaimerHeading: "## ⚠️ 免责声明",
	DisclaimerBody:    "本仓库聚合第三方志愿者**公开分享**的代理配置。我们不运营任何服务器,不保证可用性或安全性,不为使用行为负责。仅供学习和个人连接使用。请遵守所在司法管辖区的法律。",

	StarHistoryHeading: "## ⭐ Star 历史",
	FinalCTA:           "如果这个项目帮到你,点个 ⭐ —— 每一颗 star 都能帮更多人发现它。",
}
