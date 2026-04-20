# Architecture & reliability

How `free-vpn-subscriptions` fetches, verifies, and publishes working VPN nodes every hour.

## Pipeline at a glance

```
                               ┌─────────────────────────────────────┐
External scheduler (hourly +    │ runs the binary off-Actions to keep │
jitter, off-Actions)            │ source URLs and timing private;     │
                               │ ctx bounded by 30 min + SIGINT       │
                               └────────────┬────────────────────────┘
                                            │
                         ┌──────────────────┴────────────────────┐
                         │ 1. Fetch — internal/sources/fetch.go   │
                         │    17 upstream subscription feeds      │
                         │    formats: uri-list, base64, clash    │
                         │    ~4,800 raw nodes (ctx-aware HTTP)   │
                         └──────────────────┬────────────────────┘
                                            │
                         ┌──────────────────┴────────────────────┐
                         │ 2. TCP probe — internal/probe/probe.go │
                         │    concurrent net.Dialer.DialContext   │
                         │    to server:port, 3s timeout, 100-way │
                         │    parallel. Drops ~40% dead nodes.    │
                         └──────────────────┬────────────────────┘
                                            │
                         ┌──────────────────┴────────────────────┐
                         │ 3. TLS handshake — internal/probe/tls  │
                         │    tls.Dialer.DialContext for Trojan / │
                         │    Hysteria2 / VLESS+tls / VMess+tls.  │
                         │    Drops fake proxies (routers         │
                         │    accepting TCP but not TLS).         │
                         └──────────────────┬────────────────────┘
                                            │
                         ┌──────────────────┴────────────────────┐
                         │ 4. HTTP-over-proxy verify — verify/*   │
                         │    Batch sing-box subprocesses, each   │
                         │    node gets a local SOCKS5 inbound.   │
                         │    Real HTTP+HTTPS GET through the     │
                         │    proxy protocol, 2 rounds 45s apart. │
                         │    Keeps ≥50% success; sort by median. │
                         └──────────────────┬────────────────────┘
                                            │
                         ┌──────────────────┴────────────────────┐
                         │ 5. GeoIP enrich — internal/geoip       │
                         │    MaxMind GeoLite2-Country.mmdb       │
                         │    resolves each node's server IP      │
                         │    into a 2-letter country code.       │
                         └──────────────────┬────────────────────┘
                                            │
                         ┌──────────────────┴────────────────────┐
                         │ 6. Aggregate — internal/aggregate      │
                         │    - protocol whitelist                │
                         │    - RTT cap (max_rtt_ms = 4000)       │
                         │    - dedup by (proto,server,port,uuid) │
                         │    - sort by HTTP-median asc           │
                         │    - top-N (150)                       │
                         └──────────────────┬────────────────────┘
                                            │
                         ┌──────────────────┴────────────────────┐
                         │ 7. Emit — subscribe + pages + readme   │
                         │    clash.yaml / singbox.json /         │
                         │    v2ray-base64.txt + per-country      │
                         │    variants + docs/*.html + 7 READMEs  │
                         └──────────────────┬────────────────────┘
                                            │
                                  git commit && git push
```

## Step 1 — Fetch

Defined in [`internal/sources/fetch.go`](./internal/sources/fetch.go).

Each source in [`config.yaml`](./config.yaml) has a `format`:

| Format | Example | Parser |
|---|---|---|
| `uri-list` | `vless://…\nvmess://…\ntrojan://…` | split by newline, parse each URI |
| `base64` | base64 of a uri-list | decode then parse |
| `clash` | YAML with `proxies:` list | `yaml.Unmarshal` → map each proxy |

**Current sources**: 17 public, volunteer-maintained upstream feeds covering uri-list, base64 and clash-YAML formats. Exact URLs live in `config.yaml`; we don't reproduce them here because enumerating them in public docs makes them easier for upstreams to rate-limit.

Fetch errors on a single source are swallowed — they log `[skip]` and the pipeline proceeds with the remaining sources. This makes one flaky upstream not break the run.

**Freshness**: each run pulls HEAD of each source. Sources typically update once or twice a day; we re-probe hourly regardless.

## Step 2 — TCP probe

Defined in [`internal/probe/probe.go`](./internal/probe/probe.go).

