-- fnctl:statement
ALTER TABLE `raw_payloads`
  ADD COLUMN `storage_kind` VARCHAR(16) NOT NULL DEFAULT 'database' COMMENT '原始正文存储后端',
  ADD COLUMN `archive_key` VARCHAR(255) NULL COMMENT '文件归档根目录下的相对键',
  ALGORITHM=INSTANT;
