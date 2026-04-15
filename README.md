# 持续更新的免费 VPN 订阅

[中文](./README.md) | [English](./README.en.md)

[![Public Repo](https://img.shields.io/badge/repo-public-0f766e)](https://github.com/Au1rxx/free-vpn-subscriptions)
[![Updated](https://img.shields.io/badge/updates-every%206%20hours-0f766e)](https://au1rxx.github.io/free-vpn-subscriptions/updates.html)
[![Verified](https://img.shields.io/badge/validation-auto%20checked-1d221c)](https://au1rxx.github.io/free-vpn-subscriptions/verification.html)
[![Formats](https://img.shields.io/badge/formats-clash%20%7C%20sing--box%20%7C%20v2ray-cf6a32)](https://au1rxx.github.io/free-vpn-subscriptions/)
[![Status Feed](https://img.shields.io/badge/status-live-1d221c)](https://au1rxx.github.io/free-vpn-subscriptions/status.html)
[![Latest Release](https://img.shields.io/github/v/release/Au1rxx/free-vpn-subscriptions)](https://github.com/Au1rxx/free-vpn-subscriptions/releases/latest)

公开提供 Clash、sing-box、V2Ray 三种格式的订阅分发入口。状态每小时检查，公开快照按计划刷新，附带客户端教程、验证说明和故障排查页面。

[打开站点首页](https://au1rxx.github.io/free-vpn-subscriptions/) • [更新记录与快照历史](https://au1rxx.github.io/free-vpn-subscriptions/updates.html) • [状态面板](https://au1rxx.github.io/free-vpn-subscriptions/status.html) • [验证与兼容性](https://au1rxx.github.io/free-vpn-subscriptions/verification.html) • [FAQ](./docs/faq.md)

[Atom Feed](https://au1rxx.github.io/free-vpn-subscriptions/updates.xml) • [JSON Feed](https://au1rxx.github.io/free-vpn-subscriptions/updates.json) • [Discussions](https://github.com/Au1rxx/free-vpn-subscriptions/discussions)

[文档索引](./docs/README.md) • [项目细节](./docs/project-details.md) • [部署与运维](./docs/deployment.md) • [贡献指南](./CONTRIBUTING.md) • [安全策略](./SECURITY.md)

## 直接可用入口

- [Clash 远程订阅](https://raw.githubusercontent.com/Au1rxx/free-vpn-subscriptions/main/output/clash.yaml)
- [sing-box 远程订阅](https://raw.githubusercontent.com/Au1rxx/free-vpn-subscriptions/main/output/singbox.json)
- [V2Ray 远程订阅](https://raw.githubusercontent.com/Au1rxx/free-vpn-subscriptions/main/output/v2ray-base64.txt)
- [实时状态 JSON](https://raw.githubusercontent.com/Au1rxx/free-vpn-subscriptions/main/output/status.json)
- [最新快照下载页](https://github.com/Au1rxx/free-vpn-subscriptions/releases/latest)

## 为什么用户会 Star

- `不用重新找入口`：Star 之后，仓库首页就是长期入口。哪天链接、客户端、状态有变化，回来就能看到。
- `不是一条死链接`：这里不只给订阅，还持续给出实时状态、最近更新、客户端教程和故障排查。
- `越用越完整`：常见客户端页、FAQ、问题讨论会持续增加，老用户回访和新用户搜索都能吃到。

## 这个项目对用户到底有什么用

- 不知道该选 `Clash`、`sing-box` 还是 `V2Ray` 时，有格式判断页。
- 想先确认今天还能不能用时，有实时状态页。
- 想直接导入时，上面就有三个可直接使用的远程订阅入口。
- 导入失败时，有客户端教程、FAQ 和故障排查页，不用再去外面重新搜一圈。

## 项目文档

- 如果你想看仓库结构、公开产物、Pages 页面和生成文件的关系：看 [项目细节](./docs/project-details.md)。
- 如果你想看私有控制面如何同步到公共仓库、GitHub Pages / Releases / Feed / 验证链路怎么串起来：看 [部署与运维](./docs/deployment.md)。
- 如果你要总览现有技术文档和客户端指南：看 [文档索引](./docs/README.md)。

## 每次回来看什么

- `状态是否新鲜`：先看 [状态面板](https://au1rxx.github.io/free-vpn-subscriptions/status.html)。
- `今天有没有新快照`：看 [更新记录与快照历史](https://au1rxx.github.io/free-vpn-subscriptions/updates.html)。
- `分享出去的入口靠不靠谱`：看 [验证与兼容性页面](https://au1rxx.github.io/free-vpn-subscriptions/verification.html)。
- `想持续跟踪更新`：`Star` 仓库，或订阅 Atom / JSON Feed。

## 先看这三个信号

- 想确认公开侧是不是还在更新，先看 [状态面板](https://au1rxx.github.io/free-vpn-subscriptions/status.html) 和 `output/status.json` 的最近检查时间。
- 想看最近一次快照、手动下载入口和历史版本，直接看 [GitHub Releases](https://github.com/Au1rxx/free-vpn-subscriptions/releases)。
- 想确认分享出去的链接当前是否可达、哪些只是快照、哪些适合长期自动刷新，直接看 [验证与兼容性页面](https://au1rxx.github.io/free-vpn-subscriptions/verification.html)。
- 想持续跟踪更新，不要只收藏仓库首页；请优先 `Star`、到仓库 Watch 菜单里关注 `Releases`，或直接订阅 [Atom Feed](https://au1rxx.github.io/free-vpn-subscriptions/updates.xml) / [JSON Feed](https://au1rxx.github.io/free-vpn-subscriptions/updates.json)。

## 项目定位

- 面向 `Clash`、`sing-box`、`V2Ray` 用户的公开订阅分发仓库。
- 多地区节点，按小时检查健康，按计划刷新订阅文件。
- 文档和页面按 GitHub / Google 可发现性设计，兼顾搜索流量与导入转化。
- 公开仓库的完整结构和职责拆分见 [项目细节](./docs/project-details.md)。

## 获取方式

- 如果你只想快速上手：从 [站点首页](https://au1rxx.github.io/free-vpn-subscriptions/) 进入。
- 如果你要判断该选哪种格式：看 [我该使用哪种订阅格式](https://au1rxx.github.io/free-vpn-subscriptions/which-subscription-format-should-i-use.html)。
- 如果你要手动下载快照：看 [GitHub Releases](https://github.com/Au1rxx/free-vpn-subscriptions/releases)。
- 如果你要确认 `raw` 和 `release` 的区别：看 [验证与兼容性页面](https://au1rxx.github.io/free-vpn-subscriptions/verification.html)。

## 更新节奏

- 节点健康状态：每小时检查一次
- 公开订阅快照：按计划刷新
- Release 快照：跟随公开产物更新，适合历史回看和手动下载
- 更新 Feed：跟随最新 Release 自动刷新
- 端到端部署链路见 [部署与运维](./docs/deployment.md)。

## 为什么主订阅链接保持不变

- 固定订阅 URL 更适合客户端长期自动刷新，不会因为链接轮换导致旧教程、旧收藏和客户端配置失效。
- 用户回访不应依赖“换链接”，而应依赖“看更新”。这个项目把回访入口放在状态页、更新记录页和 Releases 快照历史中。
- 想确认最近有没有刷新、是否有历史快照、要不要手动下载，请直接看 [更新记录与快照历史](https://au1rxx.github.io/free-vpn-subscriptions/updates.html)。
- 想把更新接入 RSS 阅读器、自动化脚本或第三方页面，可以直接订阅 [Atom Feed](https://au1rxx.github.io/free-vpn-subscriptions/updates.xml) 或 [JSON Feed](https://au1rxx.github.io/free-vpn-subscriptions/updates.json)。

## 支持的客户端

- Clash Verge / Clash Meta / Mihomo 兼容客户端
- sing-box
- V2Ray 兼容客户端，包括 v2rayNG、NekoBox
- 可导入标准订阅链接的 iOS 客户端

## 快速开始

1. 先确认你的客户端需要哪种订阅格式。
2. 从首页或对应客户端教程进入，不要先盲贴原始链接。
3. 节点有变更时刷新订阅。
4. 排查问题前先看状态页。

## 你在这里能直接得到什么

- 可直接导入的 Clash、sing-box、V2Ray 三种远程订阅
- 当前节点状态与最近检查时间
- 最新快照、历史版本和手动下载入口
- 常见客户端教程
- 常见导入失败和连接失败的排障页面

## 客户端与教程

- [如何导入 Clash 订阅](./docs/clash-subscription.md)
- [如何导入 sing-box 订阅](./docs/sing-box-subscription.md)
- [如何导入 V2Ray 订阅](./docs/v2ray-subscription.md)
- [如何使用 Clash Verge Rev](./docs/clash-verge-rev.md)
- [如何使用 FlClash](./docs/flclash.md)
- [如何使用 Clash Meta for Android](./docs/clash-meta-android.md)
- [如何使用 Hiddify Next](./docs/hiddify-next.md)
- [如何使用 NekoBox](./docs/nekobox.md)
- [如何使用 v2rayNG](./docs/v2rayng.md)
- [如何使用 Shadowrocket](./docs/shadowrocket.md)
- [FAQ 与故障排查](./docs/faq.md)

## 搜索入口页

- [Clash 订阅不可用](https://au1rxx.github.io/free-vpn-subscriptions/clash-subscription-not-working.html)
- [V2Ray 订阅链接](https://au1rxx.github.io/free-vpn-subscriptions/v2ray-subscription-url.html)
- [Shadowrocket 订阅链接](https://au1rxx.github.io/free-vpn-subscriptions/shadowrocket-subscription-url.html)
- [Android 免费 VPN](https://au1rxx.github.io/free-vpn-subscriptions/free-vpn-for-android.html)
- [如何刷新 Clash 配置](https://au1rxx.github.io/free-vpn-subscriptions/how-to-refresh-clash-profile.html)
- [V2Ray 订阅不可用](https://au1rxx.github.io/free-vpn-subscriptions/v2ray-subscription-not-working.html)
- [Shadowrocket 无法连接](https://au1rxx.github.io/free-vpn-subscriptions/shadowrocket-not-connecting.html)
- [Clash 配置更新失败](https://au1rxx.github.io/free-vpn-subscriptions/clash-profile-update-failed.html)
- [iPhone 免费 VPN](https://au1rxx.github.io/free-vpn-subscriptions/free-vpn-for-iphone.html)
- [Android 最佳 Clash 客户端路径](https://au1rxx.github.io/free-vpn-subscriptions/best-clash-client-for-android.html)
- [故障排查总入口](https://au1rxx.github.io/free-vpn-subscriptions/troubleshooting-hub.html)
- [免费 VPN 订阅链接入口](https://au1rxx.github.io/free-vpn-subscriptions/free-vpn-subscription-links.html)
- [我该使用哪种订阅格式](https://au1rxx.github.io/free-vpn-subscriptions/which-subscription-format-should-i-use.html)
- [Clash 和 V2Ray 订阅对比](https://au1rxx.github.io/free-vpn-subscriptions/clash-vs-v2ray-subscription.html)
- [Windows 最佳客户端路径](https://au1rxx.github.io/free-vpn-subscriptions/best-vpn-client-for-windows.html)
- [Mac 最佳客户端路径](https://au1rxx.github.io/free-vpn-subscriptions/best-vpn-client-for-mac.html)

## 社区入口

- [提交安装或导入问题](https://github.com/Au1rxx/free-vpn-subscriptions/discussions/new?category=q-a)
- [建议新的客户端教程或落地页](https://github.com/Au1rxx/free-vpn-subscriptions/discussions/new?category=ideas)
- [分享可用的客户端配置](https://github.com/Au1rxx/free-vpn-subscriptions/discussions/new?category=show-and-tell)
- [报告公开链接或状态异常](https://github.com/Au1rxx/free-vpn-subscriptions/issues/new/choose)

## 常用客户端路径

- 桌面端：Clash Verge Rev、FlClash、sing-box desktop
- Android：Clash Meta for Android、v2rayNG、NekoBox、Hiddify Next、sing-box mobile
- iPhone / iPad：Shadowrocket 兼容订阅导入路径
- 多平台：Hiddify Next，覆盖 Android、iPhone、Windows、macOS、Linux

这些客户端页和问题页的目标，是同时提升搜索可见性、导入成功率和后续回访率。

## Star / Watch / 回流

如果这个仓库对你有帮助：

- 请点 `Star`，把它保存成持续更新入口，而不是只记住某个文件地址。
- 如果你想跟踪快照更新，请到仓库的 Watch 菜单里关注 `Releases`，或者直接订阅上面的 Atom / JSON Feed。
- 安装问题、客户端请求、兼容性反馈请优先走 `Discussions`。
- 只有在公开链接、Release 资产或 Pages 页面真实异常时，再提交 `Issues`。

## 说明

- 本仓库分发的是公开订阅资源，请把这些节点视为共享公共资源。
- 节点可用性和性能会随时间变化。
- 当前可用性的最准入口是状态页，而不是客户端本地缓存。

## 分享文案

- GitHub / 论坛短文案：持续更新的免费 VPN 公开入口，附带实时节点状态、更新记录、客户端教程和故障排查页面。
- 社群短文案：公开入口 + 实时状态 + 更新记录 + 多客户端教程，支持 Clash、v2rayNG、NekoBox、Shadowrocket 等常见导入路径。
- 站点入口：`https://au1rxx.github.io/free-vpn-subscriptions/`
- 详细文案文件：[社交分享文案](./docs/social-share-copy.zh-CN.md)
