# Free VPN Subscriptions

**English** · [简体中文](./README_CN.md) · [日本語](./README_JA.md) · [한국어](./README_KO.md) · [Español](./README_ES.md) · [Português](./README_PT.md) · [Русский](./README_RU.md)

<p align="center"><img src="https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/assets/hero.png" alt="Free VPN Subscriptions — hourly-refreshed free VPN subscriptions for Clash, sing-box, v2ray" width="780"></p>

![nodes](https://img.shields.io/badge/nodes-150-brightgreen) ![alive](https://img.shields.io/badge/alive-2599-blue) ![median--rtt](https://img.shields.io/badge/median--rtt-1ms-orange) ![updated](https://img.shields.io/badge/updated-2026-04-19_10:33_UTC-informational)

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
| 🇺🇸 United States (`US`) | 20 | [clash-US.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-US.yaml) | [singbox-US.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-US.json) | [v2ray-base64-US.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-US.txt) |

## 📖 Step-by-step Guides

New to VPN clients? Pick your platform and follow the tutorial:

- [**Clash Verge**](https://au1rxx.github.io/free-vpn-subscriptions/guides/clash-verge.html) · Windows / macOS / Linux
- [**v2rayNG**](https://au1rxx.github.io/free-vpn-subscriptions/guides/v2rayng.html) · Android
- [**Shadowrocket**](https://au1rxx.github.io/free-vpn-subscriptions/guides/shadowrocket.html) · iOS / iPadOS
- [**sing-box**](https://au1rxx.github.io/free-vpn-subscriptions/guides/sing-box.html) · Windows / macOS / Linux / iOS / Android

## 🧩 Supported Clients

- **Windows**: v2rayN, Clash Verge, Hiddify, NekoRay
- **macOS**: ClashX Pro, Clash Verge, sing-box, Hiddify
- **iOS**: Shadowrocket, Stash, Loon, sing-box, Hiddify
- **Android**: v2rayNG, NekoBox, Clash Meta for Android, Hiddify, sing-box
- **Linux**: mihomo (Clash.Meta), sing-box, v2ray-core

## 📊 Live Stats

- **Nodes selected**: 150
- **Alive across all sources**: 2599
- **Fastest node RTT**: 1 ms
- **Median RTT**: 1 ms
- **Last updated (UTC)**: 2026-04-19 10:33 UTC

**Protocol mix:** shadowsocks × 22 · trojan × 19 · vless × 93 · vmess × 16

**Sources used this run:** `barry-far-v2ray` × 40 · `ebrasha-v2ray` × 15 · `epodonios` × 29 · `mfuu-clash` × 3 · `ninjastrikers` × 44 · `ruking-clash` × 13 · `snakem982` × 4 · `surfboard-eternity` × 2

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
