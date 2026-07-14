# 协议公平首验补位实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking. 本项目规则禁止启动子代理，因此仅允许当前会话内联执行。

**Goal:** 在不新增 migration、不重启主验证 Worker、不重置 72 小时影子起点的前提下，让 HTTP、HTTPS、SOCKS5 等指定协议安全领取尚未首验的到期任务，并通过低并发 systemd oneshot 建立真实连通性、延迟和吞吐质量样本。

**Architecture:** `internal/store` 在现有租约事务上增加按协议、`attempts=0` 的过滤领取；`validation.SQLQueue` 根据可选字段选择旧全局领取或新过滤领取；`validate-worker` 在打开配置和数据库之前校验 `--initial-protocol`。私有运维仓库提供受资源约束的模板 unit，生产部署只原子替换磁盘二进制并运行独立 oneshot，已运行的主 Worker 保持旧 inode 和原 PID。

**Tech Stack:** Go 1.25、`database/sql`、MySQL 9.7 `FOR UPDATE OF q SKIP LOCKED`、Cobra、Bash、systemd。

**Status:** Completed
**Progress:** 5/5 tasks complete (100%)
**Updated:** 2026-07-14

## Global Constraints

- 生产数据库继续保持 8 个 migration 和 22 张业务表，不增加或修改 schema。
- 未设置 `--initial-protocol` 时，`ClaimValidationJobs` 的排序、过期租约恢复、批次上限和输出完全兼容。
- 过滤领取只允许模型中已声明的协议，只领取 `attempts=0`、`pending`、已到期且 `is_exportable=0` 的任务。
- 使用 `FOR UPDATE OF q SKIP LOCKED` 只锁 `validation_queue`；领取、批量 lease 更新和提交保持同一事务。
- 不重启 `free-vpn-validate.service`，不修改影子起点；主 Worker 的 PID、启动时间和 `NRestarts` 必须保持不变。
- transient 首验实例使用独立 validator ID、`--limit 100`、`--concurrency 20`，按 HTTP、HTTPS、SOCKS5 串行运行。
- 集成测试只创建随机临时数据库并在 cleanup 中删除，禁止对 `vpn_nodes` 写入测试夹具。

---

## Objective

目标协议在生产库中的 `attempted_nodes` 从 0 增长，且每个新质量结果都走现有快速检查、真实代理验证、性能采样和持久化路径；主 Worker、数据库、采集和影子审计无回归。

## Scope

- 公共仓库的协议枚举校验、过滤领取 SQL、SQLQueue 路由、CLI 参数和测试。
- 私有仓库的 `free-vpn-validate-initial@.service`、安装器契约测试和验收文档。
- 无重启生产部署、HTTP/HTTPS/SOCKS5 小批量首验以及链路复核。

## Non-Goals

- 不新增 `validation_queue.protocol` 列或永久比例调度 migration。
- 不调整主 Worker 的 `--limit 1000 --concurrency 150`。
- 不改变评分、TTL、导出、发布切换或 72 小时验收门槛。
- 不要求一次 oneshot 清空各协议全部积压。

## Assumptions And Decisions

- `node_configs.idx_node_configs_export(is_exportable, protocol, last_success_at)` 已在生产存在，过滤查询从该索引定位目标协议。
- `attempts=0` 是首验边界；任务一旦成功领取，现有更新会先把 attempts 加一，即使进程随后异常也由主 Worker 的过期租约路径恢复。
- 生产初始补位使用可停止的独立 oneshot；停止该实例即可回滚运行行为，不触碰主 Worker。

## Tasks

### Task 1: 过滤领取事务

