# 节点数据库运行手册

项目使用 MySQL `vpn_nodes` 数据库保存来源、原始响应、规范化配置、验证记录、分类、导出批次和 TTL 治理数据。全部表和字段都带中文备注，迁移文件位于 `db/migrations/`，由 `fnctl migrate` 按版本和 SHA-256 校验执行。

## 连接边界

- 采集机只连接本机 `127.0.0.1:13306`，由受管 SSH 隧道转发到数据库私网地址。
- 密码只从权限受限的凭据文件读取，不写入 YAML、命令参数、日志或 Git。
- `tls_mode: required` 强制数据库会话启用 TLS；`db-status` 会检查 TLS 密码套件、只读状态和 UTC 时区。
- 运行时连接池默认最多 20 个连接、10 个空闲连接，避免采集并发耗尽数据库连接。

## 初始化与验收

配置文件只引用凭据路径：

```yaml
database:
  enabled: true
  address: 127.0.0.1:13306
  name: vpn_nodes
  user: adminai
  password_file: /run/credentials/free-vpn-harvester.service/mysql-password
  tls_mode: required
  max_open_conns: 20
  max_idle_conns: 10
```

在 SSH 隧道可用且凭据文件已挂载后执行：

```bash
fnctl migrate --config /etc/free-vpn-harvester/config.yaml
fnctl db-status --config /etc/free-vpn-harvester/config.yaml
make test-migrations CONFIG=/etc/free-vpn-harvester/config.yaml
```

迁移验收会重复执行迁移，确认第二次全部跳过，并检查：6 个迁移、22 张业务表、所有表和字段中文备注完整、6 条容量/TTL 策略启用、数据库可写、TLS 与 UTC 正常。迁移采用可重复执行的 DDL；MySQL DDL 会隐式提交，因此只有一个版本的全部语句成功后才记录版本。

## 50GB 容量和 TTL

`storage_metrics.capacity_bytes` 固定以 50 GiB（53,687,091,200 字节）为容量基线。`storage_policies` 预置分层保留策略：原始响应 30 天、解析错误 90 天、验证明细 180 天、来源抓取与导出明细 365 天；达到警告水位时缩短 TTL，达到严重水位时暂停低价值冷源。节点配置删除前写入轻量墓碑，避免丢失长期去重和历史成功信息。

数据库密码轮换时只替换 systemd credential 并重启相关服务；不要修改配置文件或迁移文件。迁移文件一旦进入数据库即不可改写，校验和不一致会立即终止。
