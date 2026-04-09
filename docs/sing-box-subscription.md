# sing-box Subscription Guide

## Direct sing-box Subscription URL

Use this URL in sing-box-compatible clients:

`https://raw.githubusercontent.com/Au1rxx/free-vpn-subscriptions/main/output/singbox.json`

Backup release download:

`https://github.com/Au1rxx/free-vpn-subscriptions/releases/latest/download/singbox.json`

## Best Fit

- sing-box desktop
- sing-box mobile clients
- clients that accept remote JSON-based profile imports

## Import Steps

1. Open your sing-box client.
2. Add a remote profile or import from URL.
3. Paste the subscription URL above.
4. Sync the profile.
5. Refresh the profile manually once to confirm the remote source was saved.
6. Pick an available outbound and connect.

## Update Behavior

- The remote JSON feed is refreshed every 6 hours.
- Public node health is published hourly on the status page.
- Release snapshots act as stable downloadable checkpoints.

## Troubleshooting

- If your client requires a local file instead of a URL, download the JSON and import it manually.
- If an outbound appears stale, refresh the remote profile.
- If the profile imports but no routes work, compare with the published status feed.
- If you are deciding between sing-box and V2Ray, prefer the format your client documents as its native remote profile path.

## Related Guides

- [V2Ray guide](./v2ray-subscription.md)
- [Hiddify Next guide](./hiddify-next.md)
- [NekoBox guide](./nekobox.md)
- [FAQ](./faq.md)
