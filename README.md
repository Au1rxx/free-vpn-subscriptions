# Free VPN Subscriptions

**English** · [简体中文](./README_CN.md) · [日本語](./README_JA.md) · [한국어](./README_KO.md) · [Español](./README_ES.md) · [Português](./README_PT.md) · [Русский](./README_RU.md)

<p align="center"><img src="https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/assets/hero.png" alt="Free VPN Subscriptions — hourly-refreshed free VPN subscriptions for Clash, sing-box, v2ray" width="780"></p>

![nodes](https://img.shields.io/badge/nodes-150-brightgreen) ![alive](https://img.shields.io/badge/alive-2446-blue) ![median--rtt](https://img.shields.io/badge/median--rtt-98ms-orange) ![updated](https://img.shields.io/badge/updated-2026-04-24_00:29_UTC-informational)

> **The easiest way to get a working free VPN — copy a subscription link, paste it into your client, connect.**  
> No signup. No payment. No installation of binaries. Refreshed hourly from public sources — every published node has demonstrably forwarded real HTTP traffic through sing-box minutes ago.

> Free VPN subscriptions · free proxy list · free v2ray / clash / sing-box · VLESS / Reality / VMess / Trojan / Shadowsocks / Hysteria2 · hourly refreshed · HTTP-over-proxy verified · by country

## 💡 Why This Project?

Every "free VPN" list on GitHub is either stale, full of dead nodes, or asks you to install a sketchy binary. This repo goes further than anything else you'll find — **we don't just check that a node answers the phone, we actually push HTTP traffic through it with sing-box and confirm a 204 comes back**, minutes before publishing. You get 3 portable subscription files — drop them into Clash, sing-box, or v2rayN and go.

> 📖 How the fetch → probe → rank pipeline works: [ARCHITECTURE.md](./ARCHITECTURE.md)

## 🔬 How we verify nodes actually work

Most free-VPN lists stop at "the TCP port is open" and publish. We don't. Here is the full verification pipeline a node has to survive before it gets into the subscription.

### ✅ What we verify at aggregation time (before publishing)

1. **TCP reachability** — open a TCP connection to every `server:port`. Dead hosts, bad DNS, and blocked ports are dropped. ~40 % of raw entries fall out here.
2. **TLS handshake** — for every TLS / Reality / WS-TLS node we complete the full handshake. Expired certs, SNI mismatches, and broken Reality short-ids are dropped. ~10 % more fall out here.
3. **sing-box config validation** — every surviving node is translated into a real sing-box outbound and run through `sing-box check`. Corrupt ciphers, bad UUIDs, and unsupported flow options are dropped before they waste a probe slot.
4. **HTTP-over-proxy probe (this is the big one)** — we batch the fastest ~900 candidates into sing-box subprocesses, each node getting its own local SOCKS5 inbound, then push real HTTP + HTTPS GETs through it:
   - `http://www.gstatic.com/generate_204` (expects 204)
   - `https://www.cloudflare.com/cdn-cgi/trace` (expects 200)

   The request traverses the actual proxy protocol (VLESS / VMess / Trojan / Shadowsocks / Hysteria2), so a node that passes has demonstrably functioning auth, routing, TLS inner handshake, and exit networking.
5. **Two rounds, 45 seconds apart** — nodes that pass once but die 45 seconds later get filtered. Only nodes with ≥ 50 % success rate across all (rounds × targets) are kept.
6. **Median real-latency sort** — survivors are ranked by their median HTTP-over-proxy round-trip (not raw TCP RTT), and the top N are published.

Typical numbers from a recent run: **17 sources → ~4,800 raw → ~2,900 TCP-alive → ~2,600 TLS-OK → ~840 config-valid → ~280 HTTP-verified → top 150 published**. Every one of the 150 has actually forwarded traffic in the last ten minutes.

### ❌ What we still can't verify

- **Bandwidth / throughput** — we measure latency, not megabits. A 50 ms node may still be slow for video.
- **Geolocation precision** — GeoIP tells us the exit IP country but not the city or ISP reliably.
- **Region-specific blocks** — a node that works from our probe infra may be blocked from yours (ISP-level filtering, captive portals, etc.).
- **Staying alive past the run** — the node passed ten minutes ago; it may have died since.

### 🛡️ Runtime safety net — for the last bullet above

The `clash.yaml` we publish ships with a `url-test` proxy group that re-tests real HTTP through each node every 5 minutes on *your* device:

```yaml
proxy-groups:
  - name: AUTO
    type: url-test
    url: http://www.gstatic.com/generate_204
    interval: 300
```

Your client keeps the node list sorted by *live* HTTP-over-proxy latency from your network and auto-picks the fastest working node. sing-box and v2ray have equivalent mechanisms. If a selected node dies between hourly aggregations, the client switches to the next without intervention.

### 🧮 What this means in practice

Of the ~150 we publish each run, a typical client finds **80-120 nodes that serve HTTP cleanly from their network** at any given moment — roughly 2-3× the hit rate of lists that only do TCP probing. The url-test group rotates transparently if one drops out.

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
| 🇺🇸 United States (`US`) | 35 | [clash-US.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-US.yaml) | [singbox-US.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-US.json) | [v2ray-base64-US.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-US.txt) |
| 🇩🇪 Germany (`DE`) | 8 | [clash-DE.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-DE.yaml) | [singbox-DE.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-DE.json) | [v2ray-base64-DE.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-DE.txt) |
| 🇨🇦 Canada (`CA`) | 7 | [clash-CA.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-CA.yaml) | [singbox-CA.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-CA.json) | [v2ray-base64-CA.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-CA.txt) |
| 🇬🇧 United Kingdom (`GB`) | 3 | [clash-GB.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-GB.yaml) | [singbox-GB.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-GB.json) | [v2ray-base64-GB.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-GB.txt) |

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
- **Alive across all sources**: 2446
- **Fastest node RTT**: 28 ms
- **Median RTT**: 98 ms
- **Last updated (UTC)**: 2026-04-24 00:29 UTC

**Protocol mix:** shadowsocks × 22 · trojan × 11 · vless × 105 · vmess × 12

**Sources used this run:** `barry-far-v2ray` × 26 · `epodonios` × 7 · `lagzian-mix` × 4 · `mahdi0024` × 25 · `mahdibland-aggregator` × 8 · `mahdibland-shadowsocks` × 1 · `ninjastrikers` × 71 · `pawdroid` × 2 · `surfboard-eternity` × 6

## ❓ FAQ

<details><summary>Is this actually free?</summary>

Yes. Nodes are operated by third-party volunteers who publish their own free subscriptions. We don't run any servers ourselves — we just test, rank, and repackage what's already public.

</details>

<details><summary>How fresh is the data?</summary>

Every hour (with a small random delay to avoid hammering upstream on the `:00` mark): pulls all sources, runs the full TCP → TLS → sing-box config check → HTTP-over-proxy probe pipeline (two rounds, 45 s apart), ranks by real HTTP latency, publishes new output files. Full pipeline takes ~10 minutes. See the `Last updated` badge above.

</details>

<details><summary>Can I trust these nodes?</summary>

Free nodes see all your traffic. **Never use them for banking, login, or anything sensitive.** Fine for bypassing geo-blocks on public content. Use your own VPS / paid provider for real privacy.

</details>

<details><summary>Why do some nodes fail even though they're listed?</summary>

Even after our HTTP-over-proxy probe, nodes can die between aggregations: quota exhausted, upstream revoked the key, your ISP blocks the exit IP, or the operator took it down. The published `clash.yaml` pairs every node with a `url-test` proxy group (`http://www.gstatic.com/generate_204`, 300 s interval) — your client auto-picks the fastest node that actually serves HTTP *from your network*. If one dies, pick the next. Expect 80-120 of the 150 to work at any given moment.

</details>

## 🤝 Contributing

Know a reliable public subscription source we should add? Open an issue with the URL and format.

## ⚠️ Disclaimer

This repository aggregates **publicly shared** proxy configurations from third-party volunteers. We do not operate any servers, do not warrant availability or security, and are not responsible for how you use them. Intended for educational and personal connectivity use. Comply with all applicable laws in your jurisdiction.

## ⭐ Star History

[![Star History Chart](https://api.star-history.com/svg?repos=Au1rxx/free-vpn-subscriptions&type=Date)](https://www.star-history.com/#Au1rxx/free-vpn-subscriptions&Date)

---

If this project helped you, give it a ⭐ — every star makes it easier for others to find.
