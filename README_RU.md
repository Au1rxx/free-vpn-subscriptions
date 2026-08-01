# Free VPN Subscriptions

<div align="center">

[English](./README.md) · [简体中文](./README_CN.md) · [日本語](./README_JA.md) · [한국어](./README_KO.md) · [Español](./README_ES.md) · [Português](./README_PT.md) · **Русский**

</div>

<p align="center"><img src="https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/assets/hero.png" alt="Free VPN Subscriptions — hourly-refreshed free VPN subscriptions for Clash, sing-box, v2ray" width="780"></p>

[![GitHub stars](https://img.shields.io/github/stars/Au1rxx/free-vpn-subscriptions?style=flat&color=gold&logo=github)](https://github.com/Au1rxx/free-vpn-subscriptions/stargazers) ![отобрано](https://img.shields.io/badge/отобрано-2000-brightgreen) ![проверено](https://img.shields.io/badge/проверено-3710-blue) ![медиана--rtt](https://img.shields.io/badge/медиана--rtt-580ms-orange) ![обновлено](https://img.shields.io/badge/обновлено-2026-08-01_11:02_UTC-informational) [![License](https://img.shields.io/github/license/Au1rxx/free-vpn-subscriptions?color=blue)](https://github.com/Au1rxx/free-vpn-subscriptions/blob/main/LICENSE)

> **Самый простой способ получить рабочий бесплатный VPN — скопируйте ссылку подписки, вставьте в клиент, подключитесь.**  
> Без регистрации. Без оплаты. Без установки каких-либо бинарников. Обновляется каждый час из публичных источников — каждый публикуемый узел несколько минут назад реально пропустил HTTP-трафик через sing-box.

> бесплатный VPN · бесплатная подписка VPN · бесплатный прокси · Clash подписка · v2ray подписка · sing-box подписка · VLESS · Reality · VMess · Trojan · Shadowsocks · Hysteria2 · обновление каждый час · HTTP через прокси проверено · по стране

## 🚀 Подписка в один клик

Скопируйте URL, соответствующий вашему клиенту, и вставьте его в поле импорта подписки:

| Клиент | Формат | URL подписки |
|---|---|---|
| Clash / Clash Verge / ClashX | `clash.yaml` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/clash.yaml` |
| sing-box | `singbox.json` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/singbox.json` |
| v2rayN / v2rayNG / Shadowrocket / NekoBox | `v2ray-base64` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/v2ray-base64.txt` |

> Database-classified shards (stable / all verified / protocol / country / network): [https://github.com/Au1rxx/free-vpn-subscriptions/tree/main/output](https://github.com/Au1rxx/free-vpn-subscriptions/tree/main/output)

⭐ **Поставьте звезду репозиторию**, если он сэкономил вам время — это единственный сигнал, по которому мы решаем, что развивать дальше.

## 💡 Зачем этот проект?

Каждый список "бесплатных VPN" на GitHub либо устаревший, либо полон мёртвых узлов, либо требует установить подозрительный бинарник. Этот репозиторий идёт на шаг дальше любого другого списка — **мы не просто проверяем, что порт отвечает; мы прогоняем реальный HTTP-трафик через узел с помощью sing-box и убеждаемся, что возвращается 204**, всё это за минуты до публикации. Вы получаете 3 переносимых файла подписки — вставьте их в Clash, sing-box или v2rayN и готово.

> 📖 How the fetch → probe → rank pipeline works: [ARCHITECTURE.md](./ARCHITECTURE.md)

<details>
<summary><b>🔬 Как мы проверяем, что узлы действительно работают</b></summary>

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

</details>

## 🌍 По странам

Нужны узлы только в определённом регионе? Используйте одну из целевых URL подписок:

| Страна | Узлов | Clash | sing-box | v2ray |
|---|---|---|---|---|
| 🇺🇸 United States (`US`) | 389 | [clash-US.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-US.yaml) | [singbox-US.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-US.json) | [v2ray-base64-US.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-US.txt) |
| 🇫🇷 France (`FR`) | 281 | [clash-FR.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-FR.yaml) | [singbox-FR.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-FR.json) | [v2ray-base64-FR.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-FR.txt) |
| 🇳🇱 Netherlands (`NL`) | 258 | [clash-NL.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-NL.yaml) | [singbox-NL.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-NL.json) | [v2ray-base64-NL.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-NL.txt) |
| 🇩🇪 Germany (`DE`) | 175 | [clash-DE.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-DE.yaml) | [singbox-DE.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-DE.json) | [v2ray-base64-DE.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-DE.txt) |
| 🇭🇰 Hong Kong (`HK`) | 80 | [clash-HK.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-HK.yaml) | [singbox-HK.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-HK.json) | [v2ray-base64-HK.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-HK.txt) |
| 🇸🇬 Singapore (`SG`) | 77 | [clash-SG.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-SG.yaml) | [singbox-SG.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-SG.json) | [v2ray-base64-SG.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-SG.txt) |
| 🇬🇧 United Kingdom (`GB`) | 76 | [clash-GB.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-GB.yaml) | [singbox-GB.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-GB.json) | [v2ray-base64-GB.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-GB.txt) |
| 🇰🇷 Korea (`KR`) | 63 | [clash-KR.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-KR.yaml) | [singbox-KR.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-KR.json) | [v2ray-base64-KR.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-KR.txt) |
| 🇹🇼 Taiwan (`TW`) | 63 | [clash-TW.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-TW.yaml) | [singbox-TW.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-TW.json) | [v2ray-base64-TW.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-TW.txt) |
| 🇨🇭 Switzerland (`CH`) | 62 | [clash-CH.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-CH.yaml) | [singbox-CH.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-CH.json) | [v2ray-base64-CH.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-CH.txt) |
| 🇯🇵 Japan (`JP`) | 52 | [clash-JP.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-JP.yaml) | [singbox-JP.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-JP.json) | [v2ray-base64-JP.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-JP.txt) |
| 🇵🇱 Poland (`PL`) | 41 | [clash-PL.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-PL.yaml) | [singbox-PL.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-PL.json) | [v2ray-base64-PL.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-PL.txt) |
| 🇫🇮 Finland (`FI`) | 38 | [clash-FI.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-FI.yaml) | [singbox-FI.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-FI.json) | [v2ray-base64-FI.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-FI.txt) |
| 🇷🇺 Russia (`RU`) | 31 | [clash-RU.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-RU.yaml) | [singbox-RU.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-RU.json) | [v2ray-base64-RU.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-RU.txt) |
| 🇹🇷 Turkey (`TR`) | 31 | [clash-TR.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-TR.yaml) | [singbox-TR.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-TR.json) | [v2ray-base64-TR.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-TR.txt) |
| 🇨🇦 Canada (`CA`) | 26 | [clash-CA.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-CA.yaml) | [singbox-CA.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-CA.json) | [v2ray-base64-CA.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-CA.txt) |
| 🇨🇳 China (`CN`) | 25 | [clash-CN.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-CN.yaml) | [singbox-CN.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-CN.json) | [v2ray-base64-CN.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-CN.txt) |
| 🇸🇪 Sweden (`SE`) | 25 | [clash-SE.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-SE.yaml) | [singbox-SE.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-SE.json) | [v2ray-base64-SE.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-SE.txt) |
| 🇷🇴 Romania (`RO`) | 19 | [clash-RO.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-RO.yaml) | [singbox-RO.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-RO.json) | [v2ray-base64-RO.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-RO.txt) |
| 🇨🇾 Cyprus (`CY`) | 17 | [clash-CY.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-CY.yaml) | [singbox-CY.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-CY.json) | [v2ray-base64-CY.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-CY.txt) |
| 🇮🇹 Italy (`IT`) | 16 | [clash-IT.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-IT.yaml) | [singbox-IT.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-IT.json) | [v2ray-base64-IT.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-IT.txt) |
| 🇱🇹 Lithuania (`LT`) | 13 | [clash-LT.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-LT.yaml) | [singbox-LT.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-LT.json) | [v2ray-base64-LT.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-LT.txt) |
| 🇪🇪 Estonia (`EE`) | 12 | [clash-EE.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-EE.yaml) | [singbox-EE.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-EE.json) | [v2ray-base64-EE.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-EE.txt) |
| 🇮🇩 Indonesia (`ID`) | 11 | [clash-ID.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-ID.yaml) | [singbox-ID.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-ID.json) | [v2ray-base64-ID.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-ID.txt) |
| 🇪🇸 Spain (`ES`) | 10 | [clash-ES.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-ES.yaml) | [singbox-ES.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-ES.json) | [v2ray-base64-ES.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-ES.txt) |
| 🇹🇭 Thailand (`TH`) | 9 | [clash-TH.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-TH.yaml) | [singbox-TH.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-TH.json) | [v2ray-base64-TH.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-TH.txt) |
| 🇮🇳 India (`IN`) | 7 | [clash-IN.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-IN.yaml) | [singbox-IN.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-IN.json) | [v2ray-base64-IN.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-IN.txt) |
| 🇱🇻 Latvia (`LV`) | 7 | [clash-LV.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-LV.yaml) | [singbox-LV.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-LV.json) | [v2ray-base64-LV.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-LV.txt) |
| 🇹🇖 T1 (`T1`) | 7 | [clash-T1.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-T1.yaml) | [singbox-T1.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-T1.json) | [v2ray-base64-T1.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-T1.txt) |
| 🇧🇬 Bulgaria (`BG`) | 6 | [clash-BG.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-BG.yaml) | [singbox-BG.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-BG.json) | [v2ray-base64-BG.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-BG.txt) |
| 🇧🇷 Brazil (`BR`) | 6 | [clash-BR.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-BR.yaml) | [singbox-BR.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-BR.json) | [v2ray-base64-BR.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-BR.txt) |
| 🇨🇴 Colombia (`CO`) | 5 | [clash-CO.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-CO.yaml) | [singbox-CO.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-CO.json) | [v2ray-base64-CO.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-CO.txt) |
| 🇳🇴 Norway (`NO`) | 5 | [clash-NO.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-NO.yaml) | [singbox-NO.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-NO.json) | [v2ray-base64-NO.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-NO.txt) |
| 🇵🇰 Pakistan (`PK`) | 5 | [clash-PK.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-PK.yaml) | [singbox-PK.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-PK.json) | [v2ray-base64-PK.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-PK.txt) |
| 🇻🇳 Vietnam (`VN`) | 5 | [clash-VN.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-VN.yaml) | [singbox-VN.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-VN.json) | [v2ray-base64-VN.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-VN.txt) |
| 🇧🇩 Bangladesh (`BD`) | 4 | [clash-BD.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-BD.yaml) | [singbox-BD.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-BD.json) | [v2ray-base64-BD.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-BD.txt) |
| 🇮🇪 Ireland (`IE`) | 4 | [clash-IE.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-IE.yaml) | [singbox-IE.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-IE.json) | [v2ray-base64-IE.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-IE.txt) |
| 🇰🇭 Cambodia (`KH`) | 4 | [clash-KH.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-KH.yaml) | [singbox-KH.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-KH.json) | [v2ray-base64-KH.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-KH.txt) |
| 🇦🇺 Australia (`AU`) | 3 | [clash-AU.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-AU.yaml) | [singbox-AU.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-AU.json) | [v2ray-base64-AU.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-AU.txt) |
| 🇮🇸 Iceland (`IS`) | 3 | [clash-IS.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-IS.yaml) | [singbox-IS.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-IS.json) | [v2ray-base64-IS.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-IS.txt) |

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

- **Выбрано узлов**: 2000
- **Живых во всех источниках**: 3710
- **RTT самого быстрого узла**: 36 ms
- **Медиана RTT**: 580 ms
- **Последнее обновление (UTC)**: 2026-08-01 11:02 UTC

**Распределение протоколов:** http × 154 · https × 2 · hysteria2 × 142 · shadowsocks × 194 · socks4 × 31 · socks5 × 138 · trojan × 300 · tuic × 5 · vless × 1033 · vmess × 1

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

[![Star History Chart](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/assets/star-history.svg)](https://github.com/Au1rxx/free-vpn-subscriptions/stargazers)

---

Если этот проект вам помог, поставьте ⭐ — каждая звезда помогает другим найти его легче.
