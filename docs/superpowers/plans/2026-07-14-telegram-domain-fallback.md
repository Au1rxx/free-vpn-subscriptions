# Telegram 公开预览域名回退实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 当且仅当 `t.me` 返回 DNS NXDOMAIN 时，通过 `telegram.me` 恢复公开预览抓取和发现，同时保留原来源身份。

**Architecture:** 新增独立的 `internal/httpfallback` HTTP 执行器，先执行原请求，只对无正文 GET、精确 `t.me` 主机和 `net.DNSError.IsNotFound` 组合执行一次固定备用域名重试。来源抓取与公开发现共用该执行器，数据库 URL 和来源哈希不变。

**Tech Stack:** Go 1.25、`net/http`、`net.DNSError`、表驱动单元测试、现有 systemd 生产安装流程。

## Global Constraints

- 仅允许 `t.me` 无正文 GET 请求回退到固定 `telegram.me`，最多一次。
- HTTP 状态、DNS 超时、连接、TLS、其他域名和带正文请求不得回退。
- 保留 context、路径、查询和请求头，不修改 `sources.url` 或规范化来源哈希。
- 继续使用现有超时、重定向、正文大小、日志和稳定错误码边界。
- 严格执行 RED → GREEN → 全仓回归 → 生产只读验收。

---

### Task 1: 受限 Telegram DNS 回退器及调用方接入

**Files:**
- Create: `internal/httpfallback/httpfallback.go`
- Create: `internal/httpfallback/httpfallback_test.go`
- Modify: `internal/sources/fetch.go`
- Modify: `internal/discovery/discovery.go`

**Interfaces:**
- Consumes: `*http.Client`、`*http.Request` 和错误链中的 `*net.DNSError`。
- Produces: `func Do(client *http.Client, request *http.Request) (*http.Response, error)`。

- [x] **Step 1: 写精确回退与拒绝回退的失败测试**

测试使用自定义 `RoundTripper` 记录主机序列。核心成功用例：

```go
request, _ := http.NewRequest(http.MethodGet, "https://t.me/s/channel?q=1", nil)
request.Header.Set("If-None-Match", `"same"`)
response, err := Do(client, request)
// 断言 hosts == []string{"t.me", "telegram.me"}
// 断言第二次请求的路径、查询和 If-None-Match 不变，response.StatusCode == 200。
```

表驱动拒绝用例覆盖：非 `t.me` 的 NXDOMAIN、`t.me` 的超时、带正文 POST，以及原请求
直接返回 HTTP 404。每项断言调用次数为 1。

- [x] **Step 2: 运行测试并确认 RED**

Run: `go test ./internal/httpfallback -run 'TestDo' -count=1`

Expected: FAIL，原因是 `internal/httpfallback.Do` 尚不存在，而不是测试语法错误。

- [x] **Step 3: 写最小实现**

实现边界：

```go
func Do(client *http.Client, request *http.Request) (*http.Response, error) {
    response, err := client.Do(request)
    if err == nil || !eligible(request, err) {
        return response, err
    }
    retry := request.Clone(request.Context())
    copiedURL := *request.URL
    copiedURL.Host = "telegram.me"
    retry.URL = &copiedURL
    response, retryErr := client.Do(retry)
    if retryErr != nil {
        return nil, fmt.Errorf("t.me DNS lookup failed (%v); telegram.me fallback failed: %w", err, retryErr)
    }
    return response, nil
}
```

`eligible` 必须检查 `client != nil`、GET、`Body == nil`、`URL.Port() == ""`、
`strings.EqualFold(URL.Hostname(), "t.me")`，并通过 `errors.As` 确认
`*net.DNSError.IsNotFound`。

将 `internal/sources.FetchRaw` 和 `internal/discovery.fetch` 中的 `client.Do(req)` 精确替换
为 `httpfallback.Do(client, req)`，不改其他 HTTP 调用。

- [x] **Step 4: 运行定向和全仓验证**

Run:

```bash
go test ./internal/httpfallback ./internal/sources ./internal/discovery -count=1
go test ./... -count=1
go vet ./...
go test -race ./internal/httpfallback ./internal/sources ./internal/discovery -count=1
git diff --check
```

Expected: 全部 PASS，`go vet` 和 `git diff --check` 无输出。

