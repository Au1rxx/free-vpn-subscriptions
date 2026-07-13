-- fnctl:statement
CREATE TABLE IF NOT EXISTS `endpoints` (
  `endpoint_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '网络端点内部主键',
  `host` VARCHAR(255) NOT NULL COMMENT '规范化主机名或IP地址',
  `host_hash` BINARY(32) NOT NULL COMMENT '规范化主机SHA-256',
  `port` SMALLINT UNSIGNED NOT NULL COMMENT '服务端口',
  `address_type` VARCHAR(16) NOT NULL COMMENT '域名或IP地址类型',
  `resolved_ipv4` JSON NULL COMMENT '最近解析的IPv4地址集合',
  `resolved_ipv6` JSON NULL COMMENT '最近解析的IPv6地址集合',
  `dns_state` VARCHAR(32) NOT NULL DEFAULT 'pending' COMMENT '最近DNS解析状态',
  `dns_duration_ms` INT UNSIGNED NULL COMMENT '最近DNS解析耗时毫秒',
  `dns_expires_at` DATETIME(6) NULL COMMENT 'DNS结果过期时间',
  `first_seen_at` DATETIME(6) NOT NULL COMMENT '端点首次出现时间',
  `last_seen_at` DATETIME(6) NOT NULL COMMENT '端点最近出现时间',
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '记录创建时间',
  `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '记录更新时间',
  PRIMARY KEY (`endpoint_id`),
  UNIQUE KEY `uk_endpoints_host_port` (`host_hash`, `port`),
  KEY `idx_endpoints_dns_due` (`dns_state`, `dns_expires_at`),
  KEY `idx_endpoints_last_seen` (`last_seen_at`)
) ENGINE=InnoDB COMMENT='主机和端口构成的网络端点';

-- fnctl:statement
CREATE TABLE IF NOT EXISTS `node_configs` (
  `node_config_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '节点配置内部主键',
  `endpoint_id` BIGINT UNSIGNED NOT NULL COMMENT '关联网络端点主键',
  `config_fingerprint` BINARY(32) NOT NULL COMMENT '完整规范化配置SHA-256',
  `protocol` VARCHAR(32) NOT NULL COMMENT '代理协议',
  `protocol_version` VARCHAR(32) NULL COMMENT '协议版本',
  `transport` VARCHAR(32) NOT NULL DEFAULT 'tcp' COMMENT '传输方式',
  `security` VARCHAR(32) NOT NULL DEFAULT 'none' COMMENT '安全层类型',
  `normalized_config` JSON NOT NULL COMMENT '完整规范化配置内容',
  `canonical_uri` MEDIUMTEXT NULL COMMENT '可重建时的规范分享链接',
  `config_bytes` INT UNSIGNED NOT NULL COMMENT '规范化配置字节数',
  `parser_version` VARCHAR(64) NOT NULL COMMENT '创建该版本的解析器版本',
  `lifecycle_state` VARCHAR(32) NOT NULL DEFAULT 'pending' COMMENT '节点生命周期状态',
  `is_exportable` BOOLEAN NOT NULL DEFAULT FALSE COMMENT '当前是否允许导出',
  `first_seen_at` DATETIME(6) NOT NULL COMMENT '配置首次出现时间',
  `last_seen_at` DATETIME(6) NOT NULL COMMENT '配置最近出现时间',
  `last_success_at` DATETIME(6) NULL COMMENT '最近真实验证成功时间',
  `consecutive_failures` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '连续验证失败次数',
  `expires_at` DATETIME(6) NULL COMMENT '完整配置计划清理时间',
  `archived_at` DATETIME(6) NULL COMMENT '配置归档时间',
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '记录创建时间',
  `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '记录更新时间',
  PRIMARY KEY (`node_config_id`),
  UNIQUE KEY `uk_node_configs_fingerprint` (`config_fingerprint`),
  KEY `idx_node_configs_endpoint` (`endpoint_id`, `protocol`),
  KEY `idx_node_configs_lifecycle` (`lifecycle_state`, `last_seen_at`),
  KEY `idx_node_configs_export` (`is_exportable`, `protocol`, `last_success_at`),
  KEY `idx_node_configs_expires` (`expires_at`),
  CONSTRAINT `fk_node_configs_endpoint` FOREIGN KEY (`endpoint_id`) REFERENCES `endpoints` (`endpoint_id`) ON DELETE CASCADE
) ENGINE=InnoDB COMMENT='按完整连接参数去重的节点配置';

