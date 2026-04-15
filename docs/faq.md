# FAQ

## Are these subscriptions public?

Yes. This repository is designed to publish public subscription artifacts and public node status output.

## How often is the feed updated?

Health status is refreshed hourly. Subscription files are refreshed every 6 hours.

## How is the public repo deployed?

The public repo is fed by an upstream private control repository that syncs only the public `output/` directory into this repo. A push here then triggers GitHub Pages, release snapshots, update-feed rendering, and public validation.

See also:

- [Project details](./project-details.md)
- [Deployment and operations](./deployment.md)

## Why does a node disappear or stop working?

Public shared nodes can degrade, rotate, or go offline. Always check `output/status.json` or the live status page for the latest reachability signal.

## Which format should I use?

- Use `Clash` if your client expects a YAML subscription.
- Use `sing-box` if your client supports remote JSON imports.
- Use `V2Ray` if your client accepts Base64 subscription links.

If you are unsure, start from the client-specific guide that matches your app.

## Why keep the control plane private?

The public repo is for distribution only. Infrastructure state, cloud credentials, SSH access, and secret material stay in the private operations repository.

## Where can I read the workflow and deployment details?

Use these docs:

- [Project details](./project-details.md)
- [Deployment and operations](./deployment.md)
- [Documentation map](./README.md)

## Can I mirror or fork this repository?

Yes, but treat it as a public feed and not as a source of private or persistent credentials.

## Where should I ask for help or report a problem?

- Use GitHub Discussions for client setup help, app compatibility questions, and guide requests.
- Use GitHub Issues only for broken public links, stale release artifacts, or incorrect status output.
- Check the live status dashboard before reporting node availability problems.

## Which client guides are available?

- Clash Verge Rev
- FlClash
- Clash Meta for Android
- Hiddify Next
- NekoBox
- v2rayNG
- Shadowrocket

Use the corresponding setup page in `site/` or the matching markdown guide in `docs/`.
