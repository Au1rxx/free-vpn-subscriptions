# 存储容量快照设计

## 目标

让每日非 dry-run 维护在 TTL、汇总和容量策略执行成功后，按表持久化一组 50GiB 容量快照，
从而为增长趋势、告警和 72 小时验收提供历史证据。dry-run 必须继续只读，不能写入快照。

## 当前事实与缺口

- migration 0006 已创建 `storage_metrics`，包含采样时间、schema、表、估算行数、数据/索引
  字节、50GiB 容量和使用率，唯一键为采样时间/schema/表。
- 生产表当前为 0 行；代码只有 `ReadStorageBytes`，没有任何 `INSERT storage_metrics`。
- `maintain.Service.Run` 会在非 dry-run 时先生成日汇总，再执行 TTL、过期配置、冷源和容量
  策略；`RunMaintenance` 最后重新读取 `AfterBytes`，但结果只打印到 CLI。
- 当前数据库约 5.38GB，远低于 35GB 影子门禁和 50GiB 上限；缺失的是历史趋势，不是当前
  容量告急。

## 决策

在 `internal/store` 新增 `RecordStorageMetrics`。函数使用一次 `INSERT ... SELECT` 从
`information_schema.tables` 采样当前 schema 的每张表。`maintain.Service.Run` 仅在
`dryRun=false` 且 `RunMaintenance` 成功后调用它；快照失败会使维护命令失败，不伪报成功。

不新增 migration、表、依赖或后台 service，不在每小时影子 dry-run 中写数据。每日维护仍是
唯一持久化入口，避免每小时 22 行无必要增长。

## 数据契约

每个采样时间对当前 schema 的每张表写一行：

- `sampled_at`：调用方传入的 UTC `now`；
- `table_schema`、`table_name`：`information_schema.tables` 当前 schema；
- `table_rows_estimate`：`COALESCE(table_rows,0)`；
- `data_bytes`、`index_bytes`：对应 information_schema 字段；
- `total_bytes`：数据与索引之和；
- `capacity_bytes`：`DatabaseCapacityBytes`，固定 `50 << 30`；
- `usage_percent`：维护后的数据库总占用 / 50GiB × 100，保留三位小数；每张表重复保存同一
  数据库水位，便于任意单表时间序列同时解释全局策略档位。

唯一键冲突时更新同一采样时间的行数、字节、容量和使用率，使显式重试幂等。函数拒绝
`capacityBytes=0`，返回写入/更新影响行数和受限错误。

## 调用顺序与错误

```text
ReadStorageBytes → PolicyForUsage → RollupDailyStats
→ RunMaintenance/AfterBytes → RecordStorageMetrics → CLI report
```

dry-run 在 `RunMaintenance` 返回后直接生成报告，`StorageMetricRows=0`。非 dry-run 快照失败时，
先前已完成的 TTL 操作不回滚，但 service 返回 `record storage metrics` 错误，systemd 将本轮标为
失败并在下次计划运行重试；这与现有非事务批量维护错误语义一致。

CLI 稳定摘要新增 `storage_metric_rows=N`，用于确认每日写入。现有字段顺序和含义保持不变。

## 测试

- 真实 MySQL 集成测试以唯一微秒时间调用 `RecordStorageMetrics`，断言写入行数等于当前 schema
  表数、`node_configs` 行存在、容量为 50GiB、使用率等于传入总字节比例；测试结束删除该
  采样时间全部行。
- 同一时间重复调用后行数不增加且值更新，证明幂等。
- `capacityBytes=0` 在访问数据库前返回错误。
- maintain 服务测试/生产验收确认 dry-run 不增加行数，非 dry-run 新增完整一组并输出
  `storage_metric_rows`。
- 执行 targeted integration、`go test ./...`、`go vet ./...` 和相关 race 测试。

## 部署与影子边界

构建新 `fnctl` 到临时文件，记录 SHA 后原子替换 `/opt/free-vpn-harvester/fnctl`。正在运行的
验证 Worker 保持原 inode 和进程，不重启；下一次 CLI/service 调用使用新版本。手动执行一次
非 dry-run 维护以产生首组快照，要求 TTL 删除计数为 0 或符合既有策略、`storage_metrics`
出现当前 22 张表的完整采样、核心 PID/重启次数和影子 epoch 不变。

该变更只增加容量遥测，不改变抓取、解析、分类、验证、质量分、TTL 阈值或发布输出，因此不
重置 72 小时影子起点。

## 回滚

原子恢复部署前二进制。新增容量快照是可保留的审计数据，无需删除；回滚后每日维护不再新增
快照，但 TTL 和其他既有维护行为恢复原状。
