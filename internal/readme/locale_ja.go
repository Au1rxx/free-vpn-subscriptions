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
	Hook2:       "登録不要。支払い不要。バイナリのインストール不要。公開ソースから毎時自動更新、公開前に全ノードを TCP + TLS で検証。",
	KeywordLine: "無料VPN 無制限 · 無料 v2ray 購読 · 無料 Clash 購読 · 無料 sing-box 購読 · VLESS · Reality · VMess · Trojan · Shadowsocks · Hysteria2 · 毎時更新 · TCP+TLS プローブ済み · 国別",

	WhyHeading: "## 💡 なぜこのプロジェクト?",
	WhyBody:    "GitHub 上の「無料 VPN」リストの多くは古いデータ、死んだノードだらけ、あるいは怪しいバイナリのインストールを要求します。このリポジトリは**数分前に TCP ハンドシェイクと TLS ハンドシェイクの両方を通過したノードのみ**を、厳選された公開ソースから、レイテンシ順に発行します。Clash / sing-box / v2rayN にそのまま貼れる 3 種類の購読ファイルが手に入ります。",

	VerificationHeading: "## 🔬 ノードが本当に使えるかどうか、どう検証しているか",
	VerificationBody: `**正直に言うと: ノードが確実にトラフィックを通すことを *保証* することはできません。** 実際にトラフィックを流してみないかぎり、どんな集約プロジェクトにも不可能です。以下に、集約時に何を検証しているか、何を検証できないか、そして「本当の保証」がどこから来るかを明示します。

### ✅ 集約時 (公開前) に検証すること

1. **TCP 到達性** —— 各 ` + "`server:port`" + ` に TCP 接続を張ります。サーバーダウン、DNS エラー、ポート遮断はすべてドロップ。原始エントリの約 40 % を除外。
2. **TLS ハンドシェイク** —— TLS / Reality / WS-TLS ノードについて、完全なハンドシェイクを実行。証明書期限切れ、SNI 不一致、Reality short-id 失効はドロップ。さらに約 10 % を除外。
3. **レイテンシ順ソート** —— 生存ノードを RTT 順にソートし、上位 N を公開。

直近の典型値: **17 ソース → ~4,800 生データ → ~2,900 TCP 生存 → ~2,600 TLS OK → 上位 200 を公開**。

### ❌ 検証できないこと

- プロキシプロトコル認証。UUID / パスワード不一致は、TLS ハンドシェイク *後* に上流サーバーで拒否されるため、私たちには見えません。
- 実際の HTTP-over-proxy 成功。
- 帯域 / スループット。
- 出口 IP の GeoIP を超えた地理情報。

### 🛡️ 実行時の検証 —— 本当の保証はここから

公開している ` + "`clash.yaml`" + ` には ` + "`url-test`" + ` プロキシグループが組み込まれており、**クライアント側で 5 分ごとに各ノードへ実 HTTP** を投げます:

` + "```yaml" + `
proxy-groups:
  - name: AUTO
    type: url-test
    url: http://www.gstatic.com/generate_204
    interval: 300
` + "```" + `

クライアントは *実際の* HTTP-over-proxy レイテンシでソートし、最速の使えるノードを自動選択します。sing-box / v2ray にも同等の機能があります。選ばれたノードが落ちても、クライアントが自動で次のノードに切り替えます。

### 🧮 実際の期待値

公開される上位 200 ノードのうち、クライアント側で HTTP を通す実測済みノードは通常 30-50 個見つかります。遅くなったら url-test グループが次の候補に切り替え、ワンクリックで済みます。`,

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
	FAQ2A:      "毎時更新 (上流を `:00` ちょうどに集中して叩かないよう小さなランダム遅延あり): すべての上流ソースを取得 → TCP+TLS プローブ → 死ノード除去 → レイテンシ順ソート → 新しい出力ファイルを公開。上のバッジの更新時刻を参照してください。",
	FAQ3Q:      "これらのノードは信用できる?",
	FAQ3A:      "無料ノードは全トラフィックが運営者に見えます。**銀行・ログイン・機密情報に絶対使わないでください。**公開コンテンツのジオブロック突破には問題ありません。本当のプライバシーには自前 VPS か有料サービスを。",
	FAQ4Q:      "リストにあるのに繋がらないノードがあるのはなぜ?",
	FAQ4A:      "私たちが検証するのは TCP 到達性と TLS ハンドシェイクのみ —— 帯域上限、ルーティング異常、証明書期限切れは残ります。公開する `clash.yaml` には `url-test` グループ (`http://www.gstatic.com/generate_204` に 300 秒間隔) が組み込まれており、クライアントが実際に HTTP を通せる最速ノードを自動選択します。落ちたら次へ。",

	ContributingHeading: "## 🤝 貢献",
	ContributingBody:    "信頼できる公開購読ソースをご存知ですか?URL と形式を添えて issue を立ててください。",

	DisclaimerHeading: "## ⚠️ 免責事項",
	DisclaimerBody:    "本リポジトリは第三者ボランティアが**公開共有**しているプロキシ設定を集約するだけです。サーバーは一切運営しておらず、可用性・セキュリティは保証できず、使用結果について責任は負いません。学習と個人的な接続目的に限ります。お住まいの法域のすべての法律を遵守してください。",

	StarHistoryHeading: "## ⭐ スター履歴",
	FinalCTA:           "役に立ったら ⭐ をお願いします —— すべてのスターが、他の人がこのプロジェクトを見つけやすくします。",
}