- [x] **Task 1: 为指定协议领取未首验任务**
  - Status: completed
  - Files: `internal/store/validation_queue.go`, `internal/store/validation_queue_test.go`
  - Interfaces:
    - Consumes: `ClaimValidationJobs(ctx context.Context, db *sql.DB, owner string, limit int, lease time.Duration)` 的现有租约语义。
    - Produces: `ClaimInitialValidationJobsByProtocol(ctx context.Context, db *sql.DB, owner, protocol string, limit int, lease time.Duration) ([]ValidationJob, error)`。
  - Acceptance: 混合协议夹具只领取目标协议；排除 attempts>0、未到期、leased 和 exportable 行；并发 owner 不重复领取；旧接口测试保持通过。
  - Verification: `VPN_NODE_TEST_CONFIG=/tmp/vpn-node-validation-test.yaml go test ./internal/store -run 'TestClaimInitialValidationJobsByProtocol' -count=1 -v`
  - Evidence: `2026-07-14` RED 仅因新函数未定义而失败；GREEN 在 MySQL 9.7 临时库通过过滤和双 owner 并发测试（4 个唯一任务、0 重复），旧参数校验测试通过；cleanup 后匹配测试库数量为 0。

  - [x] **Step 1: 写 RED 集成测试和临时库夹具**

    在 `validation_queue_test.go` 增加读取 `VPN_NODE_TEST_CONFIG` 的临时库 helper，创建最小 `node_configs`、`validation_queue`、`validation_attempts` 表和精确的 `idx_node_configs_export`。测试插入目标协议与旁路协议任务，并调用尚不存在的：

    ```go
    jobs, err := ClaimInitialValidationJobsByProtocol(
        ctx, db, "initial-http-a", node.ProtoHTTP, 2, time.Minute,
    )
    ```

    断言返回任务全部为 HTTP、领取后 state 为 leased、owner 正确、attempts 从 0 变为 1。另一个测试同时启动两个 owner，各领取 2 条，断言 ID 集合无交集且总数为 4。cleanup 必须 `DROP DATABASE IF EXISTS` 随机库。

  - [x] **Step 2: 运行 RED**

    Run:

    ```bash
    VPN_NODE_TEST_CONFIG=/tmp/vpn-node-validation-test.yaml \
      go test ./internal/store -run 'TestClaimInitialValidationJobsByProtocol' -count=1 -v
    ```

    Expected: 编译失败，报告 `undefined: ClaimInitialValidationJobsByProtocol`；失败原因不得是数据库连接或夹具语法。

  - [x] **Step 3: 提取共享 lease 事务并实现过滤 SQL**

    `ClaimValidationJobs` 和新入口先执行相同 owner/limit/lease/db 校验，再把查询及参数交给共享事务 helper。新查询必须保持完整扫描字段和 `PerformanceDue` 子查询，并使用：

    ```sql
    FROM node_configs n FORCE INDEX (idx_node_configs_export)
    JOIN validation_queue q ON q.node_config_id=n.node_config_id
    WHERE n.is_exportable=0
      AND n.protocol=?
      AND q.attempts=0
      AND q.job_state='pending'
      AND q.next_attempt_at <= UTC_TIMESTAMP(6)
    ORDER BY n.last_success_at ASC, n.node_config_id ASC
    LIMIT ? FOR UPDATE OF q SKIP LOCKED
    ```

    共享 helper 继续批量把 owner、UTC lease 截止时间、state 和 attempts 写回，校验 affected rows 等于领取数，并只在提交成功后修改返回对象。

  - [x] **Step 4: 运行 GREEN 与旧接口回归**

    Run:

    ```bash
    VPN_NODE_TEST_CONFIG=/tmp/vpn-node-validation-test.yaml \
      go test ./internal/store -run 'TestClaimInitialValidationJobsByProtocol|TestValidationQueue' -count=1 -v
    ```

    Expected: PASS；临时数据库已被 cleanup 删除。

### Task 2: 协议校验、SQLQueue 路由与 CLI

