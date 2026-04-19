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
{{range .Alternates}}<link rel="alternate" hreflang="{{.Code}}" href="{{.URL}}">
{{end}}<meta property="og:type" content="website">
<meta property="og:title" content="{{.Title}}">
<meta property="og:description" content="{{.Description}}">
<meta property="og:url" content="{{.Canonical}}">
<meta property="og:image" content="{{.OGImage}}">
<meta property="og:locale" content="{{.LangAttr}}">
<meta name="twitter:card" content="summary_large_image">
<meta name="twitter:title" content="{{.Title}}">
<meta name="twitter:description" content="{{.Description}}">
<meta name="twitter:image" content="{{.OGImage}}">
<script type="application/ld+json">{{.JSONLD}}</script>
<style>` + css + `</style>
</head>
<body>
<div class="lang-switch">
  <span>{{.L10n.LanguageLabel}}</span>
  {{range .LanguageSw}}{{if .Current}}<strong>{{.Label}}</strong>{{else}}<a href="{{.URL}}">{{.Label}}</a>{{end}}
  {{end}}
</div>

<header class="hero">
  <h1>{{.Heading}}</h1>
  <p class="sub">{{.L10n.IndexSubTagline}}</p>
  <p class="meta">
    <span class="badge badge-green">{{.L10n.BadgeNodes}} · {{.Stats.TotalSelected}}</span>
    <span class="badge badge-blue">{{.L10n.BadgeAlive}} · {{.Stats.TotalAlive}}</span>
    <span class="badge badge-orange">{{.L10n.BadgeMedianRTT}} · {{.Stats.MedianLatencyMS}} ms</span>
    <span class="badge badge-gray">{{.L10n.BadgeUpdated}} · {{.UpdatedHuman}}</span>
  </p>
  <p class="repo">
    <a href="{{.RepoURL}}" class="btn-outline" target="_blank" rel="noopener">{{.L10n.StarButton}}</a>
  </p>
</header>

<section class="card">
  <h2>{{.L10n.OneClickHeading}}</h2>
  <p>{{.L10n.OneClickIntro}}</p>
  <div class="urls">
    <div class="url-row">
      <strong>{{.L10n.ColClash}}</strong>
      <code><a href="{{.URLClash}}" target="_blank" rel="noopener">{{.URLClash}}</a></code>
    </div>
    <div class="url-row">
      <strong>{{.L10n.ColSing}}</strong>
      <code><a href="{{.URLSing}}" target="_blank" rel="noopener">{{.URLSing}}</a></code>
    </div>
    <div class="url-row">
      <strong>{{.L10n.ColV2ray}}</strong>
      <code><a href="{{.URLV2ray}}" target="_blank" rel="noopener">{{.URLV2ray}}</a></code>
    </div>
  </div>
</section>

{{if .Countries}}
<section class="card">
  <h2>{{.L10n.ByCountryHeading}}</h2>
  <p>{{.L10n.ByCountryIntro}}</p>
  <div class="country-grid">
    {{range .Countries}}
    <a class="country-card" href="{{.URLPage}}">
      <div class="country-flag">{{.Flag}}</div>
      <div class="country-name">{{.Name}}</div>
      <div class="country-count">{{.Count}} {{$.L10n.NodesSuffix}}</div>
    </a>
    {{end}}
  </div>
</section>
{{end}}

{{if .Guides}}
<section class="card" id="guides">
  <h2>{{.L10n.GuidesHeading}}</h2>
  <p>{{.L10n.GuidesIntro}}</p>
  <ul class="client-list">
    {{range .Guides}}
    <li><a href="{{.URL}}"><strong>{{.Name}}</strong></a> · {{.OS}}</li>
    {{end}}
  </ul>
</section>
{{end}}

<section class="card">
  <h2>{{.L10n.ClientsHeading}}</h2>
  <ul class="client-list">
    <li>{{safe .L10n.ClientsWindows}}</li>
    <li>{{safe .L10n.ClientsMacOS}}</li>
    <li>{{safe .L10n.ClientsIOS}}</li>
    <li>{{safe .L10n.ClientsAndroid}}</li>
    <li>{{safe .L10n.ClientsLinux}}</li>
  </ul>
