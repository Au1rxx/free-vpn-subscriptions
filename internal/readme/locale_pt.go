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
	Hook2:       "Sem cadastro. Sem pagamento. Sem instalar nenhum binário. Atualizado a cada hora a partir de fontes públicas — cada nó é sondado via TCP + TLS antes da publicação.",
	KeywordLine: "VPN grátis · assinatura VPN gratuita · proxy grátis · Clash assinatura · v2ray assinatura · sing-box assinatura · VLESS · Reality · VMess · Trojan · Shadowsocks · Hysteria2 · atualizado por hora · TCP+TLS testado · por país",

	WhyHeading: "## 💡 Por que este projeto?",
	WhyBody:    "Cada lista de \"VPN gratuita\" no GitHub está desatualizada, cheia de nós mortos, ou pede para instalar um binário suspeito. Este repositório **publica apenas nós que passaram um handshake TCP E um handshake TLS minutos atrás**, a partir de fontes públicas selecionadas, ordenados por latência. Você recebe 3 arquivos de assinatura portáteis — use-os em Clash, sing-box ou v2rayN e pronto.",

	VerificationHeading: "## 🔬 Como verificamos que os nós realmente funcionam",
	VerificationBody: `**Resposta honesta primeiro: não podemos *garantir* que um nó passe o seu tráfego.** Nenhum agregador pode, sem realmente enviar tráfego através dele. Aqui está exatamente o que verificamos, o que não podemos, e de onde vem a garantia real.

### ✅ O que verificamos na agregação (antes de publicar)

1. **Acessibilidade TCP** — abrimos uma conexão TCP para cada ` + "`server:port`" + `. Hosts mortos, DNS errado, portas bloqueadas são descartados. Elimina cerca de 40 % das entradas cruas.
2. **Handshake TLS** — para cada nó TLS / Reality / WS-TLS completamos o handshake inteiro. Certificados expirados, SNI incompatíveis e short-ids Reality quebrados são descartados. Elimina mais ~10 %.
3. **Ordenação por latência** — os sobreviventes são ordenados por RTT e mantemos os N mais rápidos.

Números típicos de uma execução recente: **17 fontes → ~4,800 brutos → ~2,900 vivos via TCP → ~2,600 TLS OK → top 200 publicados**.

### ❌ O que não podemos verificar

- Autenticação do protocolo proxy. UUID / senha errados só são rejeitados *depois* do handshake TLS pelo servidor upstream.
- Sucesso real de HTTP via proxy.
- Largura de banda ou throughput.
- Geolocalização além do que o GeoIP diz sobre o IP de saída.

### 🛡️ Verificação em tempo de execução — é daqui que vem a garantia real

O ` + "`clash.yaml`" + ` que publicamos inclui um grupo ` + "`url-test`" + ` que **testa HTTP real através de cada nó** a cada 5 minutos:

` + "```yaml" + `
proxy-groups:
  - name: AUTO
    type: url-test
    url: http://www.gstatic.com/generate_204
    interval: 300
` + "```" + `

Seu cliente ordena a lista de nós pela latência *real* de HTTP via proxy e auto-seleciona o nó mais rápido que funciona. sing-box e v2ray têm mecanismos equivalentes. Se um nó selecionado morrer, o cliente muda para o próximo sem intervenção.

### 🧮 Resultado esperado

Dos top 200 publicados por execução, um cliente típico encontra 30-50 que servem HTTP limpo em qualquer momento. Rotacione se um ficar lento — o grupo url-test faz isso com um clique.`,

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

	GuidesHeading:     "## 📖 Tutoriais passo a passo",
	GuidesIntro:       "Novo nos clientes VPN? Escolha sua plataforma e siga o tutorial:",
	GuideLocaleSuffix: "",

	FAQHeading: "## ❓ Perguntas frequentes",
	FAQ1Q:      "Isso é realmente grátis?",
	FAQ1A:      "Sim. Os nós são operados por voluntários de terceiros que publicam suas próprias assinaturas gratuitas. Nós não operamos nenhum servidor — apenas testamos, classificamos e reempacotamos o que já é público.",
	FAQ2Q:      "Quão atualizados são os dados?",
	FAQ2A:      "A cada hora (com um pequeno atraso aleatório para evitar bater nas fontes upstream exatamente em `:00`): puxa todas as fontes, faz sondagem TCP+TLS em cada nó, descarta os mortos, ordena por latência e publica os novos arquivos. Veja o carimbo `Last updated` acima.",
	FAQ3Q:      "Posso confiar nesses nós?",
	FAQ3A:      "Nós gratuitos veem todo o seu tráfego. **Nunca os use para banco, login ou algo sensível.** Bom para driblar bloqueios geográficos em conteúdo público. Use seu próprio VPS / serviço pago para privacidade real.",
	FAQ4Q:      "Por que alguns nós listados falham?",
	FAQ4A:      "Verificamos apenas acessibilidade TCP e handshake TLS — um nó ainda pode ter cota esgotada, roteamento errado ou certificado expirado. O `clash.yaml` publicado inclui um grupo `url-test` (`http://www.gstatic.com/generate_204`, intervalo de 300 s); seu cliente auto-seleciona o nó mais rápido que realmente serve HTTP. Se um morrer, pegue o próximo.",

	ContributingHeading: "## 🤝 Contribuir",
	ContributingBody:    "Conhece uma fonte de assinatura pública confiável que deveríamos adicionar? Abra uma issue com a URL e o formato.",

	DisclaimerHeading: "## ⚠️ Aviso legal",
	DisclaimerBody:    "Este repositório agrega configurações de proxy **compartilhadas publicamente** por voluntários de terceiros. Não operamos nenhum servidor, não garantimos disponibilidade ou segurança, e não somos responsáveis pelo uso. Destinado a uso educacional e conectividade pessoal. Cumpra todas as leis aplicáveis em sua jurisdição.",

	StarHistoryHeading: "## ⭐ Histórico de estrelas",
	FinalCTA:           "Se este projeto te ajudou, deixe uma ⭐ — cada estrela facilita para outros o encontrarem.",
}
