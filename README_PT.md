# Free VPN Subscriptions

[English](./README.md) · [简体中文](./README_CN.md) · [日本語](./README_JA.md) · [한국어](./README_KO.md) · [Español](./README_ES.md) · **Português** · [Русский](./README_RU.md)

<p align="center"><img src="https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/assets/hero.png" alt="Free VPN Subscriptions — hourly-refreshed free VPN subscriptions for Clash, sing-box, v2ray" width="780"></p>

![nós](https://img.shields.io/badge/nós-150-brightgreen) ![ativos](https://img.shields.io/badge/ativos-2579-blue) ![rtt--mediano](https://img.shields.io/badge/rtt--mediano-8ms-orange) ![atualizado](https://img.shields.io/badge/atualizado-2026-04-19_11:05_UTC-informational)

> **A forma mais fácil de obter uma VPN gratuita funcional — copie um link de assinatura, cole no seu cliente, conecte.**  
> Sem cadastro. Sem pagamento. Sem instalar nenhum binário. Atualizado a cada hora a partir de fontes públicas — cada nó é sondado via TCP + TLS antes da publicação.

> VPN grátis · assinatura VPN gratuita · proxy grátis · Clash assinatura · v2ray assinatura · sing-box assinatura · VLESS · Reality · VMess · Trojan · Shadowsocks · Hysteria2 · atualizado por hora · TCP+TLS testado · por país

## 💡 Por que este projeto?

Cada lista de "VPN gratuita" no GitHub está desatualizada, cheia de nós mortos, ou pede para instalar um binário suspeito. Este repositório **publica apenas nós que passaram um handshake TCP E um handshake TLS minutos atrás**, a partir de fontes públicas selecionadas, ordenados por latência. Você recebe 3 arquivos de assinatura portáteis — use-os em Clash, sing-box ou v2rayN e pronto.

> 📖 How the fetch → probe → rank pipeline works: [ARCHITECTURE.md](./ARCHITECTURE.md)

## 🔬 Como verificamos que os nós realmente funcionam

**Resposta honesta primeiro: não podemos *garantir* que um nó passe o seu tráfego.** Nenhum agregador pode, sem realmente enviar tráfego através dele. Aqui está exatamente o que verificamos, o que não podemos, e de onde vem a garantia real.

### ✅ O que verificamos na agregação (antes de publicar)

1. **Acessibilidade TCP** — abrimos uma conexão TCP para cada `server:port`. Hosts mortos, DNS errado, portas bloqueadas são descartados. Elimina cerca de 40 % das entradas cruas.
2. **Handshake TLS** — para cada nó TLS / Reality / WS-TLS completamos o handshake inteiro. Certificados expirados, SNI incompatíveis e short-ids Reality quebrados são descartados. Elimina mais ~10 %.
3. **Ordenação por latência** — os sobreviventes são ordenados por RTT e mantemos os N mais rápidos.

Números típicos de uma execução recente: **17 fontes → ~4,800 brutos → ~2,900 vivos via TCP → ~2,600 TLS OK → top 200 publicados**.

### ❌ O que não podemos verificar

- Autenticação do protocolo proxy. UUID / senha errados só são rejeitados *depois* do handshake TLS pelo servidor upstream.
- Sucesso real de HTTP via proxy.
- Largura de banda ou throughput.
- Geolocalização além do que o GeoIP diz sobre o IP de saída.

### 🛡️ Verificação em tempo de execução — é daqui que vem a garantia real

O `clash.yaml` que publicamos inclui um grupo `url-test` que **testa HTTP real através de cada nó** a cada 5 minutos:

```yaml
proxy-groups:
  - name: AUTO
    type: url-test
    url: http://www.gstatic.com/generate_204
    interval: 300
```

Seu cliente ordena a lista de nós pela latência *real* de HTTP via proxy e auto-seleciona o nó mais rápido que funciona. sing-box e v2ray têm mecanismos equivalentes. Se um nó selecionado morrer, o cliente muda para o próximo sem intervenção.

### 🧮 Resultado esperado

Dos top 200 publicados por execução, um cliente típico encontra 30-50 que servem HTTP limpo em qualquer momento. Rotacione se um ficar lento — o grupo url-test faz isso com um clique.

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
| 🇺🇸 United States (`US`) | 21 | [clash-US.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-US.yaml) | [singbox-US.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-US.json) | [v2ray-base64-US.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-US.txt) |
| 🇩🇪 Germany (`DE`) | 7 | [clash-DE.yaml](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/clash-DE.yaml) | [singbox-DE.json](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/singbox-DE.json) | [v2ray-base64-DE.txt](https://github.com/Au1rxx/free-vpn-subscriptions/raw/main/output/by-country/v2ray-base64-DE.txt) |

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
- **Ativos em todas as fontes**: 2579
- **RTT do nó mais rápido**: 2 ms
- **RTT mediano**: 8 ms
- **Última atualização (UTC)**: 2026-04-19 11:05 UTC

**Mix de protocolos:** shadowsocks × 25 · trojan × 18 · vless × 85 · vmess × 22

**Fontes usadas nesta execução:** `barry-far-v2ray` × 30 · `ebrasha-v2ray` × 9 · `epodonios` × 33 · `freefq` × 1 · `mahdi0024` × 1 · `mahdibland-aggregator` × 1 · `mahdibland-shadowsocks` × 1 · `mfuu-clash` × 2 · `ninjastrikers` × 35 · `pawdroid` × 1 · `ruking-clash` × 21 · `snakem982` × 12 · `surfboard-eternity` × 2 · `vxiaov-clash` × 1

## ❓ Perguntas frequentes

<details><summary>Isso é realmente grátis?</summary>

Sim. Os nós são operados por voluntários de terceiros que publicam suas próprias assinaturas gratuitas. Nós não operamos nenhum servidor — apenas testamos, classificamos e reempacotamos o que já é público.

</details>

<details><summary>Quão atualizados são os dados?</summary>

A cada hora (com um pequeno atraso aleatório para evitar bater nas fontes upstream exatamente em `:00`): puxa todas as fontes, faz sondagem TCP+TLS em cada nó, descarta os mortos, ordena por latência e publica os novos arquivos. Veja o carimbo `Last updated` acima.

</details>

<details><summary>Posso confiar nesses nós?</summary>

Nós gratuitos veem todo o seu tráfego. **Nunca os use para banco, login ou algo sensível.** Bom para driblar bloqueios geográficos em conteúdo público. Use seu próprio VPS / serviço pago para privacidade real.

</details>

<details><summary>Por que alguns nós listados falham?</summary>

Verificamos apenas acessibilidade TCP e handshake TLS — um nó ainda pode ter cota esgotada, roteamento errado ou certificado expirado. O `clash.yaml` publicado inclui um grupo `url-test` (`http://www.gstatic.com/generate_204`, intervalo de 300 s); seu cliente auto-seleciona o nó mais rápido que realmente serve HTTP. Se um morrer, pegue o próximo.

</details>

## 🤝 Contribuir

Conhece uma fonte de assinatura pública confiável que deveríamos adicionar? Abra uma issue com a URL e o formato.

## ⚠️ Aviso legal

Este repositório agrega configurações de proxy **compartilhadas publicamente** por voluntários de terceiros. Não operamos nenhum servidor, não garantimos disponibilidade ou segurança, e não somos responsáveis pelo uso. Destinado a uso educacional e conectividade pessoal. Cumpra todas as leis aplicáveis em sua jurisdição.

## ⭐ Histórico de estrelas

[![Star History Chart](https://api.star-history.com/svg?repos=Au1rxx/free-vpn-subscriptions&type=Date)](https://www.star-history.com/#Au1rxx/free-vpn-subscriptions&Date)

---

Se este projeto te ajudou, deixe uma ⭐ — cada estrela facilita para outros o encontrarem.