- [x] **Task 2: 增加 `--initial-protocol` 并在数据库连接前拒绝非法值**
  - Status: completed
  - Files: `pkg/node/node.go`, `pkg/node/node_test.go`, `internal/validation/worker.go`, `internal/validation/worker_test.go`, `cmd/fnctl/validate_worker.go`, `cmd/fnctl/validate_worker_test.go`
  - Interfaces:
    - Produces: `node.IsSupportedProtocol(protocol string) bool`。
    - Produces: `validation.SQLQueue{DB: db, InitialProtocol: protocol}`；空字符串走旧领取，非空走 Task 1 新入口。
    - Produces: Cobra 字符串参数 `--initial-protocol`。
  - Acceptance: 所有 13 个模型协议通过校验，未知值失败；CLI 在读取不存在的 config 前返回 unsupported protocol；空参数行为不变。
  - Verification: `go test ./pkg/node ./internal/validation ./cmd/fnctl -count=1`
  - Evidence: helper 与 flag 测试先分别因 undefined 和 unknown flag 失败；实现后 `pkg/node`、`internal/validation`、`cmd/fnctl` 全部通过。13 个协议常量均被识别，非法协议在不存在的 config 被读取前返回明确错误，空 flag 默认兼容旧路径。

  - [x] **Step 1: 写协议与 CLI RED 测试**

    `pkg/node/node_test.go` 表驱动覆盖所有 `Proto*` 常量和 `invalid`；`cmd/fnctl/validate_worker_test.go` 临时把 `cfgPath` 指向不存在的文件，执行：

    ```go
    command := newValidateWorkerCmd()
    command.SetArgs([]string{"--initial-protocol", "invalid"})
    err := command.ExecuteContext(context.Background())
    ```

    断言错误包含 `unsupported initial protocol` 而不是配置文件错误，并断言 flag 存在且默认值为空。

  - [x] **Step 2: 运行 RED**

    Run: `go test ./pkg/node ./cmd/fnctl -run 'TestIsSupportedProtocol|TestValidateWorkerInitialProtocol' -count=1`

    Expected: FAIL，原因是 helper/flag 尚不存在。

  - [x] **Step 3: 写最小实现并接入 SQLQueue**

    `IsSupportedProtocol` 用 switch 覆盖 `node.go` 中所有协议常量。`validate-worker` 在 `openIngestService` 之前执行：

    ```go
    if initialProtocol != "" && !node.IsSupportedProtocol(initialProtocol) {
        return fmt.Errorf("unsupported initial protocol %q", initialProtocol)
    }
    ```

    构造 `validation.SQLQueue{DB: db, InitialProtocol: initialProtocol}`；`SQLQueue.Claim` 仅在字段非空时调用新 store 入口，否则调用旧入口。

  - [x] **Step 4: 运行 GREEN 和 Worker 回归**

    Run:

    ```bash
    go test ./pkg/node ./internal/validation ./cmd/fnctl -count=1
    gofmt -w pkg/node/node.go pkg/node/node_test.go internal/store/validation_queue.go \
      internal/store/validation_queue_test.go internal/validation/worker.go \
      internal/validation/worker_test.go cmd/fnctl/validate_worker.go \
      cmd/fnctl/validate_worker_test.go
    git diff --check
    ```

    Expected: 全部退出码 0，格式化后 diff check 无输出。

### Task 3: 公共仓库全量验证与审阅

- [x] **Task 3: 完成回归、竞态、静态检查与差异审阅**
  - Status: completed
  - Files: Task 1/2 的全部公共仓库文件、本计划文件。
  - Acceptance: 全仓测试、vet、相关 race 测试全部通过；无 migration 变化；审阅无 Critical/Important 问题。
  - Verification: 下述命令全部退出 0。
  - Evidence: `go test ./...`、`go vet ./...`、4 个相关包 race、带真实 MySQL 的过滤领取 race、`git diff --check` 全部退出 0；migration diff 为空。代码审阅无 Critical/Important 发现。生产只读 EXPLAIN 使用 `idx_node_configs_export` 与 `uk_validation_queue_node_stage` 且无 filesort；HTTP/HTTPS/SOCKS5 到期未首验候选分别为 502,196/466,710/499,694。

  - [x] **Step 1: 运行全量和竞态验证**

    ```bash
    go test ./... -count=1
    go vet ./...
    go test -race ./internal/store ./internal/validation ./cmd/fnctl ./pkg/node -count=1
    git diff --check
    git diff --name-only -- db/migrations
    ```

    Expected: 前四项退出 0；最后一项无输出。

  - [x] **Step 2: 按代码审阅技能检查完整 diff**

    检查事务失败路径、行锁范围、输入校验顺序、旧路径兼容、上下文取消和测试是否真正执行临时 MySQL。Critical/Important 发现必须修复并重新运行本任务验证。

  - [x] **Step 3: 提交公共仓库实现**

    ```bash
    git add internal/store/validation_queue.go internal/store/validation_queue_test.go \
      internal/validation/worker.go internal/validation/worker_test.go \
      cmd/fnctl/validate_worker.go cmd/fnctl/validate_worker_test.go \
      pkg/node/node.go pkg/node/node_test.go \
      docs/superpowers/plans/2026-07-14-protocol-fair-validation-bootstrap.md
    git commit -m "feat: bootstrap validation by protocol"
    ```

### Task 4: 私有运维模板与安装契约

