-- fnctl:statement
CREATE TABLE IF NOT EXISTS `storage_metrics` (
  `storage_metric_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '容量采样内部主键',
  `sampled_at` DATETIME(6) NOT NULL COMMENT '容量采样时间',
  `table_schema` VARCHAR(64) NOT NULL COMMENT '数据库名称',
  `table_name` VARCHAR(64) NOT NULL COMMENT '数据表名称',
  `table_rows_estimate` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '数据行数估计值',
  `data_bytes` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '表数据占用字节数',
  `index_bytes` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '表索引占用字节数',
  `total_bytes` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '表和索引总字节数',
  `capacity_bytes` BIGINT UNSIGNED NOT NULL DEFAULT 53687091200 COMMENT '数据库容量上限字节数',
  `usage_percent` DECIMAL(6,3) NOT NULL DEFAULT 0 COMMENT '容量使用百分比',
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '记录创建时间',
  PRIMARY KEY (`storage_metric_id`),
  UNIQUE KEY `uk_storage_metrics_sample_table` (`sampled_at`, `table_schema`, `table_name`),
  KEY `idx_storage_metrics_table_time` (`table_name`, `sampled_at`),
  KEY `idx_storage_metrics_usage` (`usage_percent`, `sampled_at`)
) ENGINE=InnoDB COMMENT='数据库表和索引容量采样';

-- fnctl:statement
CREATE TABLE IF NOT EXISTS `storage_policies` (
  `policy_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '存储策略内部主键',
  `policy_key` VARCHAR(64) NOT NULL COMMENT '存储策略唯一名称',
  `data_class` VARCHAR(64) NOT NULL COMMENT '策略适用数据类别',
  `normal_ttl_days` INT UNSIGNED NULL COMMENT '正常容量下保留天数',
  `warning_ttl_days` INT UNSIGNED NULL COMMENT '警告容量下保留天数',
  `warning_percent` DECIMAL(5,2) NOT NULL DEFAULT 70 COMMENT '策略警告水位百分比',
  `critical_percent` DECIMAL(5,2) NOT NULL DEFAULT 90 COMMENT '策略严重水位百分比',
  `action` VARCHAR(64) NOT NULL COMMENT '达到水位后的动作',
  `enabled` BOOLEAN NOT NULL DEFAULT TRUE COMMENT '策略是否启用',
  `description` VARCHAR(512) NOT NULL COMMENT '策略中文说明',
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '记录创建时间',
  `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '记录更新时间',
  PRIMARY KEY (`policy_id`),
  UNIQUE KEY `uk_storage_policies_key` (`policy_key`),
  KEY `idx_storage_policies_class` (`data_class`, `enabled`)
) ENGINE=InnoDB COMMENT='TTL和容量水位治理策略';

-- fnctl:statement
CREATE TABLE IF NOT EXISTS `node_tombstones` (
  `node_tombstone_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '节点墓碑内部主键',
  `config_fingerprint` BINARY(32) NOT NULL COMMENT '被清理完整配置SHA-256',
  `endpoint_fingerprint` BINARY(32) NOT NULL COMMENT '被清理网络端点SHA-256',
  `protocol` VARCHAR(32) NOT NULL COMMENT '被清理节点协议',
  `first_seen_at` DATETIME(6) NOT NULL COMMENT '历史首次出现时间',
  `last_seen_at` DATETIME(6) NOT NULL COMMENT '历史最近出现时间',
  `last_success_at` DATETIME(6) NULL COMMENT '历史最近成功时间',
  `ever_succeeded` BOOLEAN NOT NULL DEFAULT FALSE COMMENT '历史上是否验证成功',
  `best_quality_score` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '历史最佳质量分',
  `source_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '历史独立来源数量',
  `last_error_code` VARCHAR(64) NULL COMMENT '清理前最近错误码',
  `purge_reason` VARCHAR(64) NOT NULL COMMENT '完整配置清理原因',
  `purged_at` DATETIME(6) NOT NULL COMMENT '完整配置清理时间',
  `restore_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '墓碑恢复次数',
  `last_restored_at` DATETIME(6) NULL COMMENT '最近从墓碑恢复时间',
  `summary` JSON NULL COMMENT '长期保留的轻量历史摘要',
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '记录创建时间',
  `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '记录更新时间',
  PRIMARY KEY (`node_tombstone_id`),
  UNIQUE KEY `uk_node_tombstones_fingerprint` (`config_fingerprint`),
  KEY `idx_node_tombstones_endpoint` (`endpoint_fingerprint`, `protocol`),
  KEY `idx_node_tombstones_purged` (`purged_at`)
) ENGINE=InnoDB COMMENT='已清理节点配置的长期轻量墓碑';

-- fnctl:statement
INSERT INTO `storage_policies`
  (`policy_key`, `data_class`, `normal_ttl_days`, `warning_ttl_days`, `warning_percent`, `critical_percent`, `action`, `description`)
VALUES
  ('raw_payloads', 'raw_payload', 30, 14, 70, 90, 'delete_expired', '压缩原始响应在容量升高时缩短保留期'),
  ('parse_errors', 'parse_error', 90, 30, 70, 90, 'delete_expired', '解析错误样本保留后按期清理'),
  ('validation_attempts', 'validation_attempt', 180, 90, 80, 90, 'rollup_then_delete', '逐次验证明细先汇总再删除'),
  ('source_fetches', 'source_fetch', 365, 180, 80, 90, 'rollup_then_delete', '来源抓取明细先汇总再删除'),
  ('export_members', 'export_member', 365, 180, 80, 90, 'rollup_then_delete', '导出成员明细保留一年'),
  ('cold_source_breaker', 'source', NULL, NULL, 90, 94, 'pause_cold_sources', '容量达到熔断线时暂停低价值冷源')
ON DUPLICATE KEY UPDATE
  `data_class` = VALUES(`data_class`),
  `normal_ttl_days` = VALUES(`normal_ttl_days`),
  `warning_ttl_days` = VALUES(`warning_ttl_days`),
  `warning_percent` = VALUES(`warning_percent`),
  `critical_percent` = VALUES(`critical_percent`),
  `action` = VALUES(`action`),
  `description` = VALUES(`description`);
