# Free VPN Subscriptions

[English](./README.md) · [简体中文](./README_CN.md) · [日本語](./README_JA.md) · [한국어](./README_KO.md) · [Español](./README_ES.md) · **Português** · [Русский](./README_RU.md)

<p align="center"><img src="https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/assets/hero.png" alt="Free VPN Subscriptions — hourly-refreshed free VPN subscriptions for Clash, sing-box, v2ray" width="780"></p>

![nós](https://img.shields.io/badge/nós-150-brightgreen) ![ativos](https://img.shields.io/badge/ativos-2195-blue) ![rtt--mediano](https://img.shields.io/badge/rtt--mediano-117ms-orange) ![atualizado](https://img.shields.io/badge/atualizado-2026-04-30_10:18_UTC-informational)

> **A forma mais fácil de obter uma VPN gratuita funcional — copie um link de assinatura, cole no seu cliente, conecte.**  
> Sem cadastro. Sem pagamento. Sem instalar nenhum binário. Atualizado a cada hora a partir de fontes públicas — cada nó publicado encaminhou tráfego HTTP real através do sing-box minutos atrás.

> VPN grátis · assinatura VPN gratuita · proxy grátis · Clash assinatura · v2ray assinatura · sing-box assinatura · VLESS · Reality · VMess · Trojan · Shadowsocks · Hysteria2 · atualizado por hora · HTTP verificado sobre proxy · por país

## 💡 Por que este projeto?

Cada lista de "VPN gratuita" no GitHub está desatualizada, cheia de nós mortos, ou pede para instalar um binário suspeito. Este repositório vai um passo além de qualquer outro — **não apenas verificamos que o nó responde, mas empurramos tráfego HTTP real através dele com sing-box e confirmamos que um 204 retorna**, tudo em minutos antes de publicar. Você recebe 3 arquivos de assinatura portáteis — use-os em Clash, sing-box ou v2rayN e pronto.

> 📖 How the fetch → probe → rank pipeline works: [ARCHITECTURE.md](./ARCHITECTURE.md)

## 🔬 Como verificamos que os nós realmente funcionam

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

## 🚀 Assinatura com um clique

Copie a URL que corresponde ao seu cliente e cole no campo de importação de assinatura:

| Cliente | Formato | URL de assinatura |
|---|---|---|
| Clash / Clash Verge / ClashX | `clash.yaml` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/clash.yaml` |
| sing-box | `singbox.json` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/singbox.json` |
| v2rayN / v2rayNG / Shadowrocket / NekoBox | `v2ray-base64` | `https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/v2ray-base64.txt` |

## 🌍 Por país

Quer nós apenas em uma região específica? Use uma dessas URLs de assinatura direcionadas:

| País | Nós | Clash | sing-box | v2ray |
|---|---|---|---|---|
| 🇺🇸 United States (`US`) | 47 | [clash-US.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-US.yaml) | [singbox-US.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-US.json) | [v2ray-base64-US.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-US.txt) |
| 🇨🇦 Canada (`CA`) | 8 | [clash-CA.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-CA.yaml) | [singbox-CA.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-CA.json) | [v2ray-base64-CA.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-CA.txt) |
| 🇩🇪 Germany (`DE`) | 5 | [clash-DE.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-DE.yaml) | [singbox-DE.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-DE.json) | [v2ray-base64-DE.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-DE.txt) |
| 🇯🇵 Japan (`JP`) | 4 | [clash-JP.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-JP.yaml) | [singbox-JP.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-JP.json) | [v2ray-base64-JP.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-JP.txt) |
| 🇬🇧 United Kingdom (`GB`) | 3 | [clash-GB.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-GB.yaml) | [singbox-GB.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-GB.json) | [v2ray-base64-GB.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-GB.txt) |

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

- **Nós selecionados**: 150
- **Ativos em todas as fontes**: 2195
- **RTT do nó mais rápido**: 29 ms
- **RTT mediano**: 117 ms
- **Última atualização (UTC)**: 2026-04-30 10:18 UTC

**Mix de protocolos:** shadowsocks × 42 · trojan × 11 · vless × 89 · vmess × 8

**Fontes usadas nesta execução:** `barry-far-v2ray` × 33 · `epodonios` × 6 · `lagzian-mix` × 4 · `mahdi0024` × 15 · `mahdibland-aggregator` × 10 · `mahdibland-shadowsocks` × 7 · `mfuu-clash` × 6 · `ninjastrikers` × 59 · `pawdroid` × 1 · `ruking-clash` × 1 · `snakem982` × 1 · `surfboard-eternity` × 7

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

[![Star History Chart](https://api.star-history.com/svg?repos=Au1rxx/free-vpn-subscriptions&type=Date)](https://www.star-history.com/#Au1rxx/free-vpn-subscriptions&Date)

---

Se este projeto te ajudou, deixe uma ⭐ — cada estrela facilita para outros o encontrarem.
