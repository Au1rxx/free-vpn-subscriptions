# Free VPN Subscriptions

**English** · [简体中文](./README_CN.md) · [日本語](./README_JA.md) · [한국어](./README_KO.md) · [Español](./README_ES.md) · [Português](./README_PT.md) · [Русский](./README_RU.md)

<p align="center"><img src="https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/assets/hero.png" alt="Free VPN Subscriptions — hourly-refreshed free VPN subscriptions for Clash, sing-box, v2ray" width="780"></p>

![nodes](https://img.shields.io/badge/nodes-150-brightgreen) ![alive](https://img.shields.io/badge/alive-2579-blue) ![median--rtt](https://img.shields.io/badge/median--rtt-8ms-orange) ![updated](https://img.shields.io/badge/updated-2026-04-19_11:05_UTC-informational)

> **The easiest way to get a working free VPN — copy a subscription link, paste it into your client, connect.**  
> No signup. No payment. No installation of binaries. Refreshed hourly from public sources — every node is TCP + TLS probed before publishing.

> Free VPN subscriptions · free proxy list · free v2ray / clash / sing-box · VLESS / Reality / VMess / Trojan / Shadowsocks / Hysteria2 · hourly refreshed · TCP + TLS probed · by country

## 💡 Why This Project?

Every "free VPN" list on GitHub is either stale, full of dead nodes, or asks you to install a sketchy binary. This repo **only publishes nodes that passed a live TCP handshake AND a TLS handshake minutes ago**, from curated public sources, sorted by latency. You get 3 portable subscription files — drop them into Clash, sing-box, or v2rayN and go.

> 📖 How the fetch → probe → rank pipeline works: [ARCHITECTURE.md](./ARCHITECTURE.md)

## 🔬 How we verify nodes actually work

**Honest answer first: we cannot *guarantee* a node will pass your traffic.** No aggregator can, without running real traffic through it. Here is exactly what we verify, what we cannot, and where the real guarantee comes from.

### ✅ What we verify at aggregation time (before publishing)

1. **TCP reachability** — we open a TCP connection to every `server:port`. Dead hosts, bad DNS, and blocked ports get dropped. Drops roughly 40 % of raw entries.
2. **TLS handshake** — for every TLS / Reality / WS-TLS node we complete the full handshake. Expired certs, SNI mismatches, and broken Reality short-ids get dropped. Drops another ~10 %.
3. **Latency sort** — survivors are ranked by RTT and the top N are kept.

Typical numbers from a recent run: **17 sources → ~4,800 raw → ~2,900 TCP-alive → ~2,600 TLS-OK → top 200 published**.

### ❌ What we cannot verify

- Proxy protocol auth. A wrong UUID / password is only rejected *after* the TLS handshake by the upstream server.
- Actual HTTP-through-proxy success.
- Bandwidth or throughput.
- Geolocation beyond what GeoIP tells us about the exit IP.

### 🛡️ Runtime verification — the real guarantee

The `clash.yaml` we publish ships with a `url-test` proxy group that probes **real HTTP through each node** every 5 minutes:

```yaml
proxy-groups:
  - name: AUTO
    type: url-test
    url: http://www.gstatic.com/generate_204
    interval: 300
```

Your client ranks the node list by *actual* HTTP-through-proxy latency and auto-picks the fastest working node. sing-box and v2ray have equivalent mechanisms. If a selected node dies, the client switches to the next without intervention.

### 🧮 Expected outcome

Of the top 200 published each run, a typical client will find 30-50 that serve HTTP cleanly at any given moment. Rotate if one gets slow — the URL-test group makes that one click.

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
| 🇺🇸 United States (`US`) | 21 | [clash-US.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-US.yaml) | [singbox-US.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-US.json) | [v2ray-base64-US.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-US.txt) |
| 🇩🇪 Germany (`DE`) | 7 | [clash-DE.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-DE.yaml) | [singbox-DE.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-DE.json) | [v2ray-base64-DE.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-DE.txt) |

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
- **Alive across all sources**: 2579
- **Fastest node RTT**: 2 ms
- **Median RTT**: 8 ms
- **Last updated (UTC)**: 2026-04-19 11:05 UTC

**Protocol mix:** shadowsocks × 25 · trojan × 18 · vless × 85 · vmess × 22

**Sources used this run:** `barry-far-v2ray` × 30 · `ebrasha-v2ray` × 9 · `epodonios` × 33 · `freefq` × 1 · `mahdi0024` × 1 · `mahdibland-aggregator` × 1 · `mahdibland-shadowsocks` × 1 · `mfuu-clash` × 2 · `ninjastrikers` × 35 · `pawdroid` × 1 · `ruking-clash` × 21 · `snakem982` × 12 · `surfboard-eternity` × 2 · `vxiaov-clash` × 1

## ❓ FAQ

<details><summary>Is this actually free?</summary>

Yes. Nodes are operated by third-party volunteers who publish their own free subscriptions. We don't run any servers ourselves — we just test, rank, and repackage what's already public.

</details>

<details><summary>How fresh is the data?</summary>

Every hour (with a small random delay to avoid hammering upstream on the `:00` mark): pulls all upstream sources, TCP + TLS probes every node, drops anything dead, sorts by latency, publishes new output files. See the `Last updated` badge above.

</details>

<details><summary>Can I trust these nodes?</summary>

Free nodes see all your traffic. **Never use them for banking, login, or anything sensitive.** Fine for bypassing geo-blocks on public content. Use your own VPS / paid provider for real privacy.

</details>

<details><summary>Why do some nodes fail even though they're listed?</summary>

We verify TCP reachability and TLS handshakes only; a node can still have an expired quota, wrong routing, or GFW poisoning. The published `clash.yaml` pairs every node with a `url-test` proxy group (`http://www.gstatic.com/generate_204`, 300 s interval) — your client auto-picks the fastest node that actually serves HTTP. If one dies, pick the next.

</details>

## 🤝 Contributing

Know a reliable public subscription source we should add? Open an issue with the URL and format.

## ⚠️ Disclaimer

This repository aggregates **publicly shared** proxy configurations from third-party volunteers. We do not operate any servers, do not warrant availability or security, and are not responsible for how you use them. Intended for educational and personal connectivity use. Comply with all applicable laws in your jurisdiction.

## ⭐ Star History

[![Star History Chart](https://api.star-history.com/svg?repos=Au1rxx/free-vpn-subscriptions&type=Date)](https://www.star-history.com/#Au1rxx/free-vpn-subscriptions&Date)

---

If this project helped you, give it a ⭐ — every star makes it easier for others to find.
