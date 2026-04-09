# Free VPN Subscriptions

Free Clash, sing-box, and V2Ray subscription links with live node status, multi-region coverage, and setup guides for common clients.

[Open Live Status Site](https://au1rxx.github.io/free-vpn-subscriptions/) • [Clash Guide](./docs/clash-subscription.md) • [sing-box Guide](./docs/sing-box-subscription.md) • [V2Ray Guide](./docs/v2ray-subscription.md) • [FAQ](./docs/faq.md)

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
- [FAQ and troubleshooting](./docs/faq.md)

## Star And Follow

If this feed is useful:

- Star the repository to help discovery.
- Watch releases if you want visible update checkpoints.
- Use Discussions for client requests, setup questions, and broken-link reports.

## Notes

- This repository distributes public subscription material. Treat these nodes as shared public resources.
- Node availability and performance can change over time.
- The live status view is the best source of truth for current reachability.
