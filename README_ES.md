# Free VPN Subscriptions

[English](./README.md) · [简体中文](./README_CN.md) · [日本語](./README_JA.md) · [한국어](./README_KO.md) · **Español** · [Português](./README_PT.md) · [Русский](./README_RU.md)

<p align="center"><img src="https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/assets/hero.png" alt="Free VPN Subscriptions — hourly-refreshed free VPN subscriptions for Clash, sing-box, v2ray" width="780"></p>

![nodos](https://img.shields.io/badge/nodos-150-brightgreen) ![activos](https://img.shields.io/badge/activos-2590-blue) ![rtt--mediana](https://img.shields.io/badge/rtt--mediana-17ms-orange) ![actualizado](https://img.shields.io/badge/actualizado-2026-04-19_10:23_UTC-informational)

> **La forma más fácil de obtener una VPN gratuita que funciona — copia un enlace de suscripción, pégalo en tu cliente, conecta.**  
> Sin registro. Sin pago. Sin instalar ningún binario. Actualizado cada hora desde fuentes públicas y cada nodo es verificado.

> VPN gratis · suscripción VPN gratuita · proxy gratis · Clash suscripción · v2ray suscripción · sing-box suscripción · VLESS · Reality · VMess · Trojan · Shadowsocks · Hysteria2 · actualizado por hora · TCP+TLS verificado · por país

## 💡 ¿Por qué este proyecto?

Cada lista de "VPN gratuita" en GitHub está desactualizada, llena de nodos muertos, o te pide instalar un binario dudoso. Este repositorio **solo publica nodos que pasaron un handshake TCP y un handshake TLS hace minutos**, desde fuentes públicas curadas, ordenados por latencia. Obtienes 3 archivos de suscripción portables — úsalos en Clash, sing-box o v2rayN y listo.

> 📖 How the fetch → probe → rank pipeline works: [ARCHITECTURE.md](./ARCHITECTURE.md)

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
| 🇺🇸 United States (`US`) | 21 | [clash-US.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-US.yaml) | [singbox-US.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-US.json) | [v2ray-base64-US.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-US.txt) |

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

- **Nodos seleccionados**: 150
- **Activos en todas las fuentes**: 2590
- **RTT del nodo más rápido**: 17 ms
- **RTT mediana**: 17 ms
- **Última actualización (UTC)**: 2026-04-19 10:23 UTC

**Mezcla de protocolos:** shadowsocks × 26 · trojan × 22 · vless × 86 · vmess × 16

**Fuentes usadas en esta ejecución:** `barry-far-v2ray` × 36 · `ebrasha-v2ray` × 15 · `epodonios` × 28 · `freefq` × 1 · `lagzian-mix` × 2 · `mahdi0024` × 2 · `mahdibland-shadowsocks` × 2 · `matin-v2ray` × 1 · `mfuu-clash` × 2 · `ninjastrikers` × 45 · `ruking-clash` × 11 · `snakem982` × 3 · `surfboard-eternity` × 2

## ❓ Preguntas frecuentes

<details><summary>¿Es realmente gratis?</summary>

Sí. Los nodos son operados por voluntarios externos que publican sus propias suscripciones gratuitas. Nosotros no operamos ningún servidor — solo probamos, clasificamos y reempaquetamos lo que ya es público.

</details>

<details><summary>¿Qué tan reciente es la información?</summary>

Una GitHub Action se ejecuta cada hora: trae todas las fuentes, prueba cada nodo con TCP+TLS, elimina los muertos, ordena por latencia y comitea los archivos nuevos. Consulta la marca de tiempo `Last updated` arriba.

</details>

<details><summary>¿Puedo confiar en estos nodos?</summary>

Los nodos gratis ven todo tu tráfico. **Nunca los uses para banca, login o algo sensible.** Bien para saltar bloqueos geográficos en contenido público. Usa tu propio VPS / proveedor de pago para privacidad real.

</details>

<details><summary>¿Por qué algunos nodos listados fallan?</summary>

Verificamos accesibilidad TCP y handshake TLS, pero un nodo aún puede tener cuota expirada, ruteo incorrecto o certificado caducado. Prueba varios; el grupo selector te da alternativas.

</details>

## 🤝 Contribuir

¿Conoces una fuente de suscripción pública confiable que deberíamos agregar? Abre un issue con la URL y el formato.

## ⚠️ Aviso legal

Este repositorio agrega configuraciones de proxy **compartidas públicamente** por voluntarios externos. No operamos ningún servidor, no garantizamos disponibilidad ni seguridad, y no somos responsables del uso que hagas. Destinado a uso educativo y de conectividad personal. Cumple con todas las leyes aplicables en tu jurisdicción.

## ⭐ Historia de estrellas

[![Star History Chart](https://api.star-history.com/svg?repos=Au1rxx/free-vpn-subscriptions&type=Date)](https://www.star-history.com/#Au1rxx/free-vpn-subscriptions&Date)

---

Si este proyecto te ayudó, déjale una ⭐ — cada estrella hace más fácil que otros lo encuentren.
