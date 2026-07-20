-- fnctl:statement
-- 覆盖终审所需的当前状态聚合列，避免读取包含 JSON 和长文本的整张宽表。
ALTER TABLE `node_current_status`
  ADD INDEX `idx_node_current_validation_status` (`availability_state`, `last_validation_at`, `quality_grade`, `quality_score`)
    COMMENT '覆盖可用性时效等级和质量分聚合查询',
  ALGORITHM=INPLACE,
  LOCK=NONE;
