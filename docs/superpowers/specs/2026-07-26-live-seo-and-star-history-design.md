# 实时 SEO 与 Star History 设计

## 目标

修复 README 中永久返回 404 的 Star History 图片，并在不批量制造低价值页面的前提下，
把现有落地页升级为可持续更新的实时网络报告。最终页面应向用户提供真实数据、自然覆盖核心
搜索意图，并保持现有零 JavaScript、静态托管和每小时发布架构。

## 当前事实与约束

- `api.star-history.com` 对本仓库返回 `GitHub restricted starred-data access`；Au1rxx 的
  authenticated GitHub API 可以读取全部 stargazer 时间线。
- 当前站点已有 7 种语言、国家页、客户端指南、canonical、hreflang、sitemap 和每小时生成；
  扩大相似页面数量会增加 Google `scaled content abuse` 风险。
- Google 不使用 `meta keywords` 参与网页排名；关键词必须自然出现在可见且有帮助的正文中。
- sitemap 的 `lastmod` 只有持续准确反映显著内容变化时才有价值；当前指南页被每小时标记更新，
  与实际内容不符。
- 首页 JSON-LD 中的 `AggregateRating 4.6/47` 没有可见评价来源，属于潜在误导标记。
- 上游抓取和节点验证继续由外部 systemd 运行，不迁移到 GitHub Actions。

## 总体决策

采用分工式自动维护：

```text
现有每小时发布器 → 当前 Summary → 30 天滚动历史 → 首页、国家页、sitemap
每日 GitHub Action → authenticated stargazers API → assets/star-history.svg
```

不新增每日文章、日期归档页、分析脚本、第三方前端依赖、数据库 migration 或外部统计服务。

## Star History

在 `fnctl` 增加独立子命令，使用 `GITHUB_TOKEN` 请求
`/repos/Au1rxx/free-vpn-subscriptions/stargazers`，携带
`application/vnd.github.star+json`，完成分页、按 UTC 日期累计并生成确定性 SVG。

SVG 写入 `assets/star-history.svg`，不嵌入生成时间，只有星标数据变化时文件才变化。README
生成器的 7 种语言全部引用仓库内图片并链接到 GitHub stargazers 页面，彻底移除
`api.star-history.com`。

新增每日 scheduled + `workflow_dispatch` GitHub Action，权限限定为 `contents: write`。
Action 构建 `fnctl`、生成 SVG、仅在文件变化时提交；推送前重新获取 `main` 并 rebase，冲突时
失败而不强推，下一次计划运行重试。它不读取节点源、数据库或服务器凭据。

## 30 天滚动历史

新增 `docs/data/network-history.json`，schema 固定为：

```json
{"schema_version":1,"points":[{"generated_at":"2026-07-26T05:01:14Z","selected":2000,"verified":7866,"median_latency_ms":584,"countries":57}]}
```

`pages` 包提供独立的历史加载/更新函数；两个发布入口在调用 `pages.Generate` 前使用它，并把
结果通过 `pages.Input` 传入。函数按 UTC 小时去重追加当前点、按时间排序并保留最近
721 个端点（当前点加 720 小时前基线），写入采用同目录临时文件加 rename。文件不存在时从当前点初始化；`countries` 明确是
满足 `MinPerCountry`、实际拥有国家页的国家数量。

未知 schema、损坏 JSON 或未来时间点不会覆盖原文件：命令层输出一条有界 warning，并让页面
以当前点和“数据积累中”降级渲染。SEO 历史故障不得阻止核心订阅发布。

趋势计算选择不晚于目标时刻的最近点，输出 24 小时、7 天、30 天变化；历史不足时明确显示
“数据积累中”，不伪造零变化。

## 落地页内容

首页在订阅链接之后增加以下可见内容，全部由现有 7 语言本地化：

1. **当前网络快照**：精选数、已验证数、国家数、中位延迟和准确更新时间。
2. **实时协议分布**：协议、节点数和占比，覆盖 Clash、sing-box、v2ray、VLESS、Reality、
   Trojan、Shadowsocks、Hysteria2 等真实搜索意图。
3. **热门国家**：前 8 个国家及数量，链接到现有国家页。
4. **7/30 天趋势**：服务端生成的零 JS SVG、文字变化摘要和数据覆盖时间。
5. **验证方法与局限**：真实说明 TCP、TLS、配置检查、HTTP-over-proxy、排序和客户端复测。
6. **选择指南**：按客户端、协议和国家进入已有订阅或指南，不新建关键词门页。

国家页增加该国家当前协议组成和更新时间，使正文与其他国家页具有真实数据差异。所有动态
内容来自同一次 `Summary`/`Selected`，不得出现跨运行数字。

## SEO 与结构化数据

- 保留现有 title、description、canonical、hreflang、Open Graph 和内部链接。
- 不扩充 `meta keywords`；关键词只进入自然语言标题、说明、数据表和链接文本。
- 删除虚构的 `SoftwareApplication.aggregateRating`，改用真实的
  `WebSite + Dataset + DataDownload + FAQPage`。
- `Dataset.dateModified` 使用本轮生成时间，三种 `DataDownload` 分别指向 Clash、sing-box、
  v2ray 输出，并声明真实 MIME 类型、免费访问和 MIT 许可证。
- 修正 FAQ 中“GitHub Action 每小时运行”为“外部调度器每小时运行”。
- 首页和国家页保留准确的小时级 `lastmod`；指南 sitemap 条目省略不准确的动态 `lastmod`。
- 不增加隐藏文字、重复城市页、自动文章或仅为搜索引擎生成的内容。

## 错误处理与兼容

历史文件异常时保留原文件、记录 warning，并用当前快照降级渲染。Star History Action 失败
只保留上一张 SVG，不影响节点发布。现有订阅 URL、页面 URL、hreflang 和输出格式不变。页面
继续无 JavaScript、无遥测、无第三方运行时请求。

## 测试与验收

- Star API 分页、日期累计、XML 转义和确定性 SVG 单元测试；README 断言只引用本地 SVG。
- 历史初始化、去重、乱序、721 端点截断、损坏/schema 错误和 24h/7d/30d 计算单元测试。
- Pages 测试覆盖 7 种语言的实时区块、协议比例、国家链接和历史不足状态。
- JSON-LD 解析测试断言存在 Dataset/DataDownload、没有 AggregateRating，数据与页面一致。
- sitemap 测试断言首页/国家页 lastmod 准确、指南不伪造本轮 lastmod。
- 运行 `go build ./...`、`go vet ./...`、`go test ./... -count=1`，并本地静态服务抽查首页、
  中文首页、国家页、sitemap 和 SVG。
- 上线后验证 Star 图片 HTTP 200、PVR/社区文件不受影响、GitHub Pages 关键页面可访问。

## 发布与回滚

先发布生成器、历史初始点和静态 SVG，再启用每日 Action；下一次每小时发布自动开始积累历史。
若 Action 没有写权限，仅停用该 workflow 并保留最后 SVG。整体回滚使用普通 revert，恢复旧
模板并停止历史写入；没有数据库或订阅格式需要迁移。

## 验收条件

- README 不再出现第三方 Star History 404，图表可从仓库稳定加载并每日自动刷新。
- 首页展示真实、可读、可索引的当前数据和滚动趋势，不生成新的日期归档页。
- 结构化数据不存在虚构评分，FAQ 与实际调度架构一致。
- 所有既有 URL 和订阅契约保持兼容，完整测试与线上抽查通过。
