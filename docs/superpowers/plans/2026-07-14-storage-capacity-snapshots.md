# 存储容量快照实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 让每次成功的非 dry-run 维护按当前 22 张表持久化 50GiB 容量快照，并在 CLI 中报告采样行数。

**Architecture:** `internal/store.RecordStorageMetrics` 以单条 `INSERT ... SELECT` 采集 information_schema 并按采样时间幂等；`maintain.Service.Run` 在 TTL 成功后调用，dry-run 跳过。构建后的 fnctl 原子部署，不重启长驻服务，手动维护生成首组生产历史。

**Tech Stack:** Go、database/sql、MySQL 9.7、现有真实 MySQL 集成测试、systemd。

**Status:** Active

**Progress:** 0/2 tasks complete (0%)

**Updated:** 2026-07-14

## Global Constraints

- 不新增或修改 migration、表、索引、TTL 阈值和 50GiB 上限。
- 每次非 dry-run 成功维护写当前 schema 每张表一行；dry-run 写入 0 行。
- 快照时间使用调用方 UTC `now`，相同时间重试必须幂等。
- `usage_percent` 使用维护后数据库总占用除以 `50 << 30`，保留三位小数。
- 快照失败必须使维护命令失败，不得伪报成功。
- 测试产生的采样必须清理；测试不得运行真实非 dry-run TTL 维护。
- 部署不得重启验证 Worker、数据库隧道或其他长驻服务，不得重置影子起点。
- 不启动子代理；在现有公共功能工作树顺序实施。

## Tasks

- [ ] **Task 1: 以真实 MySQL 集成测试实现幂等表级采样**
  - Status: pending
  - Files: `internal/store/storage_metrics.go`, `internal/store/storage_metrics_test.go`
  - Acceptance: 首次采样行数等于当前 schema 表数；字段、50GiB 容量和使用率正确；相同时间重试不增加行数并更新水位；零容量在访问数据库前失败；测试清理完整。
  - Verification: `VPN_NODE_TEST_CONFIG=/etc/free-vpn-harvester/config.yaml go test ./internal/store -run TestRecordStorageMetricsIntegration -count=1 -v`
  - Evidence: 记录 RED 未定义符号、GREEN 写入行数、幂等结果和测试前后生产表行数。

  - [ ] **Step 1: 写入失败集成测试**

    加载 `VPN_NODE_TEST_CONFIG`，以 20 秒 context 打开真实 MySQL。采样时间使用
    `time.Now().UTC().Add(24*time.Hour).Truncate(time.Microsecond)`，避免与生产维护冲突；defer
    按该时间删除测试行。调用：

    ```go
    written, err := RecordStorageMetrics(ctx, db, sampledAt, DatabaseCapacityBytes, 5<<30)
    ```

    断言 `written` 和该时间数据库行数等于 information_schema 当前表数；`node_configs` 的
    schema/name/capacity/usage 为 `vpn_nodes`、`node_configs`、`50<<30`、`10.000`。再次以
    `10<<30` 调用，行数不增加且 usage 更新为 `20.000`。零容量返回错误。

  - [ ] **Step 2: 运行 RED**

    Run:

    ```bash
    CREDENTIALS_DIRECTORY=/tmp/vpn-node-test-credential \
    VPN_NODE_TEST_CONFIG=/etc/free-vpn-harvester/config.yaml \
      go test ./internal/store -run TestRecordStorageMetricsIntegration -count=1 -v
    ```

    Expected: 编译失败，`RecordStorageMetrics` 未定义。

  - [ ] **Step 3: 实现最小采样函数**

    新文件导出：

    ```go
    func RecordStorageMetrics(ctx context.Context, db *sql.DB, sampledAt time.Time,
        capacityBytes, usageBytes uint64) (int64, error)
    ```

    `capacityBytes==0` 返回错误。SQL 从 `information_schema.tables WHERE table_schema=DATABASE()`
    插入 9 个业务字段，使用 `ON DUPLICATE KEY UPDATE` 更新行数、字节、容量、使用率；错误包装
    为 `record storage metrics`。返回 `RowsAffected()`。

  - [ ] **Step 4: 运行 GREEN、重复和清理检查**

    运行 targeted integration 两次；每次测试前后查询 `storage_metrics` 总行数必须相同，确保
    defer 清理。执行 `go test ./internal/store -count=1` 和 `go vet ./internal/store`。

  - [ ] **Step 5: 提交存储层**

    ```bash
    git add internal/store/storage_metrics.go internal/store/storage_metrics_test.go
    git commit -m "feat: persist storage capacity snapshots"
    ```

