# ikuai-cli

[![CI](https://github.com/ikuaidev/ikuai-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/ikuaidev/ikuai-cli/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/ikuaidev/ikuai-cli)](https://github.com/ikuaidev/ikuai-cli/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/ikuaidev/ikuai-cli)](go.mod)
[![Go Report Card](https://goreportcard.com/badge/github.com/ikuaidev/ikuai-cli)](https://goreportcard.com/report/github.com/ikuaidev/ikuai-cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

iKuai 路由器命令行工具 — 在终端管理网络、用户、VPN、防火墙等

[English](README.md)

## 安装

**macOS / Linux：**

```bash
curl -fsSL https://raw.githubusercontent.com/ikuaidev/ikuai-cli/main/scripts/install.sh | sh
```

**Go install：**

```bash
go install github.com/ikuaidev/ikuai-cli/cmd/ikuai-cli@latest
```

<details>
<summary>其他安装方式</summary>

**预编译二进制：** 从 [Releases 页面](https://github.com/ikuaidev/ikuai-cli/releases) 下载。

```bash
VERSION=0.1.0 ARCH=linux_amd64
curl -fsSL "https://github.com/ikuaidev/ikuai-cli/releases/download/v${VERSION}/ikuai-cli_${ARCH}.tar.gz" -o ikuai-cli.tar.gz
curl -fsSL "https://github.com/ikuaidev/ikuai-cli/releases/download/v${VERSION}/checksums.txt" -o checksums.txt
sha256sum --check --ignore-missing checksums.txt
tar -xzf ikuai-cli.tar.gz
sudo mv ikuai-cli /usr/local/bin/
```

Windows：从 [Releases](https://github.com/ikuaidev/ikuai-cli/releases) 下载 `ikuai-cli_windows_amd64.zip`。

**从源码编译：**

```bash
git clone https://github.com/ikuaidev/ikuai-cli.git
cd ikuai-cli
make build
```

**Shell 补全：**

```bash
ikuai-cli completion bash > ~/.local/share/bash-completion/completions/ikuai-cli
ikuai-cli completion zsh  > ~/.zsh/completions/_ikuai-cli
ikuai-cli completion fish > ~/.config/fish/completions/ikuai-cli.fish
ikuai-cli completion powershell > ikuai-cli.ps1
```

</details>

## 从路由器中申请 Token 的方法

![第一步](docs/images/token-step1.png)

## 快速开始

### 1. 认证

**方式 A：环境变量**（推荐用于脚本和 Agent）

```bash
export IKUAI_CLI_BASE_URL=https://192.168.1.1
export IKUAI_CLI_TOKEN=<你的Token>
ikuai-cli monitor system   # 直接可用
```

**方式 B：持久化 Session**（保存到 `~/.ikuai-cli/config.json`）

```bash
ikuai-cli auth set-url https://192.168.1.1
ikuai-cli auth set-token <你的Token>
```

### 2. 验证

```bash
ikuai-cli auth status
```

### 3. 开始使用

```bash
ikuai-cli monitor system                     # CPU、内存、运行时间、WAN IP
ikuai-cli network dns get                    # DNS 配置
ikuai-cli network dns proxy create --domain example.com --dns-addr 8.8.8.8 --parse-type ipv4
ikuai-cli network pppoe set --comment maintenance --mtu 1480 --mru 1480
ikuai-cli users online                       # 在线用户
ikuai-cli security acl list                  # 安全策略
ikuai-cli log system list --human-time       # 系统日志
```

> 完整命令参考：[docs/cli-reference.md](docs/cli-reference.md)

## 功能特性

- **网络管理** — DNS、DHCP、VLAN、NAT、PPPoE、网口配置
- **系统监控** — CPU、内存、运行时间、流量、在线终端
- **安全策略** — ACL、MAC 过滤、L7 规则、URL 过滤、域名黑名单、连接数限制、终端标注
- **用户管理** — 认证账号、在线会话、踢下线、认证套餐、带宽限制
- **路由配置** — 静态路由、策略路由、多 WAN 负载均衡
- **VPN** — PPTP、L2TP、OpenVPN、IKEv2、IPSec、WireGuard
- **对象组** — IP、IPv6、MAC、端口、协议、域名、时间对象组
- **无线管理** — Wi-Fi 配置与管理
- **QoS** — 带宽控制和流量整形
- **系统维护** — 配置管理、定时任务、远程访问、VRRP、SSH 重置
- **日志** — 系统日志和审计记录
- **交互模式** — `repl` 模式，支持多级 Tab 补全
- **结构化输出** — 默认表格；支持 `--format json/yaml`、`--raw`；`--human-time` 显示可读时间；`--wide` / `--columns` 控制列显示

## 输出格式

默认表格输出。当 stdout 不是 TTY 时（管道或重定向），自动切换为 JSON。

列表命令在表格模式下只显示核心默认列。使用 `--wide` 查看全部字段，或 `--columns` 选择特定列。列宽会自适应终端宽度。

```bash
ikuai-cli monitor system                     # 表格（人类可读）
ikuai-cli monitor system --format json       # JSON（脚本、jq、Agent）
ikuai-cli monitor system --format yaml       # YAML
ikuai-cli monitor system --raw               # 完整 API 信封（调试）
ikuai-cli log system list --human-time       # 可读时间戳
ikuai-cli security acl list --wide           # 显示全部列
ikuai-cli security acl list --columns id,src_addr,action  # 自定义列
```

> `--raw` 和 `--format` 互斥。`--wide` 和 `--columns` 互斥。

## Session 存储

凭证保存在 `~/.ikuai-cli/config.json`。可通过环境变量覆盖：

```bash
export IKUAI_CLI_CONFIG_FILE=/path/to/config.json
```

**优先级：** Session 文件 > 环境变量 > 无。

环境变量（`IKUAI_CLI_BASE_URL` / `IKUAI_CLI_TOKEN`）不会写入磁盘。

## AI Agent 集成

ikuai-cli 内置 [`SKILL.md`](./SKILL.md) 和领域 [skills](./skills/)，让 AI Agent 可以直接操作 iKuai 路由器。

| 技能 | 说明 |
|------|------|
| [monitor](skills/monitor.md) | 系统状态、CPU、内存、流量、在线客户端 |
| [network](skills/network.md) | DNS、DHCP、VLAN、NAT/DNAT、WAN、LAN、PPPoE、DMZ、DNS 代理 |
| [users](skills/users.md) | 在线用户、按 ID 踢下线、账户、套餐 |
| [security](skills/security.md) | ACL、MAC 过滤、L7、URL 过滤、域名黑名单、连接数限制、终端标注 |
| [vpn](skills/vpn.md) | PPTP、L2TP、OpenVPN、IKEv2、IPSec、WireGuard |
| [objects](skills/objects.md) | IP、IPv6、MAC、端口、协议、域名、时间对象组 |
| [system](skills/system.md) | 系统配置、定时任务、远程访问、VRRP、SSH 重置 |
| [batch](skills/batch.md) | 组合工作流：初始化、批量 DHCP、配置备份 |

### 安装 Skills

**[Skills CLI](https://github.com/vercel-labs/skills)（推荐）：**

```bash
npx skills add ikuai/ikuai-cli
```

| 参数 | 说明 |
|------|------|
| `-g` | 全局安装（用户级，跨项目共享） |
| `-a claude-code` | 指定目标 Agent |
| `-y` | 非交互模式 |

**手动安装：**

```bash
mkdir -p .agents/skills
git clone https://github.com/ikuaidev/ikuai-cli.git .agents/skills/ikuai-cli
```

### Agent 输出

使用 `--format json` 获取结构化输出：

```bash
ikuai-cli monitor system --format json
ikuai-cli users online --format json | jq '.data[] | {id, ip_addr, mac, username}'
```

使用 `--format yaml` 获取节省 token 的输出（当不需要完整 JSON 精度时）。

## 开发

```bash
git clone https://github.com/ikuaidev/ikuai-cli.git
cd ikuai-cli
make test       # 运行测试（含竞态检测）
make lint       # golangci-lint 代码检查
make build      # 编译
make smoke      # 冒烟测试
```

交叉编译：

```bash
make linux-amd64    make linux-arm64
make darwin-amd64   make darwin-arm64
```

## 社区

- [Issues](https://github.com/ikuaidev/ikuai-cli/issues) — Bug 报告与功能建议
- [Contributing](CONTRIBUTING.md) — 贡献指南
- [Security](SECURITY.md) — 安全漏洞报告
- [Code of Conduct](CODE_OF_CONDUCT.md) — 行为准则

## License

MIT — see [LICENSE](LICENSE)
