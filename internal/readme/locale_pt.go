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
	Hook2:       "Sem cadastro. Sem pagamento. Sem instalar nenhum binário. Atualizado a cada hora a partir de fontes públicas — cada nó publicado encaminhou tráfego HTTP real através do sing-box minutos atrás.",
	KeywordLine: "VPN grátis · assinatura VPN gratuita · proxy grátis · Clash assinatura · v2ray assinatura · sing-box assinatura · VLESS · Reality · VMess · Trojan · Shadowsocks · Hysteria2 · atualizado por hora · HTTP verificado sobre proxy · por país",

	WhyHeading: "## 💡 Por que este projeto?",
	WhyBody:    "Cada lista de \"VPN gratuita\" no GitHub está desatualizada, cheia de nós mortos, ou pede para instalar um binário suspeito. Este repositório vai um passo além de qualquer outro — **não apenas verificamos que o nó responde, mas empurramos tráfego HTTP real através dele com sing-box e confirmamos que um 204 retorna**, tudo em minutos antes de publicar. Você recebe 3 arquivos de assinatura portáteis — use-os em Clash, sing-box ou v2rayN e pronto.",

	VerificationHeading: "## 🔬 Como verificamos que os nós realmente funcionam",
	VerificationBody: `A maioria das listas de VPN gratuita para em \"a porta TCP está aberta\" e publica. Nós não. Aqui está a pipeline completa que um nó precisa passar antes de entrar na assinatura.

### ✅ O que verificamos na agregação (antes de publicar)

1. **Acessibilidade TCP** — abrimos uma conexão TCP para cada ` + "`server:port`" + `. Hosts mortos, DNS errado, portas bloqueadas são descartados. ~40 % das entradas cruas caem aqui.
2. **Handshake TLS** — para cada nó TLS / Reality / WS-TLS completamos o handshake inteiro. Certificados expirados, SNI incompatíveis e short-ids Reality quebrados são descartados. Mais ~10 % caem aqui.
3. **Validação de configuração sing-box** — cada nó sobrevivente é traduzido em um outbound real de sing-box e passa pelo ` + "`sing-box check`" + `. Cifras corrompidas, UUIDs errados e opções flow não suportadas são descartados antes de desperdiçar um slot de sondagem.
4. **Sondagem HTTP-over-proxy (esta é a chave)** — agrupamos os ~900 candidatos mais rápidos em subprocessos sing-box, cada nó recebendo seu próprio inbound SOCKS5 local, e então enviamos GETs HTTP e HTTPS reais através dele:
   - ` + "`http://www.gstatic.com/generate_204`" + ` (espera 204)
   - ` + "`https://www.cloudflare.com/cdn-cgi/trace`" + ` (espera 200)

   A requisição atravessa o protocolo proxy real (VLESS / VMess / Trojan / Shadowsocks / Hysteria2), então um nó que passa tem demonstravelmente autenticação, roteamento, handshake TLS interno e rede de saída funcionando.
5. **Duas rodadas, 45 segundos de intervalo** — nós que passam uma vez mas morrem 45 segundos depois são filtrados. Apenas nós com ≥ 50 % de taxa de sucesso em (rodadas × alvos) ficam.
6. **Ordenar por mediana de latência real** — os sobreviventes são ordenados pela mediana do ida-e-volta HTTP-over-proxy (não RTT TCP bruto) e os top N são publicados.

Números típicos de uma execução recente: **17 fontes → ~4,800 brutos → ~2,900 vivos via TCP → ~2,600 TLS OK → ~840 configuração válida → ~280 verificados por HTTP → top 150 publicados**. Cada um dos 150 de fato encaminhou tráfego nos últimos dez minutos.

### ❌ O que ainda não podemos verificar

- **Largura de banda / throughput** — medimos latência, não megabits. Um nó de 50 ms ainda pode ser lento para vídeo.
- **Precisão de geolocalização** — GeoIP diz o país do IP de saída mas não a cidade ou ISP de forma confiável.
- **Bloqueios específicos por região** — um nó que funciona da nossa infraestrutura de sondagem pode estar bloqueado da sua (filtragem no nível do ISP, captive portals, etc.).
- **Continuar vivo depois da execução** — o nó passou dez minutos atrás; pode ter morrido desde então.

### 🛡️ Rede de segurança em tempo de execução — para o último item acima

O ` + "`clash.yaml`" + ` que publicamos inclui um grupo ` + "`url-test`" + ` que retesta HTTP real através de cada nó a cada 5 minutos no *seu* dispositivo:

` + "```yaml" + `
proxy-groups:
  - name: AUTO
    type: url-test
    url: http://www.gstatic.com/generate_204
    interval: 300
` + "```" + `

Seu cliente mantém a lista de nós ordenada por latência *ao vivo* de HTTP-over-proxy da sua rede e auto-seleciona o nó mais rápido que funciona. sing-box e v2ray têm mecanismos equivalentes. Se um nó selecionado morrer entre agregações horárias, o cliente muda para o próximo sem intervenção.

### 🧮 O que isso significa na prática

Dos ~150 publicados por execução, um cliente típico encontra **80-120 nós que servem HTTP limpo da sua rede** em qualquer momento — aproximadamente 2-3× a taxa de acerto de listas que só fazem sondagem TCP. O grupo url-test rotaciona de forma transparente se um cair.`,

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
	FAQ2A:      "A cada hora (com um pequeno atraso aleatório para evitar bater nas fontes upstream exatamente em `:00`): puxa todas as fontes → TCP → TLS → validação de configuração sing-box → sondagem HTTP-over-proxy (duas rodadas, 45 s de intervalo) → ordena por latência HTTP real → publica os novos arquivos. Pipeline completo leva ~10 minutos. Veja o carimbo `Last updated` acima.",
	FAQ3Q:      "Posso confiar nesses nós?",
	FAQ3A:      "Nós gratuitos veem todo o seu tráfego. **Nunca os use para banco, login ou algo sensível.** Bom para driblar bloqueios geográficos em conteúdo público. Use seu próprio VPS / serviço pago para privacidade real.",
	FAQ4Q:      "Por que alguns nós listados falham?",
	FAQ4A:      "Mesmo após nossa sondagem HTTP-over-proxy, os nós podem morrer entre agregações: cota esgotada, upstream revogou a chave, seu ISP bloqueia o IP de saída, ou o operador desligou. O `clash.yaml` publicado inclui um grupo `url-test` (`http://www.gstatic.com/generate_204`, intervalo de 300 s); seu cliente auto-seleciona o nó mais rápido que realmente serve HTTP *da sua rede*. Se um morrer, pegue o próximo. Espere que 80-120 dos 150 funcionem em qualquer momento.",

	ContributingHeading: "## 🤝 Contribuir",
	ContributingBody:    "Conhece uma fonte de assinatura pública confiável que deveríamos adicionar? Abra uma issue com a URL e o formato.",

	DisclaimerHeading: "## ⚠️ Aviso legal",
	DisclaimerBody:    "Este repositório agrega configurações de proxy **compartilhadas publicamente** por voluntários de terceiros. Não operamos nenhum servidor, não garantimos disponibilidade ou segurança, e não somos responsáveis pelo uso. Destinado a uso educacional e conectividade pessoal. Cumpra todas as leis aplicáveis em sua jurisdição.",

	StarHistoryHeading: "## ⭐ Histórico de estrelas",
	FinalCTA:           "Se este projeto te ajudou, deixe uma ⭐ — cada estrela facilita para outros o encontrarem.",
}