</section>

<section class="card">
  <h2>{{.L10n.FAQHeading}}</h2>
  <details><summary>{{.L10n.FAQ1Q}}</summary><p>{{.L10n.FAQ1A}}</p></details>
  <details><summary>{{.L10n.FAQ2Q}}</summary><p>{{.L10n.FAQ2A}}</p></details>
  <details><summary>{{.L10n.FAQ3Q}}</summary><p>{{.L10n.FAQ3A}}</p></details>
  <details><summary>{{.L10n.FAQ4Q}}</summary><p>{{.L10n.FAQ4A}}</p></details>
</section>

<footer>
  <p>{{.L10n.FooterLicense}}</p>
  <p class="small">{{.L10n.FooterDisclaimer}}</p>
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
{{range .Alternates}}<link rel="alternate" hreflang="{{.Code}}" href="{{.URL}}">
{{end}}<meta property="og:type" content="website">
<meta property="og:title" content="{{.Title}}">
<meta property="og:description" content="{{.Description}}">
<meta property="og:url" content="{{.Canonical}}">
<meta property="og:image" content="{{.OGImage}}">
<meta property="og:locale" content="{{.LangAttr}}">
<meta name="twitter:card" content="summary_large_image">
<meta name="twitter:title" content="{{.Title}}">
<meta name="twitter:description" content="{{.Description}}">
<meta name="twitter:image" content="{{.OGImage}}">
<script type="application/ld+json">{{.JSONLD}}</script>
<style>` + css + `</style>
</head>
<body>
<div class="lang-switch">
  <span>{{.L10n.LanguageLabel}}</span>
  {{range .LanguageSw}}{{if .Current}}<strong>{{.Label}}</strong>{{else}}<a href="{{.URL}}">{{.Label}}</a>{{end}}
  {{end}}
</div>

<nav class="breadcrumb">
  <a href="{{.HomeURL}}">{{.L10n.CountryBreadcrumb}}</a>
</nav>

<header class="hero">
  <h1>{{.Heading}}</h1>
  <p class="sub">{{.Description}}</p>
  <p class="meta">
    <span class="badge badge-green">{{.CurrentRows}} {{.L10n.NodesSuffix}}</span>
    <span class="badge badge-gray">{{.L10n.BadgeUpdated}} · {{.UpdatedHuman}}</span>
  </p>
</header>

<section class="card">
  <h2>{{.CurrentSubscribeHeading}}</h2>
  <div class="urls">
    <div class="url-row">
      <strong>{{.L10n.ColClash}}</strong>
      <code><a href="{{.URLClash}}" target="_blank" rel="noopener">{{.URLClash}}</a></code>
    </div>
    <div class="url-row">
      <strong>{{.L10n.ColSing}}</strong>
      <code><a href="{{.URLSing}}" target="_blank" rel="noopener">{{.URLSing}}</a></code>
    </div>
    <div class="url-row">
      <strong>{{.L10n.ColV2ray}}</strong>
      <code><a href="{{.URLV2ray}}" target="_blank" rel="noopener">{{.URLV2ray}}</a></code>
    </div>
  </div>
</section>

<section class="card">
  <h2>{{.L10n.CountryOtherHeading}}</h2>
  <div class="country-grid">
    {{range .Countries}}
    {{if ne .CC $.CurrentCC}}
    <a class="country-card" href="{{.URLPage}}">
      <div class="country-flag">{{.Flag}}</div>
      <div class="country-name">{{.Name}}</div>
      <div class="country-count">{{.Count}} {{$.L10n.NodesSuffix}}</div>
    </a>
    {{end}}
    {{end}}
  </div>
</section>

<footer>
  <p>{{.L10n.FooterLicense}}</p>
