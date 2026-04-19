# Free VPN Subscriptions

[English](./README.md) · **简体中文** · [日本語](./README_JA.md) · [한국어](./README_KO.md) · [Español](./README_ES.md) · [Português](./README_PT.md) · [Русский](./README_RU.md)

<p align="center"><img src="https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/assets/workflow.svg" alt="Free VPN Subscriptions — aggregation workflow" width="780"></p>

![节点](https://img.shields.io/badge/节点-150-brightgreen) ![存活](https://img.shields.io/badge/存活-394-blue) ![中位延迟](https://img.shields.io/badge/中位延迟-134ms-orange) ![更新](https://img.shields.io/badge/更新-2026-04-19_06:03_UTC-informational)

> **获取可用免费 VPN 的最简单方式 —— 复制订阅链接,粘贴到客户端,连上。**  
> 无需注册。无需付费。无需安装任何二进制。每小时从公共源自动抓取,逐个节点测试。

> 免费 VPN 订阅 · 免费机场 · 免费梯子 · 免费科学上网 · Clash 订阅 · v2ray 订阅 · sing-box 订阅 · VLESS Reality · VMess · Trojan · Shadowsocks · Hysteria2 · 每小时刷新 · TCP+TLS 探测 · 按国家分类

## 💡 为什么用这个项目?

GitHub 上几乎所有的"免费 VPN"列表都有三个问题:数据过期、全是死节点、或者要你装来路不明的二进制。本仓库**只发布几分钟前通过 TCP 握手并通过 TLS 握手的节点**,来源于筛选过的公共订阅,按延迟排序。直接给你 3 种通用订阅文件 —— 粘到 Clash / sing-box / v2rayN 即用。

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
| 🇺🇸 United States (`US`) | 27 | [clash-US.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-US.yaml) | [singbox-US.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-US.json) | [v2ray-base64-US.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-US.txt) |
| 🇬🇧 United Kingdom (`GB`) | 26 | [clash-GB.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-GB.yaml) | [singbox-GB.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-GB.json) | [v2ray-base64-GB.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-GB.txt) |
| 🇳🇱 Netherlands (`NL`) | 13 | [clash-NL.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-NL.yaml) | [singbox-NL.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-NL.json) | [v2ray-base64-NL.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-NL.txt) |
| 🇨🇦 Canada (`CA`) | 11 | [clash-CA.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-CA.yaml) | [singbox-CA.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-CA.json) | [v2ray-base64-CA.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-CA.txt) |
| 🇯🇵 Japan (`JP`) | 8 | [clash-JP.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-JP.yaml) | [singbox-JP.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-JP.json) | [v2ray-base64-JP.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-JP.txt) |
| 🇭🇰 Hong Kong (`HK`) | 6 | [clash-HK.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-HK.yaml) | [singbox-HK.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-HK.json) | [v2ray-base64-HK.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-HK.txt) |
| 🇹🇼 Taiwan (`TW`) | 6 | [clash-TW.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-TW.yaml) | [singbox-TW.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-TW.json) | [v2ray-base64-TW.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-TW.txt) |
| 🇩🇪 Germany (`DE`) | 5 | [clash-DE.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-DE.yaml) | [singbox-DE.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-DE.json) | [v2ray-base64-DE.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-DE.txt) |
| 🇰🇷 Korea (`KR`) | 4 | [clash-KR.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-KR.yaml) | [singbox-KR.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-KR.json) | [v2ray-base64-KR.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-KR.txt) |
| 🇲🇦 MA (`MA`) | 3 | [clash-MA.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-MA.yaml) | [singbox-MA.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-MA.json) | [v2ray-base64-MA.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-MA.txt) |

## 🧩 支持的客户端

- **Windows**:v2rayN、Clash Verge、Hiddify、NekoRay
- **macOS**:ClashX Pro、Clash Verge、sing-box、Hiddify
- **iOS**:Shadowrocket、Stash、Loon、sing-box、Hiddify
- **Android**:v2rayNG、NekoBox、Clash Meta for Android、Hiddify、sing-box
- **Linux**:mihomo (Clash.Meta)、sing-box、v2ray-core

## 📊 实时统计

- **精选节点数**: 150
- **全源存活总数**: 394
- **最快延迟**: 8 ms
- **中位延迟**: 134 ms
- **最后更新 (UTC)**: 2026-04-19 06:03 UTC

**协议分布:** shadowsocks × 102 · trojan × 7 · vmess × 41

**本次使用的源:** `freefq` × 2 · `mahdibland-aggregator` × 72 · `mahdibland-shadowsocks` × 58 · `pawdroid` × 3 · `vxiaov-clash` × 15

## ❓ 常见问题

<details><summary>真的完全免费吗?</summary>

是的。所有节点由第三方志愿者自己运营并公开免费订阅。本仓库不运营任何服务器,只是做测试、排名、重新打包公开内容。

</details>

<details><summary>数据多新?</summary>

GitHub Actions 每小时运行一次:拉取所有上游源 → TCP+TLS 探测每个节点 → 丢弃死节点 → 按延迟排序 → 提交新的输出文件。见顶部徽章上的更新时间。

</details>

<details><summary>这些节点可以信任吗?</summary>

免费节点能看到你所有流量。**绝不要用来登录银行、邮箱等敏感账号。**用来突破地区限制访问公开内容没问题。真正需要隐私请自建 VPS 或付费服务。

</details>

<details><summary>列表里的节点为什么有的连不上?</summary>

我们验证 TCP 可达和 TLS 握手,但节点仍可能配额用完、路由被污染、证书到期。多试几个,selector 组自带 fallback。

</details>

## 🤝 贡献

知道稳定的公共订阅源可以加入?提 issue 给我们 URL 和格式。

## ⚠️ 免责声明

本仓库聚合第三方志愿者**公开分享**的代理配置。我们不运营任何服务器,不保证可用性或安全性,不为使用行为负责。仅供学习和个人连接使用。请遵守所在司法管辖区的法律。

## ⭐ Star 历史

[![Star History Chart](https://api.star-history.com/svg?repos=Au1rxx/free-vpn-subscriptions&type=Date)](https://www.star-history.com/#Au1rxx/free-vpn-subscriptions&Date)

---

如果这个项目帮到你,点个 ⭐ —— 每一颗 star 都能帮更多人发现它。
