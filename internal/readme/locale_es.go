package readme

var ES = Locale{
	Code:        "es",
	DisplayName: "Español",
	FileName:    "README_ES.md",
	LangAttr:    "es",

	BadgeNodes:   "nodos",
	BadgeAlive:   "activos",
	BadgeMedian:  "rtt--mediana",
	BadgeUpdated: "actualizado",

	Hook1:       "**La forma más fácil de obtener una VPN gratuita que funciona — copia un enlace de suscripción, pégalo en tu cliente, conecta.**",
	Hook2:       "Sin registro. Sin pago. Sin instalar ningún binario. Actualizado cada hora desde fuentes públicas — cada nodo publicado ha reenviado tráfico HTTP real a través de sing-box hace minutos.",
	KeywordLine: "VPN gratis · suscripción VPN gratuita · proxy gratis · Clash suscripción · v2ray suscripción · sing-box suscripción · VLESS · Reality · VMess · Trojan · Shadowsocks · Hysteria2 · actualizado por hora · HTTP verificado sobre proxy · por país",

	WhyHeading: "## 💡 ¿Por qué este proyecto?",
	WhyBody:    "Cada lista de \"VPN gratuita\" en GitHub está desactualizada, llena de nodos muertos, o te pide instalar un binario dudoso. Este repositorio va un paso más allá que cualquier otro —— **no solo verificamos que el nodo responda, sino que empujamos tráfico HTTP real a través de él con sing-box y confirmamos que vuelve un 204** antes de publicar, todo en minutos. Obtienes 3 archivos de suscripción portables — úsalos en Clash, sing-box o v2rayN y listo.",

	VerificationHeading: "## 🔬 Cómo verificamos que los nodos realmente funcionan",
	VerificationBody: `La mayoría de listas de VPN gratuitas paran en \"el puerto TCP está abierto\" y publican. Nosotros no. Aquí está la tubería completa que un nodo debe superar antes de entrar en la suscripción.

### ✅ Qué verificamos en tiempo de agregación (antes de publicar)

1. **Accesibilidad TCP** — abrimos una conexión TCP a cada ` + "`server:port`" + `. Hosts caídos, DNS incorrecto y puertos bloqueados se descartan. ~40 % de las entradas crudas caen aquí.
2. **Handshake TLS** — para cada nodo TLS / Reality / WS-TLS completamos el handshake entero. Certificados expirados, SNI incorrectos y short-ids de Reality rotos se descartan. ~10 % más caen aquí.
3. **Validación de configuración sing-box** — cada nodo sobreviviente se traduce a un outbound real de sing-box y pasa por ` + "`sing-box check`" + `. Cifras corruptas, UUIDs incorrectos y opciones flow no soportadas se descartan antes de gastar un slot de sondeo.
4. **Sondeo HTTP-over-proxy (esta es la clave)** — agrupamos los ~900 candidatos más rápidos en subprocesos sing-box, cada nodo con su propio inbound SOCKS5 local, y enviamos GETs HTTP y HTTPS reales a través de él:
   - ` + "`http://www.gstatic.com/generate_204`" + ` (espera 204)
   - ` + "`https://www.cloudflare.com/cdn-cgi/trace`" + ` (espera 200)

   La solicitud atraviesa el protocolo proxy real (VLESS / VMess / Trojan / Shadowsocks / Hysteria2), así que un nodo que pasa tiene demostrablemente autenticación, enrutamiento, handshake TLS interno y red de salida funcionales.
5. **Dos rondas, 45 segundos de separación** — nodos que pasan una vez pero mueren 45 segundos después se filtran. Solo nodos con ≥ 50 % de éxito en (rondas × objetivos) se mantienen.
6. **Ordenar por mediana de latencia real** — los sobrevivientes se ordenan por la mediana del ida y vuelta HTTP-over-proxy (no RTT TCP crudo) y los top N se publican.

Números típicos de una ejecución reciente: **17 fuentes → ~4,800 crudos → ~2,900 TCP vivos → ~2,600 TLS OK → ~840 configuración válida → ~280 verificados por HTTP → top 150 publicados**. Cada uno de los 150 ha reenviado tráfico realmente en los últimos diez minutos.

### ❌ Qué todavía no podemos verificar

- **Ancho de banda / throughput** — medimos latencia, no megabits. Un nodo de 50 ms puede seguir siendo lento para vídeo.
- **Precisión de geolocalización** — GeoIP dice el país de la IP de salida pero no la ciudad o ISP confiablemente.
- **Bloqueos específicos por región** — un nodo que funciona desde nuestra infraestructura de sondeo puede estar bloqueado desde la tuya (filtrado a nivel ISP, captive portals, etc.).
- **Seguir vivo después de la ejecución** — el nodo pasó hace diez minutos; puede haber muerto desde entonces.

### 🛡️ Red de seguridad en tiempo de ejecución — para el último punto arriba

El ` + "`clash.yaml`" + ` que publicamos incluye un grupo ` + "`url-test`" + ` que retesta HTTP real a través de cada nodo cada 5 minutos en *tu* dispositivo:

` + "```yaml" + `
proxy-groups:
  - name: AUTO
    type: url-test
    url: http://www.gstatic.com/generate_204
    interval: 300
` + "```" + `

Tu cliente mantiene la lista de nodos ordenada por latencia *en vivo* de HTTP-over-proxy desde tu red y selecciona automáticamente el nodo más rápido que funciona. sing-box y v2ray tienen mecanismos equivalentes. Si un nodo seleccionado muere entre agregaciones por hora, el cliente cambia al siguiente sin intervención.

### 🧮 Qué significa en la práctica

De los ~150 que publicamos cada ejecución, un cliente típico encuentra **80-120 nodos que sirven HTTP limpiamente desde su red** en cualquier momento — aproximadamente 2-3× la tasa de acierto de listas que solo hacen sondeo TCP. El grupo url-test rota de forma transparente si uno se cae.`,

	SubscribeHeading:   "## 🚀 Suscripción con un clic",
	SubscribeIntro:     "Copia la URL que coincida con tu cliente y pégala en el campo de importación de suscripción:",
	SubscribeColClient: "Cliente",
	SubscribeColFormat: "Formato",
	SubscribeColURL:    "URL de suscripción",

	ClientsHeading: "## 🧩 Clientes compatibles",
	ClientsWindows: "**Windows**: v2rayN, Clash Verge, Hiddify, NekoRay",
	ClientsMacOS:   "**macOS**: ClashX Pro, Clash Verge, sing-box, Hiddify",
	ClientsIOS:     "**iOS**: Shadowrocket, Stash, Loon, sing-box, Hiddify",
	ClientsAndroid: "**Android**: v2rayNG, NekoBox, Clash Meta for Android, Hiddify, sing-box",
	ClientsLinux:   "**Linux**: mihomo (Clash.Meta), sing-box, v2ray-core",

	StatsHeading:     "## 📊 Estadísticas en vivo",
	StatsNodes:       "**Nodos seleccionados**",
	StatsAlive:       "**Activos en todas las fuentes**",
	StatsFastest:     "**RTT del nodo más rápido**",
	StatsMedian:      "**RTT mediana**",
	StatsUpdated:     "**Última actualización (UTC)**",
	ProtocolMixLabel: "**Mezcla de protocolos:**",
	SourcesLabel:     "**Fuentes usadas en esta ejecución:**",

	ByCountryHeading: "## 🌍 Por país",
	ByCountryIntro:   "¿Quieres nodos solo en una región específica? Usa una de estas URLs de suscripción dirigidas:",
	ByCountryColCC:   "País",
	ByCountryColN:    "Nodos",

	GuidesHeading:     "## 📖 Guías paso a paso",
	GuidesIntro:       "¿Nuevo con los clientes VPN? Elige tu plataforma y sigue el tutorial:",
	GuideLocaleSuffix: "",

	FAQHeading: "## ❓ Preguntas frecuentes",
	FAQ1Q:      "¿Es realmente gratis?",
	FAQ1A:      "Sí. Los nodos son operados por voluntarios externos que publican sus propias suscripciones gratuitas. Nosotros no operamos ningún servidor — solo probamos, clasificamos y reempaquetamos lo que ya es público.",
	FAQ2Q:      "¿Qué tan reciente es la información?",
	FAQ2A:      "Cada hora (con un pequeño retraso aleatorio para evitar golpear las fuentes upstream exactamente en `:00`): trae todas las fuentes → TCP → TLS → validación de configuración sing-box → sondeo HTTP-over-proxy (dos rondas, 45 s de separación) → ordena por latencia HTTP real → publica los archivos nuevos. La tubería completa tarda ~10 minutos. Consulta la marca de tiempo `Last updated` arriba.",
	FAQ3Q:      "¿Puedo confiar en estos nodos?",
	FAQ3A:      "Los nodos gratis ven todo tu tráfico. **Nunca los uses para banca, login o algo sensible.** Bien para saltar bloqueos geográficos en contenido público. Usa tu propio VPS / proveedor de pago para privacidad real.",
	FAQ4Q:      "¿Por qué algunos nodos listados fallan?",
	FAQ4A:      "Incluso después de nuestro sondeo HTTP-over-proxy, los nodos pueden morir entre agregaciones: cuota agotada, upstream revocó la clave, tu ISP bloquea la IP de salida, o el operador lo apagó. El `clash.yaml` publicado incluye un grupo `url-test` (`http://www.gstatic.com/generate_204`, intervalo de 300 s); tu cliente selecciona automáticamente el nodo más rápido que realmente sirve HTTP *desde tu red*. Si uno muere, toma el siguiente. Espera que 80-120 de los 150 funcionen en cualquier momento.",

	ContributingHeading: "## 🤝 Contribuir",
	ContributingBody:    "¿Conoces una fuente de suscripción pública confiable que deberíamos agregar? Abre un issue con la URL y el formato.",

	DisclaimerHeading: "## ⚠️ Aviso legal",
	DisclaimerBody:    "Este repositorio agrega configuraciones de proxy **compartidas públicamente** por voluntarios externos. No operamos ningún servidor, no garantizamos disponibilidad ni seguridad, y no somos responsables del uso que hagas. Destinado a uso educativo y de conectividad personal. Cumple con todas las leyes aplicables en tu jurisdicción.",

	StarHistoryHeading: "## ⭐ Historia de estrellas",
	FinalCTA:           "Si este proyecto te ayudó, déjale una ⭐ — cada estrella hace más fácil que otros lo encuentren.",
}
