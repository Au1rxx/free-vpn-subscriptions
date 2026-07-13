-- fnctl:statement
CREATE TABLE IF NOT EXISTS `sources` (
  `source_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '来源内部主键',
  `name` VARCHAR(255) NOT NULL COMMENT '来源显示名称',
  `kind` VARCHAR(32) NOT NULL COMMENT '来源类型',
  `url` TEXT NOT NULL COMMENT '来源原始地址',
  `canonical_url` TEXT NOT NULL COMMENT '规范化来源地址',
  `canonical_url_hash` BINARY(32) NOT NULL COMMENT '规范化地址SHA-256',
  `format_hint` VARCHAR(32) NOT NULL DEFAULT 'auto' COMMENT '内容格式提示',
  `discovery_method` VARCHAR(32) NOT NULL COMMENT '来源发现方式',
  `state` VARCHAR(32) NOT NULL DEFAULT 'active' COMMENT '来源生命周期状态',
  `enabled` BOOLEAN NOT NULL DEFAULT TRUE COMMENT '是否启用抓取',
  `priority` SMALLINT NOT NULL DEFAULT 0 COMMENT '抓取优先级',
  `depth` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '嵌套发现深度',
  `fetch_interval_seconds` INT UNSIGNED NOT NULL DEFAULT 3600 COMMENT '计划抓取间隔秒数',
  `next_fetch_at` DATETIME(6) NULL COMMENT '下次计划抓取时间',
  `etag` VARCHAR(512) NULL COMMENT '最近响应ETag',
  `last_modified` VARCHAR(255) NULL COMMENT '最近响应修改时间',
  `consecutive_failures` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '连续失败次数',
  `last_http_status` SMALLINT UNSIGNED NULL COMMENT '最近HTTP状态码',
  `last_success_at` DATETIME(6) NULL COMMENT '最近成功抓取时间',
  `last_failure_at` DATETIME(6) NULL COMMENT '最近失败抓取时间',
  `quality_score` DECIMAL(5,2) NOT NULL DEFAULT 0 COMMENT '来源质量评分',
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '记录创建时间',
  `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '记录更新时间',
  PRIMARY KEY (`source_id`),
  UNIQUE KEY `uk_sources_canonical_url` (`canonical_url_hash`),
  KEY `idx_sources_due` (`enabled`, `state`, `next_fetch_at`, `priority`),
  KEY `idx_sources_kind_state` (`kind`, `state`)
) ENGINE=InnoDB COMMENT='公开数据来源注册表';

