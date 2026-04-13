# Continuously Updated Free VPN Subscriptions

[中文](./README.md) | [English](./README.en.md)

[![Public Repo](https://img.shields.io/badge/repo-public-0f766e)](https://github.com/Au1rxx/free-vpn-subscriptions)
[![Updated](https://img.shields.io/badge/updates-every%206%20hours-0f766e)](https://au1rxx.github.io/free-vpn-subscriptions/updates.html)
[![Verified](https://img.shields.io/badge/validation-auto%20checked-1d221c)](https://au1rxx.github.io/free-vpn-subscriptions/verification.html)
[![Formats](https://img.shields.io/badge/formats-clash%20%7C%20sing--box%20%7C%20v2ray-cf6a32)](https://au1rxx.github.io/free-vpn-subscriptions/)
[![Status Feed](https://img.shields.io/badge/status-live-1d221c)](https://au1rxx.github.io/free-vpn-subscriptions/status.html)
[![Latest Release](https://img.shields.io/github/v/release/Au1rxx/free-vpn-subscriptions)](https://github.com/Au1rxx/free-vpn-subscriptions/releases/latest)

Public distribution entry for Clash, sing-box, and V2Ray formats. Health is checked hourly, public snapshots are refreshed on schedule, and the repo includes setup guides, verification notes, and troubleshooting pages.

[Open Live Status Site](https://au1rxx.github.io/free-vpn-subscriptions/) • [Updates And Snapshot History](https://au1rxx.github.io/free-vpn-subscriptions/updates.html) • [Status Dashboard](https://au1rxx.github.io/free-vpn-subscriptions/status.html) • [Verification And Compatibility](https://au1rxx.github.io/free-vpn-subscriptions/verification.html) • [FAQ](./docs/faq.md)

[Atom Feed](https://au1rxx.github.io/free-vpn-subscriptions/updates.xml) • [JSON Feed](https://au1rxx.github.io/free-vpn-subscriptions/updates.json) • [Discussions](https://github.com/Au1rxx/free-vpn-subscriptions/discussions)

## Why this repo is worth a star

- This is not just a single subscription URL. It combines `live status`, `update history`, `verification notes`, `client guides`, and `troubleshooting` into one public entry point.
- Returning traffic should come from “did it update today?”, “is the public status healthy?”, and “is there a new guide?”, not from rotating feed URLs.
- Releases, update feeds, Discussions, and public Pages accumulate visible update signals over time, which is more durable than bookmarking a single download path.

## What returning visitors should check

- `Is the public status fresh?` Open the [Status Dashboard](https://au1rxx.github.io/free-vpn-subscriptions/status.html).
- `Was there a new snapshot today?` Open [Updates And Snapshot History](https://au1rxx.github.io/free-vpn-subscriptions/updates.html).
- `Are the shared paths still valid?` Open the [Verification And Compatibility](https://au1rxx.github.io/free-vpn-subscriptions/verification.html) page.
- `Want ongoing updates?` Star the repo or subscribe to Atom / JSON feed.

## Check These Signals First

- If you want to confirm that the public side is still fresh, check the [Status Dashboard](https://au1rxx.github.io/free-vpn-subscriptions/status.html) and the latest timestamp in `output/status.json`.
- If you want the newest snapshot, clean download links, and release history, go straight to [GitHub Releases](https://github.com/Au1rxx/free-vpn-subscriptions/releases).
- If you want to know which public links are currently machine-validated, which paths are snapshot downloads, and which links are meant for long-lived auto-refresh, open the [Verification And Compatibility](https://au1rxx.github.io/free-vpn-subscriptions/verification.html) page.
- If you want ongoing updates, do not rely on bookmarking the repo homepage alone. Star the repo, use the repository Watch menu for `Releases`, or subscribe to the [Atom Feed](https://au1rxx.github.io/free-vpn-subscriptions/updates.xml) / [JSON Feed](https://au1rxx.github.io/free-vpn-subscriptions/updates.json).

## Why This Repo

- Public subscription feed for `Clash`, `sing-box`, and `V2Ray` users.
- Live health status published from the private control plane.
- Multi-region endpoints with hourly health checks and scheduled subscription refreshes.
- Public documentation designed for fast onboarding and high GitHub discoverability.

## How to get started

- If you want the fastest entry, start from the [homepage](https://au1rxx.github.io/free-vpn-subscriptions/).
- If you need help choosing a format, open [Which subscription format should I use?](https://au1rxx.github.io/free-vpn-subscriptions/which-subscription-format-should-i-use.html).
- If you want a downloadable snapshot, open [GitHub Releases](https://github.com/Au1rxx/free-vpn-subscriptions/releases).
- If you want the difference between raw auto-refresh paths and release snapshots, open the [Verification And Compatibility](https://au1rxx.github.io/free-vpn-subscriptions/verification.html) page.

## Update cadence

- Fleet health checks: every hour
- Public subscription snapshots: refreshed on schedule
- Release snapshots: follow public artifact updates for manual download and history
- Update feeds: automatically refreshed after new releases

## Why The Main Subscription URLs Stay Stable

- Stable URLs are better for client auto-refresh and prevent old guides, bookmarks, and shared posts from breaking.
- Returning traffic should come from updates, release snapshots, live status, and new documentation, not from rotating the core feed URLs.
- Use the [Updates And Snapshot History](https://au1rxx.github.io/free-vpn-subscriptions/updates.html) page when you want to check recent release activity or download checkpoints manually.
- Use the [Atom Feed](https://au1rxx.github.io/free-vpn-subscriptions/updates.xml) or [JSON Feed](https://au1rxx.github.io/free-vpn-subscriptions/updates.json) if you want to pipe release updates into RSS readers, automation, or external dashboards.

## Supported Clients

- Clash Verge / Clash Meta / Mihomo-compatible clients
- sing-box
- V2Ray-compatible clients including v2rayNG and NekoBox
- iOS clients that can import standard subscription URLs

## Quick Start

1. Pick the subscription format that matches your client.
2. Start from the homepage or the client-specific guide instead of pasting a raw link blindly.
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
- [How to use FlClash](./docs/flclash.md)
- [How to use Clash Meta for Android](./docs/clash-meta-android.md)
- [How to use Hiddify Next](./docs/hiddify-next.md)
- [How to use NekoBox](./docs/nekobox.md)
- [How to use v2rayNG](./docs/v2rayng.md)
- [How to use Shadowrocket](./docs/shadowrocket.md)
- [FAQ and troubleshooting](./docs/faq.md)

## Search Intent Pages

- [Clash subscription not working](https://au1rxx.github.io/free-vpn-subscriptions/clash-subscription-not-working.html)
- [V2Ray subscription URL](https://au1rxx.github.io/free-vpn-subscriptions/v2ray-subscription-url.html)
- [Shadowrocket subscription URL](https://au1rxx.github.io/free-vpn-subscriptions/shadowrocket-subscription-url.html)
- [Free VPN for Android](https://au1rxx.github.io/free-vpn-subscriptions/free-vpn-for-android.html)
- [How to refresh a Clash profile](https://au1rxx.github.io/free-vpn-subscriptions/how-to-refresh-clash-profile.html)
- [V2Ray subscription not working](https://au1rxx.github.io/free-vpn-subscriptions/v2ray-subscription-not-working.html)
- [Shadowrocket not connecting](https://au1rxx.github.io/free-vpn-subscriptions/shadowrocket-not-connecting.html)
- [Clash profile update failed](https://au1rxx.github.io/free-vpn-subscriptions/clash-profile-update-failed.html)
- [Free VPN for iPhone](https://au1rxx.github.io/free-vpn-subscriptions/free-vpn-for-iphone.html)
- [Best Clash client for Android](https://au1rxx.github.io/free-vpn-subscriptions/best-clash-client-for-android.html)
- [Troubleshooting hub](https://au1rxx.github.io/free-vpn-subscriptions/troubleshooting-hub.html)
- [Free VPN subscription links](https://au1rxx.github.io/free-vpn-subscriptions/free-vpn-subscription-links.html)
- [Which subscription format should I use?](https://au1rxx.github.io/free-vpn-subscriptions/which-subscription-format-should-i-use.html)
- [Clash vs V2Ray subscription](https://au1rxx.github.io/free-vpn-subscriptions/clash-vs-v2ray-subscription.html)
- [Best VPN client for Windows](https://au1rxx.github.io/free-vpn-subscriptions/best-vpn-client-for-windows.html)
- [Best VPN client for Mac](https://au1rxx.github.io/free-vpn-subscriptions/best-vpn-client-for-mac.html)

## Community

- [Open a setup question](https://github.com/Au1rxx/free-vpn-subscriptions/discussions/new?category=q-a)
- [Suggest a new client guide or landing page](https://github.com/Au1rxx/free-vpn-subscriptions/discussions/new?category=ideas)
- [Share a working client setup](https://github.com/Au1rxx/free-vpn-subscriptions/discussions/new?category=show-and-tell)
- [Report a broken public feed or stale status](https://github.com/Au1rxx/free-vpn-subscriptions/issues/new/choose)

## Popular Client Paths

- Desktop: Clash Verge Rev, FlClash, sing-box desktop
- Android: Clash Meta for Android, v2rayNG, NekoBox, Hiddify Next, sing-box mobile
- iPhone and iPad: Shadowrocket-compatible subscription import flows
- Multi-platform: Hiddify Next on Android, iPhone, Windows, macOS, and Linux

These client-specific guides exist to improve search visibility for real user queries and reduce setup friction after discovery.

## Star And Follow

If this feed is useful:

- Star the repository so you keep the full public entry point, not only a single file path.
- Use the repository Watch menu for `Releases` if you want visible update checkpoints, or subscribe to the Atom / JSON feed above.
- Use Discussions for client requests, setup questions, and broken-link reports.
- Use Issues only when a public feed artifact, release asset, or Pages route is actually broken.

## Notes

- This repository distributes public subscription material. Treat these nodes as shared public resources.
- Node availability and performance can change over time.
- The live status view is the best source of truth for current reachability.
