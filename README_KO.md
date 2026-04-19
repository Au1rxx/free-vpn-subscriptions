# Free VPN Subscriptions

[English](./README.md) · [简体中文](./README_CN.md) · [日本語](./README_JA.md) · **한국어** · [Español](./README_ES.md) · [Português](./README_PT.md) · [Русский](./README_RU.md)

<p align="center"><img src="https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/assets/workflow.svg" alt="Free VPN Subscriptions — aggregation workflow" width="780"></p>

![노드](https://img.shields.io/badge/노드-150-brightgreen) ![생존](https://img.shields.io/badge/생존-393-blue) ![중앙값--rtt](https://img.shields.io/badge/중앙값--rtt-96ms-orange) ![업데이트](https://img.shields.io/badge/업데이트-2026-04-19_07:12_UTC-informational)

> **작동하는 무료 VPN을 얻는 가장 쉬운 방법 —— 구독 링크를 복사하고 클라이언트에 붙여 넣고 연결하세요.**  
> 가입 불필요. 결제 불필요. 바이너리 설치 불필요. 공개 소스에서 매시간 자동 갱신되며 모든 노드가 검증됩니다.

> 무료 VPN 구독 · 무료 v2ray 구독 · 무료 Clash 구독 · 무료 sing-box 구독 · VLESS · Reality · VMess · Trojan · Shadowsocks · Hysteria2 · 매시간 갱신 · TCP+TLS 프로브 완료 · 국가별

## 💡 왜 이 프로젝트?

GitHub의 거의 모든 "무료 VPN" 목록은 데이터가 오래되었거나, 죽은 노드로 가득 차 있거나, 출처가 불분명한 바이너리 설치를 요구합니다. 이 저장소는 **몇 분 전에 TCP 핸드셰이크와 TLS 핸드셰이크를 모두 통과한 노드만** 선별된 공개 소스에서 레이턴시 순으로 게시합니다. Clash / sing-box / v2rayN에 바로 붙여 넣을 수 있는 3가지 범용 구독 파일을 제공합니다.

> 📖 How the fetch → probe → rank pipeline works: [ARCHITECTURE.md](./ARCHITECTURE.md)

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
| 🇬🇧 United Kingdom (`GB`) | 31 | [clash-GB.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-GB.yaml) | [singbox-GB.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-GB.json) | [v2ray-base64-GB.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-GB.txt) |
| 🇺🇸 United States (`US`) | 24 | [clash-US.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-US.yaml) | [singbox-US.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-US.json) | [v2ray-base64-US.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-US.txt) |
| 🇳🇱 Netherlands (`NL`) | 14 | [clash-NL.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-NL.yaml) | [singbox-NL.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-NL.json) | [v2ray-base64-NL.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-NL.txt) |
| 🇨🇦 Canada (`CA`) | 11 | [clash-CA.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-CA.yaml) | [singbox-CA.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-CA.json) | [v2ray-base64-CA.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-CA.txt) |
| 🇯🇵 Japan (`JP`) | 8 | [clash-JP.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-JP.yaml) | [singbox-JP.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-JP.json) | [v2ray-base64-JP.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-JP.txt) |
| 🇩🇪 Germany (`DE`) | 6 | [clash-DE.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-DE.yaml) | [singbox-DE.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-DE.json) | [v2ray-base64-DE.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-DE.txt) |
| 🇰🇷 Korea (`KR`) | 4 | [clash-KR.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-KR.yaml) | [singbox-KR.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-KR.json) | [v2ray-base64-KR.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-KR.txt) |
| 🇲🇦 MA (`MA`) | 3 | [clash-MA.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-MA.yaml) | [singbox-MA.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-MA.json) | [v2ray-base64-MA.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-MA.txt) |

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
- **전체 소스 생존 수**: 393
- **최고 속도 RTT**: 5 ms
- **중앙값 RTT**: 96 ms
- **최종 업데이트 (UTC)**: 2026-04-19 07:12 UTC

**프로토콜 분포:** shadowsocks × 102 · trojan × 12 · vmess × 36

**이번 실행에 사용된 소스:** `freefq` × 2 · `mahdibland-aggregator` × 76 · `mahdibland-shadowsocks` × 56 · `pawdroid` × 3 · `vxiaov-clash` × 13

## ❓ 자주 묻는 질문

<details><summary>정말 무료인가요?</summary>

네. 모든 노드는 제3자 자원봉사자가 운영하며 공개 구독을 스스로 게시합니다. 저희는 어떤 서버도 운영하지 않으며, 이미 공개된 것을 테스트하고 순위를 매기고 재포장할 뿐입니다.

</details>

<details><summary>데이터는 얼마나 신선한가요?</summary>

GitHub Action이 매시간 실행됩니다: 모든 상위 소스 가져오기 → 각 노드 TCP+TLS 프로브 → 죽은 것 제거 → 레이턴시 순 정렬 → 새 출력 파일 커밋. 위의 `Last updated` 타임스탬프를 확인하세요.

</details>

<details><summary>이 노드들을 신뢰할 수 있나요?</summary>

무료 노드는 모든 트래픽을 운영자가 볼 수 있습니다. **은행 거래, 로그인, 민감한 작업에는 절대 사용하지 마세요.** 공개 콘텐츠의 지역 제한 우회에는 적합합니다. 실제 프라이버시에는 자체 VPS/유료 서비스를 사용하세요.

</details>

<details><summary>목록에 있는데 작동하지 않는 노드가 있는 이유는?</summary>

TCP 도달성과 TLS 핸드셰이크를 검증하지만 노드는 여전히 할당량 소진, 잘못된 라우팅, 만료된 인증서를 가질 수 있습니다. 몇 개 시도해 보세요. selector 그룹에 대체 항목이 있습니다.

</details>

## 🤝 기여

신뢰할 수 있는 공개 구독 소스를 알고 계신가요? URL과 형식을 포함한 이슈를 열어 주세요.

## ⚠️ 면책 조항

이 저장소는 제3자 자원봉사자가 **공개 공유**한 프록시 구성을 집계합니다. 저희는 어떤 서버도 운영하지 않고, 가용성이나 보안을 보장하지 않으며, 사용 방식에 대해 책임지지 않습니다. 교육 및 개인 연결 용도로만 사용하세요. 해당 관할권의 모든 법률을 준수하세요.

## ⭐ 스타 히스토리

[![Star History Chart](https://api.star-history.com/svg?repos=Au1rxx/free-vpn-subscriptions&type=Date)](https://www.star-history.com/#Au1rxx/free-vpn-subscriptions&Date)

---

이 프로젝트가 도움이 되셨다면 ⭐을 남겨 주세요 —— 모든 스타가 다른 사람들이 이 프로젝트를 더 쉽게 발견하도록 도와줍니다.
