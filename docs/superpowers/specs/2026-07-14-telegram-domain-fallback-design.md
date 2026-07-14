# Telegram 公开预览域名回退设计

## 背景

2026-07-14 生产采集显示约 40 个 `http_failed`，与数据库中 42 个启用的
`telegram-public` 来源基本一致。主机的 Azure DNS、Cloudflare DoH 和 Google DoH
均对 `t.me` 返回权威 NXDOMAIN；公开 DNS 状态同时显示该域名被注册局置为
`serverHold`。相同公开预览路径通过 `telegram.me`、`telegram.dog` 和
`telegram.org` 均返回 HTTP 200。

现有来源以规范化 URL 哈希作为身份。批量把 `t.me` 改写为其他域名会制造新来源、
割裂抓取历史，并且无法处理以后动态发现的新 `t.me` 链接。

## 目标与非目标

目标：

- 恢复现有和未来 `t.me` 公开来源的抓取与发现；
- 保留数据库中的原始 URL、规范化身份和历史关系；
- 只在可确认的 DNS 不存在错误上回退，避免掩盖 HTTP、TLS、超时或内容错误；
- 保留现有超时、重定向、正文大小和错误码边界。

非目标：

- 不绕过 Telegram 账号、频道权限或平台认证；
- 不为任意失效域名提供通用镜像重写；
- 不批量修改生产数据库中的来源 URL；
- 不增加第二次以上的域名重试。

## 方案比较

1. 批量迁移 `sources.url` 和版本化种子到 `telegram.me`：实现简单，但改变来源身份，
   会产生重复来源并丢失未来动态 `t.me` 链接，故不采用。
2. 只在 `internal/sources` 抓取器中回退：能恢复定时抓取，但每轮主动 Telegram
   发现仍失败，覆盖不完整，故不采用。
3. 提供抓取与发现共用的受限 HTTP 回退器：保留来源身份，同时覆盖两个入口，且策略
   可以独立单元测试。采用此方案。

## 详细设计

新增内部包 `internal/httpfallback`，只暴露：

```go
func Do(client *http.Client, request *http.Request) (*http.Response, error)
```

执行流程：

1. 使用调用方提供的 `http.Client` 执行原请求；
2. 成功返回响应时不做任何处理，包括 4xx/5xx；
3. 请求失败时同时满足以下条件才允许回退：
   - 原请求主机名不区分大小写等于 `t.me`；
   - 错误链包含 `net.DNSError`，且 `IsNotFound=true`；
4. 仅对当前两个调用方使用的无正文 GET 请求执行回退；克隆请求，保留 context、
   路径、查询和请求头，仅把主机替换为 `telegram.me`，然后重试一次；
5. 回退成功时返回备用域名响应。`source_fetches.final_url` 可以记录实际访问的
   `telegram.me`，但 `sources.url` 和来源哈希保持原值；
6. 回退也失败时返回包含两次失败上下文的错误，现有上层仍将其归类为
   `http_failed`。

`internal/sources.FetchRaw` 和 `internal/discovery.fetch` 均改用该回退器。其余 HTTP
调用不受影响。

## 错误与安全边界

- 非 `t.me`、DNS 超时、连接拒绝、TLS 错误、HTTP 状态错误都只执行一次；
- 非 GET 请求或带正文请求不回退，避免重复发送不可重放的请求体；
- 备用目标是固定常量，不能由响应或采集内容控制；
- 使用 `URL.Hostname()` 判断，带用户信息、伪造后缀或非标准端口不能冒充 `t.me`；
- 备用请求继续经过调用方现有的重定向检查、超时和正文读取限制；
- 日志不增加响应正文、凭据或完整敏感配置。

## 测试与生产验收

先写失败测试，使用自定义 `RoundTripper` 验证：

- `t.me` 的 `net.DNSError{IsNotFound:true}` 会精确回退一次，且路径、查询和请求头不变；
- 非 `t.me` DNS 错误不回退；
- `t.me` 非 DNS 错误不回退；
- HTTP 404 返回原响应，不回退；
- 回退失败时错误链仍可被上层稳定归为 `http_failed`。

随后运行新包测试、`internal/sources`、`internal/discovery`、全仓测试、`go vet` 和相关
竞态测试。生产安装后手动执行一次 Telegram 发现和下一轮采集，验收标准为：

- `discover --kind telegram --url https://t.me/s/V2RayRootFree` 成功；
- 新一轮 Telegram 来源不再集中产生约 40 个 DNS `http_failed`；
- `source_fetches.final_url` 可出现 `telegram.me`，来源总数不因域名替换批量翻倍；
- 采集、验证和影子服务无新增重启或过期租约。
