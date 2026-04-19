package readme

var PT = Locale{
	Code:        "pt",
	DisplayName: "Português",
	FileName:    "README_PT.md",
	LangAttr:    "pt",

	BadgeNodes:   "nós",
	BadgeAlive:   "ativos",
	BadgeMedian:  "rtt--mediano",
	BadgeUpdated: "atualizado",

	Hook1:       "**A forma mais fácil de obter uma VPN gratuita funcional — copie um link de assinatura, cole no seu cliente, conecte.**",
	Hook2:       "Sem cadastro. Sem pagamento. Sem instalar nenhum binário. Atualizado a cada hora a partir de fontes públicas e cada nó é testado.",
	KeywordLine: "VPN grátis · assinatura VPN gratuita · proxy grátis · Clash assinatura · v2ray assinatura · sing-box assinatura · VLESS · Reality · VMess · Trojan · Shadowsocks · Hysteria2 · atualizado por hora · TCP+TLS testado · por país",

	WhyHeading: "## 💡 Por que este projeto?",
	WhyBody:    "Cada lista de \"VPN gratuita\" no GitHub está desatualizada, cheia de nós mortos, ou pede para instalar um binário suspeito. Este repositório **publica apenas nós que passaram um handshake TCP E um handshake TLS minutos atrás**, a partir de fontes públicas selecionadas, ordenados por latência. Você recebe 3 arquivos de assinatura portáteis — use-os em Clash, sing-box ou v2rayN e pronto.",

	SubscribeHeading:   "## 🚀 Assinatura com um clique",
	SubscribeIntro:     "Copie a URL que corresponde ao seu cliente e cole no campo de importação de assinatura:",
	SubscribeColClient: "Cliente",
	SubscribeColFormat: "Formato",
	SubscribeColURL:    "URL de assinatura",

	ClientsHeading: "## 🧩 Clientes suportados",
	ClientsWindows: "**Windows**: v2rayN, Clash Verge, Hiddify, NekoRay",
	ClientsMacOS:   "**macOS**: ClashX Pro, Clash Verge, sing-box, Hiddify",
	ClientsIOS:     "**iOS**: Shadowrocket, Stash, Loon, sing-box, Hiddify",
	ClientsAndroid: "**Android**: v2rayNG, NekoBox, Clash Meta for Android, Hiddify, sing-box",
	ClientsLinux:   "**Linux**: mihomo (Clash.Meta), sing-box, v2ray-core",

	StatsHeading:     "## 📊 Estatísticas ao vivo",
	StatsNodes:       "**Nós selecionados**",
	StatsAlive:       "**Ativos em todas as fontes**",
	StatsFastest:     "**RTT do nó mais rápido**",
	StatsMedian:      "**RTT mediano**",
	StatsUpdated:     "**Última atualização (UTC)**",
	ProtocolMixLabel: "**Mix de protocolos:**",
	SourcesLabel:     "**Fontes usadas nesta execução:**",

	ByCountryHeading: "## 🌍 Por país",
	ByCountryIntro:   "Quer nós apenas em uma região específica? Use uma dessas URLs de assinatura direcionadas:",
	ByCountryColCC:   "País",
	ByCountryColN:    "Nós",

	FAQHeading: "## ❓ Perguntas frequentes",
	FAQ1Q:      "Isso é realmente grátis?",
	FAQ1A:      "Sim. Os nós são operados por voluntários de terceiros que publicam suas próprias assinaturas gratuitas. Nós não operamos nenhum servidor — apenas testamos, classificamos e reempacotamos o que já é público.",
	FAQ2Q:      "Quão atualizados são os dados?",
	FAQ2A:      "Uma GitHub Action roda a cada hora: puxa todas as fontes, faz sondagem TCP+TLS em cada nó, descarta os mortos, ordena por latência e comita os novos arquivos. Veja o carimbo `Last updated` acima.",
	FAQ3Q:      "Posso confiar nesses nós?",
	FAQ3A:      "Nós gratuitos veem todo o seu tráfego. **Nunca os use para banco, login ou algo sensível.** Bom para driblar bloqueios geográficos em conteúdo público. Use seu próprio VPS / serviço pago para privacidade real.",
	FAQ4Q:      "Por que alguns nós listados falham?",
	FAQ4A:      "Verificamos acessibilidade TCP e handshake TLS, mas um nó ainda pode ter cota esgotada, roteamento errado ou certificado expirado. Tente alguns; o grupo selector oferece alternativas.",

	ContributingHeading: "## 🤝 Contribuir",
	ContributingBody:    "Conhece uma fonte de assinatura pública confiável que deveríamos adicionar? Abra uma issue com a URL e o formato.",

	DisclaimerHeading: "## ⚠️ Aviso legal",
	DisclaimerBody:    "Este repositório agrega configurações de proxy **compartilhadas publicamente** por voluntários de terceiros. Não operamos nenhum servidor, não garantimos disponibilidade ou segurança, e não somos responsáveis pelo uso. Destinado a uso educacional e conectividade pessoal. Cumpra todas as leis aplicáveis em sua jurisdição.",

	StarHistoryHeading: "## ⭐ Histórico de estrelas",
	FinalCTA:           "Se este projeto te ajudou, deixe uma ⭐ — cada estrela facilita para outros o encontrarem.",
}
