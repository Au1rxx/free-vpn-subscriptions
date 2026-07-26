# Contributing

Thanks for helping. This repository is mostly machine-operated — an aggregator
runs hourly, verifies nodes, and pushes the result — so the useful ways to
contribute are a little different from a typical project.

## Before opening an issue about a broken node

Published nodes come from public, volunteer-run sources and they die
constantly. Before reporting, please check:

1. **Re-fetch your subscription.** The feed is rewritten every hour. A node
   that failed ten minutes ago may already be gone from the current file.
2. **Try another node.** `clash.yaml` ships with a `url-test` group that picks
   the fastest node that actually works *from your network*. Expect a
   meaningful share of nodes to be unreachable from any given location — a
   node that passed our probe can still be blocked by your ISP.
3. **Check the `updated` badge** in the README. If it is many hours old, the
   publisher itself has stalled, which is a real bug worth reporting.

A single dead node is normal and not worth an issue. An entire subscription
file that is empty, malformed, or hours stale is a bug — please report that.

## Suggesting a source

Open an issue with:

- the URL,
- the format (`uri-list`, `base64`, or Clash YAML),
- roughly how often it updates,
- whether it needs a token or account (we only use sources that need neither).

Note that the live source list is not published. Enumerating upstream URLs in
public documentation makes it easier for those upstreams to be rate-limited or
blocked, which would hurt everyone using them — see
[ARCHITECTURE.md](./ARCHITECTURE.md) for the reasoning.

## Code changes

```bash
go build ./...
go vet ./...
go test ./...
```

All three must pass. A few things worth knowing before you start:

- **README and the `docs/` site are generated.** Editing them by hand does
  nothing — the next hourly publish overwrites your change. Edit the
  generators in `internal/readme/` and `internal/pages/` instead.
- **`output/` is generated too.** It is the published payload, rewritten every
  run.
- **Published node names must stay stable across runs.** They are not allowed
  to encode a node's rank or position: clients key their latency history off
  the name, and unstable names also rewrite the whole subscription file on
  every publish. There are tests covering this.
- **Never commit credentials**, and do not add upstream source URLs to files
  that ship publicly.

## Security

Please do not put exploit details or credentials in a public issue. Use
GitHub's private
["Report a vulnerability"](https://github.com/Au1rxx/free-vpn-subscriptions/security/advisories/new)
form instead.

## What this project will not do

To set expectations, so nobody spends effort on a PR that cannot be merged:

- We do not operate servers or add paid/private nodes.
- We do not accept sources that require an account, token, or payment.
- We do not guarantee availability. Every node here is third-party and
  unwarranted; see the disclaimer in the README.
