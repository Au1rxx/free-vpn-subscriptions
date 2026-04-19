# Free VPN Subscriptions

[English](./README.md) · [简体中文](./README_CN.md) · [日本語](./README_JA.md) · [한국어](./README_KO.md) · [Español](./README_ES.md) · [Português](./README_PT.md) · **Русский**

<p align="center"><img src="https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/assets/hero.png" alt="Free VPN Subscriptions — hourly-refreshed free VPN subscriptions for Clash, sing-box, v2ray" width="780"></p>

![узлы](https://img.shields.io/badge/узлы-150-brightgreen) ![живые](https://img.shields.io/badge/живые-2597-blue) ![медиана--rtt](https://img.shields.io/badge/медиана--rtt-8ms-orange) ![обновлено](https://img.shields.io/badge/обновлено-2026-04-19_10:41_UTC-informational)

> **Самый простой способ получить рабочий бесплатный VPN — скопируйте ссылку подписки, вставьте в клиент, подключитесь.**  
> Без регистрации. Без оплаты. Без установки каких-либо бинарников. Обновляется каждый час из публичных источников, каждый узел протестирован.

> бесплатный VPN · бесплатная подписка VPN · бесплатный прокси · Clash подписка · v2ray подписка · sing-box подписка · VLESS · Reality · VMess · Trojan · Shadowsocks · Hysteria2 · обновление каждый час · TCP+TLS проверка · по стране

## 💡 Зачем этот проект?

Каждый список "бесплатных VPN" на GitHub либо устаревший, либо полон мёртвых узлов, либо требует установить подозрительный бинарник. Этот репозиторий **публикует только узлы, прошедшие TCP handshake И TLS handshake несколько минут назад**, из отобранных публичных источников, отсортированных по задержке. Вы получаете 3 переносимых файла подписки — вставьте их в Clash, sing-box или v2rayN и готово.

> 📖 How the fetch → probe → rank pipeline works: [ARCHITECTURE.md](./ARCHITECTURE.md)

## 🚀 Подписка в один клик

Скопируйте URL, соответствующий вашему клиенту, и вставьте его в поле импорта подписки:

| Клиент | Формат | URL подписки |
|---|---|---|
| Clash / Clash Verge / ClashX | `clash.yaml` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/clash.yaml` |
| sing-box | `singbox.json` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/singbox.json` |
| v2rayN / v2rayNG / Shadowrocket / NekoBox | `v2ray-base64` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/v2ray-base64.txt` |

## 🌍 По странам

Нужны узлы только в определённом регионе? Используйте одну из целевых URL подписок:

| Страна | Узлов | Clash | sing-box | v2ray |
|---|---|---|---|---|
| 🇺🇸 United States (`US`) | 20 | [clash-US.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-US.yaml) | [singbox-US.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-US.json) | [v2ray-base64-US.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-US.txt) |
| 🇩🇪 Germany (`DE`) | 6 | [clash-DE.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-DE.yaml) | [singbox-DE.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-DE.json) | [v2ray-base64-DE.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-DE.txt) |
| 🇸🇪 Sweden (`SE`) | 3 | [clash-SE.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-SE.yaml) | [singbox-SE.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-SE.json) | [v2ray-base64-SE.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-SE.txt) |

## 📖 Пошаговые инструкции

Впервые настраиваете VPN-клиент? Выберите платформу и следуйте инструкции:

- [**Clash Verge**](https://au1rxx.github.io/free-vpn-subscriptions/guides/clash-verge.html) · Windows / macOS / Linux
- [**v2rayNG**](https://au1rxx.github.io/free-vpn-subscriptions/guides/v2rayng.html) · Android
- [**Shadowrocket**](https://au1rxx.github.io/free-vpn-subscriptions/guides/shadowrocket.html) · iOS / iPadOS
- [**sing-box**](https://au1rxx.github.io/free-vpn-subscriptions/guides/sing-box.html) · Windows / macOS / Linux / iOS / Android

## 🧩 Поддерживаемые клиенты

- **Windows**: v2rayN, Clash Verge, Hiddify, NekoRay
- **macOS**: ClashX Pro, Clash Verge, sing-box, Hiddify
- **iOS**: Shadowrocket, Stash, Loon, sing-box, Hiddify
- **Android**: v2rayNG, NekoBox, Clash Meta for Android, Hiddify, sing-box
- **Linux**: mihomo (Clash.Meta), sing-box, v2ray-core

## 📊 Статистика в реальном времени

- **Выбрано узлов**: 150
- **Живых во всех источниках**: 2597
- **RTT самого быстрого узла**: 2 ms
- **Медиана RTT**: 8 ms
- **Последнее обновление (UTC)**: 2026-04-19 10:41 UTC

**Распределение протоколов:** shadowsocks × 29 · trojan × 17 · vless × 77 · vmess × 27

**Источники в этом запуске:** `barry-far-v2ray` × 26 · `ebrasha-v2ray` × 15 · `epodonios` × 33 · `mahdibland-aggregator` × 3 · `mahdibland-shadowsocks` × 2 · `matin-v2ray` × 2 · `mfuu-clash` × 1 · `ninjastrikers` × 39 · `ruking-clash` × 18 · `snakem982` × 6 · `surfboard-eternity` × 4 · `vxiaov-clash` × 1

## ❓ Часто задаваемые вопросы

<details><summary>Это правда бесплатно?</summary>

Да. Узлы управляются сторонними волонтёрами, которые сами публикуют свои бесплатные подписки. Мы не управляем никакими серверами — только тестируем, ранжируем и переупаковываем то, что уже публично.

</details>

<details><summary>Насколько свежие данные?</summary>

GitHub Action запускается каждый час: получает все источники, проводит TCP+TLS проверку каждого узла, отбрасывает мёртвые, сортирует по задержке и коммитит новые файлы. Смотрите метку `Last updated` выше.

</details>

<details><summary>Можно ли доверять этим узлам?</summary>

Бесплатные узлы видят весь ваш трафик. **Никогда не используйте их для банкинга, логинов или чего-то чувствительного.** Подходит для обхода гео-блокировок на публичном контенте. Для реальной приватности используйте свой VPS / платный сервис.

</details>

<details><summary>Почему некоторые узлы из списка не работают?</summary>

Мы проверяем TCP доступность и TLS handshake, но у узла всё равно могут быть исчерпанные квоты, неверная маршрутизация или просроченный сертификат. Попробуйте несколько; selector группа даёт альтернативы.

</details>

## 🤝 Участие

Знаете надёжный публичный источник подписок, который стоит добавить? Откройте issue с URL и форматом.

## ⚠️ Отказ от ответственности

Этот репозиторий агрегирует **публично доступные** конфигурации прокси от сторонних волонтёров. Мы не управляем никакими серверами, не гарантируем доступность или безопасность и не несём ответственности за использование. Предназначено для образовательных и личных целей подключения. Соблюдайте все применимые законы вашей юрисдикции.

## ⭐ История звёзд

[![Star History Chart](https://api.star-history.com/svg?repos=Au1rxx/free-vpn-subscriptions&type=Date)](https://www.star-history.com/#Au1rxx/free-vpn-subscriptions&Date)

---

Если этот проект вам помог, поставьте ⭐ — каждая звезда помогает другим найти его легче.
