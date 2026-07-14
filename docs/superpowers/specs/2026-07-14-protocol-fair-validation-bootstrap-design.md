# 协议公平首验补位设计

**日期：** 2026-07-14

**状态：** 已按用户的自动实施授权确认

## 1. 背景与问题

生产验证 Worker 从 `validation_queue` 按优先级、到期时间和任务主键领取全局批次，
随后才在内存中对已领取任务执行协议轮询。大规模初始导入按来源和协议成批写入，导致
一个领取批次可能全部属于 SOCKS4；内存重排无法让尚未进入批次的 HTTP、HTTPS、
SOCKS5 获得最低首验份额。现场证据是全局已评分超过 36 万时，这三个协议首验仍为 0。

已批准的平台设计要求每种协议保留最低验证份额。修复必须避免新增 migration、避免
重启当前主 Worker、避免重置 72 小时影子起点，并继续复用现有租约和验证结果路径。

## 2. 目标与非目标

目标：

1. 为指定协议领取一批尚未首验且已经到期的任务。
2. 领取只锁定 `validation_queue`，与主 Worker 并发时使用 `SKIP LOCKED` 避免重复。
3. 复用现有 `validate-worker --once`、快速检查、真实代理验证、吞吐采样和结果持久化。
4. 生产用限时 transient sampler 为当前饥饿协议建立首验样本，不重启主 Worker。
5. 未设置过滤参数时保持当前全局领取行为和输出兼容。

非目标：

1. 不新增第 9 个 migration，不向队列表冗余协议字段。
2. 不替换主 Worker，不改变现有 1000 条批次和 150 并发参数。
3. 不承诺 sampler 单独完成全局 90% 覆盖；它消除协议零覆盖，主 Worker继续完成全量。
4. 不改变质量评分、TTL、公共导出或发布切换门槛。

## 3. 方案选择

采用“可选协议过滤 + 未首验过滤 + transient sampler”。另外两种方案不采用：

- 给 `validation_queue` 增加 `protocol` 和复合索引会使查询最直接，但需要新 migration、
  2.1M 行回填和生产 DDL，不符合当前 8 migration 与连续影子运行约束。
- 直接批量提高任务 `priority` 虽无需发版，但旧 Worker 持久化结果时不会复位优先级，
  可能让同一批失败任务永久压过普通任务。

## 4. 接口与数据流

`validate-worker` 增加一个可选参数：

- `--initial-protocol <name>`：只领取指定协议且 `attempts=0` 的未首验任务；值必须是
  节点模型支持的协议常量。

`validation.SQLQueue` 保存可选的首验协议。未设置时继续调用
`store.ClaimValidationJobs`；设置后调用新的
`store.ClaimInitialValidationJobsByProtocol`。Worker、Queue 接口和结果处理无需改变。

过滤领取在一个事务中执行：

1. 从 `node_configs` 强制使用现有 `idx_node_configs_export`，以
   `is_exportable=0 AND protocol=?` 做索引范围定位。
2. 通过 `node_config_id` 连接 `validation_queue`，要求 `attempts=0`、pending 且到期。
3. 以 `last_success_at,node_config_id` 稳定排序并有界 `LIMIT`。
4. 使用 `FOR UPDATE OF q SKIP LOCKED`，只锁队列表行。
5. 复用现有批量 lease 更新、受影响行数校验和事务提交逻辑。

`is_exportable=0` 对未首验任务成立：新配置默认不可导出，只有验证结果才能改变该字段。
它同时让现有复合索引按协议定位，无需生产 DDL。`attempts=0` 是最终首验边界；如果
进程在领取后崩溃，任务由主 Worker按现有过期租约路径恢复，不会被重新当成未领取任务。

## 5. 生产运行与回滚

部署采用原子替换磁盘二进制，但不重启 `free-vpn-validate.service`；其 PID、启动时间和
`NRestarts` 必须不变。新功能只由独立 transient systemd 服务调用，按
HTTP、HTTPS、SOCKS5 顺序执行 `validate-worker --once`，使用独立 validator ID、
小批次和低于主 Worker 的并发。每轮结束后重新检查目标协议首验数、过期租约、数据库
容量和主 Worker状态。

停止 transient 服务即可回滚；主 Worker和全局队列未变。若过滤领取出现错误，事务
回滚且 sampler 非零退出，不影响主 Worker。不得为部署或回滚修改影子起点。

## 6. 测试与验收

TDD 和验收覆盖：

1. 旧实现缺少过滤接口和 CLI 参数时目标测试先失败。
2. 临时 MySQL 中混合协议任务只返回目标协议。
3. `--initial-protocol` 不领取 `attempts>0`、未到期或已租约任务。
4. 两个并发 owner 使用 `SKIP LOCKED` 不领取同一任务。
5. 无过滤参数继续使用当前全局顺序和租约语义。
6. 无效的 `--initial-protocol` 必须在连接数据库前失败。
7. 全量 Go 测试、`go vet` 和相关 race 测试通过。
8. 生产 sampler 首轮使 HTTP、HTTPS、SOCKS5 的 `attempted` 均从 0 增长。
9. 主 Worker和隧道 PID、启动时间、重启数不变，过期租约保持 0。
10. 影子审计除既有时间窗口外不得出现新失败。

## 7. 已知限制

过滤查询为指定协议提供首验补位，不替代长期全局调度。初始积压清空后，新配置增量较小，
主 Worker的全局 FIFO和已有批内协议轮询可继续处理；若以后再次进行百万级批量导入，
可重用同一 sampler。若需要永久按比例调度，应在独立设计中增加队列协议列和索引，不能
在本次影子恢复中隐式扩展 schema。
