# Free VPN Subscriptions

<div align="center">

[English](./README.md) · **简体中文** · [日本語](./README_JA.md) · [한국어](./README_KO.md) · [Español](./README_ES.md) · [Português](./README_PT.md) · [Русский](./README_RU.md)

</div>

<p align="center"><img src="https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/assets/hero.png" alt="Free VPN Subscriptions — hourly-refreshed free VPN subscriptions for Clash, sing-box, v2ray" width="780"></p>

[![GitHub stars](https://img.shields.io/github/stars/Au1rxx/free-vpn-subscriptions?style=flat&color=gold&logo=github)](https://github.com/Au1rxx/free-vpn-subscriptions/stargazers) ![精选](https://img.shields.io/badge/精选-2000-brightgreen) ![已验证](https://img.shields.io/badge/已验证-3748-blue) ![中位延迟](https://img.shields.io/badge/中位延迟-574ms-orange) ![更新](https://img.shields.io/badge/更新-2026-08-01_23:28_UTC-informational) [![License](https://img.shields.io/github/license/Au1rxx/free-vpn-subscriptions?color=blue)](https://github.com/Au1rxx/free-vpn-subscriptions/blob/main/LICENSE)

> **获取可用免费 VPN 的最简单方式 —— 复制订阅链接,粘贴到客户端,连上。**  
> 无需注册。无需付费。无需安装任何二进制。每小时从公共源自动抓取 —— 每个发布的节点都在几分钟前通过 sing-box 真实转发过 HTTP 流量。

> 免费 VPN 订阅 · 免费机场 · 免费梯子 · 免费科学上网 · Clash 订阅 · v2ray 订阅 · sing-box 订阅 · VLESS Reality · VMess · Trojan · Shadowsocks · Hysteria2 · 每小时刷新 · HTTP 实测验证 · 按国家分类

## 🚀 一键订阅

复制对应客户端的 URL,粘贴到订阅导入框:

| 客户端 | 格式 | 订阅链接 |
|---|---|---|
| Clash / Clash Verge / ClashX | `clash.yaml` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/clash.yaml` |
| sing-box | `singbox.json` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/singbox.json` |
| v2rayN / v2rayNG / Shadowrocket / NekoBox | `v2ray-base64` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/v2ray-base64.txt` |

> Database-classified shards (stable / all verified / protocol / country / network): [https://github.com/Au1rxx/free-vpn-subscriptions/tree/main/output](https://github.com/Au1rxx/free-vpn-subscriptions/tree/main/output)

⭐ **点个 Star** —— 如果它帮你省了时间。这是我们判断该继续做什么的唯一信号。

## 💡 为什么用这个项目?

GitHub 上几乎所有的"免费 VPN"列表都有三个问题:数据过期、全是死节点、或者要你装来路不明的二进制。本仓库比任何其他列表都更进一步 —— **我们不只是检查节点端口能不能通,而是用 sing-box 把真实 HTTP 请求经代理打出去、确认能收到 204 才发布**,全部在几分钟内完成。直接给你 3 种通用订阅文件 —— 粘到 Clash / sing-box / v2rayN 即用。

> 📖 How the fetch → probe → rank pipeline works: [ARCHITECTURE.md](./ARCHITECTURE.md)

<details>
<summary><b>🔬 我们如何验证节点可用</b></summary>

多数免费 VPN 列表停在"端口能连上"就发布了。我们不这样。下面是节点必须跑通的完整验证管线。

### ✅ 聚合阶段(发布前)的验证

1. **TCP 可达性** —— 对每个 `server:port` 发起一次 TCP 连接。服务器宕机、DNS 错误、端口被封全部被丢。大约过滤掉 40% 的原始条目。
2. **TLS 握手** —— 对所有 TLS / Reality / WS-TLS 节点完整跑一遍 TLS 握手。证书过期、SNI 不匹配、Reality short-id 失效都会被丢。再过滤掉约 10%。
3. **sing-box 配置校验** —— 每个幸存节点都被翻译成真实的 sing-box outbound,过 `sing-box check` 验证。密码错、UUID 畸形、flow 选项不支持 —— 全部在占用探测资源之前剔除。
4. **HTTP 过代理实测(这是关键的一步)** —— 我们把最快的约 900 个候选节点分批塞进 sing-box 子进程,每个节点分配一个本地 SOCKS5 入口,然后通过它真正发 HTTP 和 HTTPS 请求:
   - `http://www.gstatic.com/generate_204`(期望 204)
   - `https://www.cloudflare.com/cdn-cgi/trace`(期望 200)

   请求会完整走一遍代理协议(VLESS / VMess / Trojan / Shadowsocks / Hysteria2),所以能过这一关的节点,就证明它的认证、路由、TLS 内握手、出口网络全部是工作的。
5. **两轮,间隔 45 秒** —— 过了一次但 45 秒后就死的节点会被筛掉。只有在(轮数 × 目标数)中成功率 ≥ 50% 的节点才会留下。
6. **按真实延迟中位数排序** —— 幸存节点按 HTTP 过代理的真实往返时间中位数(不是原始 TCP RTT)排序,取前 N 发布。

最近一次运行的典型数字:**17 个源 → ~4,800 原始 → ~2,900 TCP 存活 → ~2,600 TLS OK → ~840 配置有效 → ~280 HTTP 实测通过 → 发布前 150**。发布出去的 150 个节点,每一个都在过去十分钟内真正转发过流量。

### ❌ 我们仍然验证不了什么

- **带宽 / 吞吐** —— 我们测延迟不测 Mbps。50ms 的节点看视频可能仍然慢。
- **精确地理位置** —— GeoIP 能告诉你出口 IP 是哪国,但城市或 ISP 级别的判断不可靠。
- **特定地区的封锁** —— 我们的探测机器能通的节点,不代表你的网络也能通(ISP 层封锁、captive portal 等)。
- **发布之后是否还活着** —— 十分钟前它是活的,之后可能挂了。

### 🛡️ 运行时兜底 —— 对付上面最后一条

我们发布的 `clash.yaml` 自带 `url-test` 组,**在你本地**每 5 分钟重新跑一次 HTTP 实测:

```yaml
proxy-groups:
  - name: AUTO
    type: url-test
    url: http://www.gstatic.com/generate_204
    interval: 300
```

客户端按 *你的网络* 下实时的 HTTP 过代理延迟排序,自动挑最快可用节点。sing-box 和 v2ray 有等价机制。如果聚合到下一次运行中途某节点挂了,客户端会自动切下一个,不需要你手动处理。

### 🧮 实际效果

每次发布的约 150 个节点里,客户端通常能找到 **80-120 个在你的网络下能稳定过 HTTP**,比只做 TCP 探测的列表命中率高 2-3 倍。其中一个挂了 url-test 组会透明地切换。

</details>

## 🌍 按国家订阅

只想要特定地区的节点?选一个针对性订阅链接:

| 国家/地区 | 节点数 | Clash | sing-box | v2ray |
|---|---|---|---|---|
| 🇺🇸 United States (`US`) | 354 | [clash-US.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-US.yaml) | [singbox-US.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-US.json) | [v2ray-base64-US.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-US.txt) |
| 🇫🇷 France (`FR`) | 290 | [clash-FR.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-FR.yaml) | [singbox-FR.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-FR.json) | [v2ray-base64-FR.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-FR.txt) |
| 🇳🇱 Netherlands (`NL`) | 265 | [clash-NL.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-NL.yaml) | [singbox-NL.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-NL.json) | [v2ray-base64-NL.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-NL.txt) |
| 🇩🇪 Germany (`DE`) | 181 | [clash-DE.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-DE.yaml) | [singbox-DE.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-DE.json) | [v2ray-base64-DE.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-DE.txt) |
| 🇹🇼 Taiwan (`TW`) | 114 | [clash-TW.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-TW.yaml) | [singbox-TW.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-TW.json) | [v2ray-base64-TW.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-TW.txt) |
| 🇭🇰 Hong Kong (`HK`) | 76 | [clash-HK.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-HK.yaml) | [singbox-HK.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-HK.json) | [v2ray-base64-HK.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-HK.txt) |
| 🇨🇭 Switzerland (`CH`) | 74 | [clash-CH.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-CH.yaml) | [singbox-CH.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-CH.json) | [v2ray-base64-CH.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-CH.txt) |
| 🇸🇬 Singapore (`SG`) | 68 | [clash-SG.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-SG.yaml) | [singbox-SG.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-SG.json) | [v2ray-base64-SG.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-SG.txt) |
| 🇬🇧 United Kingdom (`GB`) | 57 | [clash-GB.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-GB.yaml) | [singbox-GB.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-GB.json) | [v2ray-base64-GB.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-GB.txt) |
| 🇯🇵 Japan (`JP`) | 47 | [clash-JP.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-JP.yaml) | [singbox-JP.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-JP.json) | [v2ray-base64-JP.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-JP.txt) |
| 🇰🇷 Korea (`KR`) | 47 | [clash-KR.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-KR.yaml) | [singbox-KR.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-KR.json) | [v2ray-base64-KR.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-KR.txt) |
| 🇵🇱 Poland (`PL`) | 34 | [clash-PL.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-PL.yaml) | [singbox-PL.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-PL.json) | [v2ray-base64-PL.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-PL.txt) |
| 🇨🇳 China (`CN`) | 33 | [clash-CN.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-CN.yaml) | [singbox-CN.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-CN.json) | [v2ray-base64-CN.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-CN.txt) |
| 🇨🇴 Colombia (`CO`) | 28 | [clash-CO.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-CO.yaml) | [singbox-CO.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-CO.json) | [v2ray-base64-CO.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-CO.txt) |
| 🇷🇺 Russia (`RU`) | 28 | [clash-RU.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-RU.yaml) | [singbox-RU.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-RU.json) | [v2ray-base64-RU.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-RU.txt) |
| 🇹🇷 Turkey (`TR`) | 28 | [clash-TR.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-TR.yaml) | [singbox-TR.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-TR.json) | [v2ray-base64-TR.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-TR.txt) |
| 🇨🇦 Canada (`CA`) | 27 | [clash-CA.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-CA.yaml) | [singbox-CA.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-CA.json) | [v2ray-base64-CA.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-CA.txt) |
| 🇸🇪 Sweden (`SE`) | 26 | [clash-SE.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-SE.yaml) | [singbox-SE.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-SE.json) | [v2ray-base64-SE.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-SE.txt) |
| 🇫🇮 Finland (`FI`) | 21 | [clash-FI.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-FI.yaml) | [singbox-FI.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-FI.json) | [v2ray-base64-FI.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-FI.txt) |
| 🇮🇩 Indonesia (`ID`) | 17 | [clash-ID.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-ID.yaml) | [singbox-ID.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-ID.json) | [v2ray-base64-ID.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-ID.txt) |
| 🇷🇴 Romania (`RO`) | 16 | [clash-RO.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-RO.yaml) | [singbox-RO.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-RO.json) | [v2ray-base64-RO.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-RO.txt) |
| 🇱🇹 Lithuania (`LT`) | 14 | [clash-LT.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-LT.yaml) | [singbox-LT.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-LT.json) | [v2ray-base64-LT.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-LT.txt) |
| 🇨🇾 Cyprus (`CY`) | 12 | [clash-CY.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-CY.yaml) | [singbox-CY.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-CY.json) | [v2ray-base64-CY.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-CY.txt) |
| 🇮🇹 Italy (`IT`) | 12 | [clash-IT.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-IT.yaml) | [singbox-IT.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-IT.json) | [v2ray-base64-IT.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-IT.txt) |
| 🇧🇷 Brazil (`BR`) | 10 | [clash-BR.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-BR.yaml) | [singbox-BR.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-BR.json) | [v2ray-base64-BR.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-BR.txt) |
| 🇱🇻 Latvia (`LV`) | 10 | [clash-LV.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-LV.yaml) | [singbox-LV.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-LV.json) | [v2ray-base64-LV.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-LV.txt) |
| 🇪🇸 Spain (`ES`) | 8 | [clash-ES.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-ES.yaml) | [singbox-ES.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-ES.json) | [v2ray-base64-ES.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-ES.txt) |
| 🇹🇭 Thailand (`TH`) | 8 | [clash-TH.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-TH.yaml) | [singbox-TH.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-TH.json) | [v2ray-base64-TH.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-TH.txt) |
| 🇪🇪 Estonia (`EE`) | 7 | [clash-EE.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-EE.yaml) | [singbox-EE.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-EE.json) | [v2ray-base64-EE.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-EE.txt) |
| 🇻🇳 Vietnam (`VN`) | 7 | [clash-VN.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-VN.yaml) | [singbox-VN.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-VN.json) | [v2ray-base64-VN.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-VN.txt) |
| 🇧🇬 Bulgaria (`BG`) | 6 | [clash-BG.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-BG.yaml) | [singbox-BG.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-BG.json) | [v2ray-base64-BG.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-BG.txt) |
| 🇪🇨 EC (`EC`) | 5 | [clash-EC.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-EC.yaml) | [singbox-EC.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-EC.json) | [v2ray-base64-EC.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-EC.txt) |
| 🇮🇳 India (`IN`) | 5 | [clash-IN.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-IN.yaml) | [singbox-IN.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-IN.json) | [v2ray-base64-IN.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-IN.txt) |
| 🇳🇴 Norway (`NO`) | 5 | [clash-NO.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-NO.yaml) | [singbox-NO.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-NO.json) | [v2ray-base64-NO.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-NO.txt) |
| 🇵🇭 Philippines (`PH`) | 5 | [clash-PH.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-PH.yaml) | [singbox-PH.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-PH.json) | [v2ray-base64-PH.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-PH.txt) |
| 🇦🇺 Australia (`AU`) | 4 | [clash-AU.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-AU.yaml) | [singbox-AU.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-AU.json) | [v2ray-base64-AU.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-AU.txt) |
| 🇧🇩 Bangladesh (`BD`) | 4 | [clash-BD.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-BD.yaml) | [singbox-BD.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-BD.json) | [v2ray-base64-BD.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-BD.txt) |
| 🇮🇪 Ireland (`IE`) | 4 | [clash-IE.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-IE.yaml) | [singbox-IE.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-IE.json) | [v2ray-base64-IE.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-IE.txt) |
| 🇲🇽 Mexico (`MX`) | 4 | [clash-MX.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-MX.yaml) | [singbox-MX.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-MX.json) | [v2ray-base64-MX.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-MX.txt) |
| 🇲🇾 Malaysia (`MY`) | 4 | [clash-MY.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-MY.yaml) | [singbox-MY.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-MY.json) | [v2ray-base64-MY.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-MY.txt) |
| 🇹🇖 T1 (`T1`) | 4 | [clash-T1.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-T1.yaml) | [singbox-T1.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-T1.json) | [v2ray-base64-T1.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-T1.txt) |
| 🇿🇦 South Africa (`ZA`) | 3 | [clash-ZA.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-ZA.yaml) | [singbox-ZA.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-ZA.json) | [v2ray-base64-ZA.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-ZA.txt) |

## 📖 客户端图文教程

新手不知道怎么配?按平台选一篇跟着做:

- [**Clash Verge**](https://au1rxx.github.io/free-vpn-subscriptions/guides/clash-verge.zh.html) · Windows / macOS / Linux
- [**v2rayNG**](https://au1rxx.github.io/free-vpn-subscriptions/guides/v2rayng.zh.html) · Android
- [**Shadowrocket**](https://au1rxx.github.io/free-vpn-subscriptions/guides/shadowrocket.zh.html) · iOS / iPadOS
- [**sing-box**](https://au1rxx.github.io/free-vpn-subscriptions/guides/sing-box.zh.html) · Windows / macOS / Linux / iOS / Android

## 🧩 支持的客户端

- **Windows**:v2rayN、Clash Verge、Hiddify、NekoRay
- **macOS**:ClashX Pro、Clash Verge、sing-box、Hiddify
- **iOS**:Shadowrocket、Stash、Loon、sing-box、Hiddify
- **Android**:v2rayNG、NekoBox、Clash Meta for Android、Hiddify、sing-box
- **Linux**:mihomo (Clash.Meta)、sing-box、v2ray-core

## 📊 实时统计

- **精选节点数**: 2000
- **全源存活总数**: 3748
- **最快延迟**: 24 ms
- **中位延迟**: 574 ms
- **最后更新 (UTC)**: 2026-08-01 23:28 UTC

**协议分布:** http × 212 · https × 2 · hysteria2 × 136 · shadowsocks × 199 · socks4 × 25 · socks5 × 118 · trojan × 313 · tuic × 4 · vless × 989 · vmess × 2

## ❓ 常见问题

<details><summary>真的完全免费吗?</summary>

是的。所有节点由第三方志愿者自己运营并公开免费订阅。本仓库不运营任何服务器,只是做测试、排名、重新打包公开内容。

</details>

<details><summary>数据多新?</summary>

每小时刷新一次(带一点随机延迟避免整点集中打上游):拉取所有源 → TCP → TLS → sing-box 配置校验 → HTTP 过代理探测(两轮,间隔 45 秒)→ 按真实 HTTP 延迟排序 → 发布新的输出文件。完整管线约 10 分钟。见顶部徽章上的更新时间。

</details>

<details><summary>这些节点可以信任吗?</summary>

免费节点能看到你所有流量。**绝不要用来登录银行、邮箱等敏感账号。**用来突破地区限制访问公开内容没问题。真正需要隐私请自建 VPS 或付费服务。

</details>

<details><summary>列表里的节点为什么有的连不上?</summary>

即使过了我们的 HTTP 实测,节点在聚合后也可能死掉:配额用尽、上游吊销了 key、你的 ISP 封锁了出口 IP,或者运营者干脆下架了。发布的 `clash.yaml` 自带 `url-test` 组(每 300 秒对 `http://www.gstatic.com/generate_204` 打一次),客户端会在 *你的网络下* 自动选真正能过 HTTP 的最快节点。挂了就换下一个。通常 150 个里随时有 80-120 个能用。

</details>

## 🤝 贡献

知道稳定的公共订阅源可以加入?提 issue 给我们 URL 和格式。

## ⚠️ 免责声明

本仓库聚合第三方志愿者**公开分享**的代理配置。我们不运营任何服务器,不保证可用性或安全性,不为使用行为负责。仅供学习和个人连接使用。请遵守所在司法管辖区的法律。

## ⭐ Star 历史

[![Star History Chart](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/assets/star-history.svg)](https://github.com/Au1rxx/free-vpn-subscriptions/stargazers)

---

如果这个项目帮到你,点个 ⭐ —— 每一颗 star 都能帮更多人发现它。
