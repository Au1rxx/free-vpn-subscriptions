package pages

// tplIndex is the landing page — global subscription + country grid + FAQ.
const tplIndex = `<!DOCTYPE html>
<html lang="{{.LangAttr}}">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>{{.Title}}</title>
<meta name="description" content="{{.Description}}">
<meta name="keywords" content="{{.Keywords}}">
<link rel="canonical" href="{{.Canonical}}">
<meta property="og:type" content="website">
<meta property="og:title" content="{{.Title}}">
<meta property="og:description" content="{{.Description}}">
<meta property="og:url" content="{{.Canonical}}">
<meta property="og:image" content="{{.OGImage}}">
<meta name="twitter:card" content="summary_large_image">
<meta name="twitter:title" content="{{.Title}}">
<meta name="twitter:description" content="{{.Description}}">
<meta name="twitter:image" content="{{.OGImage}}">
<script type="application/ld+json">{{.JSONLD}}</script>
<style>` + css + `</style>
</head>
<body>
<header class="hero">
  <h1>{{.Heading}}</h1>
  <p class="sub">The easiest way to get a working free VPN — copy a subscription link, paste it into your client, connect.</p>
  <p class="meta">
    <span class="badge badge-green">nodes · {{.Stats.TotalSelected}}</span>
    <span class="badge badge-blue">alive · {{.Stats.TotalAlive}}</span>
    <span class="badge badge-orange">median RTT · {{.Stats.MedianLatencyMS}} ms</span>
    <span class="badge badge-gray">updated · {{.UpdatedHuman}}</span>
  </p>
  <p class="repo">
    <a href="{{.RepoURL}}" class="btn-outline" target="_blank" rel="noopener">⭐ Star on GitHub</a>
  </p>
</header>

<section class="card">
  <h2>🚀 One-Click Subscribe</h2>
  <p>Copy the URL that matches your client and paste it into the subscription import field:</p>
  <div class="urls">
    <div class="url-row">
      <strong>Clash / Clash Verge / ClashX</strong>
      <code><a href="{{.URLClash}}" target="_blank" rel="noopener">{{.URLClash}}</a></code>
    </div>
    <div class="url-row">
      <strong>sing-box</strong>
      <code><a href="{{.URLSing}}" target="_blank" rel="noopener">{{.URLSing}}</a></code>
    </div>
    <div class="url-row">
      <strong>v2rayN / v2rayNG / Shadowrocket / NekoBox</strong>
      <code><a href="{{.URLV2ray}}" target="_blank" rel="noopener">{{.URLV2ray}}</a></code>
    </div>
  </div>
</section>

{{if .Countries}}
<section class="card">
  <h2>🌍 By Country</h2>
  <p>Want nodes in a specific region only? Choose a targeted subscription:</p>
  <div class="country-grid">
    {{range .Countries}}
    <a class="country-card" href="{{.URLPage}}">
      <div class="country-flag">{{.Flag}}</div>
      <div class="country-name">{{.Name}}</div>
      <div class="country-count">{{.Count}} nodes</div>
    </a>
    {{end}}
  </div>
</section>
{{end}}

<section class="card">
  <h2>🧩 Supported Clients</h2>
  <ul class="client-list">
    <li><strong>Windows</strong>: v2rayN, Clash Verge, Hiddify, NekoRay</li>
    <li><strong>macOS</strong>: ClashX Pro, Clash Verge, sing-box, Hiddify</li>
    <li><strong>iOS</strong>: Shadowrocket, Stash, Loon, sing-box, Hiddify</li>
    <li><strong>Android</strong>: v2rayNG, NekoBox, Clash Meta for Android, Hiddify, sing-box</li>
    <li><strong>Linux</strong>: mihomo (Clash.Meta), sing-box, v2ray-core</li>
  </ul>
</section>

<section class="card">
  <h2>❓ FAQ</h2>
  <details><summary>Is this actually free?</summary><p>Yes. Nodes are operated by third-party volunteers. We don't run any servers — we just test, rank, and repackage what's already public.</p></details>
  <details><summary>How fresh is the data?</summary><p>A GitHub Action runs every hour: pulls all upstream sources, TCP+TLS probes every node, drops anything dead, sorts by latency, commits new output files.</p></details>
  <details><summary>Can I trust these nodes?</summary><p>Free nodes see all your traffic. Never use them for banking, login, or anything sensitive. Fine for bypassing geo-blocks on public content. Use your own VPS or a paid provider for real privacy.</p></details>
  <details><summary>Why do some nodes fail?</summary><p>We verify TCP reachability and TLS handshakes, but a node may still have an expired quota, bad routing, or revoked certs. Try a few; the selector group gives you fallbacks.</p></details>
</section>

<footer>
  <p>Open source on <a href="{{.RepoURL}}" target="_blank" rel="noopener">GitHub</a>. MIT licensed.</p>
  <p class="small">This project aggregates publicly-shared proxy configurations. We do not operate any servers and make no guarantees. Comply with the laws of your jurisdiction.</p>
</footer>
</body>
</html>
`