</footer>
</body>
</html>
`

// tplGuide renders a single client tutorial page.
const tplGuide = `<!DOCTYPE html>
<html lang="{{.LangAttr}}">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>{{.Title}}</title>
<meta name="description" content="{{.Description}}">
<meta name="keywords" content="{{.Keywords}}">
<link rel="canonical" href="{{.Canonical}}">
{{range .Alternates}}<link rel="alternate" hreflang="{{.Code}}" href="{{.URL}}">
{{end}}<meta property="og:type" content="article">
<meta property="og:title" content="{{.Title}}">
<meta property="og:description" content="{{.Description}}">
<meta property="og:url" content="{{.Canonical}}">
<meta property="og:image" content="{{.OGImage}}">
<meta property="og:locale" content="{{.LangAttr}}">
<meta name="twitter:card" content="summary_large_image">
<meta name="twitter:title" content="{{.Title}}">
<meta name="twitter:description" content="{{.Description}}">
<meta name="twitter:image" content="{{.OGImage}}">
<script type="application/ld+json">{{.JSONLD}}</script>
<style>` + css + `</style>
</head>
<body>
<div class="lang-switch">
  {{range .LanguageSw}}{{if .Current}}<strong>{{.Label}}</strong>{{else}}<a href="{{.URL}}">{{.Label}}</a>{{end}}
  {{end}}
</div>

<nav class="breadcrumb">
  <a href="{{.HomeURL}}">{{.L10n.HomeLinkText}}</a>
</nav>

<header class="hero">
  <h1>{{.Heading}}</h1>
  <p class="sub">{{.Description}}</p>
  <p class="meta">
    <span class="badge badge-blue">{{.ClientName}}</span>
    <span class="badge badge-gray">{{.OSList}}</span>
    <span class="badge badge-green">{{.L10n.UpdatedLabel}} · {{.UpdatedHuman}}</span>
  </p>
  <p class="repo">
    <a href="{{.DownloadURL}}" class="btn-outline" target="_blank" rel="noopener">{{printf .L10n.DownloadLabelTpl .ClientName}}</a>
  </p>
</header>

<section class="card">
  <h2>{{.L10n.SubscribeHeading}}</h2>
  <p>{{.L10n.SubscribeIntro}}</p>
  <div class="urls">
    <div class="url-row">
      <strong>{{.L10n.SubscribeLabel}}</strong>
      <code><a href="{{.SubscribeURL}}" target="_blank" rel="noopener">{{.SubscribeURL}}</a></code>
    </div>
  </div>
</section>

<section class="card">
  <h2>{{.L10n.StepsHeading}}</h2>
  <ol class="steps">
    {{range .Steps}}
    <li>
      <h3>{{.Title}}</h3>
      <p>{{.Body}}</p>
    </li>
    {{end}}
  </ol>
</section>

<section class="card">
  <h2>{{.L10n.TipsHeading}}</h2>
  {{range .Tips}}
  <details><summary>{{.Q}}</summary><p>{{.A}}</p></details>
  {{end}}
</section>

<section class="card">
  <h2>{{.L10n.OtherGuidesHeading}}</h2>
  <ul class="client-list">
    {{range .OtherGuides}}
    <li><a href="{{.URL}}">{{.Name}}</a> · {{.OS}}</li>
    {{end}}
  </ul>
</section>

<footer>
  <p><a href="{{.RepoURL}}" target="_blank" rel="noopener">GitHub</a> · MIT</p>
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
  .lang-switch {
    max-width: 900px; margin: 0 auto; padding: 12px 32px 0;
    text-align: right; font-size: 13px; color: #94a3b8;
  }
  .lang-switch a { color: #93c5fd; text-decoration: none; margin: 0 4px; }
  .lang-switch a:hover { text-decoration: underline; }
  .lang-switch strong { color: #f1f5f9; margin: 0 4px; }
  .lang-switch span { margin-right: 4px; }
  .hero {
    text-align: center; padding: 32px 20px 32px;
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
  .steps { padding-left: 24px; }
  .steps li { margin-bottom: 18px; }
  .steps h3 { margin: 0 0 6px; color: #f1f5f9; font-size: 16px; }
  .steps p { margin: 0; color: #cbd5e1; }
  footer { text-align: center; padding: 40px 20px; color: #94a3b8; font-size: 14px; }
  footer a { color: #93c5fd; }
  footer .small { font-size: 12px; margin-top: 8px; max-width: 680px; margin-left: auto; margin-right: auto; }
`
