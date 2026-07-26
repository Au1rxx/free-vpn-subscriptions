# 实时 SEO 与 Star History 实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use `superpowers:executing-plans` to implement this plan task-by-task. 本仓库规则禁止子代理，因此只允许当前会话内联执行。

**状态：** 进行中
**进度：** 已完成 5/6 项（83%）
**更新时间：** 2026-07-26

**Goal:** 用仓库内每日更新的 SVG 修复 Star History，并让现有静态站点展示每小时真实快照和 30 天滚动趋势。

**Architecture:** 现有发布器继续拥有节点数据和 Pages 生成；`docs/data/network-history.json` 保存最多 721 个端点，用于覆盖完整 720 小时区间。独立 GitHub Action 使用仓库 `GITHUB_TOKEN` 每日生成 `assets/star-history.svg`，不接触节点源或数据库。

**Tech Stack:** Go 标准库、Cobra、`html/template`、GitHub REST API、GitHub Actions、静态 HTML/SVG/JSON。

## 范围

- 修复全部 README 的 Star History，新增每日自动更新的仓库内 SVG。
- 在既有首页和国家页增加当前数据、协议/国家分布及 24h/7d/30d 趋势。
- 修正 JSON-LD、FAQ 和 sitemap 的事实性问题，更新架构文档和生成产物。

## 非目标

- 不新增文章、日期归档、城市页、分析服务、数据库表或客户端 JavaScript。
- 不调整节点筛选、验证、排序、订阅格式、定时发布频率或现有 URL。

## 假设与决策

- 已批准使用“单一实时报告页 + 30 天滚动趋势”，不生成每日独立报告。
- Star 数据由 GitHub Action 每日读取；节点历史由现有每小时发布器维护。
- GitHub Pages 和 README 仍由生成器控制，生成文件不得手工修改。

## Global Constraints

- 不新增第三方 Go、前端或 Action 依赖；页面保持零 JavaScript、零遥测。
- 不改变订阅 URL、页面 URL、节点格式、hreflang 或现有外部 systemd 发布架构。
- 不生成日期归档页或批量关键词页；动态历史只保留一个 30 天滚动文件。
- 历史故障必须降级为当前快照且不得阻塞订阅发布；Action 故障保留上一张 SVG。
- 所有行为变更先记录 RED，再实施 GREEN；每项结束运行精确验证并单独提交。

## 测试策略与命令

- quick：`go test ./internal/starhistory ./internal/readme ./internal/pages ./cmd/fnctl -count=1`
- 单元测试：上述命令覆盖 API 分页、SVG、历史去重/截断/趋势、HTML、JSON-LD 和 sitemap。
- 集成测试：`go test ./... -count=1` 覆盖两个站点生成入口与现有输出契约。
- E2E：本地 `python3 -m http.server` 后用 `curl` 检查首页、中文首页、国家页、SVG、JSON 和 sitemap；上线后重复关键 HTTP 检查。
- CI：新增 workflow 通过 `go test ./internal/starhistory ./internal/readme` 后才允许生成/提交 SVG。
- full：`go build ./... && go vet ./... && go test ./... -count=1`
- 不适用：没有浏览器运行时和数据库 schema 变化，不新增浏览器自动化或数据库 migration 测试。

---

### Task 1: Star History API 与确定性 SVG

- [x] **任务状态：已完成**

**Files:**
- Create: `internal/starhistory/starhistory.go`
- Create: `internal/starhistory/starhistory_test.go`
- Create: `cmd/fnctl/star_history.go`
- Create: `cmd/fnctl/star_history_test.go`
- Modify: `cmd/fnctl/main.go`

**Interfaces:**
- Produces: `starhistory.Fetch(ctx context.Context, client *http.Client, apiBase, repo, token string) ([]time.Time, error)`
- Produces: `starhistory.RenderSVG(repo string, stars []time.Time) ([]byte, error)`
- Produces: `newStarHistoryCmd() *cobra.Command`

- [x] **Step 1: 写 API 分页和 SVG RED 测试**
  - `httptest.Server` 返回两页 `starred_at`，断言携带 `Accept: application/vnd.github.star+json` 和 Bearer token。
  - 断言重复/乱序时间按 UTC 日期累计，SVG 包含 `<svg`、仓库标题、最终星数且两次输出字节一致。
  - Run: `go test ./internal/starhistory -count=1`
  - Expected: FAIL，包或函数尚不存在。