```go
dialer := &net.Dialer{Timeout: 3 * time.Second}
conn, err := dialer.DialContext(ctx, "tcp", host:port)
```

- 100-way concurrent (configurable via `probe.concurrency`).
- Latency (time to TCP handshake) is recorded on the `Node.LatencyMS` field.
- Nodes that fail to connect are dropped entirely.
- **Context-aware**: the top-level `fnctl aggregate` wraps `context.Background()` with both `signal.NotifyContext` (SIGINT/SIGTERM) and a 30-minute deadline (`runDeadline`). Cancelling `ctx` propagates here so pending dials abort immediately instead of leaking goroutines when the host scheduler sends SIGTERM or when a run approaches the timer's `TimeoutStartSec` ceiling.

Typical result: ~60% of fetched nodes pass this stage. The rest are dead hosts, port-blocked by our runner's network, or have revoked IPs.

## Step 3 — TLS handshake

Defined in [`internal/probe/tls.go`](./internal/probe/tls.go).

This is the crucial stage that separates us from most "free VPN list" repos: we verify the node actually speaks TLS, not just "something listens on port 443."

```go
tlsDialer := &tls.Dialer{
    NetDialer: &net.Dialer{Timeout: 3 * time.Second},
    Config: &tls.Config{
        ServerName:         sni,
        InsecureSkipVerify: true,  // free nodes often use self-signed
    },
}
conn, err := tlsDialer.DialContext(ctx, "tcp", addr)
```

- `InsecureSkipVerify: true` because free nodes commonly present self-signed or expired certs. We care only that a real TLS handshake completes.
- `ServerName` uses the node's declared SNI; falls back to the bare server hostname.
- Applied to: Trojan, Hysteria2, VLESS+tls, VMess+tls.
- **Reality is intentionally skipped** — Reality nodes spoof their ClientHello to look like legitimate targets (e.g. microsoft.com), so a real handshake against the proxy tells us nothing useful.
- **Context-aware**: `tls.Dialer.DialContext` respects the same run-level ctx as the TCP probe, so SIGINT / deadline propagate through handshakes already in flight.

Typical result: ~60% of TCP-alive TLS nodes pass this stage. The rest are fake — cheap VPS routers that accept any TCP connection and silently drop anything that isn't HTTP/1.1.

## Step 4 — HTTP-over-proxy verify

Defined in [`internal/verify/`](./internal/verify/) — `verify.go` orchestrates, `singbox.go` builds config + spawns subprocess, `outbound.go` translates each `Node` into a sing-box outbound spec, `probe.go` sends the real HTTP request through a SOCKS5 dialer.

TCP + TLS tells us the transport is alive, but nothing about whether **the proxy protocol itself works**. A Trojan server can complete TLS yet reject our password. A VLESS+Reality server can handshake yet drop traffic because the short-id is stale. To close this gap, before publishing we actually run HTTP traffic through each candidate.

Pipeline per run:

1. **Candidate selection** — take the TLS-alive pool, sort by measured TCP RTT ascending, cap to `verify.candidate_pool` (default 900). Verifying ~900 nodes is the sweet spot: large enough that we still publish 150 after the funnel, small enough to finish in ~10 minutes.

2. **Config pre-filter** — run `sing-box check -c <single-node.json>` on each candidate in parallel. Corrupt SS ciphers, garbage UUIDs, unsupported flow options are dropped here. Without this, a single bad node's config error would make `sing-box check` abort the whole batch. Typical pass rate at this stage: ~90% (drops ~60–70 malformed entries).

3. **Batch start** — group the survivors into batches of `verify.batch_size` (default 40) and spawn one sing-box subprocess per batch:
   - 40 outbounds (one per node, tagged `out-0`…`out-39`)
   - 40 mixed inbounds on `127.0.0.1:base_port+i` (default base `20000`)
   - `route.rules` one-to-one mapping `in-i` → `out-i`, default outbound `direct`
   - `Setpgid: true` so we can kill the whole process group on cleanup
   - We poll `tryDial(base_port)` until sing-box binds or `startup_timeout_ms` (default 10 s) elapses.

