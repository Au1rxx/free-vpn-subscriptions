# Free VPN Subscriptions

[![Public Repo](https://img.shields.io/badge/repo-public-0f766e)](https://github.com/Au1rxx/free-vpn-subscriptions)
[![Formats](https://img.shields.io/badge/formats-clash%20%7C%20sing--box%20%7C%20v2ray-cf6a32)](https://github.com/Au1rxx/free-vpn-subscriptions/tree/main/output)
[![Status Feed](https://img.shields.io/badge/status-live-1d221c)](https://raw.githubusercontent.com/Au1rxx/free-vpn-subscriptions/main/output/status.json)
[![Latest Release](https://img.shields.io/github/v/release/Au1rxx/free-vpn-subscriptions)](https://github.com/Au1rxx/free-vpn-subscriptions/releases/latest)

Free Clash, sing-box, and V2Ray subscription links with live node status, multi-region coverage, and setup guides for common clients.

[Open Live Status Site](https://au1rxx.github.io/free-vpn-subscriptions/) • [Status Dashboard](https://au1rxx.github.io/free-vpn-subscriptions/status.html) • [Clash Guide](./docs/clash-subscription.md) • [sing-box Guide](./docs/sing-box-subscription.md) • [V2Ray Guide](./docs/v2ray-subscription.md) • [FAQ](./docs/faq.md)

## Why This Repo

- Public subscription feed for `Clash`, `sing-box`, and `V2Ray` users.
- Live health status published from the private control plane.
- Multi-region endpoints with hourly health checks and scheduled subscription refreshes.
- Public documentation designed for fast onboarding and high GitHub discoverability.

## Subscription Links

Use the raw links below in your client:

| Format | Direct Link | Update Cadence |
|--------|-------------|----------------|
| Clash | `https://raw.githubusercontent.com/Au1rxx/free-vpn-subscriptions/main/output/clash.yaml` | Every 6 hours |
| sing-box | `https://raw.githubusercontent.com/Au1rxx/free-vpn-subscriptions/main/output/singbox.json` | Every 6 hours |
| V2Ray | `https://raw.githubusercontent.com/Au1rxx/free-vpn-subscriptions/main/output/v2ray-base64.txt` | Every 6 hours |
| Status | `https://raw.githubusercontent.com/Au1rxx/free-vpn-subscriptions/main/output/status.json` | Every hour |

## Release Downloads

Use Releases as a backup distribution channel, a snapshot archive, or a clean download entry for users who do not want to browse raw files.

| Format | Latest Release Asset |
|--------|----------------------|
| Clash | `https://github.com/Au1rxx/free-vpn-subscriptions/releases/latest/download/clash.yaml` |
| sing-box | `https://github.com/Au1rxx/free-vpn-subscriptions/releases/latest/download/singbox.json` |
| V2Ray | `https://github.com/Au1rxx/free-vpn-subscriptions/releases/latest/download/v2ray-base64.txt` |
| Status | `https://github.com/Au1rxx/free-vpn-subscriptions/releases/latest/download/status.json` |

## Supported Clients

- Clash Verge / Clash Meta / Mihomo-compatible clients
- sing-box
- V2Ray-compatible clients including v2rayNG and NekoBox
- iOS clients that can import standard subscription URLs

## Quick Start

1. Pick the subscription format that matches your client.
2. Import the direct link above into your VPN client.
3. Refresh the profile when the node set changes.
4. Check the live status site before troubleshooting client-side issues.

## What Gets Published

This public repository only contains public subscription artifacts and sanitized status output:

- `output/clash.yaml`
- `output/singbox.json`
- `output/v2ray-base64.txt`
- `output/status.json`

The private control plane, infrastructure state, deployment credentials, and cloud access remain in a separate private repository.

## Guides

- [How to import Clash subscriptions](./docs/clash-subscription.md)
- [How to import sing-box subscriptions](./docs/sing-box-subscription.md)
- [How to import V2Ray subscriptions](./docs/v2ray-subscription.md)
- [How to use Clash Verge Rev](./docs/clash-verge-rev.md)
- [How to use Hiddify Next](./docs/hiddify-next.md)
- [How to use NekoBox](./docs/nekobox.md)
- [How to use v2rayNG](./docs/v2rayng.md)
- [How to use Shadowrocket](./docs/shadowrocket.md)
- [FAQ and troubleshooting](./docs/faq.md)

## Community

- [Open a setup question](https://github.com/Au1rxx/free-vpn-subscriptions/discussions/new?category=q-a)
- [Suggest a new client guide or landing page](https://github.com/Au1rxx/free-vpn-subscriptions/discussions/new?category=ideas)
- [Share a working client setup](https://github.com/Au1rxx/free-vpn-subscriptions/discussions/new?category=show-and-tell)
- [Report a broken public feed or stale status](https://github.com/Au1rxx/free-vpn-subscriptions/issues/new/choose)

## Popular Client Paths

- Desktop: Clash Verge Rev, FlClash, sing-box desktop
- Android: v2rayNG, NekoBox, Hiddify Next, sing-box mobile
- iPhone and iPad: Shadowrocket-compatible subscription import flows
- Multi-platform: Hiddify Next on Android, iPhone, Windows, macOS, and Linux

These client-specific guides exist to improve search visibility for real user queries and reduce setup friction after discovery.

## Star And Follow

If this feed is useful:

- Star the repository to help discovery.
- Watch releases if you want visible update checkpoints.
- Use Discussions for client requests, setup questions, and broken-link reports.
- Use Issues only when a public feed artifact, release asset, or Pages route is actually broken.

## Notes

- This repository distributes public subscription material. Treat these nodes as shared public resources.
- Node availability and performance can change over time.
- The live status view is the best source of truth for current reachability.
