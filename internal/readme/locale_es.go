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
	Hook2:       "Sin registro. Sin pago. Sin instalar ningún binario. Actualizado cada hora desde fuentes públicas y cada nodo es verificado.",
	KeywordLine: "VPN gratis · suscripción VPN gratuita · proxy gratis · Clash suscripción · v2ray suscripción · sing-box suscripción · VLESS · Reality · VMess · Trojan · Shadowsocks · Hysteria2 · actualizado por hora · TCP+TLS verificado · por país",

	WhyHeading: "## 💡 ¿Por qué este proyecto?",
	WhyBody:    "Cada lista de \"VPN gratuita\" en GitHub está desactualizada, llena de nodos muertos, o te pide instalar un binario dudoso. Este repositorio **solo publica nodos que pasaron un handshake TCP y un handshake TLS hace minutos**, desde fuentes públicas curadas, ordenados por latencia. Obtienes 3 archivos de suscripción portables — úsalos en Clash, sing-box o v2rayN y listo.",

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
	FAQ2A:      "Una GitHub Action se ejecuta cada hora: trae todas las fuentes, prueba cada nodo con TCP+TLS, elimina los muertos, ordena por latencia y comitea los archivos nuevos. Consulta la marca de tiempo `Last updated` arriba.",
	FAQ3Q:      "¿Puedo confiar en estos nodos?",
	FAQ3A:      "Los nodos gratis ven todo tu tráfico. **Nunca los uses para banca, login o algo sensible.** Bien para saltar bloqueos geográficos en contenido público. Usa tu propio VPS / proveedor de pago para privacidad real.",
	FAQ4Q:      "¿Por qué algunos nodos listados fallan?",
	FAQ4A:      "Verificamos accesibilidad TCP y handshake TLS, pero un nodo aún puede tener cuota expirada, ruteo incorrecto o certificado caducado. Prueba varios; el grupo selector te da alternativas.",

	ContributingHeading: "## 🤝 Contribuir",
	ContributingBody:    "¿Conoces una fuente de suscripción pública confiable que deberíamos agregar? Abre un issue con la URL y el formato.",

	DisclaimerHeading: "## ⚠️ Aviso legal",
	DisclaimerBody:    "Este repositorio agrega configuraciones de proxy **compartidas públicamente** por voluntarios externos. No operamos ningún servidor, no garantizamos disponibilidad ni seguridad, y no somos responsables del uso que hagas. Destinado a uso educativo y de conectividad personal. Cumple con todas las leyes aplicables en tu jurisdicción.",

	StarHistoryHeading: "## ⭐ Historia de estrellas",
	FinalCTA:           "Si este proyecto te ayudó, déjale una ⭐ — cada estrella hace más fácil que otros lo encuentren.",
}