// tplCountry is a per-country page — narrower, targeting long-tail queries
// like "free HK VPN subscription" or "free US VPN clash".
const tplCountry = `<!DOCTYPE html>
<html lang="{{.LangAttr}}">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>{{.Title}}</title>
<meta name="description" content="{{.Description}}">
<meta name="keywords" content="{{.Keywords}}">
<link rel="canonical" href="{{.Canonical}}">
<meta property="og:type" content="website">
<meta property="og:title" content="{{.Title}}">
<meta property="og:description" content="{{.Description}}">
<meta property="og:url" content="{{.Canonical}}">
<meta property="og:image" content="{{.OGImage}}">
<meta name="twitter:card" content="summary_large_image">
<meta name="twitter:title" content="{{.Title}}">
<meta name="twitter:description" content="{{.Description}}">
<meta name="twitter:image" content="{{.OGImage}}">
<script type="application/ld+json">{{.JSONLD}}</script>
<style>` + css + `</style>
</head>
<body>
<nav class="breadcrumb">
  <a href="{{.HomeURL}}">← All countries</a>
</nav>

<header class="hero">
  <h1>{{.Heading}}</h1>
  <p class="sub">{{.CurrentRows}} free VPN nodes in {{.CurrentName}}, TCP+TLS verified, refreshed hourly.</p>
  <p class="meta">
    <span class="badge badge-green">{{.CurrentRows}} nodes</span>
    <span class="badge badge-gray">updated · {{.UpdatedHuman}}</span>
  </p>
</header>

<section class="card">
  <h2>🚀 Subscribe to {{.CurrentName}} nodes only</h2>
  <div class="urls">
    <div class="url-row">
      <strong>Clash</strong>
      <code><a href="{{.URLClash}}" target="_blank" rel="noopener">{{.URLClash}}</a></code>
    </div>
    <div class="url-row">
      <strong>sing-box</strong>
      <code><a href="{{.URLSing}}" target="_blank" rel="noopener">{{.URLSing}}</a></code>
    </div>
    <div class="url-row">
      <strong>v2ray base64</strong>
      <code><a href="{{.URLV2ray}}" target="_blank" rel="noopener">{{.URLV2ray}}</a></code>
    </div>
  </div>
</section>

<section class="card">
  <h2>🌍 Other countries</h2>
  <div class="country-grid">
    {{range .Countries}}
    {{if ne .CC $.CurrentCC}}
    <a class="country-card" href="{{.URLPage}}">
      <div class="country-flag">{{.Flag}}</div>
      <div class="country-name">{{.Name}}</div>
      <div class="country-count">{{.Count}} nodes</div>
    </a>
    {{end}}
    {{end}}
  </div>
</section>

<footer>
  <p>Open source on <a href="{{.RepoURL}}" target="_blank" rel="noopener">GitHub</a>. MIT licensed.</p>
</footer>
</body>
</html>
`