- [x] **Task 4: 版本化低并发首验 oneshot**
  - Status: completed
  - Files: `/home/ubuntu/worktrees/vpn-lab-node-platform/ops/feed-publisher/systemd/free-vpn-validate-initial@.service`, `/home/ubuntu/worktrees/vpn-lab-node-platform/ops/feed-publisher/install-harvester.sh`, `/home/ubuntu/worktrees/vpn-lab-node-platform/ops/feed-publisher/tests/check-harvester-units_test.sh`
  - Acceptance: 模板只执行一次、加载 systemd credential、使用 `%i` 协议及独立 validator ID、limit 100/concurrency 20、有 30 分钟硬上限；安装器包含该模板；契约测试通过。
  - Verification: `bash ops/feed-publisher/tests/check-harvester-units_test.sh && ops/feed-publisher/install-harvester.sh --dry-run`
  - Evidence: 契约先因模板不存在 RED；新增模板和安装器接入后通过。`systemd-analyze verify` 随后发现 `RuntimeMaxSec` 对 oneshot 无效，强化否定测试并改用 `TimeoutStartSec=30min` 后，unit 契约、dry-run、Bash 语法、systemd verify、diff check 均退出 0 且无警告；私有提交 `abe967c`。

  - [x] **Step 1: 先扩展契约测试并确认 RED**

    在 required 数组增加模板；断言 `Type=oneshot`、`TimeoutStartSec=30min`、不存在无效的 `RuntimeMaxSec`、credential、`--once --initial-protocol=%i --limit 100 --concurrency 20`、`validator-id=ai-a1-initial-%i`。运行测试，Expected: 因 unit 不存在而 FAIL。

  - [x] **Step 2: 创建受限 unit 并加入安装器**

    unit 复用主 Worker 的 tunnel 依赖、安全沙箱和状态目录，但设置 `Type=oneshot`、有效限制 oneshot 执行期的 `TimeoutStartSec=30min`、`CPUQuota=100%`、`MemoryMax=2G`、`TasksMax=512`，不设置 Restart。安装器的 systemd 文件集合必须包括 `free-vpn-validate-initial@.service`。

  - [x] **Step 3: 运行 GREEN 并提交私有仓库变更**

    ```bash
    bash ops/feed-publisher/tests/check-harvester-units_test.sh
    ops/feed-publisher/install-harvester.sh --dry-run
    git diff --check
    git add ops/feed-publisher/systemd/free-vpn-validate-initial@.service \
      ops/feed-publisher/install-harvester.sh \
      ops/feed-publisher/tests/check-harvester-units_test.sh
    git commit -m "ops: add bounded protocol validation sampler"
    ```

### Task 5: 无重启生产部署与质量验收

