# Free VPN Subscriptions

[English](./README.md) · [简体中文](./README_CN.md) · [日本語](./README_JA.md) · **한국어** · [Español](./README_ES.md) · [Português](./README_PT.md) · [Русский](./README_RU.md)

<p align="center"><img src="https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/assets/hero.png" alt="Free VPN Subscriptions — hourly-refreshed free VPN subscriptions for Clash, sing-box, v2ray" width="780"></p>

![노드](https://img.shields.io/badge/노드-150-brightgreen) ![생존](https://img.shields.io/badge/생존-2530-blue) ![중앙값--rtt](https://img.shields.io/badge/중앙값--rtt-121ms-orange) ![업데이트](https://img.shields.io/badge/업데이트-2026-04-20_12:21_UTC-informational)

> **작동하는 무료 VPN을 얻는 가장 쉬운 방법 —— 구독 링크를 복사하고 클라이언트에 붙여 넣고 연결하세요.**  
> 가입 불필요. 결제 불필요. 바이너리 설치 불필요. 공개 소스에서 매시간 자동 갱신 —— 발행되는 모든 노드는 몇 분 전에 sing-box 를 통해 실제 HTTP 트래픽을 전달한 이력이 있습니다.

> 무료 VPN 구독 · 무료 v2ray 구독 · 무료 Clash 구독 · 무료 sing-box 구독 · VLESS · Reality · VMess · Trojan · Shadowsocks · Hysteria2 · 매시간 갱신 · HTTP 실측 검증 · 국가별

## 💡 왜 이 프로젝트?

GitHub의 거의 모든 "무료 VPN" 목록은 데이터가 오래되었거나, 죽은 노드로 가득 차 있거나, 출처가 불분명한 바이너리 설치를 요구합니다. 이 저장소는 다른 어떤 곳보다 한 단계 더 나아갑니다 —— **단지 포트가 열려 있는지 확인하는 것이 아니라, sing-box 로 실제 HTTP 트래픽을 노드를 거쳐 보내고 204 가 돌아오는지 확인한 뒤에만 발행합니다**, 모두 수 분 내에. Clash / sing-box / v2rayN에 바로 붙여 넣을 수 있는 3가지 범용 구독 파일을 제공합니다.

> 📖 How the fetch → probe → rank pipeline works: [ARCHITECTURE.md](./ARCHITECTURE.md)

## 🔬 노드가 실제로 작동하는지 어떻게 검증하나

대부분의 무료 VPN 목록은 "TCP 포트가 열려 있다"에서 멈추고 발행합니다. 우리는 다릅니다. 아래는 노드가 구독에 포함되기 전에 통과해야 하는 전체 검증 파이프라인입니다.

### ✅ 집계 단계 (발행 전) 에서 검증하는 것

1. **TCP 도달성** —— 모든 `server:port` 에 TCP 연결을 시도합니다. 죽은 호스트, 잘못된 DNS, 차단된 포트는 모두 드롭. 원시 데이터의 약 40 % 가 여기서 제거됩니다.
2. **TLS 핸드셰이크** —— TLS / Reality / WS-TLS 노드에 대해 전체 핸드셰이크를 수행합니다. 만료된 인증서, SNI 불일치, 손상된 Reality short-id 는 드롭. 추가로 약 10 % 가 제거됩니다.
3. **sing-box 설정 검증** —— 생존한 각 노드는 실제 sing-box outbound 로 변환되어 `sing-box check` 를 통과합니다. 손상된 암호화 방식, 잘못된 UUID, 지원되지 않는 flow 옵션은 프로브 자원을 소모하기 전에 제거됩니다.
4. **HTTP-over-proxy 프로브 (핵심 단계)** —— 가장 빠른 약 900 개 후보를 sing-box 서브프로세스에 배치 투입하고, 각 노드에 로컬 SOCKS5 인바운드를 할당하여 실제 HTTP 및 HTTPS 요청을 전송합니다:
   - `http://www.gstatic.com/generate_204` (204 기대)
   - `https://www.cloudflare.com/cdn-cgi/trace` (200 기대)

   요청은 실제 프록시 프로토콜 (VLESS / VMess / Trojan / Shadowsocks / Hysteria2) 을 완전히 통과하므로, 이 단계를 통과한 노드는 인증, 라우팅, 내부 TLS 핸드셰이크, 출구 네트워크가 모두 작동함이 입증됩니다.
5. **2 라운드, 45 초 간격** —— 한 번 통과했지만 45 초 후 죽는 노드는 걸러집니다. (라운드 × 타겟) 중 성공률 50 % 이상인 노드만 남습니다.
6. **실제 레이턴시 중앙값 정렬** —— 생존 노드는 HTTP-over-proxy 실측 왕복 시간 중앙값 (원시 TCP RTT 아님) 으로 정렬되어 상위 N 개를 발행합니다.

최근 실행의 전형적인 수치: **17 개 소스 → ~4,800 원시 → ~2,900 TCP 생존 → ~2,600 TLS OK → ~840 설정 유효 → ~280 HTTP 실측 통과 → 상위 150 발행**. 발행되는 150 개 노드 모두 지난 10 분 내에 실제로 트래픽을 전달한 실적이 있습니다.

### ❌ 그래도 검증할 수 없는 것

- **대역폭 / 처리량** —— 우리는 레이턴시를 측정하며 Mbps 는 아닙니다. 50ms 노드라도 동영상에서는 느릴 수 있습니다.
- **정확한 지리 위치** —— GeoIP 는 출구 IP 의 국가는 알려주지만 도시나 ISP 단위로는 신뢰할 수 없습니다.
- **지역별 차단** —— 우리 프로브에서 통하는 노드가 당신의 네트워크에서도 통한다는 보장은 없습니다 (ISP 레벨 필터링, captive portal 등).
- **발행 이후의 생존** —— 10 분 전에는 살아 있었지만 그 후 죽었을 수 있습니다.

### 🛡️ 런타임 안전망 —— 위 마지막 항목 대응

발행하는 `clash.yaml` 에는 `url-test` 프록시 그룹이 포함되어 있으며, **당신의 기기에서** 5 분마다 각 노드에 실제 HTTP 를 다시 테스트합니다:

```yaml
proxy-groups:
  - name: AUTO
    type: url-test
    url: http://www.gstatic.com/generate_204
    interval: 300
```

클라이언트는 *당신 네트워크의* 실시간 HTTP-over-proxy 레이턴시로 노드를 정렬하여 가장 빠른 작동 노드를 자동 선택합니다. sing-box / v2ray 에도 동등한 메커니즘이 있습니다. 매시간 집계 사이에 노드가 죽으면 클라이언트가 개입 없이 다음으로 전환합니다.

### 🧮 실제 기대치

발행되는 ~150 개 노드 중, 클라이언트는 보통 **80-120 개가 당신의 네트워크에서 HTTP 를 통과시키는 것을** 찾아냅니다 —— TCP 프로브만 하는 목록보다 2-3 배 높은 명중률입니다. 하나가 느려지면 url-test 그룹이 투명하게 로테이션합니다.

## 🚀 원클릭 구독

클라이언트에 맞는 URL을 복사하여 구독 가져오기 필드에 붙여 넣으세요:

| 클라이언트 | 형식 | 구독 URL |
|---|---|---|
| Clash / Clash Verge / ClashX | `clash.yaml` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/clash.yaml` |
| sing-box | `singbox.json` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/singbox.json` |
| v2rayN / v2rayNG / Shadowrocket / NekoBox | `v2ray-base64` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/v2ray-base64.txt` |

## 🌍 국가별 구독

특정 지역의 노드만 필요하신가요? 전용 구독 URL을 선택하세요:

| 국가 | 노드 수 | Clash | sing-box | v2ray |
|---|---|---|---|---|
| 🇺🇸 United States (`US`) | 40 | [clash-US.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-US.yaml) | [singbox-US.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-US.json) | [v2ray-base64-US.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-US.txt) |
| 🇩🇪 Germany (`DE`) | 7 | [clash-DE.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-DE.yaml) | [singbox-DE.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-DE.json) | [v2ray-base64-DE.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-DE.txt) |
| 🇨🇦 Canada (`CA`) | 6 | [clash-CA.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-CA.yaml) | [singbox-CA.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-CA.json) | [v2ray-base64-CA.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-CA.txt) |
| 🇬🇧 United Kingdom (`GB`) | 3 | [clash-GB.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-GB.yaml) | [singbox-GB.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-GB.json) | [v2ray-base64-GB.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-GB.txt) |

## 📖 클라이언트 설정 가이드

처음이신가요? 플랫폼에 맞는 튜토리얼을 따라 해보세요:

- [**Clash Verge**](https://au1rxx.github.io/free-vpn-subscriptions/guides/clash-verge.html) · Windows / macOS / Linux
- [**v2rayNG**](https://au1rxx.github.io/free-vpn-subscriptions/guides/v2rayng.html) · Android
- [**Shadowrocket**](https://au1rxx.github.io/free-vpn-subscriptions/guides/shadowrocket.html) · iOS / iPadOS
- [**sing-box**](https://au1rxx.github.io/free-vpn-subscriptions/guides/sing-box.html) · Windows / macOS / Linux / iOS / Android

## 🧩 지원 클라이언트

- **Windows**: v2rayN, Clash Verge, Hiddify, NekoRay
- **macOS**: ClashX Pro, Clash Verge, sing-box, Hiddify
- **iOS**: Shadowrocket, Stash, Loon, sing-box, Hiddify
- **Android**: v2rayNG, NekoBox, Clash Meta for Android, Hiddify, sing-box
- **Linux**: mihomo (Clash.Meta), sing-box, v2ray-core

## 📊 실시간 통계

- **선정된 노드**: 150
- **전체 소스 생존 수**: 2530
- **최고 속도 RTT**: 28 ms
- **중앙값 RTT**: 121 ms
- **최종 업데이트 (UTC)**: 2026-04-20 12:21 UTC

**프로토콜 분포:** shadowsocks × 22 · trojan × 15 · vless × 101 · vmess × 12

**이번 실행에 사용된 소스:** `barry-far-v2ray` × 34 · `epodonios` × 10 · `lagzian-mix` × 4 · `mahdi0024` × 26 · `mahdibland-aggregator` × 8 · `mahdibland-shadowsocks` × 2 · `matin-v2ray` × 1 · `ninjastrikers` × 52 · `pawdroid` × 1 · `ruking-clash` × 7 · `surfboard-eternity` × 5

## ❓ 자주 묻는 질문

<details><summary>정말 무료인가요?</summary>

네. 모든 노드는 제3자 자원봉사자가 운영하며 공개 구독을 스스로 게시합니다. 저희는 어떤 서버도 운영하지 않으며, 이미 공개된 것을 테스트하고 순위를 매기고 재포장할 뿐입니다.

</details>

<details><summary>데이터는 얼마나 신선한가요?</summary>

매시간 갱신 (상위 소스를 `:00` 에 집중적으로 때리지 않도록 작은 무작위 지연 포함): 모든 소스 가져오기 → TCP → TLS → sing-box 설정 검증 → HTTP-over-proxy 프로브 (2 라운드, 45 초 간격) → 실 HTTP 레이턴시 정렬 → 새 출력 파일 발행. 전체 파이프라인 약 10 분. 위의 `Last updated` 타임스탬프를 확인하세요.

</details>

<details><summary>이 노드들을 신뢰할 수 있나요?</summary>

무료 노드는 모든 트래픽을 운영자가 볼 수 있습니다. **은행 거래, 로그인, 민감한 작업에는 절대 사용하지 마세요.** 공개 콘텐츠의 지역 제한 우회에는 적합합니다. 실제 프라이버시에는 자체 VPS/유료 서비스를 사용하세요.

</details>

<details><summary>목록에 있는데 작동하지 않는 노드가 있는 이유는?</summary>

HTTP-over-proxy 프로브를 통과한 후에도 노드는 집계 사이에 죽을 수 있습니다: 할당량 소진, 상위 키 취소, ISP 가 출구 IP 차단, 운영자 중단 등. 발행하는 `clash.yaml` 에는 `url-test` 그룹 (`http://www.gstatic.com/generate_204`, 300초 간격) 이 포함되어 있어 클라이언트가 *당신의 네트워크에서* 실제로 HTTP 를 통과시키는 가장 빠른 노드를 자동 선택합니다. 죽으면 다음으로. 150 개 중 80-120 개가 수시로 작동합니다.

</details>

## 🤝 기여

신뢰할 수 있는 공개 구독 소스를 알고 계신가요? URL과 형식을 포함한 이슈를 열어 주세요.

## ⚠️ 면책 조항

이 저장소는 제3자 자원봉사자가 **공개 공유**한 프록시 구성을 집계합니다. 저희는 어떤 서버도 운영하지 않고, 가용성이나 보안을 보장하지 않으며, 사용 방식에 대해 책임지지 않습니다. 교육 및 개인 연결 용도로만 사용하세요. 해당 관할권의 모든 법률을 준수하세요.

## ⭐ 스타 히스토리

[![Star History Chart](https://api.star-history.com/svg?repos=Au1rxx/free-vpn-subscriptions&type=Date)](https://www.star-history.com/#Au1rxx/free-vpn-subscriptions&Date)

---

이 프로젝트가 도움이 되셨다면 ⭐을 남겨 주세요 —— 모든 스타가 다른 사람들이 이 프로젝트를 더 쉽게 발견하도록 도와줍니다.
