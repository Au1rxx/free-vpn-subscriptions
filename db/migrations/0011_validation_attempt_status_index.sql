-- fnctl:statement
-- 覆盖终审验证结果与吞吐聚合，避免并发采集期间回表扫描完整验证明细。
ALTER TABLE `validation_attempts`
  ADD INDEX `idx_validation_attempts_status` (`passed`, `partial_success`, `performance_bytes`, `performance_error_code`, `bytes_per_second`)
    COMMENT '覆盖验证结果和吞吐状态聚合查询',
  ALGORITHM=INPLACE,
  LOCK=NONE;