- [x] **Task 5: 补测 HTTP/HTTPS/SOCKS5 并证明全链路无回归**
  - Status: completed
  - Files: Runtime `/opt/free-vpn-harvester/fnctl`, `/etc/systemd/system/free-vpn-validate-initial@.service`; evidence `/home/ubuntu/worktrees/vpn-lab-node-platform/docs/vpn-node-data-platform-acceptance.md` 与本计划。
  - Acceptance: 三个协议 attempted_nodes 均增长；主 Worker PID/启动时间/NRestarts 不变；过期租约 0；数据库健康、采集服务和影子审计无新增失败。
  - Verification: 生产只读聚合 SQL、`systemctl show`、`db-status`、`validation-status`、shadow report。
  - Evidence: `12:00:17Z` 原子部署后主 Worker PID `1009053`、`10:14:13Z` 启动时间、`NRestarts=0`、运行 inode `525102` 和旧进程哈希全部不变；磁盘新哈希 `03c5d52c2215517757c74ae41b2b49ff54cf34bdc132318ad0c796400d0fc9d6`。HTTP/HTTPS/SOCKS5 oneshot 各领取并持久化 100 条，均退出 0 且 persist_errors=0；300/300 状态有评分明细和验证时间，HTTP 的 2 个 degraded 均有延迟且 2 次吞吐成功。过期租约 0，数据库 8/22、TLS、可写、约 5.01GB；`12:01:19Z` 影子审计唯一失败仍为 `shadow_window_lt_72h`。

  - [x] **Step 1: 记录部署前基线并构建**

    记录 UTC、主 Worker/tunnel PID、启动时间、NRestarts、二进制 SHA-256、三个协议 attempted 数、过期租约、数据库字节数和最新影子失败项。构建 `/tmp/fnctl-protocol-bootstrap` 并运行 `sha256sum`。

  - [x] **Step 2: 原子安装二进制和模板但不重启主 Worker**

    先把新二进制安装为 `/opt/free-vpn-harvester/.fnctl.new`，再同文件系统 rename 到 `fnctl`；只安装新模板并 daemon-reload。随后再次比较主 Worker PID、启动时间、NRestarts 和 `/proc/<pid>/exe` inode，必须不变。

  - [x] **Step 3: 串行后台启动三个 oneshot**

    按 `http`、`https`、`socks5` 顺序执行 `systemctl start --no-block free-vpn-validate-initial@<protocol>.service`。每个实例必须自然结束且 `Result=success, ExecMainStatus=0` 后才能启动下一个；完整日志留在 journal，只用 `journalctl -n 80 --no-pager` 查看有界尾部。

  - [x] **Step 4: 验证质量与系统健康**

    聚合查询每协议总数、distinct attempted、passed、available/degraded、latency 非空和 throughput 非空；至少 attempted>0，且成功节点的延迟字段非空、性能成功有 bytes_per_second。确认 expired leases=0、隧道和主 Worker active、数据库 migration=8/tables=22/TLS/read_only=0/低于 35GB；采集解析继续前进；影子审计除 `shadow_window_lt_72h` 外无失败。

  - [x] **Step 5: 记录证据、提交并继续 72 小时门禁**

    将 UTC、前后计数、哈希、PID、不变式和 shadow 结果写入私有验收文档及本计划 Evidence，提交文档。Task 5 完成只表示协议补位完成，不表示总目标完成；72 小时硬指标仍由主生产计划持续跟踪。

## Risks And Blockers

- 过滤查询若未使用目标索引会扩大数据库负载；部署前用生产 `EXPLAIN` 复核 key、rows 和 filesort，异常时停止，不运行 sampler。
- oneshot 中途退出会留下已增加 attempts 的 lease；现有过期租约恢复可接管，验收必须确认 expired leases 回到 0。
- 原子替换磁盘文件不影响已运行进程，但安装器若被误用可能覆盖其他 unit；本任务只安装新模板和二进制，不调用会重启服务的流程。
- 若目标协议没有到期 attempts=0 任务，claimed=0 不算功能失败；先用只读 SQL确认候选数，必要时等待 `next_attempt_at`，不得篡改生产优先级。

## Change Log

- 2026-07-14: 从已批准的 `docs/superpowers/specs/2026-07-14-protocol-fair-validation-bootstrap-design.md` 创建；明确公共代码、私有 systemd 和生产无重启验收边界。
- 2026-07-14: 完成公共仓库 TDD、MySQL 并发验证、全仓回归和生产 EXPLAIN；无 schema 变化。
- 2026-07-14: `systemd-analyze verify` 证明 `RuntimeMaxSec` 对 oneshot 无效；计划改为由 `TimeoutStartSec=30min` 提供可执行硬上限，并强化否定契约断言。
- 2026-07-14: 完成无重启生产部署和三个饥饿协议的真实首验；质量字段、租约、数据库、采集和影子审计均通过。

## Rollback

1. `systemctl stop free-vpn-validate-initial@http.service free-vpn-validate-initial@https.service free-vpn-validate-initial@socks5.service`。
2. 删除或停用新模板不会影响 `free-vpn-validate.service`；主 Worker持续使用启动时的旧 inode。
3. 若磁盘新二进制存在缺陷，从部署前备份原子恢复；不重启主 Worker，先运行 CLI 只读状态命令验证。
4. 不修改 validation_queue 数据、不回滚 migration、不重置影子起点；已完成验证结果作为正常历史保留。

## Completion

- 5 个顶层任务均已通过 TDD、真实 MySQL 集成、race、静态检查、代码审阅、systemd 契约和生产验收。
- 公共实现提交 `de8859491`，私有运维提交 `abe967c`；生产磁盘二进制和模板已安装，主 Worker 未重启。
- 残余风险是目标协议的初始样本相对百万级积压仍小；后续仅以同一受限 sampler 扩大覆盖，不改变本计划接口和 schema。
- 总平台目标仍要求 72 小时影子硬指标；本计划完成不作为发布切换依据。
