-- fnctl:statement
-- 覆盖 Node API 按协议和最近发现时间执行的游标分页，避免每页扫描并排序整张配置表。
ALTER TABLE `node_configs`
  ADD INDEX `idx_node_configs_protocol_seen_id` (`protocol`, `last_seen_at`, `node_config_id`)
    COMMENT '覆盖协议和最近发现时间的节点接口游标分页',
  ALGORITHM=INPLACE,
  LOCK=NONE;