4. **HTTP probe** — for each outbound, dial its local SOCKS5 inbound via `golang.org/x/net/proxy` and send real requests:
   - `http://www.gstatic.com/generate_204` — expects 204
   - `https://www.cloudflare.com/cdn-cgi/trace` — expects 200

   Both requests traverse the full proxy stack: auth, transport (TLS/WS/gRPC/Reality), and the egress network. A status in `[200, 400)` counts as success. Latency is the HTTP round-trip, not TCP RTT.

5. **Stability rounds** — repeat the probe `verify.rounds` (default 2) times with `verify.round_gap_ms` (default 45 s) between rounds. Nodes that pass once but die 45 seconds later are filtered out by the round gap. Nodes whose `successes ≥ (rounds × targets) / 2` (so ≥2 out of 4 for defaults) survive.

6. **Rank** — on survivors, `Node.LatencyMS` is overwritten with the **median HTTP-through-proxy latency** across all successful attempts. The original TCP RTT is preserved on `Node.TCPLatencyMS` for display. Aggregation then sorts by the HTTP median.

Typical funnel on a recent run: **17 sources → ~4,800 raw → ~2,900 TCP-alive → ~2,600 TLS-OK → ~840 config-valid → ~280 HTTP-verified → top 150 published**. Total runtime ~10 min, of which step 4 takes ~7 min (15 batches × 12 s × 2 rounds + 45 s gap + cleanup).

**Why sing-box subprocess, not Go library?** sing-box's internal Go API is unstable and its dependency tree is large. Embedding it would couple our release cadence to theirs. The subprocess approach lets us pin sing-box independently (`/usr/local/bin/sing-box`, currently v1.13.8) and replace it without rebuilding the binary.

**Why not wireshark-level diagnostics?** We want evidence that the proxy works end-to-end, not protocol-level debugging. A successful HTTP GET through the full stack is that evidence.

## Step 5 — GeoIP

Defined in [`internal/geoip/geoip.go`](./internal/geoip/geoip.go).

Uses MaxMind's free **GeoLite2-Country** database mirrored at `P3TERX/GeoLite.mmdb`. On each run the binary:
1. Downloads the latest mmdb (cached at `output/.cache/` — skipped if fresh).
2. Resolves every node's server hostname to an IPv4 address (with a small in-process DNS cache).
3. Looks up the country code, stored on `Node.Country`.

Soft-failure design: if the GeoIP database can't be downloaded, the pipeline still produces global outputs — only the per-country filter is skipped that run.

## Step 6 — Aggregate

Defined in [`internal/aggregate/aggregate.go`](./internal/aggregate/aggregate.go).

Applied in strict order:

1. **Protocol filter** — keep only `vless | vmess | trojan | shadowsocks | hysteria2` (configurable).
2. **RTT cap** — drop anything above `aggregate.max_rtt_ms` (default 4000 ms). Note: `Node.LatencyMS` is now the HTTP-over-proxy median from step 4, typically 2–3× higher than raw TCP RTT, which is why the cap is set higher than the older TCP-only 1500 ms.
3. **Dedup** — key is `(protocol, server, port, uuid_or_password)`. When the same endpoint appears from multiple sources, we keep the one with lowest measured latency.
4. **Sort** ascending by HTTP-median latency.
5. **Top-N** — slice to `aggregate.top_n` (default 150). Shorter lists are easier for clients to process and keep the selector group responsive.

The median-latency statistic reported in `status.json` and READMEs uses the true median: the middle element for odd-length lists, and the mean of the two middle elements for even-length lists. Implemented as `medianLatency` in `aggregate.go` — previously we returned the upper-middle value, biasing the number up by ~5–10 ms on typical runs.

A country variant is emitted per ISO-2 country when that country has at least `geoip.min_per_country` nodes (default 3) — avoids publishing a `clash-BZ.yaml` with only 1 node.

## Step 7 — Emit

The binary writes:

| Output | Path | Consumer |
|---|---|---|
| Clash | `output/clash.yaml` | Clash Verge, ClashX, mihomo |
| sing-box | `output/singbox.json` | sing-box CLI + mobile apps |
| v2ray-base64 | `output/v2ray-base64.txt` | v2rayN, v2rayNG, Shadowrocket |
| Per-country variants | `output/by-country/{clash,singbox,v2ray-base64}-XX.{yaml,json,txt}` | targeted subscriptions |
| Status | `output/status.json` | summary for dashboards |
| READMEs | `README.md`, `README_CN.md`, …, `README_RU.md` | GitHub repo front page |
| Pages site | `docs/index.{en,zh,ja,ko,es,pt,ru}.html`, `docs/XX.{en,zh,ja,ko,es,pt,ru}.html` per qualifying country, `docs/guides/{slug}.{en,zh}.html`, `docs/sitemap.xml`, `docs/robots.txt` | SEO landing for au1rxx.github.io |