- [x] **Step 2: 实现最小 API 与 SVG**
  - 校验 `owner/repo`、token 和 HTTP 状态；跟随 RFC 5988 `Link rel="next"`，限制最多 100 页。
  - SVG 仅使用转义后的文本、折线、坐标轴和 `<title>/<desc>`；不嵌入当前时间。
  - Run: `go test ./internal/starhistory -count=1`
  - Expected: PASS。

- [x] **Step 3: 写 CLI RED 测试并实现**
  - 命令：`fnctl star-history --repo owner/repo --output assets/star-history.svg`。
  - token 只从 `GITHUB_TOKEN` 读取；缺失 token、无效 repo、API 非 2xx 必须返回错误；输出通过临时文件加 rename。
  - 在 `newRootCmd` 注册 `newStarHistoryCmd()`。
  - Run: `go test ./cmd/fnctl -run StarHistory -count=1`
  - Expected: 先 FAIL，实施后 PASS。

- [x] **Step 4: 验证并提交**
  - Run: `go test ./internal/starhistory ./cmd/fnctl -count=1`
  - Commit: `feat: add deterministic star history generator`
  - 证据：2026-07-26 targeted tests 退出 0；另验证分页跨 origin 会在携带 token 前失败。

### Task 2: README 本地图片与每日 Action

- [x] **任务状态：已完成**

**Files:**
- Modify: `internal/readme/generator.go`
- Modify: `internal/readme/generator_test.go`
- Create: `.github/workflows/update-star-history.yml`
- Create: `assets/star-history.svg`

**Interfaces:**
- Consumes: `fnctl star-history`。
- Produces: 7 个 README 均引用 `<RepoURL>/raw/main/assets/star-history.svg`。

- [x] **Step 1: 写 README RED 测试**
  - 对所有 locale 断言存在仓库内 SVG 和 `/stargazers` 链接，不含 `api.star-history.com`。
  - Run: `go test ./internal/readme -run Star -count=1`
  - Expected: FAIL，当前仍引用第三方 API。

- [x] **Step 2: 修改生成器并生成首张 SVG**
  - README Star History 图片使用仓库 raw URL，非 GitHub repo 继续完全省略该区块。
  - 使用当前 Au1rxx 凭据以环境变量方式运行新命令，凭据不进入参数、日志或文件。
  - Run: `go test ./internal/readme -run Star -count=1`
  - Expected: PASS，`file assets/star-history.svg` 识别为 SVG。

- [x] **Step 3: 新增每日 workflow**
  - `schedule: "17 3 * * *"` + `workflow_dispatch`，`permissions: contents: write`，同名 `concurrency`。
  - 使用 `actions/checkout@v4`、`actions/setup-go@v5`，先运行 targeted tests，再以 `${{ github.token }}` 生成 SVG。
  - 仅有 diff 时提交；`git pull --rebase origin main` 后普通 push，禁止 force。
  - 用 PyYAML 解析 workflow 并断言 schedule、权限、测试和生成步骤存在。

- [x] **Step 4: 验证并提交**
  - Run: `go test ./internal/readme ./internal/starhistory ./cmd/fnctl -count=1`
  - Commit: `fix: self-host the star history chart`
  - 证据：2026-07-26 targeted tests 退出 0；workflow 结构断言通过，SVG XML 可解析且重复生成字节一致。

### Task 3: 30 天历史与趋势计算

- [x] **任务状态：已完成**

**Files:**
- Create: `internal/pages/history.go`
- Create: `internal/pages/history_test.go`
- Create: `cmd/fnctl/page_history.go`
- Create: `cmd/fnctl/page_history_test.go`
- Modify: `cmd/fnctl/main.go`
- Modify: `cmd/fnctl/export_db.go`
- Modify: `internal/pages/pages.go`

**Interfaces:**
- Produces: `pages.HistoryPoint`，字段为 `GeneratedAt time.Time`、`Selected`、`Verified`、`MedianLatencyMS`、`Countries`。
- Produces: `pages.TrendDelta{Available bool; Selected int; Verified int; MedianLatencyMS int}`，整数为当前值减基线值。
- Produces: `pages.TrendSummary{Hours24 TrendDelta; Days7 TrendDelta; Days30 TrendDelta}`。
- Produces: `pages.UpdateHistory(path string, current HistoryPoint) ([]HistoryPoint, error)`。
- Produces: `pages.BuildTrends(points []HistoryPoint, now time.Time) TrendSummary`。
- Extends: `pages.Input.History []HistoryPoint`。