-- fnctl:statement
CREATE TABLE IF NOT EXISTS `raw_payloads` (
  `payload_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '原始响应内部主键',
  `content_sha256` BINARY(32) NOT NULL COMMENT '未压缩正文SHA-256',
  `content_type` VARCHAR(255) NULL COMMENT '响应内容类型',
  `content_encoding` VARCHAR(64) NULL COMMENT '上游内容编码',
  `compression` VARCHAR(16) NOT NULL DEFAULT 'gzip' COMMENT '数据库压缩算法',
  `original_bytes` BIGINT UNSIGNED NOT NULL COMMENT '原始正文大小字节数',
  `compressed_bytes` BIGINT UNSIGNED NOT NULL COMMENT '压缩正文大小字节数',
  `compressed_body` MEDIUMBLOB NULL COMMENT '压缩后的原始正文',
  `first_seen_at` DATETIME(6) NOT NULL COMMENT '正文首次出现时间',
  `last_seen_at` DATETIME(6) NOT NULL COMMENT '正文最近出现时间',
  `expires_at` DATETIME(6) NULL COMMENT '正文计划清理时间',
  `reference_count` BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '抓取引用次数',
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '记录创建时间',
  PRIMARY KEY (`payload_id`),
  UNIQUE KEY `uk_raw_payloads_sha256` (`content_sha256`),
  KEY `idx_raw_payloads_expires` (`expires_at`)
) ENGINE=InnoDB COMMENT='按内容去重的压缩原始响应';

-- fnctl:statement
CREATE TABLE IF NOT EXISTS `source_fetches` (
  `fetch_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '抓取记录内部主键',
  `source_id` BIGINT UNSIGNED NOT NULL COMMENT '所属来源主键',
  `payload_id` BIGINT UNSIGNED NULL COMMENT '关联原始响应主键',
  `started_at` DATETIME(6) NOT NULL COMMENT '抓取开始时间',
  `finished_at` DATETIME(6) NULL COMMENT '抓取结束时间',
  `http_status` SMALLINT UNSIGNED NULL COMMENT 'HTTP响应状态码',
  `final_url` TEXT NULL COMMENT '重定向后的最终地址',
  `etag` VARCHAR(512) NULL COMMENT '本次响应ETag',
  `last_modified` VARCHAR(255) NULL COMMENT '本次响应修改时间',
  `content_type` VARCHAR(255) NULL COMMENT '本次响应内容类型',
  `content_encoding` VARCHAR(64) NULL COMMENT '本次响应内容编码',
  `response_bytes` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '读取响应字节数',
  `duration_ms` INT UNSIGNED NULL COMMENT '抓取总耗时毫秒',
  `response_headers` JSON NULL COMMENT '允许保存的响应头摘要',
  `fetch_state` VARCHAR(32) NOT NULL COMMENT '抓取终态',
  `parse_state` VARCHAR(32) NOT NULL DEFAULT 'pending' COMMENT '解析处理状态',
  `error_code` VARCHAR(64) NULL COMMENT '稳定抓取错误码',
  `error_summary` VARCHAR(1024) NULL COMMENT '受限长度错误摘要',
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '记录创建时间',
  PRIMARY KEY (`fetch_id`),
  KEY `idx_source_fetches_source_time` (`source_id`, `started_at`),
  KEY `idx_source_fetches_parse` (`parse_state`, `finished_at`),
  KEY `idx_source_fetches_payload` (`payload_id`),
  CONSTRAINT `fk_source_fetches_source` FOREIGN KEY (`source_id`) REFERENCES `sources` (`source_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_source_fetches_payload` FOREIGN KEY (`payload_id`) REFERENCES `raw_payloads` (`payload_id`) ON DELETE SET NULL
) ENGINE=InnoDB COMMENT='每次来源抓取结果';

-- fnctl:statement
CREATE TABLE IF NOT EXISTS `source_links` (
  `source_link_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '来源关系内部主键',
  `parent_source_id` BIGINT UNSIGNED NOT NULL COMMENT '父来源主键',
  `child_source_id` BIGINT UNSIGNED NOT NULL COMMENT '子来源主键',
  `discovery_location` VARCHAR(1024) NULL COMMENT '发现位置摘要',
  `evidence` VARCHAR(1024) NULL COMMENT '发现证据摘要',
  `depth` TINYINT UNSIGNED NOT NULL COMMENT '子来源发现深度',
  `first_seen_at` DATETIME(6) NOT NULL COMMENT '关系首次出现时间',
  `last_seen_at` DATETIME(6) NOT NULL COMMENT '关系最近出现时间',
  `seen_count` BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '关系累计出现次数',
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '记录创建时间',
  `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '记录更新时间',
  PRIMARY KEY (`source_link_id`),
  UNIQUE KEY `uk_source_links_pair` (`parent_source_id`, `child_source_id`),
  KEY `idx_source_links_child` (`child_source_id`, `last_seen_at`),
  CONSTRAINT `fk_source_links_parent` FOREIGN KEY (`parent_source_id`) REFERENCES `sources` (`source_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_source_links_child` FOREIGN KEY (`child_source_id`) REFERENCES `sources` (`source_id`) ON DELETE CASCADE
) ENGINE=InnoDB COMMENT='来源发现关系图';

-- fnctl:statement
CREATE TABLE IF NOT EXISTS `parse_runs` (
  `parse_run_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '解析批次内部主键',
  `fetch_id` BIGINT UNSIGNED NOT NULL COMMENT '关联抓取记录主键',
  `parser_version` VARCHAR(64) NOT NULL COMMENT '解析器版本',
  `detected_format` VARCHAR(32) NOT NULL COMMENT '实际识别内容格式',
  `started_at` DATETIME(6) NOT NULL COMMENT '解析开始时间',
  `finished_at` DATETIME(6) NULL COMMENT '解析结束时间',
  `input_entries` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '输入条目数量',
  `success_entries` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '成功解析数量',
  `error_entries` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '解析失败数量',
  `discovered_urls` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '发现的新地址数量',
  `parse_state` VARCHAR(32) NOT NULL COMMENT '解析批次状态',
  `error_summary` VARCHAR(1024) NULL COMMENT '批次错误摘要',
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '记录创建时间',
  PRIMARY KEY (`parse_run_id`),
  UNIQUE KEY `uk_parse_runs_fetch_version` (`fetch_id`, `parser_version`),
  KEY `idx_parse_runs_state_time` (`parse_state`, `started_at`),
  CONSTRAINT `fk_parse_runs_fetch` FOREIGN KEY (`fetch_id`) REFERENCES `source_fetches` (`fetch_id`) ON DELETE CASCADE
) ENGINE=InnoDB COMMENT='原始响应解析批次';

-- fnctl:statement
CREATE TABLE IF NOT EXISTS `parse_errors` (
  `parse_error_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '解析错误内部主键',
  `parse_run_id` BIGINT UNSIGNED NOT NULL COMMENT '所属解析批次主键',
  `source_id` BIGINT UNSIGNED NOT NULL COMMENT '所属来源主键',
  `fetch_id` BIGINT UNSIGNED NOT NULL COMMENT '所属抓取记录主键',
  `line_number` INT UNSIGNED NULL COMMENT '错误所在行号',
  `scheme_hint` VARCHAR(32) NULL COMMENT '推测的协议类型',
  `error_code` VARCHAR(64) NOT NULL COMMENT '稳定解析错误码',
  `sample_sha256` BINARY(32) NOT NULL COMMENT '错误样本SHA-256',
  `sample_excerpt` VARCHAR(512) NULL COMMENT '脱敏且截断的错误样本',
  `error_message` VARCHAR(1024) NOT NULL COMMENT '受限长度错误说明',
  `first_seen_at` DATETIME(6) NOT NULL COMMENT '错误首次出现时间',
  `last_seen_at` DATETIME(6) NOT NULL COMMENT '错误最近出现时间',
  `expires_at` DATETIME(6) NULL COMMENT '错误样本计划清理时间',
  `seen_count` BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '错误累计出现次数',
  `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '记录创建时间',
  PRIMARY KEY (`parse_error_id`),
  UNIQUE KEY `uk_parse_errors_run_sample` (`parse_run_id`, `sample_sha256`, `error_code`),
  KEY `idx_parse_errors_expires` (`expires_at`),
  KEY `idx_parse_errors_source_code` (`source_id`, `error_code`),
  CONSTRAINT `fk_parse_errors_run` FOREIGN KEY (`parse_run_id`) REFERENCES `parse_runs` (`parse_run_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_parse_errors_source` FOREIGN KEY (`source_id`) REFERENCES `sources` (`source_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_parse_errors_fetch` FOREIGN KEY (`fetch_id`) REFERENCES `source_fetches` (`fetch_id`) ON DELETE CASCADE
) ENGINE=InnoDB COMMENT='无法解析条目的结构化错误';
