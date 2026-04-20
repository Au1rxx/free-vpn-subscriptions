package readme

var JA = Locale{
	Code:        "ja",
	DisplayName: "日本語",
	FileName:    "README_JA.md",
	LangAttr:    "ja",

	BadgeNodes:   "ノード",
	BadgeAlive:   "生存",
	BadgeMedian:  "中央値--rtt",
	BadgeUpdated: "更新",

	Hook1:       "**動作する無料 VPN を手に入れる一番かんたんな方法 —— 購読リンクをコピーしてクライアントに貼るだけ。**",
	Hook2:       "登録不要。支払い不要。バイナリのインストール不要。公開ソースから毎時自動更新 —— 公開される全ノードは、数分前に sing-box 経由で実 HTTP トラフィックを転送した実績があります。",
	KeywordLine: "無料VPN 無制限 · 無料 v2ray 購読 · 無料 Clash 購読 · 無料 sing-box 購読 · VLESS · Reality · VMess · Trojan · Shadowsocks · Hysteria2 · 毎時更新 · HTTP 実測検証 · 国別",

	WhyHeading: "## 💡 なぜこのプロジェクト?",
	WhyBody:    "GitHub 上の「無料 VPN」リストの多くは古いデータ、死んだノードだらけ、あるいは怪しいバイナリのインストールを要求します。このリポジトリはどこよりも一歩進んでいます —— **ポートが開いているかをチェックするだけでなく、sing-box を使って実際に HTTP トラフィックをノード経由で流し、204 が返ってくることを確認してから公開します**、全て数分以内に。Clash / sing-box / v2rayN にそのまま貼れる 3 種類の購読ファイルが手に入ります。",

	VerificationHeading: "## 🔬 ノードが本当に使えるかどうか、どう検証しているか",
	VerificationBody: `多くの無料 VPN リストは「TCP ポートが開いている」だけで公開します。私たちは違います。以下が、公開前にノードが通過しなければならない完全な検証パイプラインです。

### ✅ 集約時 (公開前) に検証すること

1. **TCP 到達性** —— 各 ` + "`server:port`" + ` に TCP 接続を張る。サーバーダウン、DNS エラー、ポート遮断はすべてドロップ。原始エントリの約 40 % を除外。
2. **TLS ハンドシェイク** —— TLS / Reality / WS-TLS ノードについて完全なハンドシェイクを実行。証明書期限切れ、SNI 不一致、Reality short-id 失効はドロップ。さらに約 10 % を除外。
3. **sing-box 設定検証** —— 生存した各ノードを実際の sing-box outbound に変換し、` + "`sing-box check`" + ` で検証。破損した暗号化方式、不正な UUID、未対応の flow オプションは、プローブスロットを浪費する前にドロップ。
4. **HTTP-over-proxy 実測 (これが肝)** —— 最速の約 900 候補を sing-box サブプロセスにバッチで投入し、各ノードに専用のローカル SOCKS5 inbound を割り当てて、実際の HTTP と HTTPS リクエストを流します:
   - ` + "`http://www.gstatic.com/generate_204`" + ` (204 を期待)
   - ` + "`https://www.cloudflare.com/cdn-cgi/trace`" + ` (200 を期待)

   リクエストは実際のプロキシプロトコル (VLESS / VMess / Trojan / Shadowsocks / Hysteria2) を完全に通過するため、ここを通過するノードは認証・ルーティング・内側 TLS ハンドシェイク・出口ネットワークすべてが機能していると実証されています。
5. **2 ラウンド、45 秒間隔** —— 一度通っても 45 秒後に死ぬノードは除外。(ラウンド数 × ターゲット数) のうち成功率 50 % 以上のノードのみ残ります。
6. **実レイテンシ中央値でソート** —— 生存ノードは HTTP-over-proxy 実測往復時間の中央値 (生の TCP RTT ではなく) でソートされ、上位 N を公開。

直近の典型値: **17 ソース → ~4,800 生データ → ~2,900 TCP 生存 → ~2,600 TLS OK → ~840 設定有効 → ~280 HTTP 実測通過 → 上位 150 を公開**。公開される 150 ノードは、すべてこの 10 分以内に実際にトラフィックを転送した実績があります。

### ❌ それでも検証できないこと

- **帯域 / スループット** —— 測っているのはレイテンシで Mbps ではありません。50ms でも動画は遅いかもしれません。
- **精密な地理位置** —— GeoIP は出口 IP の国を教えてくれますが、都市や ISP レベルでは信頼できません。
- **地域特有のブロック** —— 私たちのプローブから通るノードが、あなたの環境から通るとは限りません (ISP 層のフィルタ、captive portal など)。
- **公開後の生存** —— 10 分前には生きていた、でもその後死んだかもしれません。

### 🛡️ 実行時のセーフティネット —— 上記最後の項目対策

公開している ` + "`clash.yaml`" + ` には ` + "`url-test`" + ` プロキシグループが組み込まれており、**あなたの端末上で** 5 分ごとに各ノードへ実 HTTP を再投入します:

` + "```yaml" + `
proxy-groups:
  - name: AUTO
    type: url-test
    url: http://www.gstatic.com/generate_204
    interval: 300
` + "```" + `

クライアントは *あなたのネットワーク上での* リアルタイム HTTP-over-proxy レイテンシでソートし、最速の使えるノードを自動選択します。sing-box / v2ray にも同等の機能があります。毎時集約の合間にノードが落ちても、クライアントが自動で次に切り替えます。

### 🧮 実際の期待値

毎回公開する ~150 ノードのうち、クライアント側で **80-120 ノードが HTTP を通す実測済み** で、TCP プローブだけのリストと比べて 2-3 倍の命中率です。一つ落ちても url-test グループが透過的にローテーションします。`,

	SubscribeHeading:   "## 🚀 ワンクリック購読",
	SubscribeIntro:     "クライアントに合う URL をコピーして購読インポート欄に貼り付けてください:",
	SubscribeColClient: "クライアント",
	SubscribeColFormat: "形式",
	SubscribeColURL:    "購読 URL",

	ClientsHeading: "## 🧩 対応クライアント",
	ClientsWindows: "**Windows**: v2rayN、Clash Verge、Hiddify、NekoRay",
	ClientsMacOS:   "**macOS**: ClashX Pro、Clash Verge、sing-box、Hiddify",
	ClientsIOS:     "**iOS**: Shadowrocket、Stash、Loon、sing-box、Hiddify",
	ClientsAndroid: "**Android**: v2rayNG、NekoBox、Clash Meta for Android、Hiddify、sing-box",
	ClientsLinux:   "**Linux**: mihomo (Clash.Meta)、sing-box、v2ray-core",

	StatsHeading:     "## 📊 リアルタイム統計",
	StatsNodes:       "**選定ノード**",
	StatsAlive:       "**全ソース生存数**",
	StatsFastest:     "**最速 RTT**",
	StatsMedian:      "**中央値 RTT**",
	StatsUpdated:     "**最終更新 (UTC)**",
	ProtocolMixLabel: "**プロトコル構成:**",
	SourcesLabel:     "**今回使用したソース:**",

	ByCountryHeading: "## 🌍 国別購読",
	ByCountryIntro:   "特定地域のノードだけが欲しい?専用の購読 URL を選んでください:",
	ByCountryColCC:   "国",
	ByCountryColN:    "ノード数",

	GuidesHeading:     "## 📖 クライアント設定ガイド",
	GuidesIntro:       "初めての方はプラットフォームに合わせてチュートリアルをどうぞ:",
	GuideLocaleSuffix: "",

	FAQHeading: "## ❓ よくある質問",
	FAQ1Q:      "本当に無料?",
	FAQ1A:      "はい。全ノードはサードパーティのボランティアが運営し、自分で公開購読を出しています。当リポジトリはサーバーを持たず、既に公開されているものをテスト・順位付け・再パッケージしているだけです。",
	FAQ2Q:      "データはどれくらい新しい?",
	FAQ2A:      "毎時更新 (上流を `:00` ちょうどに集中して叩かないよう小さなランダム遅延あり): すべてのソースを取得 → TCP → TLS → sing-box 設定検証 → HTTP-over-proxy 実測 (2 ラウンド、45 秒間隔) → 実 HTTP レイテンシ順ソート → 新しい出力ファイルを公開。完全パイプラインで約 10 分。上のバッジの更新時刻を参照してください。",
	FAQ3Q:      "これらのノードは信用できる?",
	FAQ3A:      "無料ノードは全トラフィックが運営者に見えます。**銀行・ログイン・機密情報に絶対使わないでください。**公開コンテンツのジオブロック突破には問題ありません。本当のプライバシーには自前 VPS か有料サービスを。",
	FAQ4Q:      "リストにあるのに繋がらないノードがあるのはなぜ?",
	FAQ4A:      "HTTP-over-proxy 実測を通過した後でも、ノードは集約間に死ぬことがあります: 帯域制限、上流キーの失効、ISP が出口 IP を遮断、運営者の廃止など。公開する `clash.yaml` には `url-test` グループ (`http://www.gstatic.com/generate_204` に 300 秒間隔) が組み込まれており、クライアントが *あなたのネットワークから* 実際に HTTP を通せる最速ノードを自動選択します。落ちたら次へ。150 のうち 80-120 が随時使えるはずです。",

	ContributingHeading: "## 🤝 貢献",
	ContributingBody:    "信頼できる公開購読ソースをご存知ですか?URL と形式を添えて issue を立ててください。",

	DisclaimerHeading: "## ⚠️ 免責事項",
	DisclaimerBody:    "本リポジトリは第三者ボランティアが**公開共有**しているプロキシ設定を集約するだけです。サーバーは一切運営しておらず、可用性・セキュリティは保証できず、使用結果について責任は負いません。学習と個人的な接続目的に限ります。お住まいの法域のすべての法律を遵守してください。",

	StarHistoryHeading: "## ⭐ スター履歴",
	FinalCTA:           "役に立ったら ⭐ をお願いします —— すべてのスターが、他の人がこのプロジェクトを見つけやすくします。",
}
