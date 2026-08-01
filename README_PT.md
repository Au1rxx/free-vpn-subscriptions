# Free VPN Subscriptions

<div align="center">

[English](./README.md) · [简体中文](./README_CN.md) · [日本語](./README_JA.md) · [한국어](./README_KO.md) · [Español](./README_ES.md) · **Português** · [Русский](./README_RU.md)

</div>

<p align="center"><img src="https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/assets/hero.png" alt="Free VPN Subscriptions — hourly-refreshed free VPN subscriptions for Clash, sing-box, v2ray" width="780"></p>

[![GitHub stars](https://img.shields.io/github/stars/Au1rxx/free-vpn-subscriptions?style=flat&color=gold&logo=github)](https://github.com/Au1rxx/free-vpn-subscriptions/stargazers) ![selecionados](https://img.shields.io/badge/selecionados-2000-brightgreen) ![verificados](https://img.shields.io/badge/verificados-3868-blue) ![rtt--mediano](https://img.shields.io/badge/rtt--mediano-553ms-orange) ![atualizado](https://img.shields.io/badge/atualizado-2026-08-01_05:26_UTC-informational) [![License](https://img.shields.io/github/license/Au1rxx/free-vpn-subscriptions?color=blue)](https://github.com/Au1rxx/free-vpn-subscriptions/blob/main/LICENSE)

> **A forma mais fácil de obter uma VPN gratuita funcional — copie um link de assinatura, cole no seu cliente, conecte.**  
> Sem cadastro. Sem pagamento. Sem instalar nenhum binário. Atualizado a cada hora a partir de fontes públicas — cada nó publicado encaminhou tráfego HTTP real através do sing-box minutos atrás.

> VPN grátis · assinatura VPN gratuita · proxy grátis · Clash assinatura · v2ray assinatura · sing-box assinatura · VLESS · Reality · VMess · Trojan · Shadowsocks · Hysteria2 · atualizado por hora · HTTP verificado sobre proxy · por país

## 🚀 Assinatura com um clique

Copie a URL que corresponde ao seu cliente e cole no campo de importação de assinatura:

| Cliente | Formato | URL de assinatura |
|---|---|---|
| Clash / Clash Verge / ClashX | `clash.yaml` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/clash.yaml` |
| sing-box | `singbox.json` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/singbox.json` |
| v2rayN / v2rayNG / Shadowrocket / NekoBox | `v2ray-base64` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/v2ray-base64.txt` |

> Database-classified shards (stable / all verified / protocol / country / network): [https://github.com/Au1rxx/free-vpn-subscriptions/tree/main/output](https://github.com/Au1rxx/free-vpn-subscriptions/tree/main/output)

⭐ **Dê uma estrela ao repositório** se ele te poupou tempo — é o único sinal que usamos para decidir o que continuar construindo.

## 💡 Por que este projeto?

Cada lista de "VPN gratuita" no GitHub está desatualizada, cheia de nós mortos, ou pede para instalar um binário suspeito. Este repositório vai um passo além de qualquer outro — **não apenas verificamos que o nó responde, mas empurramos tráfego HTTP real através dele com sing-box e confirmamos que um 204 retorna**, tudo em minutos antes de publicar. Você recebe 3 arquivos de assinatura portáteis — use-os em Clash, sing-box ou v2rayN e pronto.

> 📖 How the fetch → probe → rank pipeline works: [ARCHITECTURE.md](./ARCHITECTURE.md)

<details>
<summary><b>🔬 Como verificamos que os nós realmente funcionam</b></summary>

A maioria das listas de VPN gratuita para em \"a porta TCP está aberta\" e publica. Nós não. Aqui está a pipeline completa que um nó precisa passar antes de entrar na assinatura.

### ✅ O que verificamos na agregação (antes de publicar)

1. **Acessibilidade TCP** — abrimos uma conexão TCP para cada `server:port`. Hosts mortos, DNS errado, portas bloqueadas são descartados. ~40 % das entradas cruas caem aqui.
2. **Handshake TLS** — para cada nó TLS / Reality / WS-TLS completamos o handshake inteiro. Certificados expirados, SNI incompatíveis e short-ids Reality quebrados são descartados. Mais ~10 % caem aqui.
3. **Validação de configuração sing-box** — cada nó sobrevivente é traduzido em um outbound real de sing-box e passa pelo `sing-box check`. Cifras corrompidas, UUIDs errados e opções flow não suportadas são descartados antes de desperdiçar um slot de sondagem.
4. **Sondagem HTTP-over-proxy (esta é a chave)** — agrupamos os ~900 candidatos mais rápidos em subprocessos sing-box, cada nó recebendo seu próprio inbound SOCKS5 local, e então enviamos GETs HTTP e HTTPS reais através dele:
   - `http://www.gstatic.com/generate_204` (espera 204)
   - `https://www.cloudflare.com/cdn-cgi/trace` (espera 200)

   A requisição atravessa o protocolo proxy real (VLESS / VMess / Trojan / Shadowsocks / Hysteria2), então um nó que passa tem demonstravelmente autenticação, roteamento, handshake TLS interno e rede de saída funcionando.
5. **Duas rodadas, 45 segundos de intervalo** — nós que passam uma vez mas morrem 45 segundos depois são filtrados. Apenas nós com ≥ 50 % de taxa de sucesso em (rodadas × alvos) ficam.
6. **Ordenar por mediana de latência real** — os sobreviventes são ordenados pela mediana do ida-e-volta HTTP-over-proxy (não RTT TCP bruto) e os top N são publicados.

Números típicos de uma execução recente: **17 fontes → ~4,800 brutos → ~2,900 vivos via TCP → ~2,600 TLS OK → ~840 configuração válida → ~280 verificados por HTTP → top 150 publicados**. Cada um dos 150 de fato encaminhou tráfego nos últimos dez minutos.

### ❌ O que ainda não podemos verificar

- **Largura de banda / throughput** — medimos latência, não megabits. Um nó de 50 ms ainda pode ser lento para vídeo.
- **Precisão de geolocalização** — GeoIP diz o país do IP de saída mas não a cidade ou ISP de forma confiável.
- **Bloqueios específicos por região** — um nó que funciona da nossa infraestrutura de sondagem pode estar bloqueado da sua (filtragem no nível do ISP, captive portals, etc.).
- **Continuar vivo depois da execução** — o nó passou dez minutos atrás; pode ter morrido desde então.

### 🛡️ Rede de segurança em tempo de execução — para o último item acima

O `clash.yaml` que publicamos inclui um grupo `url-test` que retesta HTTP real através de cada nó a cada 5 minutos no *seu* dispositivo:

```yaml
proxy-groups:
  - name: AUTO
    type: url-test
    url: http://www.gstatic.com/generate_204
    interval: 300
```

Seu cliente mantém a lista de nós ordenada por latência *ao vivo* de HTTP-over-proxy da sua rede e auto-seleciona o nó mais rápido que funciona. sing-box e v2ray têm mecanismos equivalentes. Se um nó selecionado morrer entre agregações horárias, o cliente muda para o próximo sem intervenção.

### 🧮 O que isso significa na prática

Dos ~150 publicados por execução, um cliente típico encontra **80-120 nós que servem HTTP limpo da sua rede** em qualquer momento — aproximadamente 2-3× a taxa de acerto de listas que só fazem sondagem TCP. O grupo url-test rotaciona de forma transparente se um cair.

</details>

## 🌍 Por país

Quer nós apenas em uma região específica? Use uma dessas URLs de assinatura direcionadas:

| País | Nós | Clash | sing-box | v2ray |
|---|---|---|---|---|
| 🇺🇸 United States (`US`) | 403 | [clash-US.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-US.yaml) | [singbox-US.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-US.json) | [v2ray-base64-US.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-US.txt) |
| 🇫🇷 France (`FR`) | 314 | [clash-FR.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-FR.yaml) | [singbox-FR.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-FR.json) | [v2ray-base64-FR.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-FR.txt) |
| 🇳🇱 Netherlands (`NL`) | 248 | [clash-NL.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-NL.yaml) | [singbox-NL.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-NL.json) | [v2ray-base64-NL.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-NL.txt) |
| 🇩🇪 Germany (`DE`) | 187 | [clash-DE.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-DE.yaml) | [singbox-DE.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-DE.json) | [v2ray-base64-DE.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-DE.txt) |
| 🇹🇼 Taiwan (`TW`) | 85 | [clash-TW.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-TW.yaml) | [singbox-TW.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-TW.json) | [v2ray-base64-TW.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-TW.txt) |
| 🇭🇰 Hong Kong (`HK`) | 76 | [clash-HK.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-HK.yaml) | [singbox-HK.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-HK.json) | [v2ray-base64-HK.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-HK.txt) |
| 🇬🇧 United Kingdom (`GB`) | 69 | [clash-GB.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-GB.yaml) | [singbox-GB.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-GB.json) | [v2ray-base64-GB.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-GB.txt) |
| 🇰🇷 Korea (`KR`) | 66 | [clash-KR.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-KR.yaml) | [singbox-KR.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-KR.json) | [v2ray-base64-KR.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-KR.txt) |
| 🇸🇬 Singapore (`SG`) | 66 | [clash-SG.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-SG.yaml) | [singbox-SG.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-SG.json) | [v2ray-base64-SG.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-SG.txt) |
| 🇨🇭 Switzerland (`CH`) | 52 | [clash-CH.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-CH.yaml) | [singbox-CH.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-CH.json) | [v2ray-base64-CH.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-CH.txt) |
| 🇫🇮 Finland (`FI`) | 48 | [clash-FI.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-FI.yaml) | [singbox-FI.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-FI.json) | [v2ray-base64-FI.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-FI.txt) |
| 🇯🇵 Japan (`JP`) | 44 | [clash-JP.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-JP.yaml) | [singbox-JP.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-JP.json) | [v2ray-base64-JP.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-JP.txt) |
| 🇨🇦 Canada (`CA`) | 36 | [clash-CA.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-CA.yaml) | [singbox-CA.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-CA.json) | [v2ray-base64-CA.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-CA.txt) |
| 🇵🇱 Poland (`PL`) | 32 | [clash-PL.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-PL.yaml) | [singbox-PL.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-PL.json) | [v2ray-base64-PL.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-PL.txt) |
| 🇨🇳 China (`CN`) | 31 | [clash-CN.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-CN.yaml) | [singbox-CN.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-CN.json) | [v2ray-base64-CN.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-CN.txt) |
| 🇹🇷 Turkey (`TR`) | 26 | [clash-TR.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-TR.yaml) | [singbox-TR.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-TR.json) | [v2ray-base64-TR.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-TR.txt) |
| 🇨🇾 Cyprus (`CY`) | 21 | [clash-CY.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-CY.yaml) | [singbox-CY.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-CY.json) | [v2ray-base64-CY.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-CY.txt) |
| 🇸🇪 Sweden (`SE`) | 21 | [clash-SE.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-SE.yaml) | [singbox-SE.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-SE.json) | [v2ray-base64-SE.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-SE.txt) |
| 🇮🇹 Italy (`IT`) | 18 | [clash-IT.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-IT.yaml) | [singbox-IT.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-IT.json) | [v2ray-base64-IT.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-IT.txt) |
| 🇪🇪 Estonia (`EE`) | 11 | [clash-EE.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-EE.yaml) | [singbox-EE.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-EE.json) | [v2ray-base64-EE.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-EE.txt) |
| 🇮🇩 Indonesia (`ID`) | 11 | [clash-ID.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-ID.yaml) | [singbox-ID.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-ID.json) | [v2ray-base64-ID.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-ID.txt) |
| 🇮🇳 India (`IN`) | 10 | [clash-IN.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-IN.yaml) | [singbox-IN.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-IN.json) | [v2ray-base64-IN.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-IN.txt) |
| 🇪🇸 Spain (`ES`) | 9 | [clash-ES.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-ES.yaml) | [singbox-ES.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-ES.json) | [v2ray-base64-ES.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-ES.txt) |
| 🇷🇴 Romania (`RO`) | 9 | [clash-RO.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-RO.yaml) | [singbox-RO.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-RO.json) | [v2ray-base64-RO.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-RO.txt) |
| 🇹🇭 Thailand (`TH`) | 9 | [clash-TH.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-TH.yaml) | [singbox-TH.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-TH.json) | [v2ray-base64-TH.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-TH.txt) |
| 🇷🇺 Russia (`RU`) | 8 | [clash-RU.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-RU.yaml) | [singbox-RU.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-RU.json) | [v2ray-base64-RU.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-RU.txt) |
| 🇱🇻 Latvia (`LV`) | 6 | [clash-LV.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-LV.yaml) | [singbox-LV.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-LV.json) | [v2ray-base64-LV.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-LV.txt) |
| 🇻🇳 Vietnam (`VN`) | 5 | [clash-VN.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-VN.yaml) | [singbox-VN.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-VN.json) | [v2ray-base64-VN.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-VN.txt) |
| 🇦🇷 Argentina (`AR`) | 4 | [clash-AR.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-AR.yaml) | [singbox-AR.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-AR.json) | [v2ray-base64-AR.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-AR.txt) |
| 🇦🇺 Australia (`AU`) | 4 | [clash-AU.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-AU.yaml) | [singbox-AU.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-AU.json) | [v2ray-base64-AU.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-AU.txt) |
| 🇧🇷 Brazil (`BR`) | 4 | [clash-BR.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-BR.yaml) | [singbox-BR.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-BR.json) | [v2ray-base64-BR.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-BR.txt) |
| 🇨🇴 Colombia (`CO`) | 4 | [clash-CO.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-CO.yaml) | [singbox-CO.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-CO.json) | [v2ray-base64-CO.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-CO.txt) |
| 🇪🇨 EC (`EC`) | 4 | [clash-EC.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-EC.yaml) | [singbox-EC.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-EC.json) | [v2ray-base64-EC.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-EC.txt) |
| 🇲🇾 Malaysia (`MY`) | 4 | [clash-MY.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-MY.yaml) | [singbox-MY.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-MY.json) | [v2ray-base64-MY.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-MY.txt) |
| 🇳🇴 Norway (`NO`) | 4 | [clash-NO.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-NO.yaml) | [singbox-NO.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-NO.json) | [v2ray-base64-NO.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-NO.txt) |
| 🇹🇖 T1 (`T1`) | 4 | [clash-T1.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-T1.yaml) | [singbox-T1.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-T1.json) | [v2ray-base64-T1.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-T1.txt) |
| 🇧🇩 Bangladesh (`BD`) | 3 | [clash-BD.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-BD.yaml) | [singbox-BD.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-BD.json) | [v2ray-base64-BD.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-BD.txt) |
| 🇨🇿 Czechia (`CZ`) | 3 | [clash-CZ.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-CZ.yaml) | [singbox-CZ.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-CZ.json) | [v2ray-base64-CZ.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-CZ.txt) |
| 🇲🇽 Mexico (`MX`) | 3 | [clash-MX.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-MX.yaml) | [singbox-MX.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-MX.json) | [v2ray-base64-MX.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-MX.txt) |

## 📖 Tutoriais passo a passo

Novo nos clientes VPN? Escolha sua plataforma e siga o tutorial:

- [**Clash Verge**](https://au1rxx.github.io/free-vpn-subscriptions/guides/clash-verge.html) · Windows / macOS / Linux
- [**v2rayNG**](https://au1rxx.github.io/free-vpn-subscriptions/guides/v2rayng.html) · Android
- [**Shadowrocket**](https://au1rxx.github.io/free-vpn-subscriptions/guides/shadowrocket.html) · iOS / iPadOS
- [**sing-box**](https://au1rxx.github.io/free-vpn-subscriptions/guides/sing-box.html) · Windows / macOS / Linux / iOS / Android

## 🧩 Clientes suportados

- **Windows**: v2rayN, Clash Verge, Hiddify, NekoRay
- **macOS**: ClashX Pro, Clash Verge, sing-box, Hiddify
- **iOS**: Shadowrocket, Stash, Loon, sing-box, Hiddify
- **Android**: v2rayNG, NekoBox, Clash Meta for Android, Hiddify, sing-box
- **Linux**: mihomo (Clash.Meta), sing-box, v2ray-core

## 📊 Estatísticas ao vivo

- **Nós selecionados**: 2000
- **Ativos em todas as fontes**: 3868
- **RTT do nó mais rápido**: 47 ms
- **RTT mediano**: 553 ms
- **Última atualização (UTC)**: 2026-08-01 05:26 UTC

**Mix de protocolos:** http × 169 · https × 2 · hysteria2 × 132 · shadowsocks × 203 · socks4 × 33 · socks5 × 126 · trojan × 335 · tuic × 5 · vless × 994 · vmess × 1

## ❓ Perguntas frequentes

<details><summary>Isso é realmente grátis?</summary>

Sim. Os nós são operados por voluntários de terceiros que publicam suas próprias assinaturas gratuitas. Nós não operamos nenhum servidor — apenas testamos, classificamos e reempacotamos o que já é público.

</details>

<details><summary>Quão atualizados são os dados?</summary>

A cada hora (com um pequeno atraso aleatório para evitar bater nas fontes upstream exatamente em `:00`): puxa todas as fontes → TCP → TLS → validação de configuração sing-box → sondagem HTTP-over-proxy (duas rodadas, 45 s de intervalo) → ordena por latência HTTP real → publica os novos arquivos. Pipeline completo leva ~10 minutos. Veja o carimbo `Last updated` acima.

</details>

<details><summary>Posso confiar nesses nós?</summary>

Nós gratuitos veem todo o seu tráfego. **Nunca os use para banco, login ou algo sensível.** Bom para driblar bloqueios geográficos em conteúdo público. Use seu próprio VPS / serviço pago para privacidade real.

</details>

<details><summary>Por que alguns nós listados falham?</summary>

Mesmo após nossa sondagem HTTP-over-proxy, os nós podem morrer entre agregações: cota esgotada, upstream revogou a chave, seu ISP bloqueia o IP de saída, ou o operador desligou. O `clash.yaml` publicado inclui um grupo `url-test` (`http://www.gstatic.com/generate_204`, intervalo de 300 s); seu cliente auto-seleciona o nó mais rápido que realmente serve HTTP *da sua rede*. Se um morrer, pegue o próximo. Espere que 80-120 dos 150 funcionem em qualquer momento.

</details>

## 🤝 Contribuir

Conhece uma fonte de assinatura pública confiável que deveríamos adicionar? Abra uma issue com a URL e o formato.

## ⚠️ Aviso legal

Este repositório agrega configurações de proxy **compartilhadas publicamente** por voluntários de terceiros. Não operamos nenhum servidor, não garantimos disponibilidade ou segurança, e não somos responsáveis pelo uso. Destinado a uso educacional e conectividade pessoal. Cumpra todas as leis aplicáveis em sua jurisdição.

## ⭐ Histórico de estrelas

[![Star History Chart](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/assets/star-history.svg)](https://github.com/Au1rxx/free-vpn-subscriptions/stargazers)

---

Se este projeto te ajudou, deixe uma ⭐ — cada estrela facilita para outros o encontrarem.