- [x] **Step 5: 提交实现**

```bash
git add internal/httpfallback internal/sources/fetch.go internal/discovery/discovery.go
git commit -m "fix: fall back from unavailable telegram preview domain"
```

### Task 2: 生产安装与真实来源验收

**Files:**
- Modify: `docs/superpowers/plans/2026-07-14-telegram-domain-fallback.md`
- Runtime: `/opt/free-vpn-harvester/fnctl`
- Runtime: `/var/lib/free-vpn-harvester/harvest-last-run.log`

**Interfaces:**
- Consumes: Task 1 提交后的 `fnctl` 和私有仓库 `install-harvester.sh`。
- Produces: 真实 `t.me` 发现成功、Telegram 定时抓取恢复和有界生产证据。

- [x] **Step 1: 安装新二进制并保持现有服务参数**

Run:

```bash
sudo /home/ubuntu/worktrees/vpn-lab-node-platform/ops/feed-publisher/install-harvester.sh
sha256sum /opt/free-vpn-harvester/fnctl
systemctl is-active free-vpn-db-tunnel.service free-vpn-validate.service
```

Expected: 安装完成；隧道和验证服务均为 `active`。不重启正在稳定运行且代码路径未变化的
验证 Worker。

- [x] **Step 2: 通过 transient systemd credential 运行真实发现**

Run:

```bash
sudo systemd-run --wait --pipe --collect \
  -p User=ubuntu -p Group=ubuntu \
  -p LoadCredential=mysql-password:/etc/free-vpn-harvester/mysql-password.credential \
  /opt/free-vpn-harvester/fnctl -c /etc/free-vpn-harvester/config.yaml \
  discover --kind telegram --url https://t.me/s/V2RayRootFree --limit 200
```

Expected: 退出码 0，不再包含 `lookup t.me ... no such host`，并输出有界候选/来源摘要。

- [x] **Step 3: 触发一轮采集并验证来源失败分布**

Run: `sudo systemctl start free-vpn-harvest-fetch.service`；服务超过 60 秒时只轮询
`systemctl show`，完整输出保存在状态日志。

Expected: 服务最终 `Result=success`、`ExecMainStatus=0`；该轮 Telegram 来源不再集中增加
约 40 个 `http_failed`，`pending_fetches` 和配置数继续前进。

- [x] **Step 4: 验证整条链路没有回归**

Run: `db-status`、`ingest-status`、`validation-status`、服务状态和影子报告。

Expected: migration 8、表 22、TLS 和可写状态正常；验证服务零新增重启、过期租约 0；
数据库低于 35GB；影子验收除时间窗口外无新失败项。

- [x] **Step 5: 同步计划状态并提交证据**

勾选上述已完成步骤，记录生产 UTC 时间、发现结果和失败计数差值，然后运行：

```bash
git add docs/superpowers/plans/2026-07-14-telegram-domain-fallback.md
git commit -m "docs: record telegram fallback production acceptance"
```

## 生产验收证据

- `2026-07-14T04:10:40Z` 安装新二进制，SHA-256 为
  `1cfa6e699a0dcea7be674dbb04adbaa30d56bbdf3a5526747b122e0563b7b0f0`；
- 真实 `https://t.me/s/V2RayRootFree` 发现退出码为 0，得到 49 个候选并 upsert
  40 个来源；
- `2026-07-14T04:11:13Z` 至 `04:21:05Z` 完整采集轮次成功，抓取 11 个到期来源，
  成功 10、失败 1；唯一失败为 `body_too_large`；
- 轮次前后 `http_failed` 均为 84，没有新增 Telegram DNS 失败；总抓取从 538 增至
  549，成功从 315 增至 325，失败仅从 223 增至 224；
- 来源总数保持 346，说明备用域名没有改变来源身份或制造批量重复。
- `2026-07-14T04:23:34Z` 分类器补齐 26,343 条新增配置，未分类剩余为 0；
- `2026-07-14T04:24:11Z` 影子复核显示 migration 8、表 22、数据库可写且 TLS
  正常、分配量 4,129,226,752 字节、过期租约 0、可用节点 1,572、影子导出
  `singbox_failures=0`；唯一未满足项仍为 `shadow_window_lt_72h`。
