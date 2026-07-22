-- fnctl:statement
-- 标记验证明细删除前已完成且后续不可被部分历史覆盖的每日汇总。
ALTER TABLE `node_daily_stats`
  ADD COLUMN `finalized_at` DATETIME(6) NULL COMMENT '验证明细删除前完成汇总的时间' AFTER `error_counts`,
  ALGORITHM=INPLACE,
  LOCK=NONE;

-- fnctl:statement
-- 为按统计日期生成最终每日汇总提供范围索引，避免扫描全部验证明细。
ALTER TABLE `validation_attempts`
  ADD INDEX `idx_validation_attempts_rollup` (`started_at`, `node_config_id`)
    COMMENT '按开始时间和节点生成删除前每日汇总',
  ALGORITHM=INPLACE,
  LOCK=NONE;

-- fnctl:statement
-- 为验证批次的有界 TTL 清理提供按时间和主键排序的覆盖索引。
ALTER TABLE `validation_batches`
  ADD INDEX `idx_validation_batches_cleanup` (`started_at`, `validation_batch_id`)
    COMMENT '按开始时间和主键有界清理验证批次',
  ALGORITHM=INPLACE,
  LOCK=NONE;

-- fnctl:statement
-- 将高速增长明细调整为适合 50GiB 上限的平衡分层 TTL，并纳入验证批次。
INSERT INTO `storage_policies`
  (`policy_key`, `data_class`, `normal_ttl_days`, `warning_ttl_days`, `warning_percent`, `critical_percent`, `action`, `description`)
VALUES
  ('raw_payloads', 'raw_payload', 30, 14, 70, 90, 'delete_expired', '原始响应按容量水位保留30、14、7、3、1天'),
  ('parse_errors', 'parse_error', 30, 14, 70, 90, 'delete_expired', '解析错误按容量水位保留30、14、7、3、1天'),
  ('validation_attempts', 'validation_attempt', 14, 7, 70, 90, 'rollup_then_delete', '验证尝试先汇总再按容量水位保留14、7、3、2、1天'),
  ('validation_batches', 'validation_batch', 14, 7, 70, 90, 'rollup_then_delete', '无明细引用的验证批次按容量水位保留14、7、3、2、1天'),
  ('source_fetches', 'source_fetch', 90, 60, 70, 90, 'rollup_then_delete', '来源抓取先汇总再按容量水位保留90、60、30、14、7天'),
  ('export_members', 'export_member', 30, 14, 70, 90, 'rollup_then_delete', '导出成员按容量水位保留30、14、7、3、1天'),
  ('cold_source_breaker', 'source', NULL, NULL, 90, 94, 'pause_cold_sources', '容量达到90%时暂停低价值冷源并加速明细TTL')
ON DUPLICATE KEY UPDATE
  `data_class` = VALUES(`data_class`),
  `normal_ttl_days` = VALUES(`normal_ttl_days`),
  `warning_ttl_days` = VALUES(`warning_ttl_days`),
  `warning_percent` = VALUES(`warning_percent`),
  `critical_percent` = VALUES(`critical_percent`),
  `action` = VALUES(`action`),
  `enabled` = TRUE,
  `description` = VALUES(`description`);