The Clash emitter builds a `proxy-groups` selector with a URL-test probe (`http://www.gstatic.com/generate_204`, 300s interval) — clients auto-pick the fastest node in real use.

## How the Pages site is actually served

This trips people up, so it gets its own section: **GitHub Pages does not run Go, Node, or any backend**. It is a static CDN that only serves whatever pre-built files exist under `docs/` on the `main` branch. Two completely separate runtime environments are at play:

```
┌──────────────────────────────────┐         ┌────────────────────────────┐
│ External scheduler (off-Actions) │  push   │ GitHub Pages (static CDN)  │
│ hourly + ±30 min jitter          │────────▶│                            │
│                                  │         │ au1rxx.github.io/...       │
│ go build → fnctl aggregate       │         │   ├─ index.html            │
│   ├─ fetch + probe + rank        │         │   ├─ index.zh.html         │
│   └─ internal/pages.Generate()   │         │   ├─ us.html / us.zh.html  │
│        writes docs/*.html        │         │   ├─ guides/*.html         │
│ git add docs/ && git push        │         │   ├─ sitemap.xml           │
└──────────────────────────────────┘         │   └─ robots.txt            │
                                             │                            │
                                             │  ← browser requests static │
                                             └────────────────────────────┘
```

`internal/pages/*.go`, `cmd/fnctl/*.go`, the whole Go source tree — **none of that ships to Pages**. The scheduler runs the binary, emits static HTML, pushes, exits. Pages serves the HTML verbatim. There is no server-side rendering, no edge functions, no runtime. Aggregation runs **off** GitHub Actions on purpose: the public Actions log would otherwise expose every upstream source URL, error message, and the exact cron timing — operational metadata that upstream maintainers can use to rate-limit or block us.

Practical consequences:
- All internationalization (i18n) must be baked into distinct URLs (`index.html` vs `index.zh.html`) — there is no `Accept-Language` negotiation.
- All dynamic values (node counts, timestamps, RTT medians) are re-computed and re-written into the HTML on each Actions run. A change in the live stats requires a new push.
- `.nojekyll` exists in `docs/` so Pages serves files verbatim and doesn't try to run Jekyll over them.

## Multilingual rendering

The site supports **two tiers of localization** so we can reach non-English searchers without being on the hook for translating 6 × 4 = 24 hand-written guide tutorials.

| Tier | Locales | What gets rendered |
|---|---|---|
| **Full** (`supportedLocales`) | `en`, `zh`, `ja`, `ko`, `es`, `pt`, `ru` — 7 total | Index + per-country pages |
| **Guides** (`supportedGuideLocales`) | `en`, `zh` — 2 total | Step-by-step client setup pages under `/guides/` |

Source of truth: [`internal/pages/l10n.go`](./internal/pages/l10n.go) (index/country strings, all 7 locales) and [`internal/pages/guides.go`](./internal/pages/guides.go) (guide content, en+zh only).

| Page type | English (canonical) | Non-English locales |
|---|---|---|
| Index | `/index.html` served as `/` | `/index.{loc}.html` for each of `zh`, `ja`, `ko`, `es`, `pt`, `ru` |
| Country | `/us.html`, `/hk.html`, … | `/us.{loc}.html`, `/hk.{loc}.html`, … — one per full locale |
| Guide | `/guides/clash-verge.html`, … | `/guides/clash-verge.zh.html` only (no `.ja/.ko/.es/.pt/.ru`) |

**Guide fallback**: an index page rendered in e.g. Japanese links to the **English** guide URL (`guides/clash-verge.html`, not `.ja.html`) because the Japanese file does not exist. This keeps the click working — better UX than a 404 — and aligns with Google's hreflang fallback expectations (a locale with no translated resource legitimately points at the default one). The guide's own `<link rel="alternate" hreflang>` tags advertise only `en` + `zh` + `x-default`, so search engines won't synthesize a Japanese guide URL.

