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
	Hook2:       "登録不要。支払い不要。バイナリのインストール不要。公開ソースから毎時自動更新、各ノード疎通テスト済み。",
	KeywordLine: "無料VPN 無制限 · 無料 v2ray 購読 · 無料 Clash 購読 · 無料 sing-box 購読 · VLESS · Reality · VMess · Trojan · Shadowsocks · Hysteria2 · 毎時更新 · TCP+TLS プローブ済み · 国別",

	WhyHeading: "## 💡 なぜこのプロジェクト?",
	WhyBody:    "GitHub 上の「無料 VPN」リストの多くは古いデータ、死んだノードだらけ、あるいは怪しいバイナリのインストールを要求します。このリポジトリは**数分前に TCP ハンドシェイクと TLS ハンドシェイクの両方を通過したノードのみ**を、厳選された公開ソースから、レイテンシ順に発行します。Clash / sing-box / v2rayN にそのまま貼れる 3 種類の購読ファイルが手に入ります。",

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

	FAQHeading: "## ❓ よくある質問",
	FAQ1Q:      "本当に無料?",
	FAQ1A:      "はい。全ノードはサードパーティのボランティアが運営し、自分で公開購読を出しています。当リポジトリはサーバーを持たず、既に公開されているものをテスト・順位付け・再パッケージしているだけです。",
	FAQ2Q:      "データはどれくらい新しい?",
	FAQ2A:      "GitHub Actions が毎時実行されます: すべての上流ソースを取得 → TCP+TLS プローブ → 死ノード除去 → レイテンシ順ソート → 新しい出力ファイルをコミット。上のバッジの更新時刻を参照してください。",
	FAQ3Q:      "これらのノードは信用できる?",
	FAQ3A:      "無料ノードは全トラフィックが運営者に見えます。**銀行・ログイン・機密情報に絶対使わないでください。**公開コンテンツのジオブロック突破には問題ありません。本当のプライバシーには自前 VPS か有料サービスを。",
	FAQ4Q:      "リストにあるのに繋がらないノードがあるのはなぜ?",
	FAQ4A:      "TCP 到達性と TLS ハンドシェイクは確認していますが、帯域上限、ルーティング異常、証明書期限切れなど実際に繋がらない要因は残ります。数個試してください。selector グループにフォールバックがあります。",

	ContributingHeading: "## 🤝 貢献",
	ContributingBody:    "信頼できる公開購読ソースをご存知ですか?URL と形式を添えて issue を立ててください。",

	DisclaimerHeading: "## ⚠️ 免責事項",
	DisclaimerBody:    "本リポジトリは第三者ボランティアが**公開共有**しているプロキシ設定を集約するだけです。サーバーは一切運営しておらず、可用性・セキュリティは保証できず、使用結果について責任は負いません。学習と個人的な接続目的に限ります。お住まいの法域のすべての法律を遵守してください。",

	StarHistoryHeading: "## ⭐ スター履歴",
	FinalCTA:           "役に立ったら ⭐ をお願いします —— すべてのスターが、他の人がこのプロジェクトを見つけやすくします。",
}
