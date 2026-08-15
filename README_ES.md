# Free VPN Subscriptions

<div align="center">

[English](./README.md) · [简体中文](./README_CN.md) · [日本語](./README_JA.md) · [한국어](./README_KO.md) · **Español** · [Português](./README_PT.md) · [Русский](./README_RU.md)

</div>

<p align="center"><img src="https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/assets/hero.png" alt="Free VPN Subscriptions — hourly-refreshed free VPN subscriptions for Clash, sing-box, v2ray" width="780"></p>

[![GitHub stars](https://img.shields.io/github/stars/Au1rxx/free-vpn-subscriptions?style=flat&color=gold&logo=github)](https://github.com/Au1rxx/free-vpn-subscriptions/stargazers) ![seleccionados](https://img.shields.io/badge/seleccionados-2000-brightgreen) ![verificados](https://img.shields.io/badge/verificados-9572-blue) ![rtt--mediana](https://img.shields.io/badge/rtt--mediana-645ms-orange) ![actualizado](https://img.shields.io/badge/actualizado-2026-08-15_22:27_UTC-informational) [![License](https://img.shields.io/github/license/Au1rxx/free-vpn-subscriptions?color=blue)](https://github.com/Au1rxx/free-vpn-subscriptions/blob/main/LICENSE)

> **La forma más fácil de obtener una VPN gratuita que funciona — copia un enlace de suscripción, pégalo en tu cliente, conecta.**  
> Sin registro. Sin pago. Sin instalar ningún binario. Actualizado cada hora desde fuentes públicas — cada nodo publicado ha reenviado tráfico HTTP real a través de sing-box hace minutos.

> VPN gratis · suscripción VPN gratuita · proxy gratis · Clash suscripción · v2ray suscripción · sing-box suscripción · VLESS · Reality · VMess · Trojan · Shadowsocks · Hysteria2 · actualizado por hora · HTTP verificado sobre proxy · por país

## 🚀 Suscripción con un clic

Copia la URL que coincida con tu cliente y pégala en el campo de importación de suscripción:

| Cliente | Formato | URL de suscripción |
|---|---|---|
| Clash / Clash Verge / ClashX | `clash.yaml` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/clash.yaml` |
| sing-box | `singbox.json` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/singbox.json` |
| v2rayN / v2rayNG / Shadowrocket / NekoBox | `v2ray-base64` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/v2ray-base64.txt` |

> Database-classified shards (stable / all verified / protocol / country / network): [https://github.com/Au1rxx/free-vpn-subscriptions/tree/main/output](https://github.com/Au1rxx/free-vpn-subscriptions/tree/main/output)

⭐ **Dale una estrella al repositorio** si te ahorró tiempo: es la única señal que usamos para decidir qué seguir construyendo.

## 💡 ¿Por qué este proyecto?

Cada lista de "VPN gratuita" en GitHub está desactualizada, llena de nodos muertos, o te pide instalar un binario dudoso. Este repositorio va un paso más allá que cualquier otro —— **no solo verificamos que el nodo responda, sino que empujamos tráfico HTTP real a través de él con sing-box y confirmamos que vuelve un 204** antes de publicar, todo en minutos. Obtienes 3 archivos de suscripción portables — úsalos en Clash, sing-box o v2rayN y listo.

<details>
<summary><b>🔬 Cómo verificamos que los nodos realmente funcionan</b></summary>

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

</details>

## 🌍 Por país

¿Quieres nodos solo en una región específica? Usa una de estas URLs de suscripción dirigidas:

| País | Nodos | Clash | sing-box | v2ray |
|---|---|---|---|---|
| 🇯🇵 Japan (`JP`) | 502 | [clash-JP.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-JP.yaml) | [singbox-JP.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-JP.json) | [v2ray-base64-JP.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-JP.txt) |
| 🇺🇸 United States (`US`) | 403 | [clash-US.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-US.yaml) | [singbox-US.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-US.json) | [v2ray-base64-US.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-US.txt) |
| 🇰🇷 Korea (`KR`) | 338 | [clash-KR.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-KR.yaml) | [singbox-KR.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-KR.json) | [v2ray-base64-KR.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-KR.txt) |
| 🇳🇱 Netherlands (`NL`) | 153 | [clash-NL.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-NL.yaml) | [singbox-NL.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-NL.json) | [v2ray-base64-NL.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-NL.txt) |
| 🇸🇬 Singapore (`SG`) | 120 | [clash-SG.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-SG.yaml) | [singbox-SG.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-SG.json) | [v2ray-base64-SG.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-SG.txt) |
| 🇭🇰 Hong Kong (`HK`) | 89 | [clash-HK.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-HK.yaml) | [singbox-HK.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-HK.json) | [v2ray-base64-HK.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-HK.txt) |
| 🇩🇪 Germany (`DE`) | 88 | [clash-DE.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-DE.yaml) | [singbox-DE.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-DE.json) | [v2ray-base64-DE.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-DE.txt) |
| 🇫🇷 France (`FR`) | 68 | [clash-FR.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-FR.yaml) | [singbox-FR.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-FR.json) | [v2ray-base64-FR.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-FR.txt) |
| 🇫🇮 Finland (`FI`) | 35 | [clash-FI.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-FI.yaml) | [singbox-FI.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-FI.json) | [v2ray-base64-FI.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-FI.txt) |
| 🇵🇱 Poland (`PL`) | 32 | [clash-PL.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-PL.yaml) | [singbox-PL.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-PL.json) | [v2ray-base64-PL.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-PL.txt) |
| 🇨🇦 Canada (`CA`) | 30 | [clash-CA.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-CA.yaml) | [singbox-CA.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-CA.json) | [v2ray-base64-CA.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-CA.txt) |
| 🇮🇪 Ireland (`IE`) | 27 | [clash-IE.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-IE.yaml) | [singbox-IE.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-IE.json) | [v2ray-base64-IE.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-IE.txt) |
| 🇷🇺 Russia (`RU`) | 16 | [clash-RU.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-RU.yaml) | [singbox-RU.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-RU.json) | [v2ray-base64-RU.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-RU.txt) |
| 🇦🇺 Australia (`AU`) | 12 | [clash-AU.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-AU.yaml) | [singbox-AU.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-AU.json) | [v2ray-base64-AU.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-AU.txt) |
| 🇬🇧 United Kingdom (`GB`) | 12 | [clash-GB.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-GB.yaml) | [singbox-GB.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-GB.json) | [v2ray-base64-GB.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-GB.txt) |
| 🇮🇳 India (`IN`) | 12 | [clash-IN.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-IN.yaml) | [singbox-IN.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-IN.json) | [v2ray-base64-IN.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-IN.txt) |
| 🇹🇼 Taiwan (`TW`) | 11 | [clash-TW.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-TW.yaml) | [singbox-TW.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-TW.json) | [v2ray-base64-TW.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-TW.txt) |
| 🇪🇪 Estonia (`EE`) | 7 | [clash-EE.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-EE.yaml) | [singbox-EE.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-EE.json) | [v2ray-base64-EE.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-EE.txt) |
| 🇱🇻 Latvia (`LV`) | 6 | [clash-LV.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-LV.yaml) | [singbox-LV.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-LV.json) | [v2ray-base64-LV.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-LV.txt) |
| 🇦🇹 Austria (`AT`) | 4 | [clash-AT.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-AT.yaml) | [singbox-AT.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-AT.json) | [v2ray-base64-AT.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-AT.txt) |
| 🇪🇸 Spain (`ES`) | 4 | [clash-ES.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-ES.yaml) | [singbox-ES.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-ES.json) | [v2ray-base64-ES.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-ES.txt) |
| 🇧🇬 Bulgaria (`BG`) | 3 | [clash-BG.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-BG.yaml) | [singbox-BG.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-BG.json) | [v2ray-base64-BG.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-BG.txt) |
| 🇨🇿 Czechia (`CZ`) | 3 | [clash-CZ.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-CZ.yaml) | [singbox-CZ.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-CZ.json) | [v2ray-base64-CZ.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-CZ.txt) |
| 🇱🇹 Lithuania (`LT`) | 3 | [clash-LT.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-LT.yaml) | [singbox-LT.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-LT.json) | [v2ray-base64-LT.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-LT.txt) |
| 🇵🇹 Portugal (`PT`) | 3 | [clash-PT.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-PT.yaml) | [singbox-PT.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-PT.json) | [v2ray-base64-PT.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-PT.txt) |

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

- **Nodos seleccionados**: 2000
- **Activos en todas las fuentes**: 9572
- **RTT del nodo más rápido**: 14 ms
- **RTT mediana**: 645 ms
- **Última actualización (UTC)**: 2026-08-15 22:27 UTC

**Mezcla de protocolos:** https × 1 · hysteria2 × 76 · shadowsocks × 137 · socks4 × 1 · trojan × 1016 · tuic × 2 · vless × 591 · vmess × 160 · wireguard × 16

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

[![Star History Chart](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/assets/star-history.svg)](https://github.com/Au1rxx/free-vpn-subscriptions/stargazers)

---

Si este proyecto te ayudó, déjale una ⭐ — cada estrella hace más fácil que otros lo encuentren.
