# 来源解析背压实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 阻止同一来源在旧正文尚未解析时重复加入大型快照，使生产解析积压有界并可排空。

**Architecture:** 由 `internal/store.ClaimDueSources` 作为唯一领取边界，用相关 `NOT EXISTS` 子查询跳过已有待解析成功正文的来源。解析终态自然解除背压，不增加新状态或 migration。

**Tech Stack:** Go、database/sql、MySQL 9.7、现有集成测试、systemd。

## Global Constraints

- 不改变来源优先级、`next_fetch_at` 或 `FOR UPDATE SKIP LOCKED` 并发语义。
- 不增加 migration，生产仍保持 8 个已验收 migration。
- 不重启数据库隧道或验证 Worker，不改变权威影子起点。
- 集成测试使用独立临时数据库，不改动生产来源的 `next_fetch_at`。

---

### Task 1: 来源级待解析背压

**Files:**
- Create: `internal/store/sources_backpressure_test.go`
- Modify: `internal/store/sources.go`

**Interfaces:**
- Consumes: `ClaimDueSources(ctx context.Context, db *sql.DB, limit int) ([]SourceRecord, error)`
- Produces: 待解析成功正文存在时的来源级领取背压。

- [x] **Step 1: 写隔离 MySQL 集成测试**

测试创建随机临时数据库及最小 `sources`、`source_fetches` 表，写入一个到期来源和一条
`fetch_state='success', parse_state='pending'` 记录。断言首次领取返回 0 个；把
`parse_state` 更新为 `success` 后，断言第二次领取返回该来源。测试通过
`VPN_NODE_TEST_CONFIG` 读取连接信息，`defer` 删除临时数据库。

- [x] **Step 2: 运行 RED**

Run: `go test ./internal/store -run TestClaimDueSourcesDefersSourceWithPendingParse -count=1 -v`

Expected: FAIL，报告有待解析正文时仍领取了 1 个来源。

- [x] **Step 3: 实现最小 SQL 修复**

把 `ClaimDueSources` 的外层表设为别名 `s`，在到期条件后增加：

```sql
AND NOT EXISTS (
  SELECT 1 FROM source_fetches f
  WHERE f.source_id=s.source_id
    AND f.fetch_state='success'
    AND f.parse_state='pending'
)
```

保留现有排序、批次上限、行锁和领取后 `next_fetch_at` 更新。

- [x] **Step 4: 运行 GREEN 和全量回归**

Run:

```bash
go test ./internal/store -run TestClaimDueSourcesDefersSourceWithPendingParse -count=1 -v
go test ./...
go vet ./...
go test -race ./internal/store ./internal/sources
```

Expected: 全部退出码 0，无失败或 race。

- [x] **Step 5: 提交修复**

```bash
git add internal/store/sources.go internal/store/sources_backpressure_test.go
git commit -m "fix: backpressure sources awaiting parse"
```

### Task 2: 生产部署与积压验收

**Files:**
- Modify: `/opt/free-vpn-harvester/fnctl`
- Modify: `/home/ubuntu/worktrees/vpn-lab-node-platform/docs/vpn-node-data-platform-acceptance.md`

**Interfaces:**
- Consumes: Task 1 已验证的 `fnctl` 二进制。
- Produces: 不重启长驻 Worker 的原子部署，以及至少两个自动采集轮次的积压证据。

- [x] **Step 1: 构建并原子安装二进制**

Run:

```bash
go build -o /tmp/fnctl-source-backpressure ./cmd/fnctl
sha256sum /tmp/fnctl-source-backpressure /opt/free-vpn-harvester/fnctl
sudo install -o root -g root -m 0755 /tmp/fnctl-source-backpressure /opt/free-vpn-harvester/.fnctl.new
sudo mv -f /opt/free-vpn-harvester/.fnctl.new /opt/free-vpn-harvester/fnctl
```

Expected: 新旧哈希已记录，原子 rename 成功；已运行的验证 Worker 仍持有旧 inode，不执行重启。

- [x] **Step 2: 手工触发一轮并验证背压**

等当前轮次自然结束，再手工启动 `free-vpn-harvest-fetch.service`。轮次结束后运行
`ingest-status`，并查询每个来源的 pending 计数。

Expected: 轮次失败 0；已有 pending 正文的来源不再新增抓取；单来源
pending 最大值不再增长。

- [x] **Step 3: 等待第二个自动轮次并确认趋势**

Expected: 总 `pending_fetches` 低于部署基线，`body_too_large` 和抓取失败无新增，
数据库隧道与验证 Worker 仍 active/running、`NRestarts=0`。

- [x] **Step 4: 记录验收并提交**

```bash
git -C /home/ubuntu/worktrees/vpn-lab-node-platform add docs/vpn-node-data-platform-acceptance.md
git -C /home/ubuntu/worktrees/vpn-lab-node-platform commit -m "docs: record source parse backpressure recovery"
```
