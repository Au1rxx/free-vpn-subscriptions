# Free VPN Subscriptions

[English](./README.md) · [简体中文](./README_CN.md) · [日本語](./README_JA.md) · **한국어** · [Español](./README_ES.md) · [Português](./README_PT.md) · [Русский](./README_RU.md)

<p align="center"><img src="https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/assets/hero.png" alt="Free VPN Subscriptions — hourly-refreshed free VPN subscriptions for Clash, sing-box, v2ray" width="780"></p>

![노드](https://img.shields.io/badge/노드-150-brightgreen) ![생존](https://img.shields.io/badge/생존-2579-blue) ![중앙값--rtt](https://img.shields.io/badge/중앙값--rtt-8ms-orange) ![업데이트](https://img.shields.io/badge/업데이트-2026-04-19_11:05_UTC-informational)

> **작동하는 무료 VPN을 얻는 가장 쉬운 방법 —— 구독 링크를 복사하고 클라이언트에 붙여 넣고 연결하세요.**  
> 가입 불필요. 결제 불필요. 바이너리 설치 불필요. 공개 소스에서 매시간 자동 갱신 — 발행 전 모든 노드를 TCP + TLS 로 검증.

> 무료 VPN 구독 · 무료 v2ray 구독 · 무료 Clash 구독 · 무료 sing-box 구독 · VLESS · Reality · VMess · Trojan · Shadowsocks · Hysteria2 · 매시간 갱신 · TCP+TLS 프로브 완료 · 국가별

## 💡 왜 이 프로젝트?

GitHub의 거의 모든 "무료 VPN" 목록은 데이터가 오래되었거나, 죽은 노드로 가득 차 있거나, 출처가 불분명한 바이너리 설치를 요구합니다. 이 저장소는 **몇 분 전에 TCP 핸드셰이크와 TLS 핸드셰이크를 모두 통과한 노드만** 선별된 공개 소스에서 레이턴시 순으로 게시합니다. Clash / sing-box / v2rayN에 바로 붙여 넣을 수 있는 3가지 범용 구독 파일을 제공합니다.

> 📖 How the fetch → probe → rank pipeline works: [ARCHITECTURE.md](./ARCHITECTURE.md)

## 🔬 노드가 실제로 작동하는지 어떻게 검증하나

**솔직히 말하면: 어떤 노드가 반드시 트래픽을 통과시킨다고 *보장* 할 수는 없습니다.** 실제로 트래픽을 흘려보내지 않는 한, 어떤 집계 프로젝트도 불가능합니다. 아래에 "집계 단계에서 무엇을 검증하고, 무엇을 못 하며, 진짜 보장은 어디에서 오는지"를 모두 밝힙니다.

### ✅ 집계 단계 (발행 전) 에서 검증하는 것

1. **TCP 도달성** —— 모든 `server:port` 에 TCP 연결을 시도합니다. 죽은 호스트, 잘못된 DNS, 차단된 포트는 모두 드롭. 원시 데이터의 약 40 % 가 여기서 제거됩니다.
2. **TLS 핸드셰이크** —— TLS / Reality / WS-TLS 노드에 대해 전체 핸드셰이크를 수행합니다. 만료된 인증서, SNI 불일치, 손상된 Reality short-id 는 드롭. 추가로 약 10 % 가 제거됩니다.
3. **레이턴시 정렬** —— 생존 노드를 RTT 순으로 정렬하여 상위 N 개를 발행합니다.

최근 실행의 전형적인 수치: **17 개 소스 → ~4,800 원시 → ~2,900 TCP 생존 → ~2,600 TLS OK → 상위 200 발행**.

### ❌ 검증할 수 없는 것

- 프록시 프로토콜 인증. UUID / 비밀번호 불일치는 TLS 핸드셰이크 *후* 에 상위 서버에서 거부되므로 우리에게는 보이지 않습니다.
- 실제 HTTP-over-proxy 성공 여부.
- 대역폭 / 처리량.
- 출구 IP 의 GeoIP 를 넘어선 정확한 지리 정보.

### 🛡️ 런타임 검증 —— 진짜 보장은 여기서

발행하는 `clash.yaml` 에는 `url-test` 프록시 그룹이 포함되어 있으며, **클라이언트가 5 분마다 각 노드에 실제 HTTP** 를 보냅니다:

```yaml
proxy-groups:
  - name: AUTO
    type: url-test
    url: http://www.gstatic.com/generate_204
    interval: 300
```

클라이언트는 *실제* HTTP-over-proxy 레이턴시로 노드를 정렬하여 가장 빠른 작동 노드를 자동 선택합니다. sing-box / v2ray 에도 동등한 메커니즘이 있습니다. 선택된 노드가 죽으면 클라이언트가 개입 없이 다음으로 전환합니다.

### 🧮 실제 기대치

발행되는 상위 200 개 노드 중, 클라이언트는 보통 30-50 개의 실제로 HTTP 를 통과시키는 노드를 찾아냅니다. 느려지면 url-test 그룹이 다음 후보로 전환하여 한 번의 클릭이면 충분합니다.

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
| 🇺🇸 United States (`US`) | 21 | [clash-US.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-US.yaml) | [singbox-US.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-US.json) | [v2ray-base64-US.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-US.txt) |
| 🇩🇪 Germany (`DE`) | 7 | [clash-DE.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-DE.yaml) | [singbox-DE.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-DE.json) | [v2ray-base64-DE.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-DE.txt) |

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
- **전체 소스 생존 수**: 2579
- **최고 속도 RTT**: 2 ms
- **중앙값 RTT**: 8 ms
- **최종 업데이트 (UTC)**: 2026-04-19 11:05 UTC

**프로토콜 분포:** shadowsocks × 25 · trojan × 18 · vless × 85 · vmess × 22

**이번 실행에 사용된 소스:** `barry-far-v2ray` × 30 · `ebrasha-v2ray` × 9 · `epodonios` × 33 · `freefq` × 1 · `mahdi0024` × 1 · `mahdibland-aggregator` × 1 · `mahdibland-shadowsocks` × 1 · `mfuu-clash` × 2 · `ninjastrikers` × 35 · `pawdroid` × 1 · `ruking-clash` × 21 · `snakem982` × 12 · `surfboard-eternity` × 2 · `vxiaov-clash` × 1

## ❓ 자주 묻는 질문

<details><summary>정말 무료인가요?</summary>

네. 모든 노드는 제3자 자원봉사자가 운영하며 공개 구독을 스스로 게시합니다. 저희는 어떤 서버도 운영하지 않으며, 이미 공개된 것을 테스트하고 순위를 매기고 재포장할 뿐입니다.

</details>

<details><summary>데이터는 얼마나 신선한가요?</summary>

매시간 갱신 (상위 소스를 `:00` 에 집중적으로 때리지 않도록 작은 무작위 지연 포함): 모든 상위 소스 가져오기 → 각 노드 TCP+TLS 프로브 → 죽은 것 제거 → 레이턴시 순 정렬 → 새 출력 파일 발행. 위의 `Last updated` 타임스탬프를 확인하세요.

</details>

<details><summary>이 노드들을 신뢰할 수 있나요?</summary>

무료 노드는 모든 트래픽을 운영자가 볼 수 있습니다. **은행 거래, 로그인, 민감한 작업에는 절대 사용하지 마세요.** 공개 콘텐츠의 지역 제한 우회에는 적합합니다. 실제 프라이버시에는 자체 VPS/유료 서비스를 사용하세요.

</details>

<details><summary>목록에 있는데 작동하지 않는 노드가 있는 이유는?</summary>

TCP 도달성과 TLS 핸드셰이크만 검증하므로 노드는 여전히 할당량 소진, 잘못된 라우팅, 만료된 인증서를 가질 수 있습니다. 발행하는 `clash.yaml` 에는 `url-test` 그룹 (`http://www.gstatic.com/generate_204`, 300초 간격) 이 포함되어 있어 클라이언트가 실제로 HTTP 를 통과시키는 가장 빠른 노드를 자동 선택합니다. 죽으면 다음으로.

</details>

## 🤝 기여

신뢰할 수 있는 공개 구독 소스를 알고 계신가요? URL과 형식을 포함한 이슈를 열어 주세요.

## ⚠️ 면책 조항

이 저장소는 제3자 자원봉사자가 **공개 공유**한 프록시 구성을 집계합니다. 저희는 어떤 서버도 운영하지 않고, 가용성이나 보안을 보장하지 않으며, 사용 방식에 대해 책임지지 않습니다. 교육 및 개인 연결 용도로만 사용하세요. 해당 관할권의 모든 법률을 준수하세요.

## ⭐ 스타 히스토리

[![Star History Chart](https://api.star-history.com/svg?repos=Au1rxx/free-vpn-subscriptions&type=Date)](https://www.star-history.com/#Au1rxx/free-vpn-subscriptions&Date)

---

이 프로젝트가 도움이 되셨다면 ⭐을 남겨 주세요 —— 모든 스타가 다른 사람들이 이 프로젝트를 더 쉽게 발견하도록 도와줍니다.
