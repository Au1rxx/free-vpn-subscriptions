# 来源解析背压设计

## 背景与根因

生产每 15 分钟领取到期来源，但单个 20–35MB 正文的解析和持久化可能跨越多个
9 分钟解析预算。`ClaimDueSources` 只检查 `next_fetch_at`，没有检查同一来源是否已有
待解析的成功正文。因此同一大型来源会在旧快照完成前反复加入新快照，造成 FIFO
队列增长和后续小正文阻塞。

生产证据显示：队头 `sevcator-trojan` 正文为 11,091,782 字节、70,486 个有效条目，
检查点为 `processed_nodes=25000`；积压中同时存在多份同源 `ebrasha-all`、
`sevcator-vmess` 的 20–35MB 快照。

## 方案

在 `internal/store.ClaimDueSources` 的数据库领取边界增加来源级背压：

- 如果同一 `source_id` 存在 `fetch_state='success' AND parse_state='pending'`，不再领取该来源。
- 正文解析成功或进入其他终态后，来源按原 `next_fetch_at`、优先级和 `SKIP LOCKED`
  规则恢复领取。
- 背压仅限当前来源；其他渠道继续抓取，不设全局积压阈值。
- 不增加 migration。现有 `idx_source_fetches_source_time(source_id, started_at)` 可先定位来源，
  而 `source_fetches` 又受 30 天 TTL 约束；当前数量下不需要改变已验收的 8 个 migration。

## 一致性与失败处理

背压只依赖已持久化状态，不引入内存锁或额外租约。解析进程被超时终止时，
`source_fetches.parse_state` 保持 `pending`，因而继续阻止同源新快照；下轮解析从
`parse_runs.error_summary=processed_nodes=N` 检查点续跑。不影响抓取失败或 HTTP 304 记录，
因为它们的 `parse_state` 为 `skipped`。

## 验收

- 隔离测试数据库中，存在成功且待解析正文时，到期来源不被领取。
- 把该正文标记为 `success` 后，同一来源立即可再次领取。
- 全量 Go 测试、`go vet ./...` 和目标 race 测试通过。
- 部署后新抓取轮次失败为 0，同一来源待解析记录不再增长，总积压转为下降。
- 数据库隧道和验证 Worker 不重启，不重置 72 小时影子起点。