// css is the inline stylesheet — kept minimal so page weight stays under 20 KB
// and no external fetch can block render.
const css = `
  :root { color-scheme: dark; }
  * { box-sizing: border-box; }
  body {
    margin: 0; font-family: -apple-system,Segoe UI,Roboto,sans-serif;
    background: linear-gradient(180deg,#0f172a,#1e293b 400px);
    color: #e2e8f0; line-height: 1.55;
  }
  .hero {
    text-align: center; padding: 64px 20px 32px;
    max-width: 900px; margin: 0 auto;
  }
  h1 { font-size: clamp(28px, 5vw, 44px); margin: 0 0 12px; }
  .sub { font-size: 18px; color: #cbd5e1; max-width: 680px; margin: 0 auto 16px; }
  .meta { display: flex; justify-content: center; flex-wrap: wrap; gap: 8px; margin: 16px 0; }
  .badge {
    display: inline-block; padding: 4px 12px; border-radius: 999px;
    font-size: 13px; font-weight: 600;
  }
  .badge-green { background: #065f46; color: #d1fae5; }
  .badge-blue { background: #1e40af; color: #dbeafe; }
  .badge-orange { background: #92400e; color: #fed7aa; }
  .badge-gray { background: #374151; color: #d1d5db; }
  .repo { margin-top: 24px; }
  .btn-outline {
    border: 1px solid #3b82f6; color: #93c5fd; padding: 10px 20px;
    border-radius: 8px; text-decoration: none; font-weight: 600;
  }
  .btn-outline:hover { background: #1e3a8a; }
  .breadcrumb { padding: 16px 32px; max-width: 900px; margin: 0 auto; }
  .breadcrumb a { color: #93c5fd; text-decoration: none; font-size: 14px; }
  .card {
    background: #1e293b; border: 1px solid #334155; border-radius: 12px;
    padding: 28px 32px; margin: 20px auto; max-width: 900px;
  }
  .card h2 { margin-top: 0; font-size: 22px; color: #f1f5f9; }
  .urls { display: flex; flex-direction: column; gap: 12px; margin-top: 16px; }
  .url-row { display: flex; flex-direction: column; gap: 4px; }
  .url-row strong { color: #f1f5f9; font-size: 14px; }
  .url-row code {
    background: #0f172a; padding: 10px 14px; border-radius: 6px;
    font-size: 13px; overflow-x: auto; border: 1px solid #334155;
  }
  .url-row code a { color: #60a5fa; text-decoration: none; word-break: break-all; }
  .country-grid {
    display: grid; gap: 12px; margin-top: 16px;
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  }
  .country-card {
    display: flex; flex-direction: column; align-items: center; gap: 4px;
    padding: 20px 12px; background: #0f172a; border: 1px solid #334155;
    border-radius: 10px; text-decoration: none; color: inherit;
    transition: border-color .2s, transform .1s;
  }
  .country-card:hover { border-color: #3b82f6; transform: translateY(-2px); }
  .country-flag { font-size: 32px; }
  .country-name { font-weight: 600; color: #f1f5f9; font-size: 14px; }
  .country-count { color: #94a3b8; font-size: 13px; }
  .client-list { padding-left: 20px; }
  .client-list li { margin-bottom: 6px; color: #cbd5e1; }
  details { margin: 10px 0; padding: 12px 16px; background: #0f172a; border-radius: 8px; border: 1px solid #334155; }
  details summary { cursor: pointer; font-weight: 600; color: #f1f5f9; }
  details p { margin: 12px 0 0; color: #cbd5e1; }
  footer { text-align: center; padding: 40px 20px; color: #94a3b8; font-size: 14px; }
  footer a { color: #93c5fd; }
  footer .small { font-size: 12px; margin-top: 8px; max-width: 680px; margin-left: auto; margin-right: auto; }
`
