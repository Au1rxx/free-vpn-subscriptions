# Free VPN Subscriptions

[English](./README.md) · [简体中文](./README_CN.md) · [日本語](./README_JA.md) · [한국어](./README_KO.md) · [Español](./README_ES.md) · [Português](./README_PT.md) · **Русский**

<p align="center"><img src="https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/assets/hero.png" alt="Free VPN Subscriptions — hourly-refreshed free VPN subscriptions for Clash, sing-box, v2ray" width="780"></p>

![узлы](https://img.shields.io/badge/узлы-150-brightgreen) ![живые](https://img.shields.io/badge/живые-2073-blue) ![медиана--rtt](https://img.shields.io/badge/медиана--rtt-102ms-orange) ![обновлено](https://img.shields.io/badge/обновлено-2026-05-02_02:30_UTC-informational)

> **Самый простой способ получить рабочий бесплатный VPN — скопируйте ссылку подписки, вставьте в клиент, подключитесь.**  
> Без регистрации. Без оплаты. Без установки каких-либо бинарников. Обновляется каждый час из публичных источников — каждый публикуемый узел несколько минут назад реально пропустил HTTP-трафик через sing-box.

> бесплатный VPN · бесплатная подписка VPN · бесплатный прокси · Clash подписка · v2ray подписка · sing-box подписка · VLESS · Reality · VMess · Trojan · Shadowsocks · Hysteria2 · обновление каждый час · HTTP через прокси проверено · по стране

## 💡 Зачем этот проект?

Каждый список "бесплатных VPN" на GitHub либо устаревший, либо полон мёртвых узлов, либо требует установить подозрительный бинарник. Этот репозиторий идёт на шаг дальше любого другого списка — **мы не просто проверяем, что порт отвечает; мы прогоняем реальный HTTP-трафик через узел с помощью sing-box и убеждаемся, что возвращается 204**, всё это за минуты до публикации. Вы получаете 3 переносимых файла подписки — вставьте их в Clash, sing-box или v2rayN и готово.

> 📖 How the fetch → probe → rank pipeline works: [ARCHITECTURE.md](./ARCHITECTURE.md)

## 🔬 Как мы проверяем, что узлы действительно работают

Большинство списков бесплатных VPN останавливаются на \"TCP-порт открыт\" и публикуют. Мы — нет. Ниже полная пайплайн-проверка, которую узел должен пройти, прежде чем попасть в подписку.

### ✅ Что мы проверяем при агрегации (перед публикацией)

1. **Доступность TCP** — открываем TCP-соединение к каждому `server:port`. Мёртвые хосты, неверный DNS, заблокированные порты отбрасываются. Отсюда отсеивается примерно 40 % исходных записей.
2. **TLS-handshake** — для каждого TLS / Reality / WS-TLS узла выполняем полный handshake. Просроченные сертификаты, несовпадения SNI, сломанные Reality short-id отбрасываются. Ещё ~10 % отсеивается.
3. **Валидация конфига sing-box** — каждый выживший узел переводится в реальный outbound sing-box и проходит `sing-box check`. Битые cipher, неправильные UUID, неподдерживаемые flow-опции отбрасываются до того, как займут слот проверки.
4. **HTTP через прокси (это самое важное)** — самые быстрые ~900 кандидатов пакетно загружаются в sing-box-подпроцессы, каждому узлу даётся свой локальный SOCKS5 inbound, и через него выполняются реальные HTTP и HTTPS GET-запросы:
   - `http://www.gstatic.com/generate_204` (ожидается 204)
   - `https://www.cloudflare.com/cdn-cgi/trace` (ожидается 200)

   Запрос проходит полностью через сам прокси-протокол (VLESS / VMess / Trojan / Shadowsocks / Hysteria2), так что узел, прошедший проверку, на деле имеет работающую аутентификацию, маршрутизацию, внутренний TLS handshake и выходную сеть.
5. **Два раунда, 45 секунд между ними** — узлы, прошедшие один раз и умершие через 45 секунд, фильтруются. Остаются только узлы с коэффициентом успеха ≥ 50 % по (раунды × цели).
6. **Сортировка по медиане реальной задержки** — выжившие сортируются по медиане HTTP-through-proxy round-trip (не по сырой TCP RTT), и публикуются top N.

Типичные цифры последнего запуска: **17 источников → ~4,800 сырых → ~2,900 живых по TCP → ~2,600 OK по TLS → ~840 с валидным конфигом → ~280 прошедших HTTP-проверку → top 150 опубликовано**. Каждый из 150 реально пропустил трафик за последние десять минут.

### ❌ Чего мы всё ещё не можем проверить

- **Пропускную способность / throughput** — мы измеряем задержку, а не мегабиты. Узел с 50 ms может быть всё ещё медленным для видео.
- **Точность геолокации** — GeoIP говорит про страну выходного IP, но не надёжен на уровне города или ISP.
- **Региональные блокировки** — узел, работающий с нашей инфраструктуры проверки, может быть заблокирован для вас (ISP-фильтрация, captive portal и т.п.).
- **Останется ли узел живым после запуска** — узел прошёл десять минут назад; с тех пор он мог умереть.

### 🛡️ Страховка в runtime — для последнего пункта выше

Публикуемый `clash.yaml` включает группу `url-test`, которая перепроверяет реальный HTTP через каждый узел каждые 5 минут на *вашем* устройстве:

```yaml
proxy-groups:
  - name: AUTO
    type: url-test
    url: http://www.gstatic.com/generate_204
    interval: 300
```

Ваш клиент держит список узлов отсортированным по *живой* задержке HTTP через прокси из вашей сети и автоматически выбирает самый быстрый рабочий узел. В sing-box и v2ray есть аналогичные механизмы. Если выбранный узел умирает между часовыми агрегациями, клиент переключается на следующий без вмешательства.

### 🧮 Что это значит на практике

Из ~150, публикуемых каждый запуск, типичный клиент находит **80-120 узлов, чисто пропускающих HTTP из его сети** в любой момент — примерно в 2-3 раза выше, чем у списков, делающих только TCP-проверку. Группа url-test прозрачно ротирует, если один выпал.

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
| 🇺🇸 United States (`US`) | 36 | [clash-US.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-US.yaml) | [singbox-US.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-US.json) | [v2ray-base64-US.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-US.txt) |
| 🇨🇦 Canada (`CA`) | 8 | [clash-CA.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-CA.yaml) | [singbox-CA.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-CA.json) | [v2ray-base64-CA.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-CA.txt) |
| 🇩🇪 Germany (`DE`) | 5 | [clash-DE.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-DE.yaml) | [singbox-DE.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-DE.json) | [v2ray-base64-DE.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-DE.txt) |
| 🇯🇵 Japan (`JP`) | 5 | [clash-JP.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-JP.yaml) | [singbox-JP.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-JP.json) | [v2ray-base64-JP.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-JP.txt) |
| 🇬🇧 United Kingdom (`GB`) | 4 | [clash-GB.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-GB.yaml) | [singbox-GB.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-GB.json) | [v2ray-base64-GB.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-GB.txt) |

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
- **Живых во всех источниках**: 2073
- **RTT самого быстрого узла**: 28 ms
- **Медиана RTT**: 102 ms
- **Последнее обновление (UTC)**: 2026-05-02 02:30 UTC

**Распределение протоколов:** shadowsocks × 38 · trojan × 8 · vless × 94 · vmess × 10

**Источники в этом запуске:** `barry-far-v2ray` × 20 · `ebrasha-v2ray` × 1 · `epodonios` × 4 · `lagzian-mix` × 3 · `mahdi0024` × 32 · `mahdibland-aggregator` × 9 · `mahdibland-shadowsocks` × 8 · `matin-v2ray` × 1 · `ninjastrikers` × 61 · `pawdroid` × 3 · `surfboard-eternity` × 8

## ❓ Часто задаваемые вопросы

<details><summary>Это правда бесплатно?</summary>

Да. Узлы управляются сторонними волонтёрами, которые сами публикуют свои бесплатные подписки. Мы не управляем никакими серверами — только тестируем, ранжируем и переупаковываем то, что уже публично.

</details>

<details><summary>Насколько свежие данные?</summary>

Каждый час (с небольшой случайной задержкой, чтобы не бить по upstream строго в `:00`): получает все источники → TCP → TLS → валидация конфига sing-box → HTTP через прокси (два раунда, 45 секунд между ними) → сортирует по реальной HTTP-задержке → публикует новые файлы. Полный пайплайн ~10 минут. Смотрите метку `Last updated` выше.

</details>

<details><summary>Можно ли доверять этим узлам?</summary>

Бесплатные узлы видят весь ваш трафик. **Никогда не используйте их для банкинга, логинов или чего-то чувствительного.** Подходит для обхода гео-блокировок на публичном контенте. Для реальной приватности используйте свой VPS / платный сервис.

</details>

<details><summary>Почему некоторые узлы из списка не работают?</summary>

Даже после нашей HTTP-проверки через прокси узлы могут умирать между агрегациями: квота исчерпана, upstream отозвал ключ, ваш ISP блокирует выходной IP, или оператор закрыл. В публикуемом `clash.yaml` есть группа `url-test` (`http://www.gstatic.com/generate_204`, интервал 300 с), клиент сам выбирает самый быстрый узел, реально пропускающий HTTP *из вашей сети*. Умер — берите следующий. Ожидайте, что 80-120 из 150 работают в любой момент.

</details>

## 🤝 Участие

Знаете надёжный публичный источник подписок, который стоит добавить? Откройте issue с URL и форматом.

## ⚠️ Отказ от ответственности

Этот репозиторий агрегирует **публично доступные** конфигурации прокси от сторонних волонтёров. Мы не управляем никакими серверами, не гарантируем доступность или безопасность и не несём ответственности за использование. Предназначено для образовательных и личных целей подключения. Соблюдайте все применимые законы вашей юрисдикции.

## ⭐ История звёзд

[![Star History Chart](https://api.star-history.com/svg?repos=Au1rxx/free-vpn-subscriptions&type=Date)](https://www.star-history.com/#Au1rxx/free-vpn-subscriptions&Date)

---

Если этот проект вам помог, поставьте ⭐ — каждая звезда помогает другим найти его легче.
