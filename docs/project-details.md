# Project Details

## Purpose

`Au1rxx/free-vpn-subscriptions` is the public distribution layer for the VPN feed. It is meant to solve four user problems in one place:

- provide directly usable public subscription URLs
- expose current public node status
- publish update history and snapshot releases
- reduce import failure with guides, FAQs, and troubleshooting pages

This repository is public by design. It should contain only public artifacts, public docs, public Pages content, and workflow code that publishes those public assets.

## Direct public entry points

- Clash remote subscription: `https://raw.githubusercontent.com/Au1rxx/free-vpn-subscriptions/main/output/clash.yaml`
- sing-box remote subscription: `https://raw.githubusercontent.com/Au1rxx/free-vpn-subscriptions/main/output/singbox.json`
- V2Ray remote subscription: `https://raw.githubusercontent.com/Au1rxx/free-vpn-subscriptions/main/output/v2ray-base64.txt`
- Live status JSON: `https://raw.githubusercontent.com/Au1rxx/free-vpn-subscriptions/main/output/status.json`
- Latest release page: `https://github.com/Au1rxx/free-vpn-subscriptions/releases/latest`
- GitHub Pages homepage: `https://au1rxx.github.io/free-vpn-subscriptions/`

## Repository layout

- `output/`
  Public feed artifacts that clients import directly.
- `site/`
  GitHub Pages content, search landing pages, status page, update page, and verification page.
- `docs/`
  Client guides plus technical documentation for maintainers.
- `scripts/`
  Small automation helpers used to build release notes, render update feeds, and validate public assets.
- `.github/workflows/`
  Public automation for Pages, release snapshots, release amplification, and validation.

## Public artifacts

### `output/clash.yaml`

- Format: YAML
- Primary use: Clash, Mihomo, and compatible remote profile import flows
- Distribution role: long-lived remote subscription URL

### `output/singbox.json`

- Format: JSON
- Primary use: sing-box remote profile imports
- Distribution role: long-lived remote subscription URL

### `output/v2ray-base64.txt`

- Format: Base64 subscription payload
- Primary use: V2Ray-compatible clients such as v2rayNG, NekoBox, Hiddify Next, and many Shadowrocket-style flows
- Distribution role: long-lived remote subscription URL

### `output/status.json`

- Format: JSON array
- Purpose: sanitized public health output for current nodes
- Expected fields: node name, region, protocol, public IP, port, status, latency, and last check timestamp
- Distribution role: source for the public status page and machine validation

## Public site structure

### Core pages

- `site/index.html`
  Main public entry point. Combines direct links, client paths, status, updates, and star/watch calls to action.
- `site/free-vpn-subscription-links.html`
  Direct-link hub for users who already know they need a raw URL or latest snapshot asset.
- `site/status.html`
  Human-readable status page backed by `output/status.json`.
- `site/updates.html`
  Release history, latest snapshots, and update-following paths.
- `site/verification.html`
  Documents what is machine-validated and how raw paths differ from release snapshots.

### Search and troubleshooting pages

The rest of `site/` contains search-intent landing pages such as:

- format choice pages
- client-specific setup pages
- troubleshooting pages for common failures
- platform entry pages for Android, iPhone, Windows, and macOS

Those pages exist to reduce setup friction and to capture recurring search queries without forcing users to leave the project.

## Generated outputs

Some files in the public repo are generated or refreshed by automation:

- `site/updates.json`
  JSON Feed rendered from recent GitHub Releases
- `site/updates.xml`
  Atom feed rendered from recent GitHub Releases
- release notes
  Generated from the current and previous `status.json` snapshots via `scripts/build_release_notes.py`

## Public/private boundary

This repository is intentionally the public side of the project.

Safe to keep here:

- public subscription artifacts
- sanitized node status output
- public Pages content
- setup guides and troubleshooting docs
- workflow code that publishes public assets

Do not add here:

- cloud credentials or access tokens
- SSH keys
- Terraform state
- private node inventory
- unsanitized control-plane state
- anything copied from the private operations repository that is not intended for public distribution

## Where to look when something is wrong

- Raw feed looks empty or stale:
  check `output/` and the upstream public sync step documented in [deployment.md](./deployment.md)
- Release assets are stale:
  check `Release Feed Snapshot`
- `updates.json` or `updates.xml` is stale:
  check `Amplify Release Updates`
- Pages content is stale:
  check `Deploy Pages`
- Pages and release are fresh but a client still fails:
  check the client guide, `status.json`, and `verification.html`
