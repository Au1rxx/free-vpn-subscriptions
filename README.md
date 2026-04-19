# Free VPN Subscriptions

**English** · [简体中文](./README_CN.md) · [日本語](./README_JA.md) · [한국어](./README_KO.md) · [Español](./README_ES.md) · [Português](./README_PT.md) · [Русский](./README_RU.md)

<p align="center"><img src="https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/assets/workflow.svg" alt="Free VPN Subscriptions — aggregation workflow" width="780"></p>

![nodes](https://img.shields.io/badge/nodes-150-brightgreen) ![alive](https://img.shields.io/badge/alive-393-blue) ![median--rtt](https://img.shields.io/badge/median--rtt-119ms-orange) ![updated](https://img.shields.io/badge/updated-2026-04-19_06:42_UTC-informational)

> **The easiest way to get a working free VPN — copy a subscription link, paste it into your client, connect.**  
> No signup. No payment. No installation of binaries. Refreshed hourly from public sources with every node tested.

> Free VPN subscriptions · free proxy list · free v2ray / clash / sing-box · VLESS / Reality / VMess / Trojan / Shadowsocks / Hysteria2 · hourly refreshed · TCP + TLS probed · by country

## 💡 Why This Project?

Every "free VPN" list on GitHub is either stale, full of dead nodes, or asks you to install a sketchy binary. This repo **only publishes nodes that passed a live TCP handshake AND a TLS handshake minutes ago**, from curated public sources, sorted by latency. You get 3 portable subscription files — drop them into Clash, sing-box, or v2rayN and go.

> 📖 How the fetch → probe → rank pipeline works: [ARCHITECTURE.md](./ARCHITECTURE.md)

## 🚀 One-Click Subscribe

Copy the URL that matches your client and paste it into the subscription import field:

| Client | Format | Subscribe URL |
|---|---|---|
| Clash / Clash Verge / ClashX | `clash.yaml` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/clash.yaml` |
| sing-box | `singbox.json` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/singbox.json` |
| v2rayN / v2rayNG / Shadowrocket / NekoBox | `v2ray-base64` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/v2ray-base64.txt` |

## 🌍 By Country

Want nodes in a specific region only? Use one of these targeted subscription URLs:

| Country | Nodes | Clash | sing-box | v2ray |
|---|---|---|---|---|
| 🇺🇸 United States (`US`) | 28 | [clash-US.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-US.yaml) | [singbox-US.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-US.json) | [v2ray-base64-US.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-US.txt) |
| 🇬🇧 United Kingdom (`GB`) | 26 | [clash-GB.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-GB.yaml) | [singbox-GB.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-GB.json) | [v2ray-base64-GB.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-GB.txt) |
| 🇳🇱 Netherlands (`NL`) | 14 | [clash-NL.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-NL.yaml) | [singbox-NL.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-NL.json) | [v2ray-base64-NL.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-NL.txt) |
| 🇨🇦 Canada (`CA`) | 11 | [clash-CA.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-CA.yaml) | [singbox-CA.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-CA.json) | [v2ray-base64-CA.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-CA.txt) |
| 🇯🇵 Japan (`JP`) | 8 | [clash-JP.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-JP.yaml) | [singbox-JP.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-JP.json) | [v2ray-base64-JP.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-JP.txt) |
| 🇭🇰 Hong Kong (`HK`) | 6 | [clash-HK.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-HK.yaml) | [singbox-HK.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-HK.json) | [v2ray-base64-HK.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-HK.txt) |
| 🇰🇷 Korea (`KR`) | 4 | [clash-KR.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-KR.yaml) | [singbox-KR.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-KR.json) | [v2ray-base64-KR.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-KR.txt) |
| 🇹🇼 Taiwan (`TW`) | 4 | [clash-TW.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-TW.yaml) | [singbox-TW.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-TW.json) | [v2ray-base64-TW.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-TW.txt) |

## 🧩 Supported Clients

- **Windows**: v2rayN, Clash Verge, Hiddify, NekoRay
- **macOS**: ClashX Pro, Clash Verge, sing-box, Hiddify
- **iOS**: Shadowrocket, Stash, Loon, sing-box, Hiddify
- **Android**: v2rayNG, NekoBox, Clash Meta for Android, Hiddify, sing-box
- **Linux**: mihomo (Clash.Meta), sing-box, v2ray-core

## 📊 Live Stats

- **Nodes selected**: 150
- **Alive across all sources**: 393
- **Fastest node RTT**: 8 ms
- **Median RTT**: 119 ms
- **Last updated (UTC)**: 2026-04-19 06:42 UTC

**Protocol mix:** shadowsocks × 104 · trojan × 7 · vmess × 39

**Sources used this run:** `freefq` × 2 · `mahdibland-aggregator` × 67 · `mahdibland-shadowsocks` × 65 · `vxiaov-clash` × 16

## ❓ FAQ

<details><summary>Is this actually free?</summary>

Yes. Nodes are operated by third-party volunteers who publish their own free subscriptions. We don't run any servers ourselves — we just test, rank, and repackage what's already public.

</details>

<details><summary>How fresh is the data?</summary>

A GitHub Action runs every hour: pulls all upstream sources, TCP + TLS probes every node, drops anything dead, sorts by latency, and commits new output files. Check the `Last updated` timestamp above.

</details>

<details><summary>Can I trust these nodes?</summary>

Free nodes see all your traffic. **Never use them for banking, login, or anything sensitive.** Fine for bypassing geo-blocks on public content. Use your own VPS / paid provider for real privacy.

</details>

<details><summary>Why do some nodes fail even though they're listed?</summary>

We verify TCP reachability and TLS handshakes, but a node can still have an expired quota, wrong routing, or GFW poisoning. Try a few; the selector group gives you fallbacks.

</details>

## 🤝 Contributing

Know a reliable public subscription source we should add? Open an issue with the URL and format.

## ⚠️ Disclaimer

This repository aggregates **publicly shared** proxy configurations from third-party volunteers. We do not operate any servers, do not warrant availability or security, and are not responsible for how you use them. Intended for educational and personal connectivity use. Comply with all applicable laws in your jurisdiction.

## ⭐ Star History

[![Star History Chart](https://api.star-history.com/svg?repos=Au1rxx/free-vpn-subscriptions&type=Date)](https://www.star-history.com/#Au1rxx/free-vpn-subscriptions&Date)

---

If this project helped you, give it a ⭐ — every star makes it easier for others to find.
