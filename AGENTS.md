# 公共产物仓规则

## 仓库职责

- 本仓库 `Au1rxx/free-vpn-subscriptions` 是 PUBLIC、自动生成的订阅产物仓，不是采集程序源码仓、构建仓、数据库仓或生产配置仓。
- 当前分支只用于向终端用户提供订阅文件、生成 README、生成站点、必要静态资产和问题反馈入口。
- harvester 源码、测试、migration、采集源、构建、systemd 和部署文档只存在于私有运维仓库，不得复制到本仓库。
- 历史提交和 `v0.1.0`–`v0.5.2` 标签为既有依赖保留，不代表当前分支继续提供 Go module；不得强推或重写历史。
- 远程分支只允许 `main`；禁止创建源码开发分支。上述 8 个历史标签已经冻结，禁止移动、删除或新增公共 `v*` 标签。

## 允许内容

当前分支只允许以下 tracked 内容：

- `AGENTS.md`、`.gitignore`、`LICENSE`、`CODE_OF_CONDUCT.md`、`SECURITY.md`、`CONTRIBUTING.md`；
- `README.md`、`README_*.md`；
- `output/**`；
- 自动生成的 `docs/**`；
- `assets/hero.png`、`assets/star-history.svg`；
- `.github/ISSUE_TEMPLATE/**` 中面向订阅使用者的反馈模板。

## 禁止内容

- 禁止 `cmd/`、`internal/`、`pkg/`、`db/`、`scripts/`；
- 禁止 `go.mod`、`go.sum`、`Makefile`、`config.yaml`、`config.example.yaml`；
- 禁止 `.github/workflows/**` 和面向源码 PR 的模板；
- 禁止 Go、Shell、SQL、Terraform、Python、JavaScript、TypeScript、Rust、C/C++、Java、Kotlin、Ruby、PHP 等程序或运维实现文件，即使把它们放进 `docs/` 或 `output/`；
- 禁止采集源清单、数据库连接配置、归档配置、凭据、内部日志和服务器路径细节；
- 禁止在本仓库构建、测试或运行 harvester。
- 禁止通过 branch、tag、release asset 或生成目录重新引入 harvester 源码；既有消费者必须固定到不高于 `v0.5.2` 的历史版本，不得依赖 `main`。

## 修改与发布

- `output/`、README、站点和星标图由私有发布链路生成；不得手工编辑生成产物来掩盖生产问题。
- 不接受采集源码 PR。订阅失效、格式或新鲜度问题通过 Issue 模板反馈。
- 每次提交和推送前必须通过私有发布链路安装的 allow-list 检查；任何越界文件都应停止发布，而不是放宽规则。
- 公共提交不得修改账号、deploy key、仓库可见性、Pages 设置或其他外部状态，除非有针对该操作的单独明确授权。