Every page advertises its alternates to search engines in two places:

1. **`<link rel="alternate" hreflang="…">` in `<head>`** — one tag per locale the page itself exists in, plus `x-default` → English.
2. **`<xhtml:link rel="alternate" hreflang="…">` in `sitemap.xml`** — same alternates declared at the sitemap level so Google discovers them even if it hits the sitemap before any page.

Index + country sitemap entries carry all 7 `hreflang` alternates; guide entries carry only `en` + `zh-Hans` + `x-default`.

A visible language switcher at the top of every page lets users toggle manually without relying on crawler-only hreflang. Switchers on index/country pages offer all 7 locales; switchers on guide pages show only en + zh.

## SEO surface

Every HTML page carries a consistent set of metadata. Implementation: [`internal/pages/pages.go`](./internal/pages/pages.go) + [`internal/pages/templates.go`](./internal/pages/templates.go).

| Signal | Purpose |
|---|---|
| Per-locale `<title>` and `<meta description>` | Direct search snippet content |
| `<link rel="canonical">` | Points at the locale-specific URL to avoid dup-content flagging |
| `<link rel="alternate" hreflang>` (× N locales + `x-default`) | Tells Google which version to serve per user locale |
| `og:type / og:locale / og:image` + Twitter card tags | Link-preview cards on Reddit, Slack, Twitter, Discord |
| JSON-LD `WebSite` | Sitelinks search box + site name |
| JSON-LD `SoftwareApplication` with `AggregateRating` | Rich-result eligibility |
| JSON-LD `FAQPage` | FAQ rich snippet on the index page |
| JSON-LD `WebPage` + `BreadcrumbList` | Breadcrumb trail on country pages |
| JSON-LD `HowTo` | Step-by-step rich snippet on guide pages (one `HowToStep` per guide step) |
| `sitemap.xml` + `robots.txt` | Crawl discovery; hourly `changefreq` on home, weekly on guides |

`inLanguage` is set on every JSON-LD entity so Google can serve the right version per query locale. All page weight stays under 20 KB (inline CSS, no JS, no external fetches) — Core Web Vitals green by construction.

### Hero image

`assets/hero.png` (525 KB, 1536×1024, 70% smaller than the original 1.77 MB GPT-generated output, compressed with `pngquant --quality=75-92 --strip` to drop metadata while keeping the visible fidelity identical) is the `og:image` referenced on every page. It's also the README hero. Smaller weight = faster link-preview fetches on Reddit/Twitter/Slack and lighter crawls by ImageBot.

## Adding a new locale

Two paths depending on how much translation effort you're committing.

### Path A — index + country pages only (cheap, ~45 strings)

This is how `ja`, `ko`, `es`, `pt`, `ru` were added. Use this when you don't want to hand-translate the multi-paragraph client setup tutorials.

1. Add an entry to `pageLocales` in [`internal/pages/l10n.go`](./internal/pages/l10n.go) — one full `pageL10n{…}` struct (~45 fields). Reuse translations from the corresponding `internal/readme/locale_*.go` where possible.
2. Append the new code to `supportedLocales` in [`internal/pages/guides.go`](./internal/pages/guides.go).
3. Add a case for the code in `localeLangAttr` if the HTML `lang` attribute differs from the locale code (e.g. `zh` → `zh-Hans`).
4. Add a case in `hreflangCode` if the Google-expected hreflang differs from the locale code.

The new locale automatically: gets its own index + country pages, shows up in every language switcher and every hreflang list, and appears in the sitemap with all 7 alternates. The index page's **Guides** section in the new locale falls back to the English guide pages (no 404).

### Path B — also render full guide translations (expensive, ~50 HTML-bearing strings per guide × 4 guides)

Do this when you have a native speaker willing to write each step-by-step tutorial body.

1. Do everything in Path A.
2. For each `guideSpec` in [`internal/pages/guides.go`](./internal/pages/guides.go), add the new code to its `L10n` map — a full `guideL10n{…}` with `Steps[]` and `Tips[]` bodies (HTML allowed).
3. Append the new code to `supportedGuideLocales` (same file). This is what switches the guide loop in `Generate()` from "English only" to "also render this locale".

Once the code is in `supportedGuideLocales`, the guide sitemap entries automatically gain its `hreflang`, the guide language switcher adds it, and index pages in that locale link to the translated guide URL instead of the English fallback.

