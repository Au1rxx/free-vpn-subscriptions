# Free VPN Subscriptions

[English](./README.md) · [简体中文](./README_CN.md) · **日本語** · [한국어](./README_KO.md) · [Español](./README_ES.md) · [Português](./README_PT.md) · [Русский](./README_RU.md)

<p align="center"><img src="https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/assets/hero.png" alt="Free VPN Subscriptions — hourly-refreshed free VPN subscriptions for Clash, sing-box, v2ray" width="780"></p>

![ノード](https://img.shields.io/badge/ノード-150-brightgreen) ![生存](https://img.shields.io/badge/生存-2599-blue) ![中央値--rtt](https://img.shields.io/badge/中央値--rtt-9ms-orange) ![更新](https://img.shields.io/badge/更新-2026-04-19_11:06_UTC-informational)

> **動作する無料 VPN を手に入れる一番かんたんな方法 —— 購読リンクをコピーしてクライアントに貼るだけ。**  
> 登録不要。支払い不要。バイナリのインストール不要。公開ソースから毎時自動更新、公開前に全ノードを TCP + TLS で検証。

> 無料VPN 無制限 · 無料 v2ray 購読 · 無料 Clash 購読 · 無料 sing-box 購読 · VLESS · Reality · VMess · Trojan · Shadowsocks · Hysteria2 · 毎時更新 · TCP+TLS プローブ済み · 国別

## 💡 なぜこのプロジェクト?

GitHub 上の「無料 VPN」リストの多くは古いデータ、死んだノードだらけ、あるいは怪しいバイナリのインストールを要求します。このリポジトリは**数分前に TCP ハンドシェイクと TLS ハンドシェイクの両方を通過したノードのみ**を、厳選された公開ソースから、レイテンシ順に発行します。Clash / sing-box / v2rayN にそのまま貼れる 3 種類の購読ファイルが手に入ります。

> 📖 How the fetch → probe → rank pipeline works: [ARCHITECTURE.md](./ARCHITECTURE.md)

## 🔬 ノードが本当に使えるかどうか、どう検証しているか

**正直に言うと: ノードが確実にトラフィックを通すことを *保証* することはできません。** 実際にトラフィックを流してみないかぎり、どんな集約プロジェクトにも不可能です。以下に、集約時に何を検証しているか、何を検証できないか、そして「本当の保証」がどこから来るかを明示します。

### ✅ 集約時 (公開前) に検証すること

1. **TCP 到達性** —— 各 `server:port` に TCP 接続を張ります。サーバーダウン、DNS エラー、ポート遮断はすべてドロップ。原始エントリの約 40 % を除外。
2. **TLS ハンドシェイク** —— TLS / Reality / WS-TLS ノードについて、完全なハンドシェイクを実行。証明書期限切れ、SNI 不一致、Reality short-id 失効はドロップ。さらに約 10 % を除外。
3. **レイテンシ順ソート** —— 生存ノードを RTT 順にソートし、上位 N を公開。

直近の典型値: **17 ソース → ~4,800 生データ → ~2,900 TCP 生存 → ~2,600 TLS OK → 上位 200 を公開**。

### ❌ 検証できないこと

- プロキシプロトコル認証。UUID / パスワード不一致は、TLS ハンドシェイク *後* に上流サーバーで拒否されるため、私たちには見えません。
- 実際の HTTP-over-proxy 成功。
- 帯域 / スループット。
- 出口 IP の GeoIP を超えた地理情報。

### 🛡️ 実行時の検証 —— 本当の保証はここから

公開している `clash.yaml` には `url-test` プロキシグループが組み込まれており、**クライアント側で 5 分ごとに各ノードへ実 HTTP** を投げます:

```yaml
proxy-groups:
  - name: AUTO
    type: url-test
    url: http://www.gstatic.com/generate_204
    interval: 300
```

クライアントは *実際の* HTTP-over-proxy レイテンシでソートし、最速の使えるノードを自動選択します。sing-box / v2ray にも同等の機能があります。選ばれたノードが落ちても、クライアントが自動で次のノードに切り替えます。

### 🧮 実際の期待値

公開される上位 200 ノードのうち、クライアント側で HTTP を通す実測済みノードは通常 30-50 個見つかります。遅くなったら url-test グループが次の候補に切り替え、ワンクリックで済みます。

## 🚀 ワンクリック購読

クライアントに合う URL をコピーして購読インポート欄に貼り付けてください:

| クライアント | 形式 | 購読 URL |
|---|---|---|
| Clash / Clash Verge / ClashX | `clash.yaml` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/clash.yaml` |
| sing-box | `singbox.json` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/singbox.json` |
| v2rayN / v2rayNG / Shadowrocket / NekoBox | `v2ray-base64` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/v2ray-base64.txt` |

## 🌍 国別購読

特定地域のノードだけが欲しい?専用の購読 URL を選んでください:

| 国 | ノード数 | Clash | sing-box | v2ray |
|---|---|---|---|---|
| 🇺🇸 United States (`US`) | 19 | [clash-US.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-US.yaml) | [singbox-US.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-US.json) | [v2ray-base64-US.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-US.txt) |
| 🇩🇪 Germany (`DE`) | 4 | [clash-DE.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-DE.yaml) | [singbox-DE.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-DE.json) | [v2ray-base64-DE.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-DE.txt) |

## 📖 クライアント設定ガイド

初めての方はプラットフォームに合わせてチュートリアルをどうぞ:

- [**Clash Verge**](https://au1rxx.github.io/free-vpn-subscriptions/guides/clash-verge.html) · Windows / macOS / Linux
- [**v2rayNG**](https://au1rxx.github.io/free-vpn-subscriptions/guides/v2rayng.html) · Android
- [**Shadowrocket**](https://au1rxx.github.io/free-vpn-subscriptions/guides/shadowrocket.html) · iOS / iPadOS
- [**sing-box**](https://au1rxx.github.io/free-vpn-subscriptions/guides/sing-box.html) · Windows / macOS / Linux / iOS / Android

## 🧩 対応クライアント

- **Windows**: v2rayN、Clash Verge、Hiddify、NekoRay
- **macOS**: ClashX Pro、Clash Verge、sing-box、Hiddify
- **iOS**: Shadowrocket、Stash、Loon、sing-box、Hiddify
- **Android**: v2rayNG、NekoBox、Clash Meta for Android、Hiddify、sing-box
- **Linux**: mihomo (Clash.Meta)、sing-box、v2ray-core

## 📊 リアルタイム統計

- **選定ノード**: 150
- **全ソース生存数**: 2599
- **最速 RTT**: 5 ms
- **中央値 RTT**: 9 ms
- **最終更新 (UTC)**: 2026-04-19 11:06 UTC

**プロトコル構成:** shadowsocks × 22 · trojan × 24 · vless × 75 · vmess × 29

**今回使用したソース:** `barry-far-v2ray` × 38 · `ebrasha-v2ray` × 11 · `epodonios` × 36 · `lagzian-mix` × 2 · `mahdibland-aggregator` × 2 · `mahdibland-shadowsocks` × 1 · `matin-v2ray` × 3 · `mfuu-clash` × 5 · `ninjastrikers` × 26 · `pawdroid` × 1 · `ruking-clash` × 19 · `snakem982` × 1 · `surfboard-eternity` × 3 · `vxiaov-clash` × 2

## ❓ よくある質問

<details><summary>本当に無料?</summary>

はい。全ノードはサードパーティのボランティアが運営し、自分で公開購読を出しています。当リポジトリはサーバーを持たず、既に公開されているものをテスト・順位付け・再パッケージしているだけです。

</details>

<details><summary>データはどれくらい新しい?</summary>

毎時更新 (上流を `:00` ちょうどに集中して叩かないよう小さなランダム遅延あり): すべての上流ソースを取得 → TCP+TLS プローブ → 死ノード除去 → レイテンシ順ソート → 新しい出力ファイルを公開。上のバッジの更新時刻を参照してください。

</details>

<details><summary>これらのノードは信用できる?</summary>

無料ノードは全トラフィックが運営者に見えます。**銀行・ログイン・機密情報に絶対使わないでください。**公開コンテンツのジオブロック突破には問題ありません。本当のプライバシーには自前 VPS か有料サービスを。

</details>

<details><summary>リストにあるのに繋がらないノードがあるのはなぜ?</summary>

私たちが検証するのは TCP 到達性と TLS ハンドシェイクのみ —— 帯域上限、ルーティング異常、証明書期限切れは残ります。公開する `clash.yaml` には `url-test` グループ (`http://www.gstatic.com/generate_204` に 300 秒間隔) が組み込まれており、クライアントが実際に HTTP を通せる最速ノードを自動選択します。落ちたら次へ。

</details>

## 🤝 貢献

信頼できる公開購読ソースをご存知ですか?URL と形式を添えて issue を立ててください。

## ⚠️ 免責事項

本リポジトリは第三者ボランティアが**公開共有**しているプロキシ設定を集約するだけです。サーバーは一切運営しておらず、可用性・セキュリティは保証できず、使用結果について責任は負いません。学習と個人的な接続目的に限ります。お住まいの法域のすべての法律を遵守してください。

## ⭐ スター履歴

[![Star History Chart](https://api.star-history.com/svg?repos=Au1rxx/free-vpn-subscriptions&type=Date)](https://www.star-history.com/#Au1rxx/free-vpn-subscriptions&Date)

---

役に立ったら ⭐ をお願いします —— すべてのスターが、他の人がこのプロジェクトを見つけやすくします。
