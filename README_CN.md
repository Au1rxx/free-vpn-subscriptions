# Free VPN Subscriptions

[English](./README.md) · **简体中文** · [日本語](./README_JA.md) · [한국어](./README_KO.md) · [Español](./README_ES.md) · [Português](./README_PT.md) · [Русский](./README_RU.md)

<p align="center"><img src="https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/assets/hero.png" alt="Free VPN Subscriptions — hourly-refreshed free VPN subscriptions for Clash, sing-box, v2ray" width="780"></p>

![节点](https://img.shields.io/badge/节点-150-brightgreen) ![存活](https://img.shields.io/badge/存活-2579-blue) ![中位延迟](https://img.shields.io/badge/中位延迟-8ms-orange) ![更新](https://img.shields.io/badge/更新-2026-04-19_11:05_UTC-informational)

> **获取可用免费 VPN 的最简单方式 —— 复制订阅链接,粘贴到客户端,连上。**  
> 无需注册。无需付费。无需安装任何二进制。每小时从公共源自动抓取,发布前每个节点都经过 TCP + TLS 探测。

> 免费 VPN 订阅 · 免费机场 · 免费梯子 · 免费科学上网 · Clash 订阅 · v2ray 订阅 · sing-box 订阅 · VLESS Reality · VMess · Trojan · Shadowsocks · Hysteria2 · 每小时刷新 · TCP+TLS 探测 · 按国家分类

## 💡 为什么用这个项目?

GitHub 上几乎所有的"免费 VPN"列表都有三个问题:数据过期、全是死节点、或者要你装来路不明的二进制。本仓库**只发布几分钟前通过 TCP 握手并通过 TLS 握手的节点**,来源于筛选过的公共订阅,按延迟排序。直接给你 3 种通用订阅文件 —— 粘到 Clash / sing-box / v2rayN 即用。

> 📖 How the fetch → probe → rank pipeline works: [ARCHITECTURE.md](./ARCHITECTURE.md)

## 🔬 我们如何验证节点可用

**先说实话:我们无法 *保证* 每个节点一定能通流量。** 任何聚合项目做不到这一点,除非真正把流量打过去跑。下面把"我们做了什么、做不了什么、真正的保证来自哪里"全部讲清楚。

### ✅ 聚合阶段(发布前)的验证

1. **TCP 可达性** —— 对每个 `server:port` 发起一次 TCP 连接。服务器宕机、DNS 错误、端口被封全部会被丢掉。大约过滤掉 40% 的原始条目。
2. **TLS 握手** —— 对所有 TLS / Reality / WS-TLS 节点,完整跑一遍 TLS 握手。证书过期、SNI 不匹配、Reality short-id 失效都会被丢。再过滤掉大约 10%。
3. **按延迟排序** —— 幸存节点按 RTT 排序,取最快的前 N 个发布。

最近一次运行的典型数字:**17 个源 → ~4,800 原始 → ~2,900 TCP 存活 → ~2,600 TLS OK → 发布前 200**。

### ❌ 我们验证不了什么

- 代理协议的身份验证。UUID / 密码错误,上游服务器是在 TLS 握手 *之后* 才拒绝的,我们看不到。
- 节点真正跑 HTTP 代理的效果。
- 带宽 / 吞吐。
- 精确地理位置(只能通过出口 IP 的 GeoIP 粗略判断)。

### 🛡️ 运行时验证 —— 真正的保证在这里

我们发布的 `clash.yaml` 自带一个 `url-test` 组,**客户端每 5 分钟** 对每个节点真正打一次 HTTP:

```yaml
proxy-groups:
  - name: AUTO
    type: url-test
    url: http://www.gstatic.com/generate_204
    interval: 300
```

客户端会按 *真实* 的 HTTP 过代理延迟排序,自动挑最快可用节点。sing-box 和 v2ray 有等价机制。选中的节点挂了,客户端会自动切下一个,不需要你手动处理。

### 🧮 实际效果

每次发布的前 200 个节点里,客户端通常能找到 30-50 个当前能稳定过 HTTP 的。慢了就让 url-test 组换下一个,一键切换。

## 🚀 一键订阅

复制对应客户端的 URL,粘贴到订阅导入框:

| 客户端 | 格式 | 订阅链接 |
|---|---|---|
| Clash / Clash Verge / ClashX | `clash.yaml` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/clash.yaml` |
| sing-box | `singbox.json` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/singbox.json` |
| v2rayN / v2rayNG / Shadowrocket / NekoBox | `v2ray-base64` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/v2ray-base64.txt` |

## 🌍 按国家订阅

只想要特定地区的节点?选一个针对性订阅链接:

| 国家/地区 | 节点数 | Clash | sing-box | v2ray |
|---|---|---|---|---|
| 🇺🇸 United States (`US`) | 21 | [clash-US.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-US.yaml) | [singbox-US.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-US.json) | [v2ray-base64-US.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-US.txt) |
| 🇩🇪 Germany (`DE`) | 7 | [clash-DE.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-DE.yaml) | [singbox-DE.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-DE.json) | [v2ray-base64-DE.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-DE.txt) |

## 📖 客户端图文教程

新手不知道怎么配?按平台选一篇跟着做:

- [**Clash Verge**](https://au1rxx.github.io/free-vpn-subscriptions/guides/clash-verge.zh.html) · Windows / macOS / Linux
- [**v2rayNG**](https://au1rxx.github.io/free-vpn-subscriptions/guides/v2rayng.zh.html) · Android
- [**Shadowrocket**](https://au1rxx.github.io/free-vpn-subscriptions/guides/shadowrocket.zh.html) · iOS / iPadOS
- [**sing-box**](https://au1rxx.github.io/free-vpn-subscriptions/guides/sing-box.zh.html) · Windows / macOS / Linux / iOS / Android

## 🧩 支持的客户端

- **Windows**:v2rayN、Clash Verge、Hiddify、NekoRay
- **macOS**:ClashX Pro、Clash Verge、sing-box、Hiddify
- **iOS**:Shadowrocket、Stash、Loon、sing-box、Hiddify
- **Android**:v2rayNG、NekoBox、Clash Meta for Android、Hiddify、sing-box
- **Linux**:mihomo (Clash.Meta)、sing-box、v2ray-core

## 📊 实时统计

- **精选节点数**: 150
- **全源存活总数**: 2579
- **最快延迟**: 2 ms
- **中位延迟**: 8 ms
- **最后更新 (UTC)**: 2026-04-19 11:05 UTC

**协议分布:** shadowsocks × 25 · trojan × 18 · vless × 85 · vmess × 22

**本次使用的源:** `barry-far-v2ray` × 30 · `ebrasha-v2ray` × 9 · `epodonios` × 33 · `freefq` × 1 · `mahdi0024` × 1 · `mahdibland-aggregator` × 1 · `mahdibland-shadowsocks` × 1 · `mfuu-clash` × 2 · `ninjastrikers` × 35 · `pawdroid` × 1 · `ruking-clash` × 21 · `snakem982` × 12 · `surfboard-eternity` × 2 · `vxiaov-clash` × 1

## ❓ 常见问题

<details><summary>真的完全免费吗?</summary>

是的。所有节点由第三方志愿者自己运营并公开免费订阅。本仓库不运营任何服务器,只是做测试、排名、重新打包公开内容。

</details>

<details><summary>数据多新?</summary>

每小时刷新一次(带一点随机延迟避免整点集中打上游):拉取所有上游源 → TCP+TLS 探测每个节点 → 丢弃死节点 → 按延迟排序 → 发布新的输出文件。见顶部徽章上的更新时间。

</details>

<details><summary>这些节点可以信任吗?</summary>

免费节点能看到你所有流量。**绝不要用来登录银行、邮箱等敏感账号。**用来突破地区限制访问公开内容没问题。真正需要隐私请自建 VPS 或付费服务。

</details>

<details><summary>列表里的节点为什么有的连不上?</summary>

我们只验证 TCP 可达和 TLS 握手 —— 节点仍可能配额用完、路由被污染、证书到期。发布的 `clash.yaml` 自带 `url-test` 组(每 300 秒对 `http://www.gstatic.com/generate_204` 打一次),客户端会自动选真正能过 HTTP 的最快节点。挂了就换下一个。

</details>

## 🤝 贡献

知道稳定的公共订阅源可以加入?提 issue 给我们 URL 和格式。

## ⚠️ 免责声明

本仓库聚合第三方志愿者**公开分享**的代理配置。我们不运营任何服务器,不保证可用性或安全性,不为使用行为负责。仅供学习和个人连接使用。请遵守所在司法管辖区的法律。

## ⭐ Star 历史

[![Star History Chart](https://api.star-history.com/svg?repos=Au1rxx/free-vpn-subscriptions&type=Date)](https://www.star-history.com/#Au1rxx/free-vpn-subscriptions&Date)

---

如果这个项目帮到你,点个 ⭐ —— 每一颗 star 都能帮更多人发现它。
