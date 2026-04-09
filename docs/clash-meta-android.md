# Clash Meta for Android Guide

## Best Subscription Format

Use the Clash YAML subscription feed:

`https://raw.githubusercontent.com/Au1rxx/free-vpn-subscriptions/main/output/clash.yaml`

Backup release download:

`https://github.com/Au1rxx/free-vpn-subscriptions/releases/latest/download/clash.yaml`

## Clash Meta for Android Setup

1. Open Clash Meta for Android and go to profiles.
2. Create a new profile from URL or add a remote subscription.
3. Paste the Clash YAML link above.
4. Import the profile and refresh it once manually.
5. Select a node and test connectivity.

## Why This Path Works Well

- It matches the profile style most Clash-based Android users expect.
- It keeps the import path consistent with desktop Clash clients.
- It supports remote refresh better than a one-time local import.

## Common Issues

- If the app imports a profile but nothing connects, switch nodes before rewriting DNS or routing settings.
- If the remote profile does not refresh, confirm the URL was saved as a subscription source.
- If all nodes fail together, compare with the public status page first.

