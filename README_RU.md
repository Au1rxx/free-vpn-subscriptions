# Free VPN Subscriptions

[English](./README.md) · [简体中文](./README_CN.md) · [日本語](./README_JA.md) · [한국어](./README_KO.md) · [Español](./README_ES.md) · [Português](./README_PT.md) · **Русский**

<p align="center"><img src="https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/assets/hero.png" alt="Free VPN Subscriptions — hourly-refreshed free VPN subscriptions for Clash, sing-box, v2ray" width="780"></p>

![узлы](https://img.shields.io/badge/узлы-150-brightgreen) ![живые](https://img.shields.io/badge/живые-2599-blue) ![медиана--rtt](https://img.shields.io/badge/медиана--rtt-9ms-orange) ![обновлено](https://img.shields.io/badge/обновлено-2026-04-19_11:06_UTC-informational)

> **Самый простой способ получить рабочий бесплатный VPN — скопируйте ссылку подписки, вставьте в клиент, подключитесь.**  
> Без регистрации. Без оплаты. Без установки каких-либо бинарников. Обновляется каждый час из публичных источников — перед публикацией каждый узел проверяется по TCP + TLS.

> бесплатный VPN · бесплатная подписка VPN · бесплатный прокси · Clash подписка · v2ray подписка · sing-box подписка · VLESS · Reality · VMess · Trojan · Shadowsocks · Hysteria2 · обновление каждый час · TCP+TLS проверка · по стране

## 💡 Зачем этот проект?

Каждый список "бесплатных VPN" на GitHub либо устаревший, либо полон мёртвых узлов, либо требует установить подозрительный бинарник. Этот репозиторий **публикует только узлы, прошедшие TCP handshake И TLS handshake несколько минут назад**, из отобранных публичных источников, отсортированных по задержке. Вы получаете 3 переносимых файла подписки — вставьте их в Clash, sing-box или v2rayN и готово.

> 📖 How the fetch → probe → rank pipeline works: [ARCHITECTURE.md](./ARCHITECTURE.md)

## 🔬 Как мы проверяем, что узлы действительно работают

**Сразу честный ответ: мы не можем *гарантировать*, что узел пропустит ваш трафик.** Ни один агрегатор не может этого, не пропустив через узел реальный трафик. Ниже подробно: что мы проверяем на этапе агрегации, чего не можем, и откуда реально берётся гарантия.

### ✅ Что мы проверяем при агрегации (перед публикацией)

1. **Доступность TCP** — открываем TCP-соединение к каждому `server:port`. Мёртвые хосты, неверный DNS, заблокированные порты отбрасываются. Отбрасывается примерно 40 % исходных записей.
2. **TLS-handshake** — для каждого TLS / Reality / WS-TLS узла выполняем полный handshake. Просроченные сертификаты, несовпадения SNI, сломанные Reality short-id отбрасываются. Ещё ~10 % отбрасывается.
3. **Сортировка по задержке** — выжившие узлы сортируются по RTT, публикуются top N.

Типичные цифры последнего запуска: **17 источников → ~4,800 сырых → ~2,900 живых по TCP → ~2,600 OK по TLS → top 200 опубликовано**.

### ❌ Что мы проверить не можем

- Аутентификацию прокси-протокола. Неверный UUID / пароль отвергается уже *после* TLS-handshake со стороны upstream-сервера, и мы этого не видим.
- Реальный успех HTTP через прокси.
- Пропускную способность / throughput.
- Геолокацию точнее, чем GeoIP по выходному IP.

### 🛡️ Проверка во время работы — реальная гарантия приходит отсюда

В публикуемом `clash.yaml` есть группа `url-test`, которая **каждые 5 минут проверяет реальный HTTP через каждый узел**:

```yaml
proxy-groups:
  - name: AUTO
    type: url-test
    url: http://www.gstatic.com/generate_204
    interval: 300
```

Ваш клиент сортирует список узлов по *реальной* задержке HTTP через прокси и автоматически выбирает самый быстрый рабочий узел. В sing-box и v2ray есть аналогичные механизмы. Если выбранный узел умирает, клиент без вмешательства переключается на следующий.

### 🧮 Что это значит на практике

Из top 200, публикуемых каждый запуск, типичный клиент находит 30-50 узлов, стабильно пропускающих HTTP в любой момент. Если один замедлился — группа url-test переключает на следующий в один клик.

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
| 🇺🇸 United States (`US`) | 19 | [clash-US.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-US.yaml) | [singbox-US.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-US.json) | [v2ray-base64-US.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-US.txt) |
| 🇩🇪 Germany (`DE`) | 4 | [clash-DE.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-DE.yaml) | [singbox-DE.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-DE.json) | [v2ray-base64-DE.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-DE.txt) |

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
- **Живых во всех источниках**: 2599
- **RTT самого быстрого узла**: 5 ms
- **Медиана RTT**: 9 ms
- **Последнее обновление (UTC)**: 2026-04-19 11:06 UTC

**Распределение протоколов:** shadowsocks × 22 · trojan × 24 · vless × 75 · vmess × 29

**Источники в этом запуске:** `barry-far-v2ray` × 38 · `ebrasha-v2ray` × 11 · `epodonios` × 36 · `lagzian-mix` × 2 · `mahdibland-aggregator` × 2 · `mahdibland-shadowsocks` × 1 · `matin-v2ray` × 3 · `mfuu-clash` × 5 · `ninjastrikers` × 26 · `pawdroid` × 1 · `ruking-clash` × 19 · `snakem982` × 1 · `surfboard-eternity` × 3 · `vxiaov-clash` × 2

## ❓ Часто задаваемые вопросы

<details><summary>Это правда бесплатно?</summary>

Да. Узлы управляются сторонними волонтёрами, которые сами публикуют свои бесплатные подписки. Мы не управляем никакими серверами — только тестируем, ранжируем и переупаковываем то, что уже публично.

</details>

<details><summary>Насколько свежие данные?</summary>

Каждый час (с небольшой случайной задержкой, чтобы не бить по upstream строго в `:00`): получает все источники, проводит TCP+TLS проверку каждого узла, отбрасывает мёртвые, сортирует по задержке и публикует новые файлы. Смотрите метку `Last updated` выше.

</details>

<details><summary>Можно ли доверять этим узлам?</summary>

Бесплатные узлы видят весь ваш трафик. **Никогда не используйте их для банкинга, логинов или чего-то чувствительного.** Подходит для обхода гео-блокировок на публичном контенте. Для реальной приватности используйте свой VPS / платный сервис.

</details>

<details><summary>Почему некоторые узлы из списка не работают?</summary>

Мы проверяем только TCP доступность и TLS handshake — у узла всё равно могут быть исчерпанные квоты, неверная маршрутизация или просроченный сертификат. В публикуемом `clash.yaml` есть группа `url-test` (`http://www.gstatic.com/generate_204`, интервал 300 с), клиент сам выбирает самый быстрый узел, реально пропускающий HTTP. Умер — берите следующий.

</details>

## 🤝 Участие

Знаете надёжный публичный источник подписок, который стоит добавить? Откройте issue с URL и форматом.

## ⚠️ Отказ от ответственности

Этот репозиторий агрегирует **публично доступные** конфигурации прокси от сторонних волонтёров. Мы не управляем никакими серверами, не гарантируем доступность или безопасность и не несём ответственности за использование. Предназначено для образовательных и личных целей подключения. Соблюдайте все применимые законы вашей юрисдикции.

## ⭐ История звёзд

[![Star History Chart](https://api.star-history.com/svg?repos=Au1rxx/free-vpn-subscriptions&type=Date)](https://www.star-history.com/#Au1rxx/free-vpn-subscriptions&Date)

---

Если этот проект вам помог, поставьте ⭐ — каждая звезда помогает другим найти его легче.
