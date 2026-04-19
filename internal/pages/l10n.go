package pages

// pageL10n holds every user-facing string used by the index and per-country
// page templates. Field values are looked up by locale code.
type pageL10n struct {
	// <html lang=...>
	LangAttr string
	// Displayed in language switcher for this locale (shown in this locale).
	NativeName string

	// <title> / heading
	IndexTitleTpl        string // %d selected
	IndexDescriptionTpl  string // %d selected
	IndexKeywords        string
	IndexHeading         string
	IndexSubTagline      string // short pitch beneath hero

	BadgeNodes        string
	BadgeAlive        string
	BadgeMedianRTT    string
	BadgeUpdated      string
	BadgeUpdatedUnit  string // "ms" suffix if needed — leave empty normally

	OneClickHeading  string
	OneClickIntro    string
	ColClash         string // label for Clash row
	ColSing          string
	ColV2ray         string

	ByCountryHeading string
	ByCountryIntro   string
	NodesSuffix      string // "nodes" / "节点" rendered after count

	GuidesHeading string
	GuidesIntro   string

	ClientsHeading string
	ClientsWindows string
	ClientsMacOS   string
	ClientsIOS     string
	ClientsAndroid string
	ClientsLinux   string

	FAQHeading string
	FAQ1Q      string
	FAQ1A      string
	FAQ2Q      string
	FAQ2A      string
	FAQ3Q      string
	FAQ3A      string
	FAQ4Q      string
	FAQ4A      string

	StarButton       string
	FooterLicense    string // "MIT licensed. Open source on GitHub."
	FooterDisclaimer string

	// Country page strings
	CountryTitleTpl        string // %s name, %d count
	CountryDescriptionTpl  string
	CountryKeywordsTpl     string // %s name (6 insertions)
	CountryHeadingTpl      string // %s flag, %s name
	CountrySubTpl          string // %d count, %s name
	CountryBreadcrumb      string // "← All countries"
	CountryOtherHeading    string // "Other countries"
	CountrySubscribeHeadingTpl string // "Subscribe to %s nodes only"
	CountryStarHistoryNote string

	// Language switcher preamble
	LanguageLabel string // "Language:" / "语言:"
}