-- fnctl:statement
CREATE TABLE IF NOT EXISTS `node_source_stats` (
  `node_source_stat_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '节点来源关系内部主键',
  `node_config_id` BIGINT UNSIGNED NOT NULL COMMENT '节点配置主键',
  `source_id` BIGINT UNSIGNED NOT NULL COMMENT '来源主键',
  `last_fetch_id` BIGINT UNSIGNED NULL COMMENT '最近发现该节点的抓取主键',
  `first_seen_at` DATETIME(6) NOT NULL COMMENT '在该来源首次出现时间',
  `last_seen_at` DATETIME(6) NOT NULL COMMENT '在该来源最近出现时间',
  `seen_count` BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '在该来源累计出现次数',
  `is_active` BOOLEAN NOT NULL DEFAULT TRUE COMMENT '该来源当前是否仍发布节点',
  `source_quality` DECIMAL(5,2) NOT NULL DEFAULT 0 COMMENT '该来源对节点的质量贡献',
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '记录创建时间',
  `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '记录更新时间',
  PRIMARY KEY (`node_source_stat_id`),
  UNIQUE KEY `uk_node_source_stats_pair` (`node_config_id`, `source_id`),
  KEY `idx_node_source_stats_source_active` (`source_id`, `is_active`, `last_seen_at`),
  KEY `idx_node_source_stats_node_active` (`node_config_id`, `is_active`),
  CONSTRAINT `fk_node_source_stats_node` FOREIGN KEY (`node_config_id`) REFERENCES `node_configs` (`node_config_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_node_source_stats_source` FOREIGN KEY (`source_id`) REFERENCES `sources` (`source_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_node_source_stats_fetch` FOREIGN KEY (`last_fetch_id`) REFERENCES `source_fetches` (`fetch_id`) ON DELETE SET NULL
) ENGINE=InnoDB COMMENT='节点与来源的长期聚合关系';

-- fnctl:statement
CREATE TABLE IF NOT EXISTS `node_source_daily` (
  `node_source_daily_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '每日来源统计内部主键',
  `stat_date` DATE NOT NULL COMMENT '统计日期',
  `node_config_id` BIGINT UNSIGNED NOT NULL COMMENT '节点配置主键',
  `source_id` BIGINT UNSIGNED NOT NULL COMMENT '来源主键',
  `seen_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '当日出现次数',
  `first_seen_at` DATETIME(6) NULL COMMENT '当日首次出现时间',
  `last_seen_at` DATETIME(6) NULL COMMENT '当日最近出现时间',
  `was_new` BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否为当日新发现',
  `was_removed` BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否在当日判定消失',
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '记录创建时间',
  `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '记录更新时间',
  PRIMARY KEY (`node_source_daily_id`),
  UNIQUE KEY `uk_node_source_daily_triplet` (`stat_date`, `node_config_id`, `source_id`),
  KEY `idx_node_source_daily_source_date` (`source_id`, `stat_date`),
  KEY `idx_node_source_daily_node_date` (`node_config_id`, `stat_date`),
  CONSTRAINT `fk_node_source_daily_node` FOREIGN KEY (`node_config_id`) REFERENCES `node_configs` (`node_config_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_node_source_daily_source` FOREIGN KEY (`source_id`) REFERENCES `sources` (`source_id`) ON DELETE CASCADE
) ENGINE=InnoDB COMMENT='节点来源关系每日变化统计';