- [x] **Step 1: 写历史 RED 测试**
  - 覆盖缺失文件初始化、同一 UTC 小时去重、乱序排序、未来点拒绝、未知 schema、损坏 JSON、721 端点截断和原子写入。
  - 覆盖 24h/7d/30d 最近不晚于目标时刻的点，以及历史不足的 `Available=false`。
  - Run: `go test ./internal/pages -run 'History|Trend' -count=1`
  - Expected: FAIL，类型和函数尚不存在。

- [x] **Step 2: 实现历史模型与计算**
  - JSON 顶层固定 `schema_version: 1`；UTC RFC3339；`Countries` 只统计达到 `MinPerCountry` 的国家。
  - 临时文件权限 `0644`，成功 `Sync/Close/Rename`；错误绝不覆盖旧文件。
  - Run: `go test ./internal/pages -run 'History|Trend' -count=1`
  - Expected: PASS。

- [x] **Step 3: 接入两个发布入口并验证降级**
  - `preparePageHistory` 捕获历史错误，向 stderr 输出单行 `warning: page history: ...`，并返回仅含当前点的 slice。
  - legacy `writeOutputs` 与 DB `renderSite` 均在调用 `pages.Generate` 前准备历史。
  - 测试损坏文件时站点仍生成、原文件未覆盖、趋势显示数据积累中。

- [x] **Step 4: 验证并提交**
  - Run: `go test ./internal/pages ./cmd/fnctl -count=1`
  - Commit: `feat: retain rolling network history`
  - 证据：2026-07-26 targeted tests 退出 0；覆盖 721 端点上限、同小时去重、趋势边界和损坏文件降级后继续生成页面。

### Task 4: 实时落地页与国家页内容

- [x] **任务状态：已完成**

**Files:**
- Modify: `internal/pages/pages.go`
- Modify: `internal/pages/templates.go`
- Modify: `internal/pages/l10n.go`
- Modify: `internal/pages/pages_test.go`

**Interfaces:**
- Produces: `metricRow{Name string; Count int; Percent string}`、`TrendSummary` 和安全的内联趋势 SVG。
- Extends: `pageCtx`，增加 `ProtocolRows`、`TopCountries`、`Trends`、`TrendSVG`、`CountryProtocols`。

- [x] **Step 1: 写 7 语言页面 RED 测试**
  - 首页必须展示当前快照、真实协议分布、前 8 国家、趋势文字、更新时间、验证方法、限制和选择指南。
  - 国家页必须展示该国家协议组成；所有百分比合计允许四舍五入误差不超过 0.2%。
  - Run: `go test ./internal/pages -run 'Live|Protocol|Locale' -count=1`
  - Expected: FAIL，模板尚无这些区块。

- [x] **Step 2: 实现数据派生与零 JS 图表**
  - 从同一次 `Summary/Selected` 排序协议和国家；相同数量按稳定字典序。
  - 趋势 SVG 使用 `template.HTML` 前必须只由数值和固定标签构造，并包含可访问 title/desc。
  - 历史不足显示本地化“数据积累中”，不得渲染虚假 delta。

- [x] **Step 3: 完成全部本地化和样式**
  - 为 `en/zh/ja/ko/es/pt/ru` 填满新增字段，复用 README 已审阅的验证与限制文案。
  - 扩展现有 CSS 的 card/table/grid，保持移动端可读、无外部资源、首屏订阅链接位置不后移。
  - Run: `go test ./internal/pages -count=1`
  - Expected: PASS。

- [x] **Step 4: 验证并提交**
  - Run: `go test ./internal/pages ./internal/readme -count=1`
  - Commit: `feat: add live network report to pages`
  - 证据：2026-07-26 targeted tests 退出 0；七种语言均覆盖实时报告，协议占比精确合计 100.0%，订阅卡仍先于报告区块。

### Task 5: 合规结构化数据与 sitemap

- [x] **任务状态：已完成**

**Files:**
- Modify: `internal/pages/pages.go`
- Modify: `internal/pages/pages_test.go`
- Modify: `internal/pages/l10n.go`
- Modify: `ARCHITECTURE.md`