var pageLocales = map[string]pageL10n{
	"en": {
		LangAttr:   "en",
		NativeName: "English",

		IndexTitleTpl:       "Free VPN Subscriptions · %d nodes · hourly · Clash / sing-box / v2ray",
		IndexDescriptionTpl: "%d TCP+TLS-verified free VPN nodes from public sources. Hourly refresh. Copy a Clash, sing-box or v2ray URL and paste into your client.",
		IndexKeywords:       "free vpn, vpn subscription, clash, sing-box, v2ray, vless, reality, trojan, shadowsocks, hysteria2, proxy list, free proxy",
		IndexHeading:        "Free VPN Subscriptions",
		IndexSubTagline:     "The easiest way to get a working free VPN — copy a subscription link, paste it into your client, connect.",

		BadgeNodes:     "nodes",
		BadgeAlive:     "alive",
		BadgeMedianRTT: "median RTT",
		BadgeUpdated:   "updated",

		OneClickHeading: "🚀 One-Click Subscribe",
		OneClickIntro:   "Copy the URL that matches your client and paste it into the subscription import field:",
		ColClash:        "Clash / Clash Verge / ClashX",
		ColSing:         "sing-box",
		ColV2ray:        "v2rayN / v2rayNG / Shadowrocket / NekoBox",

		ByCountryHeading: "🌍 By Country",
		ByCountryIntro:   "Want nodes in a specific region only? Choose a targeted subscription:",
		NodesSuffix:      "nodes",

		GuidesHeading: "📖 Step-by-step Guides",
		GuidesIntro:   "New to VPN clients? Pick your platform:",

		ClientsHeading: "🧩 Supported Clients",
		ClientsWindows: "<strong>Windows</strong>: v2rayN, Clash Verge, Hiddify, NekoRay",
		ClientsMacOS:   "<strong>macOS</strong>: ClashX Pro, Clash Verge, sing-box, Hiddify",
		ClientsIOS:     "<strong>iOS</strong>: Shadowrocket, Stash, Loon, sing-box, Hiddify",
		ClientsAndroid: "<strong>Android</strong>: v2rayNG, NekoBox, Clash Meta for Android, Hiddify, sing-box",
		ClientsLinux:   "<strong>Linux</strong>: mihomo (Clash.Meta), sing-box, v2ray-core",

		FAQHeading: "❓ FAQ",
		FAQ1Q:      "Is this actually free?",
		FAQ1A:      "Yes. Nodes are operated by third-party volunteers. We don't run any servers — we just test, rank, and repackage what's already public.",
		FAQ2Q:      "How fresh is the data?",
		FAQ2A:      "A GitHub Action runs every hour: pulls all upstream sources, TCP+TLS probes every node, drops anything dead, sorts by latency, commits new output files.",
		FAQ3Q:      "Can I trust these nodes?",
		FAQ3A:      "Free nodes see all your traffic. Never use them for banking, login, or anything sensitive. Fine for bypassing geo-blocks on public content. Use your own VPS or a paid provider for real privacy.",
		FAQ4Q:      "Why do some nodes fail?",
		FAQ4A:      "We verify TCP reachability and TLS handshakes, but a node may still have an expired quota, bad routing, or revoked certs. Try a few; the selector group gives you fallbacks.",

		StarButton:       "⭐ Star on GitHub",
		FooterLicense:    "Open source on GitHub. MIT licensed.",
		FooterDisclaimer: "This project aggregates publicly-shared proxy configurations. We do not operate any servers and make no guarantees. Comply with the laws of your jurisdiction.",

		CountryTitleTpl:            "Free %s VPN Subscription · %d nodes · Clash / sing-box / v2ray",
		CountryDescriptionTpl:      "%d TCP+TLS-verified free VPN nodes in %s, refreshed hourly. Copy a Clash, sing-box, or v2ray subscription URL.",
		CountryKeywordsTpl:         "free %s vpn, %s vpn subscription, %s clash, %s v2ray, %s proxy, %s free vpn",
		CountryHeadingTpl:          "Free %s %s VPN Subscription",
		CountrySubTpl:              "%d free VPN nodes in %s, TCP+TLS verified, refreshed hourly.",
		CountryBreadcrumb:          "← All countries",
		CountryOtherHeading:        "🌍 Other countries",
		CountrySubscribeHeadingTpl: "🚀 Subscribe to %s nodes only",

		LanguageLabel: "Language:",
	},
	"zh": {
		LangAttr:   "zh-Hans",
		NativeName: "简体中文",

		IndexTitleTpl:       "免费 VPN 订阅 · %d 节点 · 每小时刷新 · Clash / sing-box / v2ray",
		IndexDescriptionTpl: "来自公共源的 %d 个 TCP+TLS 实测节点,每小时自动刷新。复制订阅链接粘贴到客户端即可,支持 Clash、sing-box、v2ray。",
		IndexKeywords:       "免费 VPN 订阅,免费机场,免费梯子,clash 订阅,v2ray 订阅,sing-box 订阅,VLESS,Reality,Trojan,Shadowsocks,Hysteria2,每小时刷新,免费节点",
		IndexHeading:        "免费 VPN 订阅",
		IndexSubTagline:     "获取可用免费 VPN 的最简单方式——复制订阅链接,粘贴到客户端,连上就用。",

		BadgeNodes:     "节点",
		BadgeAlive:     "存活",
		BadgeMedianRTT: "中位延迟",
		BadgeUpdated:   "更新",

		OneClickHeading: "🚀 一键订阅",
		OneClickIntro:   "复制对应客户端的 URL,粘贴到订阅导入框:",
		ColClash:        "Clash / Clash Verge / ClashX",
		ColSing:         "sing-box",
		ColV2ray:        "v2rayN / v2rayNG / Shadowrocket / NekoBox",

		ByCountryHeading: "🌍 按国家订阅",
		ByCountryIntro:   "只想要特定地区的节点?选一个针对性订阅链接:",
		NodesSuffix:      "节点",

		GuidesHeading: "📖 客户端图文教程",
		GuidesIntro:   "新手?按平台挑一篇跟着做:",

		ClientsHeading: "🧩 支持的客户端",
		ClientsWindows: "<strong>Windows</strong>:v2rayN、Clash Verge、Hiddify、NekoRay",
		ClientsMacOS:   "<strong>macOS</strong>:ClashX Pro、Clash Verge、sing-box、Hiddify",
		ClientsIOS:     "<strong>iOS</strong>:Shadowrocket、Stash、Loon、sing-box、Hiddify",
		ClientsAndroid: "<strong>Android</strong>:v2rayNG、NekoBox、Clash Meta for Android、Hiddify、sing-box",
		ClientsLinux:   "<strong>Linux</strong>:mihomo (Clash.Meta)、sing-box、v2ray-core",

		FAQHeading: "❓ 常见问题",
		FAQ1Q:      "真的完全免费吗?",
		FAQ1A:      "是的。所有节点由第三方志愿者运营并公开免费订阅。本项目不运营任何服务器,只做测试、排名和重新打包。",
		FAQ2Q:      "数据多新?",
		FAQ2A:      "GitHub Actions 每小时运行一次:拉取所有上游源,对每个节点做 TCP+TLS 探测,丢弃死节点,按延迟排序,提交新文件。",
		FAQ3Q:      "这些节点可以信任吗?",
		FAQ3A:      "免费节点能看到你所有流量。绝不要用来登录银行、邮箱等敏感账号。用来突破地区限制访问公开内容没问题。真正需要隐私请自建 VPS 或付费服务。",
		FAQ4Q:      "列表里有些节点连不上?",
		FAQ4A:      "我们验证 TCP 可达和 TLS 握手,但节点仍可能配额用完、路由被污染或证书到期。多试几个,selector 组自带 fallback。",

		StarButton:       "⭐ GitHub 上点 Star",
		FooterLicense:    "开源于 GitHub,MIT 许可证。",
		FooterDisclaimer: "本项目聚合公开分享的代理配置,不运营任何服务器,不作任何保证。请遵守所在司法管辖区的法律。",

		CountryTitleTpl:            "%s 免费 VPN 订阅 · %d 节点 · Clash / sing-box / v2ray",
		CountryDescriptionTpl:      "%d 个 %s 免费 VPN 节点,每小时 TCP+TLS 实测过滤。复制 Clash、sing-box 或 v2ray 订阅链接即可使用。",
		CountryKeywordsTpl:         "%s 免费 vpn,%s 免费机场,%s clash 订阅,%s v2ray,%s 免费节点,%s 梯子",
		CountryHeadingTpl:          "%s %s 免费 VPN 订阅",
		CountrySubTpl:              "%d 个 %s 节点,TCP+TLS 实测过滤,每小时刷新。",
		CountryBreadcrumb:          "← 所有国家",
		CountryOtherHeading:        "🌍 其他国家",
		CountrySubscribeHeadingTpl: "🚀 仅订阅 %s 节点",

		LanguageLabel: "语言:",
	},
	"ja": {
		LangAttr:   "ja",
		NativeName: "日本語",

		IndexTitleTpl:       "無料 VPN 購読 · %d ノード · 毎時更新 · Clash / sing-box / v2ray",
		IndexDescriptionTpl: "公開ソースから集めた %d 個の TCP+TLS 検証済み無料 VPN ノード。毎時自動更新。Clash / sing-box / v2ray の購読 URL をコピーしてクライアントに貼るだけ。",
		IndexKeywords:       "無料vpn, 無料 vpn 購読, clash 購読, v2ray 購読, sing-box 購読, VLESS, Reality, Trojan, Shadowsocks, Hysteria2, 毎時更新, 無料 プロキシ",
		IndexHeading:        "無料 VPN 購読",
		IndexSubTagline:     "動作する無料 VPN を手に入れる一番かんたんな方法 —— 購読リンクをコピーしてクライアントに貼るだけ。",

		BadgeNodes:     "ノード",
		BadgeAlive:     "生存",
		BadgeMedianRTT: "中央値 RTT",
		BadgeUpdated:   "更新",

		OneClickHeading: "🚀 ワンクリック購読",
		OneClickIntro:   "クライアントに合う URL をコピーして購読インポート欄に貼り付けてください:",
		ColClash:        "Clash / Clash Verge / ClashX",
		ColSing:         "sing-box",
		ColV2ray:        "v2rayN / v2rayNG / Shadowrocket / NekoBox",

		ByCountryHeading: "🌍 国別購読",
		ByCountryIntro:   "特定地域のノードだけ欲しい?専用の購読リンクを選んでください:",
		NodesSuffix:      "ノード",

		GuidesHeading: "📖 ステップバイステップのチュートリアル",
		GuidesIntro:   "VPN クライアント初心者の方は、プラットフォームを選んでください:",

		ClientsHeading: "🧩 対応クライアント",
		ClientsWindows: "<strong>Windows</strong>: v2rayN、Clash Verge、Hiddify、NekoRay",
		ClientsMacOS:   "<strong>macOS</strong>: ClashX Pro、Clash Verge、sing-box、Hiddify",
		ClientsIOS:     "<strong>iOS</strong>: Shadowrocket、Stash、Loon、sing-box、Hiddify",
		ClientsAndroid: "<strong>Android</strong>: v2rayNG、NekoBox、Clash Meta for Android、Hiddify、sing-box",
		ClientsLinux:   "<strong>Linux</strong>: mihomo (Clash.Meta)、sing-box、v2ray-core",

		FAQHeading: "❓ よくある質問",
		FAQ1Q:      "本当に無料?",
		FAQ1A:      "はい。全ノードはサードパーティのボランティアが運営し、自分で公開購読を出しています。当リポジトリはサーバーを持たず、既に公開されているものをテスト・順位付け・再パッケージしているだけです。",
		FAQ2Q:      "データはどれくらい新しい?",
		FAQ2A:      "GitHub Actions が毎時実行されます: すべての上流ソースを取得 → TCP+TLS プローブ → 死ノード除去 → レイテンシ順ソート → 新しい出力ファイルをコミット。",
		FAQ3Q:      "これらのノードは信用できる?",
		FAQ3A:      "無料ノードは全トラフィックが運営者に見えます。銀行・ログイン・機密情報に絶対使わないでください。公開コンテンツのジオブロック突破には問題ありません。本当のプライバシーには自前 VPS か有料サービスを。",
		FAQ4Q:      "リストにあるのに繋がらないノードがあるのはなぜ?",
		FAQ4A:      "TCP 到達性と TLS ハンドシェイクは確認していますが、帯域上限、ルーティング異常、証明書期限切れは残ります。いくつか試してください。selector グループにフォールバックがあります。",

		StarButton:       "⭐ GitHub でスター",
		FooterLicense:    "GitHub でオープンソース、MIT ライセンス。",
		FooterDisclaimer: "本プロジェクトは公開共有されているプロキシ設定を集約しているだけです。サーバーは一切運営せず、可用性・セキュリティは保証できません。お住まいの法域の法律を遵守してください。",

		CountryTitleTpl:            "%s 無料 VPN 購読 · %d ノード · Clash / sing-box / v2ray",
		CountryDescriptionTpl:      "%d 個の %s 無料 VPN ノード、毎時 TCP+TLS 検証。Clash / sing-box / v2ray の購読 URL をコピーするだけ。",
		CountryKeywordsTpl:         "%s 無料 vpn, %s vpn 購読, %s clash, %s v2ray, %s プロキシ, %s 無料",
		CountryHeadingTpl:          "%s %s 無料 VPN 購読",
		CountrySubTpl:              "%d 個の %s ノード、TCP+TLS 検証済み、毎時更新。",
		CountryBreadcrumb:          "← 国一覧",
		CountryOtherHeading:        "🌍 他の国",
		CountrySubscribeHeadingTpl: "🚀 %s のノードのみ購読",

		LanguageLabel: "言語:",
	},
	"ko": {
		LangAttr:   "ko",
		NativeName: "한국어",

		IndexTitleTpl:       "무료 VPN 구독 · %d 노드 · 매시간 갱신 · Clash / sing-box / v2ray",
		IndexDescriptionTpl: "공개 소스에서 수집한 TCP+TLS 검증 무료 VPN 노드 %d 개. 매시간 자동 갱신. Clash / sing-box / v2ray 구독 URL을 복사해 클라이언트에 붙여넣으세요.",
		IndexKeywords:       "무료 vpn, 무료 vpn 구독, clash 구독, v2ray 구독, sing-box 구독, VLESS, Reality, Trojan, Shadowsocks, Hysteria2, 매시간 갱신, 무료 프록시",
		IndexHeading:        "무료 VPN 구독",
		IndexSubTagline:     "작동하는 무료 VPN을 얻는 가장 쉬운 방법 —— 구독 링크를 복사해 클라이언트에 붙여넣고 연결하세요.",

		BadgeNodes:     "노드",
		BadgeAlive:     "생존",
		BadgeMedianRTT: "중앙값 RTT",
		BadgeUpdated:   "업데이트",

		OneClickHeading: "🚀 원클릭 구독",
		OneClickIntro:   "클라이언트에 맞는 URL을 복사하여 구독 가져오기 필드에 붙여넣으세요:",
		ColClash:        "Clash / Clash Verge / ClashX",
		ColSing:         "sing-box",
		ColV2ray:        "v2rayN / v2rayNG / Shadowrocket / NekoBox",

		ByCountryHeading: "🌍 국가별 구독",
		ByCountryIntro:   "특정 지역 노드만 원하시나요? 전용 구독 링크를 선택하세요:",
		NodesSuffix:      "노드",

		GuidesHeading: "📖 단계별 설정 가이드",
		GuidesIntro:   "VPN 클라이언트가 처음이신가요? 플랫폼을 선택하세요:",

		ClientsHeading: "🧩 지원 클라이언트",
		ClientsWindows: "<strong>Windows</strong>: v2rayN, Clash Verge, Hiddify, NekoRay",
		ClientsMacOS:   "<strong>macOS</strong>: ClashX Pro, Clash Verge, sing-box, Hiddify",
		ClientsIOS:     "<strong>iOS</strong>: Shadowrocket, Stash, Loon, sing-box, Hiddify",
		ClientsAndroid: "<strong>Android</strong>: v2rayNG, NekoBox, Clash Meta for Android, Hiddify, sing-box",
		ClientsLinux:   "<strong>Linux</strong>: mihomo (Clash.Meta), sing-box, v2ray-core",

		FAQHeading: "❓ 자주 묻는 질문",
		FAQ1Q:      "정말 무료인가요?",
		FAQ1A:      "네. 모든 노드는 제3자 자원봉사자가 운영하며 공개 구독을 스스로 게시합니다. 저희는 어떤 서버도 운영하지 않으며, 이미 공개된 것을 테스트하고 순위를 매기고 재포장할 뿐입니다.",
		FAQ2Q:      "데이터는 얼마나 신선한가요?",
		FAQ2A:      "GitHub Action이 매시간 실행됩니다: 모든 상위 소스 가져오기 → 각 노드 TCP+TLS 프로브 → 죽은 것 제거 → 레이턴시 순 정렬 → 새 출력 파일 커밋.",
		FAQ3Q:      "이 노드들을 신뢰할 수 있나요?",
		FAQ3A:      "무료 노드는 모든 트래픽을 운영자가 볼 수 있습니다. 은행 거래, 로그인, 민감한 작업에는 절대 사용하지 마세요. 공개 콘텐츠의 지역 제한 우회에는 적합합니다. 실제 프라이버시에는 자체 VPS / 유료 서비스를 사용하세요.",
		FAQ4Q:      "목록에 있는데 작동하지 않는 노드가 있는 이유는?",
		FAQ4A:      "TCP 도달성과 TLS 핸드셰이크를 검증하지만 노드는 여전히 할당량 소진, 잘못된 라우팅, 만료된 인증서를 가질 수 있습니다. 몇 개 시도해 보세요. selector 그룹에 대체 항목이 있습니다.",

		StarButton:       "⭐ GitHub에서 Star",
		FooterLicense:    "GitHub 오픈소스, MIT 라이선스.",
		FooterDisclaimer: "이 프로젝트는 공개 공유된 프록시 구성을 집계할 뿐입니다. 서버를 운영하지 않으며 어떤 보장도 하지 않습니다. 귀하의 관할권 법률을 준수하세요.",

		CountryTitleTpl:            "%s 무료 VPN 구독 · %d 노드 · Clash / sing-box / v2ray",
		CountryDescriptionTpl:      "%s의 TCP+TLS 검증 무료 VPN 노드 %d 개, 매시간 갱신. Clash / sing-box / v2ray 구독 URL을 복사하세요.",
		CountryKeywordsTpl:         "%s 무료 vpn, %s vpn 구독, %s clash, %s v2ray, %s 프록시, %s 무료",
		CountryHeadingTpl:          "%s %s 무료 VPN 구독",
		CountrySubTpl:              "%s의 무료 VPN 노드 %d 개, TCP+TLS 검증, 매시간 갱신.",
		CountryBreadcrumb:          "← 모든 국가",
		CountryOtherHeading:        "🌍 다른 국가",
		CountrySubscribeHeadingTpl: "🚀 %s 노드만 구독",

		LanguageLabel: "언어:",
	},
	"es": {
		LangAttr:   "es",
		NativeName: "Español",

		IndexTitleTpl:       "Suscripciones VPN gratuitas · %d nodos · por hora · Clash / sing-box / v2ray",
		IndexDescriptionTpl: "%d nodos VPN gratuitos verificados con TCP+TLS desde fuentes públicas. Actualizado cada hora. Copia una URL de Clash, sing-box o v2ray y pégala en tu cliente.",
		IndexKeywords:       "vpn gratis, suscripción vpn gratuita, clash, sing-box, v2ray, vless, reality, trojan, shadowsocks, hysteria2, lista de proxies, proxy gratis",
		IndexHeading:        "Suscripciones VPN gratuitas",
		IndexSubTagline:     "La forma más fácil de obtener una VPN gratuita que funciona — copia un enlace de suscripción, pégalo en tu cliente, conecta.",

		BadgeNodes:     "nodos",
		BadgeAlive:     "activos",
		BadgeMedianRTT: "RTT mediana",
		BadgeUpdated:   "actualizado",

		OneClickHeading: "🚀 Suscripción con un clic",
		OneClickIntro:   "Copia la URL que coincida con tu cliente y pégala en el campo de importación de suscripción:",
		ColClash:        "Clash / Clash Verge / ClashX",
		ColSing:         "sing-box",
		ColV2ray:        "v2rayN / v2rayNG / Shadowrocket / NekoBox",

		ByCountryHeading: "🌍 Por país",
		ByCountryIntro:   "¿Quieres nodos solo en una región específica? Elige una suscripción dirigida:",
		NodesSuffix:      "nodos",

		GuidesHeading: "📖 Guías paso a paso",
		GuidesIntro:   "¿Nuevo con los clientes VPN? Elige tu plataforma:",

		ClientsHeading: "🧩 Clientes compatibles",
		ClientsWindows: "<strong>Windows</strong>: v2rayN, Clash Verge, Hiddify, NekoRay",
		ClientsMacOS:   "<strong>macOS</strong>: ClashX Pro, Clash Verge, sing-box, Hiddify",
		ClientsIOS:     "<strong>iOS</strong>: Shadowrocket, Stash, Loon, sing-box, Hiddify",
		ClientsAndroid: "<strong>Android</strong>: v2rayNG, NekoBox, Clash Meta for Android, Hiddify, sing-box",
		ClientsLinux:   "<strong>Linux</strong>: mihomo (Clash.Meta), sing-box, v2ray-core",

		FAQHeading: "❓ Preguntas frecuentes",
		FAQ1Q:      "¿Es realmente gratis?",
		FAQ1A:      "Sí. Los nodos son operados por voluntarios externos que publican sus propias suscripciones gratuitas. Nosotros no operamos ningún servidor — solo probamos, clasificamos y reempaquetamos lo que ya es público.",
		FAQ2Q:      "¿Qué tan reciente es la información?",
		FAQ2A:      "Una GitHub Action se ejecuta cada hora: trae todas las fuentes, prueba cada nodo con TCP+TLS, elimina los muertos, ordena por latencia y comitea los archivos nuevos.",
		FAQ3Q:      "¿Puedo confiar en estos nodos?",
		FAQ3A:      "Los nodos gratis ven todo tu tráfico. Nunca los uses para banca, login o algo sensible. Bien para saltar bloqueos geográficos en contenido público. Usa tu propio VPS o un proveedor de pago para privacidad real.",
		FAQ4Q:      "¿Por qué algunos nodos listados fallan?",
		FAQ4A:      "Verificamos accesibilidad TCP y handshake TLS, pero un nodo aún puede tener cuota expirada, ruteo incorrecto o certificado caducado. Prueba varios; el grupo selector te da alternativas.",

		StarButton:       "⭐ Estrella en GitHub",
		FooterLicense:    "Código abierto en GitHub. Licencia MIT.",
		FooterDisclaimer: "Este proyecto agrega configuraciones de proxy compartidas públicamente. No operamos ningún servidor y no ofrecemos garantías. Cumple con las leyes de tu jurisdicción.",

		CountryTitleTpl:            "Suscripción VPN %s gratis · %d nodos · Clash / sing-box / v2ray",
		CountryDescriptionTpl:      "%d nodos VPN gratuitos en %s verificados con TCP+TLS, actualizados cada hora. Copia una URL de Clash, sing-box o v2ray.",
		CountryKeywordsTpl:         "vpn gratis %s, suscripción vpn %s, %s clash, %s v2ray, %s proxy, %s vpn gratis",
		CountryHeadingTpl:          "Suscripción VPN %s %s gratis",
		CountrySubTpl:              "%d nodos VPN gratuitos en %s, verificados con TCP+TLS, actualizados cada hora.",
		CountryBreadcrumb:          "← Todos los países",
		CountryOtherHeading:        "🌍 Otros países",
		CountrySubscribeHeadingTpl: "🚀 Suscribirse solo a nodos de %s",

		LanguageLabel: "Idioma:",
	},
	"pt": {
		LangAttr:   "pt",
		NativeName: "Português",

		IndexTitleTpl:       "Assinaturas VPN gratuitas · %d nós · por hora · Clash / sing-box / v2ray",
		IndexDescriptionTpl: "%d nós VPN gratuitos verificados com TCP+TLS a partir de fontes públicas. Atualizado a cada hora. Copie uma URL Clash, sing-box ou v2ray e cole no seu cliente.",
		IndexKeywords:       "vpn grátis, assinatura vpn grátis, clash, sing-box, v2ray, vless, reality, trojan, shadowsocks, hysteria2, lista de proxies, proxy grátis",
		IndexHeading:        "Assinaturas VPN gratuitas",
		IndexSubTagline:     "A forma mais fácil de obter uma VPN gratuita funcional — copie um link de assinatura, cole no seu cliente, conecte.",

		BadgeNodes:     "nós",
		BadgeAlive:     "ativos",
		BadgeMedianRTT: "RTT mediano",
		BadgeUpdated:   "atualizado",

		OneClickHeading: "🚀 Assinatura com um clique",
		OneClickIntro:   "Copie a URL que corresponde ao seu cliente e cole no campo de importação de assinatura:",
		ColClash:        "Clash / Clash Verge / ClashX",
		ColSing:         "sing-box",
		ColV2ray:        "v2rayN / v2rayNG / Shadowrocket / NekoBox",

		ByCountryHeading: "🌍 Por país",
		ByCountryIntro:   "Quer nós apenas em uma região específica? Escolha uma assinatura direcionada:",
		NodesSuffix:      "nós",

		GuidesHeading: "📖 Tutoriais passo a passo",
		GuidesIntro:   "Novo nos clientes VPN? Escolha sua plataforma:",

		ClientsHeading: "🧩 Clientes suportados",
		ClientsWindows: "<strong>Windows</strong>: v2rayN, Clash Verge, Hiddify, NekoRay",
		ClientsMacOS:   "<strong>macOS</strong>: ClashX Pro, Clash Verge, sing-box, Hiddify",
		ClientsIOS:     "<strong>iOS</strong>: Shadowrocket, Stash, Loon, sing-box, Hiddify",
		ClientsAndroid: "<strong>Android</strong>: v2rayNG, NekoBox, Clash Meta for Android, Hiddify, sing-box",
		ClientsLinux:   "<strong>Linux</strong>: mihomo (Clash.Meta), sing-box, v2ray-core",

		FAQHeading: "❓ Perguntas frequentes",
		FAQ1Q:      "Isso é realmente grátis?",
		FAQ1A:      "Sim. Os nós são operados por voluntários de terceiros que publicam suas próprias assinaturas gratuitas. Nós não operamos nenhum servidor — apenas testamos, classificamos e reempacotamos o que já é público.",
		FAQ2Q:      "Quão atualizados são os dados?",
		FAQ2A:      "Uma GitHub Action roda a cada hora: puxa todas as fontes, faz sondagem TCP+TLS em cada nó, descarta os mortos, ordena por latência e comita os novos arquivos.",
		FAQ3Q:      "Posso confiar nesses nós?",
		FAQ3A:      "Nós gratuitos veem todo o seu tráfego. Nunca os use para banco, login ou algo sensível. Bom para driblar bloqueios geográficos em conteúdo público. Use seu próprio VPS ou serviço pago para privacidade real.",
		FAQ4Q:      "Por que alguns nós listados falham?",
		FAQ4A:      "Verificamos acessibilidade TCP e handshake TLS, mas um nó ainda pode ter cota esgotada, roteamento errado ou certificado expirado. Tente alguns; o grupo selector oferece alternativas.",

		StarButton:       "⭐ Estrela no GitHub",
		FooterLicense:    "Código aberto no GitHub. Licença MIT.",
		FooterDisclaimer: "Este projeto agrega configurações de proxy compartilhadas publicamente. Não operamos nenhum servidor e não oferecemos garantias. Cumpra as leis da sua jurisdição.",

		CountryTitleTpl:            "Assinatura VPN %s grátis · %d nós · Clash / sing-box / v2ray",
		CountryDescriptionTpl:      "%d nós VPN gratuitos em %s verificados com TCP+TLS, atualizados a cada hora. Copie uma URL Clash, sing-box ou v2ray.",
		CountryKeywordsTpl:         "vpn grátis %s, assinatura vpn %s, %s clash, %s v2ray, %s proxy, %s vpn grátis",
		CountryHeadingTpl:          "Assinatura VPN %s %s grátis",
		CountrySubTpl:              "%d nós VPN gratuitos em %s, verificados com TCP+TLS, atualizados a cada hora.",
		CountryBreadcrumb:          "← Todos os países",
		CountryOtherHeading:        "🌍 Outros países",
		CountrySubscribeHeadingTpl: "🚀 Assinar apenas nós de %s",

		LanguageLabel: "Idioma:",
	},
	"ru": {
		LangAttr:   "ru",
		NativeName: "Русский",

		IndexTitleTpl:       "Бесплатные VPN подписки · %d узлов · ежечасно · Clash / sing-box / v2ray",
		IndexDescriptionTpl: "%d проверенных TCP+TLS бесплатных VPN узлов из публичных источников. Обновление каждый час. Скопируйте Clash, sing-box или v2ray URL и вставьте в клиент.",
		IndexKeywords:       "бесплатный vpn, бесплатная подписка vpn, clash, sing-box, v2ray, vless, reality, trojan, shadowsocks, hysteria2, список прокси, бесплатный прокси",
		IndexHeading:        "Бесплатные VPN подписки",
		IndexSubTagline:     "Самый простой способ получить рабочий бесплатный VPN — скопируйте ссылку подписки, вставьте в клиент, подключитесь.",

		BadgeNodes:     "узлы",
		BadgeAlive:     "живые",
		BadgeMedianRTT: "медиана RTT",
		BadgeUpdated:   "обновлено",

		OneClickHeading: "🚀 Подписка в один клик",
		OneClickIntro:   "Скопируйте URL, соответствующий вашему клиенту, и вставьте его в поле импорта подписки:",
		ColClash:        "Clash / Clash Verge / ClashX",
		ColSing:         "sing-box",
		ColV2ray:        "v2rayN / v2rayNG / Shadowrocket / NekoBox",

		ByCountryHeading: "🌍 По странам",
		ByCountryIntro:   "Нужны узлы только в определённом регионе? Выберите целевую подписку:",
		NodesSuffix:      "узлов",

		GuidesHeading: "📖 Пошаговые инструкции",
		GuidesIntro:   "Впервые настраиваете VPN-клиент? Выберите платформу:",

		ClientsHeading: "🧩 Поддерживаемые клиенты",
		ClientsWindows: "<strong>Windows</strong>: v2rayN, Clash Verge, Hiddify, NekoRay",
		ClientsMacOS:   "<strong>macOS</strong>: ClashX Pro, Clash Verge, sing-box, Hiddify",
		ClientsIOS:     "<strong>iOS</strong>: Shadowrocket, Stash, Loon, sing-box, Hiddify",
		ClientsAndroid: "<strong>Android</strong>: v2rayNG, NekoBox, Clash Meta for Android, Hiddify, sing-box",
		ClientsLinux:   "<strong>Linux</strong>: mihomo (Clash.Meta), sing-box, v2ray-core",

		FAQHeading: "❓ Часто задаваемые вопросы",
		FAQ1Q:      "Это правда бесплатно?",
		FAQ1A:      "Да. Узлы управляются сторонними волонтёрами, которые сами публикуют свои бесплатные подписки. Мы не управляем никакими серверами — только тестируем, ранжируем и переупаковываем то, что уже публично.",
		FAQ2Q:      "Насколько свежие данные?",
		FAQ2A:      "GitHub Action запускается каждый час: получает все источники, проводит TCP+TLS проверку каждого узла, отбрасывает мёртвые, сортирует по задержке и коммитит новые файлы.",
		FAQ3Q:      "Можно ли доверять этим узлам?",
		FAQ3A:      "Бесплатные узлы видят весь ваш трафик. Никогда не используйте их для банкинга, логинов или чего-то чувствительного. Подходит для обхода гео-блокировок на публичном контенте. Для реальной приватности используйте свой VPS или платный сервис.",
		FAQ4Q:      "Почему некоторые узлы из списка не работают?",
		FAQ4A:      "Мы проверяем TCP доступность и TLS handshake, но у узла всё равно могут быть исчерпанные квоты, неверная маршрутизация или просроченный сертификат. Попробуйте несколько; selector группа даёт альтернативы.",

		StarButton:       "⭐ Звезду на GitHub",
		FooterLicense:    "Открытый исходный код на GitHub. Лицензия MIT.",
		FooterDisclaimer: "Этот проект агрегирует публично доступные конфигурации прокси. Мы не управляем серверами и не даём гарантий. Соблюдайте законы вашей юрисдикции.",

		CountryTitleTpl:            "Бесплатная VPN подписка %s · %d узлов · Clash / sing-box / v2ray",
		CountryDescriptionTpl:      "%d бесплатных VPN узлов в %s, проверены TCP+TLS, обновление каждый час. Скопируйте Clash, sing-box или v2ray URL подписки.",
		CountryKeywordsTpl:         "бесплатный vpn %s, подписка vpn %s, %s clash, %s v2ray, %s прокси, %s бесплатный vpn",
		CountryHeadingTpl:          "Бесплатная VPN подписка %s %s",
		CountrySubTpl:              "%d бесплатных VPN узлов в %s, проверены TCP+TLS, обновление каждый час.",
		CountryBreadcrumb:          "← Все страны",
		CountryOtherHeading:        "🌍 Другие страны",
		CountrySubscribeHeadingTpl: "🚀 Подписаться только на узлы %s",

		LanguageLabel: "Язык:",
	},
}
