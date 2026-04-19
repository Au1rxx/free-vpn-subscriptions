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
	Hook2:       "Sin registro. Sin pago. Sin instalar ningún binario. Actualizado cada hora desde fuentes públicas — cada nodo se verifica con sondeo TCP + TLS antes de publicarse.",
	KeywordLine: "VPN gratis · suscripción VPN gratuita · proxy gratis · Clash suscripción · v2ray suscripción · sing-box suscripción · VLESS · Reality · VMess · Trojan · Shadowsocks · Hysteria2 · actualizado por hora · TCP+TLS verificado · por país",

	WhyHeading: "## 💡 ¿Por qué este proyecto?",
	WhyBody:    "Cada lista de \"VPN gratuita\" en GitHub está desactualizada, llena de nodos muertos, o te pide instalar un binario dudoso. Este repositorio **solo publica nodos que pasaron un handshake TCP y un handshake TLS hace minutos**, desde fuentes públicas curadas, ordenados por latencia. Obtienes 3 archivos de suscripción portables — úsalos en Clash, sing-box o v2rayN y listo.",

	VerificationHeading: "## 🔬 Cómo verificamos que los nodos realmente funcionan",
	VerificationBody: `**Respuesta honesta primero: no podemos *garantizar* que un nodo pase tu tráfico.** Ningún agregador puede sin enviar tráfico real a través de él. Aquí está exactamente lo que verificamos, lo que no, y de dónde viene la garantía real.

### ✅ Qué verificamos en tiempo de agregación (antes de publicar)

1. **Accesibilidad TCP** — abrimos una conexión TCP a cada ` + "`server:port`" + `. Hosts caídos, DNS incorrecto y puertos bloqueados se descartan. Elimina aproximadamente el 40 % de las entradas crudas.
2. **Handshake TLS** — para cada nodo TLS / Reality / WS-TLS completamos el handshake entero. Certificados expirados, SNI incorrectos y short-ids de Reality rotos se descartan. Elimina otro ~10 %.
3. **Orden por latencia** — los supervivientes se ordenan por RTT y guardamos los N más rápidos.

Números típicos de una ejecución reciente: **17 fuentes → ~4,800 crudos → ~2,900 TCP vivos → ~2,600 TLS OK → top 200 publicados**.

### ❌ Qué no podemos verificar

- Autenticación del protocolo proxy. UUID / contraseña incorrectos solo son rechazados *después* del handshake TLS por el servidor upstream.
- Éxito real de HTTP a través del proxy.
- Ancho de banda o throughput.
- Geolocalización más allá de lo que GeoIP dice sobre la IP de salida.

### 🛡️ Verificación en tiempo de ejecución — de aquí viene la garantía real

El ` + "`clash.yaml`" + ` que publicamos incluye un grupo ` + "`url-test`" + ` que **prueba HTTP real a través de cada nodo** cada 5 minutos:

` + "```yaml" + `
proxy-groups:
  - name: AUTO
    type: url-test
    url: http://www.gstatic.com/generate_204
    interval: 300
` + "```" + `

Tu cliente ordena la lista de nodos por latencia *real* de HTTP a través del proxy y selecciona automáticamente el nodo más rápido que funciona. sing-box y v2ray tienen mecanismos equivalentes. Si un nodo seleccionado muere, el cliente cambia al siguiente sin intervención.

### 🧮 Resultado esperado

De los top 200 publicados cada ejecución, un cliente típico encontrará 30-50 que sirven HTTP limpiamente en cualquier momento. Rota si uno se pone lento — el grupo url-test hace eso con un clic.`,

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
	FAQ2A:      "Cada hora (con un pequeño retraso aleatorio para evitar golpear las fuentes upstream exactamente en `:00`): trae todas las fuentes, prueba cada nodo con TCP+TLS, elimina los muertos, ordena por latencia y publica los archivos nuevos. Consulta la marca de tiempo `Last updated` arriba.",
	FAQ3Q:      "¿Puedo confiar en estos nodos?",
	FAQ3A:      "Los nodos gratis ven todo tu tráfico. **Nunca los uses para banca, login o algo sensible.** Bien para saltar bloqueos geográficos en contenido público. Usa tu propio VPS / proveedor de pago para privacidad real.",
	FAQ4Q:      "¿Por qué algunos nodos listados fallan?",
	FAQ4A:      "Solo verificamos accesibilidad TCP y handshake TLS — un nodo aún puede tener cuota expirada, ruteo incorrecto o certificado caducado. El `clash.yaml` publicado incluye un grupo `url-test` (`http://www.gstatic.com/generate_204`, intervalo de 300 s); tu cliente selecciona automáticamente el nodo más rápido que realmente sirve HTTP. Si uno muere, toma el siguiente.",

	ContributingHeading: "## 🤝 Contribuir",
	ContributingBody:    "¿Conoces una fuente de suscripción pública confiable que deberíamos agregar? Abre un issue con la URL y el formato.",

	DisclaimerHeading: "## ⚠️ Aviso legal",
	DisclaimerBody:    "Este repositorio agrega configuraciones de proxy **compartidas públicamente** por voluntarios externos. No operamos ningún servidor, no garantizamos disponibilidad ni seguridad, y no somos responsables del uso que hagas. Destinado a uso educativo y de conectividad personal. Cumple con todas las leyes aplicables en tu jurisdicción.",

	StarHistoryHeading: "## ⭐ Historia de estrellas",
	FinalCTA:           "Si este proyecto te ayudó, déjale una ⭐ — cada estrella hace más fácil que otros lo encuentren.",
}
