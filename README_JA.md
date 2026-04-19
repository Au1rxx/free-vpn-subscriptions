# Free VPN Subscriptions

[English](./README.md) · [简体中文](./README_CN.md) · **日本語** · [한국어](./README_KO.md) · [Español](./README_ES.md) · [Português](./README_PT.md) · [Русский](./README_RU.md)

<p align="center"><img src="https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/assets/hero.png" alt="Free VPN Subscriptions — hourly-refreshed free VPN subscriptions for Clash, sing-box, v2ray" width="780"></p>

![ノード](https://img.shields.io/badge/ノード-150-brightgreen) ![生存](https://img.shields.io/badge/生存-393-blue) ![中央値--rtt](https://img.shields.io/badge/中央値--rtt-107ms-orange) ![更新](https://img.shields.io/badge/更新-2026-04-19_08:02_UTC-informational)

> **動作する無料 VPN を手に入れる一番かんたんな方法 —— 購読リンクをコピーしてクライアントに貼るだけ。**  
> 登録不要。支払い不要。バイナリのインストール不要。公開ソースから毎時自動更新、各ノード疎通テスト済み。

> 無料VPN 無制限 · 無料 v2ray 購読 · 無料 Clash 購読 · 無料 sing-box 購読 · VLESS · Reality · VMess · Trojan · Shadowsocks · Hysteria2 · 毎時更新 · TCP+TLS プローブ済み · 国別

## 💡 なぜこのプロジェクト?

GitHub 上の「無料 VPN」リストの多くは古いデータ、死んだノードだらけ、あるいは怪しいバイナリのインストールを要求します。このリポジトリは**数分前に TCP ハンドシェイクと TLS ハンドシェイクの両方を通過したノードのみ**を、厳選された公開ソースから、レイテンシ順に発行します。Clash / sing-box / v2rayN にそのまま貼れる 3 種類の購読ファイルが手に入ります。

> 📖 How the fetch → probe → rank pipeline works: [ARCHITECTURE.md](./ARCHITECTURE.md)

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
| 🇬🇧 United Kingdom (`GB`) | 26 | [clash-GB.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-GB.yaml) | [singbox-GB.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-GB.json) | [v2ray-base64-GB.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-GB.txt) |
| 🇺🇸 United States (`US`) | 26 | [clash-US.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-US.yaml) | [singbox-US.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-US.json) | [v2ray-base64-US.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-US.txt) |
| 🇳🇱 Netherlands (`NL`) | 14 | [clash-NL.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-NL.yaml) | [singbox-NL.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-NL.json) | [v2ray-base64-NL.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-NL.txt) |
| 🇨🇦 Canada (`CA`) | 11 | [clash-CA.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-CA.yaml) | [singbox-CA.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-CA.json) | [v2ray-base64-CA.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-CA.txt) |
| 🇯🇵 Japan (`JP`) | 8 | [clash-JP.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-JP.yaml) | [singbox-JP.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-JP.json) | [v2ray-base64-JP.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-JP.txt) |
| 🇩🇪 Germany (`DE`) | 6 | [clash-DE.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-DE.yaml) | [singbox-DE.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-DE.json) | [v2ray-base64-DE.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-DE.txt) |
| 🇹🇼 Taiwan (`TW`) | 5 | [clash-TW.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-TW.yaml) | [singbox-TW.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-TW.json) | [v2ray-base64-TW.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-TW.txt) |
| 🇰🇷 Korea (`KR`) | 4 | [clash-KR.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-KR.yaml) | [singbox-KR.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-KR.json) | [v2ray-base64-KR.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-KR.txt) |
| 🇲🇦 MA (`MA`) | 3 | [clash-MA.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-MA.yaml) | [singbox-MA.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-MA.json) | [v2ray-base64-MA.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-MA.txt) |

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
- **全ソース生存数**: 393
- **最速 RTT**: 1 ms
- **中央値 RTT**: 107 ms
- **最終更新 (UTC)**: 2026-04-19 08:02 UTC

**プロトコル構成:** shadowsocks × 106 · trojan × 8 · vmess × 36

**今回使用したソース:** `freefq` × 2 · `mahdibland-aggregator` × 72 · `mahdibland-shadowsocks` × 59 · `pawdroid` × 3 · `vxiaov-clash` × 14

## ❓ よくある質問

<details><summary>本当に無料?</summary>

はい。全ノードはサードパーティのボランティアが運営し、自分で公開購読を出しています。当リポジトリはサーバーを持たず、既に公開されているものをテスト・順位付け・再パッケージしているだけです。

</details>

<details><summary>データはどれくらい新しい?</summary>

GitHub Actions が毎時実行されます: すべての上流ソースを取得 → TCP+TLS プローブ → 死ノード除去 → レイテンシ順ソート → 新しい出力ファイルをコミット。上のバッジの更新時刻を参照してください。

</details>

<details><summary>これらのノードは信用できる?</summary>

無料ノードは全トラフィックが運営者に見えます。**銀行・ログイン・機密情報に絶対使わないでください。**公開コンテンツのジオブロック突破には問題ありません。本当のプライバシーには自前 VPS か有料サービスを。

</details>

<details><summary>リストにあるのに繋がらないノードがあるのはなぜ?</summary>

TCP 到達性と TLS ハンドシェイクは確認していますが、帯域上限、ルーティング異常、証明書期限切れなど実際に繋がらない要因は残ります。数個試してください。selector グループにフォールバックがあります。

</details>

## 🤝 貢献

信頼できる公開購読ソースをご存知ですか?URL と形式を添えて issue を立ててください。

## ⚠️ 免責事項

本リポジトリは第三者ボランティアが**公開共有**しているプロキシ設定を集約するだけです。サーバーは一切運営しておらず、可用性・セキュリティは保証できず、使用結果について責任は負いません。学習と個人的な接続目的に限ります。お住まいの法域のすべての法律を遵守してください。

## ⭐ スター履歴

[![Star History Chart](https://api.star-history.com/svg?repos=Au1rxx/free-vpn-subscriptions&type=Date)](https://www.star-history.com/#Au1rxx/free-vpn-subscriptions&Date)

---

役に立ったら ⭐ をお願いします —— すべてのスターが、他の人がこのプロジェクトを見つけやすくします。
