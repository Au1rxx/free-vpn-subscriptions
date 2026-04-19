package readme

var KO = Locale{
	Code:        "ko",
	DisplayName: "한국어",
	FileName:    "README_KO.md",
	LangAttr:    "ko",

	BadgeNodes:   "노드",
	BadgeAlive:   "생존",
	BadgeMedian:  "중앙값--rtt",
	BadgeUpdated: "업데이트",

	Hook1:       "**작동하는 무료 VPN을 얻는 가장 쉬운 방법 —— 구독 링크를 복사하고 클라이언트에 붙여 넣고 연결하세요.**",
	Hook2:       "가입 불필요. 결제 불필요. 바이너리 설치 불필요. 공개 소스에서 매시간 자동 갱신되며 모든 노드가 검증됩니다.",
	KeywordLine: "무료 VPN 구독 · 무료 v2ray 구독 · 무료 Clash 구독 · 무료 sing-box 구독 · VLESS · Reality · VMess · Trojan · Shadowsocks · Hysteria2 · 매시간 갱신 · TCP+TLS 프로브 완료 · 국가별",

	WhyHeading: "## 💡 왜 이 프로젝트?",
	WhyBody:    "GitHub의 거의 모든 \"무료 VPN\" 목록은 데이터가 오래되었거나, 죽은 노드로 가득 차 있거나, 출처가 불분명한 바이너리 설치를 요구합니다. 이 저장소는 **몇 분 전에 TCP 핸드셰이크와 TLS 핸드셰이크를 모두 통과한 노드만** 선별된 공개 소스에서 레이턴시 순으로 게시합니다. Clash / sing-box / v2rayN에 바로 붙여 넣을 수 있는 3가지 범용 구독 파일을 제공합니다.",

	SubscribeHeading:   "## 🚀 원클릭 구독",
	SubscribeIntro:     "클라이언트에 맞는 URL을 복사하여 구독 가져오기 필드에 붙여 넣으세요:",
	SubscribeColClient: "클라이언트",
	SubscribeColFormat: "형식",
	SubscribeColURL:    "구독 URL",

	ClientsHeading: "## 🧩 지원 클라이언트",
	ClientsWindows: "**Windows**: v2rayN, Clash Verge, Hiddify, NekoRay",
	ClientsMacOS:   "**macOS**: ClashX Pro, Clash Verge, sing-box, Hiddify",
	ClientsIOS:     "**iOS**: Shadowrocket, Stash, Loon, sing-box, Hiddify",
	ClientsAndroid: "**Android**: v2rayNG, NekoBox, Clash Meta for Android, Hiddify, sing-box",
	ClientsLinux:   "**Linux**: mihomo (Clash.Meta), sing-box, v2ray-core",

	StatsHeading:     "## 📊 실시간 통계",
	StatsNodes:       "**선정된 노드**",
	StatsAlive:       "**전체 소스 생존 수**",
	StatsFastest:     "**최고 속도 RTT**",
	StatsMedian:      "**중앙값 RTT**",
	StatsUpdated:     "**최종 업데이트 (UTC)**",
	ProtocolMixLabel: "**프로토콜 분포:**",
	SourcesLabel:     "**이번 실행에 사용된 소스:**",

	ByCountryHeading: "## 🌍 국가별 구독",
	ByCountryIntro:   "특정 지역의 노드만 필요하신가요? 전용 구독 URL을 선택하세요:",
	ByCountryColCC:   "국가",
	ByCountryColN:    "노드 수",

	GuidesHeading:     "## 📖 클라이언트 설정 가이드",
	GuidesIntro:       "처음이신가요? 플랫폼에 맞는 튜토리얼을 따라 해보세요:",
	GuideLocaleSuffix: "",

	FAQHeading: "## ❓ 자주 묻는 질문",
	FAQ1Q:      "정말 무료인가요?",
	FAQ1A:      "네. 모든 노드는 제3자 자원봉사자가 운영하며 공개 구독을 스스로 게시합니다. 저희는 어떤 서버도 운영하지 않으며, 이미 공개된 것을 테스트하고 순위를 매기고 재포장할 뿐입니다.",
	FAQ2Q:      "데이터는 얼마나 신선한가요?",
	FAQ2A:      "GitHub Action이 매시간 실행됩니다: 모든 상위 소스 가져오기 → 각 노드 TCP+TLS 프로브 → 죽은 것 제거 → 레이턴시 순 정렬 → 새 출력 파일 커밋. 위의 `Last updated` 타임스탬프를 확인하세요.",
	FAQ3Q:      "이 노드들을 신뢰할 수 있나요?",
	FAQ3A:      "무료 노드는 모든 트래픽을 운영자가 볼 수 있습니다. **은행 거래, 로그인, 민감한 작업에는 절대 사용하지 마세요.** 공개 콘텐츠의 지역 제한 우회에는 적합합니다. 실제 프라이버시에는 자체 VPS/유료 서비스를 사용하세요.",
	FAQ4Q:      "목록에 있는데 작동하지 않는 노드가 있는 이유는?",
	FAQ4A:      "TCP 도달성과 TLS 핸드셰이크를 검증하지만 노드는 여전히 할당량 소진, 잘못된 라우팅, 만료된 인증서를 가질 수 있습니다. 몇 개 시도해 보세요. selector 그룹에 대체 항목이 있습니다.",

	ContributingHeading: "## 🤝 기여",
	ContributingBody:    "신뢰할 수 있는 공개 구독 소스를 알고 계신가요? URL과 형식을 포함한 이슈를 열어 주세요.",

	DisclaimerHeading: "## ⚠️ 면책 조항",
	DisclaimerBody:    "이 저장소는 제3자 자원봉사자가 **공개 공유**한 프록시 구성을 집계합니다. 저희는 어떤 서버도 운영하지 않고, 가용성이나 보안을 보장하지 않으며, 사용 방식에 대해 책임지지 않습니다. 교육 및 개인 연결 용도로만 사용하세요. 해당 관할권의 모든 법률을 준수하세요.",

	StarHistoryHeading: "## ⭐ 스타 히스토리",
	FinalCTA:           "이 프로젝트가 도움이 되셨다면 ⭐을 남겨 주세요 —— 모든 스타가 다른 사람들이 이 프로젝트를 더 쉽게 발견하도록 도와줍니다.",
}
