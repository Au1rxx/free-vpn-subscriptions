-- fnctl:statement
CREATE DATABASE IF NOT EXISTS `vpn_nodes`
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;

-- fnctl:statement
CREATE TABLE IF NOT EXISTS `schema_migrations` (
  `version` VARCHAR(32) NOT NULL COMMENT '迁移版本号',
  `name` VARCHAR(255) NOT NULL COMMENT '迁移文件名称',
  `checksum` BINARY(32) NOT NULL COMMENT '迁移文件SHA-256校验值',
  `applied_at` DATETIME(6) NOT NULL COMMENT '迁移执行时间',
  PRIMARY KEY (`version`)
) ENGINE=InnoDB COMMENT='数据库迁移版本记录';
