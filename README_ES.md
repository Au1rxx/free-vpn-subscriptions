# Free VPN Subscriptions

[English](./README.md) · [简体中文](./README_CN.md) · [日本語](./README_JA.md) · [한국어](./README_KO.md) · **Español** · [Português](./README_PT.md) · [Русский](./README_RU.md)

<p align="center"><img src="https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/assets/hero.png" alt="Free VPN Subscriptions — hourly-refreshed free VPN subscriptions for Clash, sing-box, v2ray" width="780"></p>

![nodos](https://img.shields.io/badge/nodos-144-brightgreen) ![activos](https://img.shields.io/badge/activos-2082-blue) ![rtt--mediana](https://img.shields.io/badge/rtt--mediana-193ms-orange) ![actualizado](https://img.shields.io/badge/actualizado-2026-05-03_17:22_UTC-informational)

> **La forma más fácil de obtener una VPN gratuita que funciona — copia un enlace de suscripción, pégalo en tu cliente, conecta.**  
> Sin registro. Sin pago. Sin instalar ningún binario. Actualizado cada hora desde fuentes públicas — cada nodo publicado ha reenviado tráfico HTTP real a través de sing-box hace minutos.

> VPN gratis · suscripción VPN gratuita · proxy gratis · Clash suscripción · v2ray suscripción · sing-box suscripción · VLESS · Reality · VMess · Trojan · Shadowsocks · Hysteria2 · actualizado por hora · HTTP verificado sobre proxy · por país

## 💡 ¿Por qué este proyecto?

Cada lista de "VPN gratuita" en GitHub está desactualizada, llena de nodos muertos, o te pide instalar un binario dudoso. Este repositorio va un paso más allá que cualquier otro —— **no solo verificamos que el nodo responda, sino que empujamos tráfico HTTP real a través de él con sing-box y confirmamos que vuelve un 204** antes de publicar, todo en minutos. Obtienes 3 archivos de suscripción portables — úsalos en Clash, sing-box o v2rayN y listo.

> 📖 How the fetch → probe → rank pipeline works: [ARCHITECTURE.md](./ARCHITECTURE.md)

## 🔬 Cómo verificamos que los nodos realmente funcionan

La mayoría de listas de VPN gratuitas paran en \"el puerto TCP está abierto\" y publican. Nosotros no. Aquí está la tubería completa que un nodo debe superar antes de entrar en la suscripción.

### ✅ Qué verificamos en tiempo de agregación (antes de publicar)

1. **Accesibilidad TCP** — abrimos una conexión TCP a cada `server:port`. Hosts caídos, DNS incorrecto y puertos bloqueados se descartan. ~40 % de las entradas crudas caen aquí.
2. **Handshake TLS** — para cada nodo TLS / Reality / WS-TLS completamos el handshake entero. Certificados expirados, SNI incorrectos y short-ids de Reality rotos se descartan. ~10 % más caen aquí.
3. **Validación de configuración sing-box** — cada nodo sobreviviente se traduce a un outbound real de sing-box y pasa por `sing-box check`. Cifras corruptas, UUIDs incorrectos y opciones flow no soportadas se descartan antes de gastar un slot de sondeo.
4. **Sondeo HTTP-over-proxy (esta es la clave)** — agrupamos los ~900 candidatos más rápidos en subprocesos sing-box, cada nodo con su propio inbound SOCKS5 local, y enviamos GETs HTTP y HTTPS reales a través de él:
   - `http://www.gstatic.com/generate_204` (espera 204)
   - `https://www.cloudflare.com/cdn-cgi/trace` (espera 200)

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

El `clash.yaml` que publicamos incluye un grupo `url-test` que retesta HTTP real a través de cada nodo cada 5 minutos en *tu* dispositivo:

```yaml
proxy-groups:
  - name: AUTO
    type: url-test
    url: http://www.gstatic.com/generate_204
    interval: 300
```

Tu cliente mantiene la lista de nodos ordenada por latencia *en vivo* de HTTP-over-proxy desde tu red y selecciona automáticamente el nodo más rápido que funciona. sing-box y v2ray tienen mecanismos equivalentes. Si un nodo seleccionado muere entre agregaciones por hora, el cliente cambia al siguiente sin intervención.

### 🧮 Qué significa en la práctica

De los ~150 que publicamos cada ejecución, un cliente típico encuentra **80-120 nodos que sirven HTTP limpiamente desde su red** en cualquier momento — aproximadamente 2-3× la tasa de acierto de listas que solo hacen sondeo TCP. El grupo url-test rota de forma transparente si uno se cae.

## 🚀 Suscripción con un clic

Copia la URL que coincida con tu cliente y pégala en el campo de importación de suscripción:

| Cliente | Formato | URL de suscripción |
|---|---|---|
| Clash / Clash Verge / ClashX | `clash.yaml` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/clash.yaml` |
| sing-box | `singbox.json` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/singbox.json` |
| v2rayN / v2rayNG / Shadowrocket / NekoBox | `v2ray-base64` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/v2ray-base64.txt` |

## 🌍 Por país

¿Quieres nodos solo en una región específica? Usa una de estas URLs de suscripción dirigidas:

| País | Nodos | Clash | sing-box | v2ray |
|---|---|---|---|---|
| 🇺🇸 United States (`US`) | 45 | [clash-US.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-US.yaml) | [singbox-US.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-US.json) | [v2ray-base64-US.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-US.txt) |
| 🇨🇦 Canada (`CA`) | 8 | [clash-CA.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-CA.yaml) | [singbox-CA.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-CA.json) | [v2ray-base64-CA.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-CA.txt) |
| 🇩🇪 Germany (`DE`) | 8 | [clash-DE.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-DE.yaml) | [singbox-DE.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-DE.json) | [v2ray-base64-DE.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-DE.txt) |
| 🇬🇧 United Kingdom (`GB`) | 4 | [clash-GB.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-GB.yaml) | [singbox-GB.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-GB.json) | [v2ray-base64-GB.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-GB.txt) |
| 🇳🇱 Netherlands (`NL`) | 3 | [clash-NL.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-NL.yaml) | [singbox-NL.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-NL.json) | [v2ray-base64-NL.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-NL.txt) |

## 📖 Guías paso a paso

¿Nuevo con los clientes VPN? Elige tu plataforma y sigue el tutorial:

- [**Clash Verge**](https://au1rxx.github.io/free-vpn-subscriptions/guides/clash-verge.html) · Windows / macOS / Linux
- [**v2rayNG**](https://au1rxx.github.io/free-vpn-subscriptions/guides/v2rayng.html) · Android
- [**Shadowrocket**](https://au1rxx.github.io/free-vpn-subscriptions/guides/shadowrocket.html) · iOS / iPadOS
- [**sing-box**](https://au1rxx.github.io/free-vpn-subscriptions/guides/sing-box.html) · Windows / macOS / Linux / iOS / Android

## 🧩 Clientes compatibles

- **Windows**: v2rayN, Clash Verge, Hiddify, NekoRay
- **macOS**: ClashX Pro, Clash Verge, sing-box, Hiddify
- **iOS**: Shadowrocket, Stash, Loon, sing-box, Hiddify
- **Android**: v2rayNG, NekoBox, Clash Meta for Android, Hiddify, sing-box
- **Linux**: mihomo (Clash.Meta), sing-box, v2ray-core

## 📊 Estadísticas en vivo

- **Nodos seleccionados**: 144
- **Activos en todas las fuentes**: 2082
- **RTT del nodo más rápido**: 29 ms
- **RTT mediana**: 193 ms
- **Última actualización (UTC)**: 2026-05-03 17:22 UTC

**Mezcla de protocolos:** shadowsocks × 27 · trojan × 23 · vless × 82 · vmess × 12

**Fuentes usadas en esta ejecución:** `barry-far-v2ray` × 28 · `ebrasha-v2ray` × 1 · `epodonios` × 9 · `lagzian-mix` × 2 · `mahdi0024` × 27 · `mahdibland-aggregator` × 5 · `mahdibland-shadowsocks` × 8 · `mfuu-clash` × 3 · `ninjastrikers` × 54 · `ruking-clash` × 1 · `surfboard-eternity` × 6

## ❓ Preguntas frecuentes

<details><summary>¿Es realmente gratis?</summary>

Sí. Los nodos son operados por voluntarios externos que publican sus propias suscripciones gratuitas. Nosotros no operamos ningún servidor — solo probamos, clasificamos y reempaquetamos lo que ya es público.

</details>

<details><summary>¿Qué tan reciente es la información?</summary>

Cada hora (con un pequeño retraso aleatorio para evitar golpear las fuentes upstream exactamente en `:00`): trae todas las fuentes → TCP → TLS → validación de configuración sing-box → sondeo HTTP-over-proxy (dos rondas, 45 s de separación) → ordena por latencia HTTP real → publica los archivos nuevos. La tubería completa tarda ~10 minutos. Consulta la marca de tiempo `Last updated` arriba.

</details>

<details><summary>¿Puedo confiar en estos nodos?</summary>

Los nodos gratis ven todo tu tráfico. **Nunca los uses para banca, login o algo sensible.** Bien para saltar bloqueos geográficos en contenido público. Usa tu propio VPS / proveedor de pago para privacidad real.

</details>

<details><summary>¿Por qué algunos nodos listados fallan?</summary>

Incluso después de nuestro sondeo HTTP-over-proxy, los nodos pueden morir entre agregaciones: cuota agotada, upstream revocó la clave, tu ISP bloquea la IP de salida, o el operador lo apagó. El `clash.yaml` publicado incluye un grupo `url-test` (`http://www.gstatic.com/generate_204`, intervalo de 300 s); tu cliente selecciona automáticamente el nodo más rápido que realmente sirve HTTP *desde tu red*. Si uno muere, toma el siguiente. Espera que 80-120 de los 150 funcionen en cualquier momento.

</details>

## 🤝 Contribuir

¿Conoces una fuente de suscripción pública confiable que deberíamos agregar? Abre un issue con la URL y el formato.

## ⚠️ Aviso legal

Este repositorio agrega configuraciones de proxy **compartidas públicamente** por voluntarios externos. No operamos ningún servidor, no garantizamos disponibilidad ni seguridad, y no somos responsables del uso que hagas. Destinado a uso educativo y de conectividad personal. Cumple con todas las leyes aplicables en tu jurisdicción.

## ⭐ Historia de estrellas

[![Star History Chart](https://api.star-history.com/svg?repos=Au1rxx/free-vpn-subscriptions&type=Date)](https://www.star-history.com/#Au1rxx/free-vpn-subscriptions&Date)

---

Si este proyecto te ayudó, déjale una ⭐ — cada estrella hace más fácil que otros lo encuentren.
