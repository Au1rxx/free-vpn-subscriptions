-- 为只发布裸 IP:端口 的协议分片保存显式协议，避免猜测和错误分类。

-- fnctl:statement
ALTER TABLE `sources`
  ADD COLUMN `protocol_hint` VARCHAR(32) NULL COMMENT '裸端点列表协议提示'
  AFTER `format_hint`;