## Reliability mechanisms summary

| Risk | Mitigation |
|---|---|
| Upstream source goes offline | Per-source try/catch; other sources still contribute |
| Stale nodes in upstream feeds | Hourly cron + TCP probe + TLS probe + HTTP-over-proxy verify drop dead/fake ones |
| "Open port but dead proxy" trap | TLS handshake filter catches routers that accept TCP but can't speak TLS |
| "Handshake succeeds but proxy rejects traffic" trap | HTTP-over-proxy verify actually forwards a GET through each proxy protocol before publishing |
| Flaky nodes (pass once, die seconds later) | Verify runs 2 rounds 45 s apart; only nodes with ≥50 % success rate across rounds×targets are kept |
| Broken config (bad cipher, wrong UUID, unsupported flow) | `sing-box check` pre-filter at the verify stage drops them before a probe slot is wasted |
| Node goes offline between probe and user | Clash selector group + URL-test fallback; per-country has multiple nodes |
| Single country dominates | Top-N cap ensures geographic spread; per-country files for targeted subs |
| GeoIP DB outage | Soft-fail — global outputs still produced |
| Host runner throttled | Concurrency bounded to 100; timeout_ms keeps one stuck probe from blocking |
| Run exceeds scheduler budget | `runDeadline = 30m` in `cmd/fnctl/main.go` caps the whole aggregate via `context.WithTimeout`; the host timer's `TimeoutStartSec=45m` is defense-in-depth |
| Runner sends SIGINT mid-run | `signal.NotifyContext` catches SIGINT/SIGTERM; ctx propagates into fetch + probe + TLS, aborting pending work instead of leaking goroutines |
| Race between scheduled bot and human commits | systemd `Type=oneshot` + a single timer ensures only one publish run at a time; the script also refuses to run on a dirty tree |

## What we deliberately do *not* do

- **No bandwidth / throughput measurement** — we verify that a small HTTP GET completes (204 / 200), which proves the proxy forwards traffic, but we don't measure Mbps. A 50 ms node might still be slow for video streaming.
- **No manual curation** — every node comes from a public upstream; we don't edit the list.
- **No analytics / telemetry** — the static site has zero JS and zero third-party resources.
- **No link to non-free providers** — affiliate links would compromise our neutrality.

## Tuning knobs

All in [`config.yaml`](./config.yaml):

```yaml
probe:
  timeout_ms: 3000       # per-node TCP dial timeout
  concurrency: 100       # parallel probe goroutines
  max_nodes_per_source: 500
  tls_verify: true       # toggle step 3

verify:                  # step 4 — HTTP-over-proxy
  enabled: true
  candidate_pool: 900    # top-N by TCP RTT fed into verify
  batch_size: 40         # outbounds per sing-box subprocess
  base_port: 20000       # first SOCKS5 inbound port
  concurrency: 20        # parallel probes within a batch
  timeout_ms: 6000       # per-target HTTP timeout
  rounds: 2              # repeat probe this many times
  round_gap_ms: 45000    # sleep between rounds
  targets:
    - http://www.gstatic.com/generate_204
    - https://www.cloudflare.com/cdn-cgi/trace
  sing_box_bin: /usr/local/bin/sing-box
  startup_timeout_ms: 10000  # wait for sing-box to bind inbounds

aggregate:
  top_n: 150             # final list size
  max_rtt_ms: 4000       # HTTP-median latency cap (was 1500 for raw TCP)
  protocols: [vless, vmess, trojan, shadowsocks, hysteria2]

geoip:
  enabled: true
  min_per_country: 3     # threshold for per-country output
```

**Prerequisite**: `sing-box` must be available at `verify.sing_box_bin` (default looks up `sing-box` in `PATH`). The production runner uses `/usr/local/bin/sing-box` v1.13.8. If the binary is missing, set `verify.enabled: false` to skip step 4 entirely (the pipeline falls back to TCP+TLS-only, same quality bar as before this stage was added).

## Monitoring

- **Latest status**: `output/status.json` — totals and breakdowns per run.
- **Workflow log**: Actions tab → `Aggregate & Publish` → latest run.
- **Output diff**: `git log -p output/` shows exactly which nodes changed hour over hour.
