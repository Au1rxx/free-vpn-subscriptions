-- fnctl:statement
CREATE TABLE IF NOT EXISTS `node_classifications` (
  `node_config_id` BIGINT UNSIGNED NOT NULL COMMENT '节点配置主键',
  `protocol` VARCHAR(32) NOT NULL COMMENT '协议分类',
  `transport` VARCHAR(32) NOT NULL COMMENT '传输分类',
  `security` VARCHAR(32) NOT NULL COMMENT '安全层分类',
  `ip_version` VARCHAR(16) NULL COMMENT '入口IP版本分类',
  `entry_country` CHAR(2) NULL COMMENT '入口国家代码',
  `entry_region` VARCHAR(128) NULL COMMENT '入口地区名称',
  `entry_city` VARCHAR(128) NULL COMMENT '入口城市名称',
  `entry_timezone` VARCHAR(64) NULL COMMENT '入口时区',
  `entry_asn` VARCHAR(32) NULL COMMENT '入口自治系统编号',
  `entry_organization` VARCHAR(255) NULL COMMENT '入口网络组织名称',
  `provider_class` VARCHAR(32) NULL COMMENT '云厂商或网络类型分类',
  `exit_country` CHAR(2) NULL COMMENT '出口国家代码',
  `exit_asn` VARCHAR(32) NULL COMMENT '出口自治系统编号',
  `freshness_class` VARCHAR(32) NOT NULL DEFAULT 'new' COMMENT '节点新鲜度分类',
  `stability_class` VARCHAR(32) NOT NULL DEFAULT 'unknown' COMMENT '节点稳定性分类',
  `risk_tags` JSON NULL COMMENT '风险标签集合',
  `classifier_version` VARCHAR(64) NOT NULL COMMENT '分类器版本',
  `classified_at` DATETIME(6) NOT NULL COMMENT '最近分类时间',
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '记录创建时间',
  `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '记录更新时间',
  PRIMARY KEY (`node_config_id`),
  KEY `idx_node_classifications_country` (`entry_country`, `protocol`),
  KEY `idx_node_classifications_asn` (`entry_asn`, `provider_class`),
  KEY `idx_node_classifications_freshness` (`freshness_class`, `stability_class`),
  CONSTRAINT `fk_node_classifications_node` FOREIGN KEY (`node_config_id`) REFERENCES `node_configs` (`node_config_id`) ON DELETE CASCADE
) ENGINE=InnoDB COMMENT='节点地理网络协议质量和风险分类';

-- fnctl:statement
CREATE TABLE IF NOT EXISTS `export_runs` (
  `export_run_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '导出批次内部主键',
  `run_uuid` CHAR(36) NOT NULL COMMENT '导出批次全局标识',
  `rules_version` VARCHAR(64) NOT NULL COMMENT '导出规则版本',
  `started_at` DATETIME(6) NOT NULL COMMENT '导出开始时间',
  `finished_at` DATETIME(6) NULL COMMENT '导出结束时间',
  `candidate_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '导出候选节点数量',
  `selected_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '实际导出节点数量',
  `file_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '生成文件数量',
  `output_bytes` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '生成文件总字节数',
  `export_state` VARCHAR(32) NOT NULL COMMENT '导出批次状态',
  `summary` JSON NULL COMMENT '导出分类统计摘要',
  `error_summary` VARCHAR(1024) NULL COMMENT '导出受限错误摘要',
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '记录创建时间',
  PRIMARY KEY (`export_run_id`),
  UNIQUE KEY `uk_export_runs_uuid` (`run_uuid`),
  KEY `idx_export_runs_time_state` (`started_at`, `export_state`)
) ENGINE=InnoDB COMMENT='公共订阅生成批次';

-- fnctl:statement
CREATE TABLE IF NOT EXISTS `export_members` (
  `export_member_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '导出成员内部主键',
  `export_run_id` BIGINT UNSIGNED NOT NULL COMMENT '导出批次主键',
  `node_config_id` BIGINT UNSIGNED NOT NULL COMMENT '被导出节点配置主键',
  `collection_name` VARCHAR(128) NOT NULL COMMENT '所属订阅集合名称',
  `rank_number` INT UNSIGNED NOT NULL COMMENT '集合内排序序号',
  `quality_score` TINYINT UNSIGNED NOT NULL COMMENT '导出时质量分',
  `quality_grade` CHAR(1) NOT NULL COMMENT '导出时质量等级',
  `selection_reason` VARCHAR(255) NOT NULL COMMENT '节点入选原因',
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '记录创建时间',
  PRIMARY KEY (`export_member_id`),
  UNIQUE KEY `uk_export_members_pair` (`export_run_id`, `node_config_id`, `collection_name`),
  KEY `idx_export_members_collection_rank` (`export_run_id`, `collection_name`, `rank_number`),
  KEY `idx_export_members_node_time` (`node_config_id`, `created_at`),
  CONSTRAINT `fk_export_members_run` FOREIGN KEY (`export_run_id`) REFERENCES `export_runs` (`export_run_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_export_members_node` FOREIGN KEY (`node_config_id`) REFERENCES `node_configs` (`node_config_id`) ON DELETE CASCADE
) ENGINE=InnoDB COMMENT='每次公共订阅包含的节点';
