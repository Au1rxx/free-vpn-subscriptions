-- 保存受控吞吐采样失败原因，使“未采样”和“采样失败”可以明确区分。

-- fnctl:statement
ALTER TABLE `validation_attempts`
  ADD COLUMN `performance_error_code` VARCHAR(64) NULL COMMENT '性能采样错误分类代码'
  AFTER `bytes_per_second`;