- [ ] **Task 2: 接入维护服务、部署并生成首组生产快照**
  - Status: pending
  - Files: `internal/maintain/service.go`, `cmd/fnctl/maintain.go`, `cmd/fnctl/maintain_test.go`, `docs/superpowers/plans/2026-07-14-storage-capacity-snapshots.md`
  - Acceptance: dry-run 报告采样 0；非 dry-run 在维护成功后记录完整表组并报告行数；全量测试通过；生产首组恰好覆盖 22 张表且核心运行态不变。
  - Verification: `go test ./... -count=1 && go vet ./...`
  - Evidence: 记录全量测试、二进制 SHA、生产维护输出、快照时间/行数/水位、PID 和影子 epoch。

  - [ ] **Step 1: 写入 service/CLI 失败合同**

    新建 `cmd/fnctl/maintain_test.go`，构造 `maintain.Report{StorageMetricRows: 22}`，调用
    `writeMaintenanceReport(&buffer, report)`，要求输出含 `storage_metric_rows=22`。先运行
    `go test ./cmd/fnctl -run TestWriteMaintenanceReport -count=1`，旧代码因字段和函数未定义失败。

  - [ ] **Step 2: 在非 dry-run 成功路径记录快照**

    `RunMaintenance` 成功后：

    ```go
    metricRows := int64(0)
    if !dryRun {
        metricRows, err = store.RecordStorageMetrics(
            ctx, s.DB, now, store.DatabaseCapacityBytes, stored.AfterBytes)
        if err != nil {
            return Report{}, err
        }
    }
    ```

    把 `metricRows` 写入 Report。把 CLI 现有 JSON/单行输出提取为
    `writeMaintenanceReport(io.Writer, maintain.Report) error`，单行末尾追加
    `storage_metric_rows=%d`，命令调用该 helper。

  - [ ] **Step 3: 运行全量验证**

    Run:

    ```bash
    gofmt -w internal/store/storage_metrics.go internal/store/storage_metrics_test.go \
      internal/maintain/service.go cmd/fnctl/maintain.go cmd/fnctl/maintain_test.go
    go test ./... -count=1
    go vet ./...
    go test -race ./internal/store ./internal/maintain ./cmd/fnctl -count=1
    ```

    Expected: 全部退出 0。

  - [ ] **Step 4: 构建并原子部署 fnctl**

    构建临时二进制、记录 SHA；备份当前 `/opt/free-vpn-harvester/fnctl`，安装 `.fnctl.new` 后
    同文件系统 `mv -f`。部署前后记录隧道/验证 Worker PID、`NRestarts` 和影子 epoch，不重启
    任何长驻 service。

  - [ ] **Step 5: 验证 dry-run 只读并执行一次生产维护**

    先记录 `storage_metrics` 行数，运行 `fnctl maintain --dry-run`，确认行数不变且
    `storage_metric_rows=0`。随后 `systemctl start --no-block free-vpn-maintain.service`，后台等待
    结果；要求 Result=success、TTL/异常计数符合策略、`storage_metric_rows=22`。

  - [ ] **Step 6: 验证首组容量历史和运行态**

    SQL 要求最新 `sampled_at` 恰好 22 行、22 个不同表、capacity 全为 53,687,091,200、使用率
    一致且约 10%；`node_configs` 等主要表字段非零。核心 PID/重启次数、影子 epoch、严格影子
    门禁和 legacy 发布状态不变。

  - [ ] **Step 7: 提交服务接入和计划证据**

    ```bash
    git add internal/maintain/service.go cmd/fnctl/maintain.go \
      docs/superpowers/plans/2026-07-14-storage-capacity-snapshots.md
    git commit -m "feat: record capacity after maintenance"
    ```

## Risks And Blockers

- `information_schema.table_rows` 是估算值；容量字节仍来自引擎统计，趋势用途足够，不能作为
  精确业务行数。
- `RowsAffected` 在幂等更新时可能按更新语义大于表数；生产使用唯一时间首次写入应为 22，测试
  的幂等断言以持久行数而非第二次 RowsAffected 为准。
- 非 dry-run 维护会执行既有 TTL；部署前 dry-run 必须为 0 异常行，否则不手动启动。
- 72 小时时间条件尚未到达；容量快照就绪不批准数据库发布。

## Change Log

- 2026-07-14: 从已批准容量快照设计创建；初始进度 0/2。

## Rollback

恢复部署前 fnctl 二进制。已经写入的容量快照保留作为审计证据，不修改 schema 或删除生产数据。

## Completion

完成标准是生产首组容量历史、dry-run 不变量、全量测试和核心运行态全部验证。总体目标继续等待
48/72 小时门槛，公共发布保持 legacy。
