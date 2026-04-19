# Free VPN Subscriptions

![Nodes](https://img.shields.io/badge/nodes-150-brightgreen) ![Alive](https://img.shields.io/badge/alive-622-blue) ![Median RTT](https://img.shields.io/badge/median--rtt-10ms-orange) ![Updated](https://img.shields.io/badge/updated-2026-04-19_04:49_UTC-informational)

> **The easiest way to get a working free VPN — copy a subscription link, paste it into your client, connect.**  
> No signup. No payment. No installation of binaries. Refreshed hourly from public sources with every node tested.

> 免费 VPN 订阅 · 免费梯子 · 免费科学上网 · free proxy · v2ray/clash/sing-box · VLESS / Reality / VMess / Trojan / Shadowsocks / Hysteria2

## 💡 Why This Project?

Every "free VPN" list on GitHub is either stale, full of dead nodes, or asks you to install a sketchy binary. This repo **only publishes nodes that passed a live TCP health check minutes ago**, from curated public sources, sorted by latency. You get 3 portable subscription files — drop them into Clash, sing-box, or v2rayN and go.

## 🚀 One-Click Subscribe

Copy the URL that matches your client and paste it into the subscription import field:

| Client | Format | Subscribe URL |
|---|---|---|
| Clash / Clash Verge / ClashX | `clash.yaml` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/clash.yaml` |
| sing-box | `singbox.json` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/singbox.json` |
| v2rayN / v2rayNG / Shadowrocket / NekoBox | `v2ray-base64` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/v2ray-base64.txt` |

## 🧩 Supported Clients

- **Windows**: v2rayN, Clash Verge, Hiddify, NekoRay
- **macOS**: ClashX Pro, Clash Verge, sing-box, Hiddify
- **iOS**: Shadowrocket, Stash, Loon, sing-box, Hiddify
- **Android**: v2rayNG, NekoBox, Clash Meta for Android, Hiddify, sing-box
- **Linux**: mihomo (Clash.Meta), sing-box, v2ray-core

## 📊 Live Stats

- **Nodes selected**: 150
- **Alive across all sources**: 622
- **Fastest node RTT**: 3 ms
- **Median RTT**: 10 ms
- **Last updated (UTC)**: 2026-04-19 04:49 UTC

**Protocol mix:** shadowsocks × 18 · trojan × 1 · vmess × 131

**Sources used this run:** `freefq` × 2 · `mahdibland-aggregator` × 77 · `mahdibland-shadowsocks` × 55 · `pawdroid` × 2 · `vxiaov-clash` × 14

## ❓ FAQ

<details><summary>Is this actually free?</summary>

Yes. Nodes are operated by third-party volunteers who publish their own free subscriptions. We don't run any servers ourselves — we just test, rank, and repackage what's already public.

</details>

<details><summary>How fresh is the data?</summary>

A GitHub Action runs every hour: pulls all upstream sources, TCP-probes every node, drops anything dead, sorts by latency, and commits new output files. Check the `Last updated` timestamp above.

</details>

<details><summary>Can I trust these nodes?</summary>

Free nodes see all your traffic. **Never use them for banking, login, or anything sensitive.** Fine for bypassing geo-blocks on public content. Use your own VPS / paid provider for real privacy.

</details>

<details><summary>Why do some nodes fail even though they're listed?</summary>

We only do TCP reachability checks. A node that handshakes may still have an expired cert, full bandwidth quota, or GFW-poisoned routes. Try a few; that's why the selector group gives you fallbacks.

</details>

## 🤝 Contributing

Know a reliable public subscription source we should add? Open an issue with the URL and format.

## ⚠️ Disclaimer

This repository aggregates **publicly shared** proxy configurations from third-party volunteers. We do not operate any servers, do not warrant availability or security, and are not responsible for how you use them. Intended for educational and personal connectivity use. Comply with all applicable laws in your jurisdiction.

## ⭐ Star History

[![Star History Chart](https://api.star-history.com/svg?repos=Au1rxx/free-vpn-subscriptions&type=Date)](https://www.star-history.com/#Au1rxx/free-vpn-subscriptions&Date)

---

If this project helped you, give it a ⭐ — every star makes it easier for others to find.