- [x] **Step 1: 写 JSON-LD/sitemap RED 测试**
  - 解析 JSON-LD，断言 `WebSite`、`Dataset`、3 个 `DataDownload`、真实 `dateModified`，且不存在 `aggregateRating`。
  - 断言首页/国家页 sitemap 有准确 lastmod，指南条目没有本轮动态 lastmod。
  - 断言 7 语言 FAQ 不再声称 GitHub Actions 聚合。
  - Run: `go test ./internal/pages -run 'JSONLD|Sitemap|FAQ' -count=1`
  - Expected: FAIL。

- [x] **Step 2: 实现并更新架构文档**
  - 用 `Dataset` 替换 `SoftwareApplication`；下载 URL/MIME、免费状态、许可证和语言与页面一致。
  - `writeSitemapEntry` 仅在 lastmod 非空时输出标签；指南传空值。
  - `ARCHITECTURE.md` 更新 SEO 表、滚动历史和 Star Action，删除评分富结果说法。

- [x] **Step 3: 验证并提交**
  - Run: `go test ./internal/pages -count=1`
  - Commit: `fix: make page metadata evidence-based`
  - 证据：2026-07-26 pages 全包测试退出 0；JSON-LD 可解析且含 3 个真实下载，静态指南 sitemap 条目不再伪造动态 `lastmod`。

### Task 6: 全量验证、生成产物与上线

- [ ] **任务状态：进行中**

**Files:**
- Generated: `README*.md`
- Generated: `docs/*.html`
- Generated: `docs/sitemap.xml`
- Generated: `docs/data/network-history.json`
- Update: 本计划的状态、证据、变更记录和完成情况。

- [x] **Step 1: 全量本地门禁**
  - Run: `go build ./... && go vet ./... && go test ./... -count=1`
  - Expected: exit 0，无失败或跳过的确定性测试。
  - 证据：合并结果 `42b892651ce0` 提交前，build、vet、全测试、工作树与暂存 diff-check 均退出 0。

- [x] **Step 2: 生成并做本地 E2E**
  - 使用现有 DB export/site-root 发布入口生成所有 README 和 Pages，不直接编辑生成文件。
  - 启动临时静态服务器，检查 `/`、`/index.zh.html`、一个国家页、`/sitemap.xml`、
    `/data/network-history.json` 和仓库 SVG；检查 JSON-LD 可解析且 HTML 无第三方运行时请求。
  - 证据：真实数据库导出成功；6 个目标均 HTTP 200，旧国家页为 404，JSON-LD/XML/JSON 解析及零第三方脚本断言通过。

- [x] **Step 3: 完整 diff 与独立审阅**
  - 检查生成文件规模、无凭据/上游 URL、无无关改动；执行 `code-review` 并解决 Critical/Important。
  - 再运行受影响 targeted tests 和 full 门禁。
  - 证据：审阅发现并修复 30 天窗口需 721 个端点、生成 HTML 尾随空格、同小时重复点和陈旧国家页四项问题；复审无剩余 Critical/Important。

- [ ] **Step 4: 推送与线上验证**
  - 普通 push `main`，禁止 force；等待 GitHub Pages 和一次 publisher 更新。
  - 验证 Star SVG、首页动态区块、中文页、国家页、sitemap、history JSON 均 HTTP 200。
  - 手动触发一次 Star workflow，确认无变化时 no-op、有变化时只提交 SVG。

- [ ] **Step 5: 完成计划**
  - 填写每项验证证据，进度更新为 6/6；记录 workflow、Pages、publisher 状态和剩余风险。

## 风险与阻塞项

- GitHub Action 与外部发布器可能同时 push：workflow 只普通 rebase/push，冲突即失败，不覆盖 publisher。
- 30 天趋势从上线时开始积累：24h/7d/30d 依次出现，页面明确显示数据积累中。
- `contents: write` 可能被仓库 Actions 设置限制：若 push 被拒绝，需要 Au1rxx 明确启用 Workflow
  read/write 权限；在此之前保留首张 SVG，不影响主站。
- 翻译字段增加容易漏项：沿用反射完整性测试，任一空字符串使测试失败。

## 变更记录

- 2026-07-26：根据已批准的实时报告方案与设计文档创建。
- 2026-07-26：完成 Tasks 1–5、本地验证、真实生成、E2E 与合并自审；上线验证进行中。

## 回滚

普通 revert 各任务提交；停用 `update-star-history.yml`；恢复旧模板后重新运行生成器。历史 JSON 和
SVG 均为可删除的派生产物，没有数据库、订阅格式或远程服务迁移。

## 完成情况

本地实现与合并验证已完成；等待普通推送、GitHub Pages、Star workflow 和外部 publisher 线上验证。
