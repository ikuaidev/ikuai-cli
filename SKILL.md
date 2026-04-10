---
name: ikuai-cli
description: |
  CLI for managing iKuai routers via local v4.0 API.
  Use when the user asks to check router status, configure network/DNS/DHCP/VPN/firewall,
  manage users, view logs, or automate any iKuai router operation.
---

# ikuai-cli — AI Agent Skill

Go CLI for managing iKuai routers through the local v4.0 HTTP API.

## When to Use

- Router status, CPU, memory, uptime, traffic
- Network: DNS, DHCP, VLAN, NAT, PPPoE, VPN
- Users: online sessions, kick, auth accounts
- Security: ACL, MAC filter, URL filter, domain blacklist
- System: logs, schedules, remote access, VRRP

## Prerequisites

```bash
# Session (persistent)
ikuai-cli auth set-url https://192.168.1.1
ikuai-cli auth set-token <TOKEN>

# OR environment variables (ephemeral)
export IKUAI_CLI_BASE_URL=https://192.168.1.1
export IKUAI_CLI_TOKEN=<TOKEN>

# Verify
ikuai-cli auth status --format json
```

> **Priority:** Session file > Environment variables > None.

## Output Modes

| Flag | Mode | Use case |
|------|------|----------|
| (default) | table (TTY) / JSON (pipe) | human / script auto-detect |
| `--format json` | compact JSON | agent parsing, pipe to `jq` |
| `--format yaml` | YAML | token-efficient for LLM |
| `--raw` | full API envelope | debug |
| `--human-time` | timestamp → local time | log reading |
| `--dry-run` | preview, no execute | safety check |

## Domain Skills

| Skill | File | Scope |
|-------|------|-------|
| Monitor | [monitor.md](skills/monitor.md) | CPU, memory, traffic, clients |
| Network | [network.md](skills/network.md) | DNS, DHCP, VLAN, NAT, WAN/LAN, PPPoE |
| Security | [security.md](skills/security.md) | ACL, MAC, L7, URL, domain blacklist, peerconn |
| VPN | [vpn.md](skills/vpn.md) | PPTP, L2TP, OpenVPN, IKEv2, IPSec, WireGuard |
| Users | [users.md](skills/users.md) | Online users, accounts, packages |
| System | [system.md](skills/system.md) | Config, schedules, remote access, VRRP, ALG, kernel |
| Auth | [auth.md](skills/auth.md) | Login, PPPoE auth accounts |
| Auth Server | [auth-server.md](skills/auth-server.md) | Web portal auth config |
| Log | [log.md](skills/log.md) | System/DHCP/PPPoE/auth/ARP logs |
| Objects | [objects.md](skills/objects.md) | IP/MAC/port/domain/time objects |
| QoS | [qos.md](skills/qos.md) | IP/MAC bandwidth control |
| Routing | [routing.md](skills/routing.md) | Static routes, stream shunting |
| Wireless | [wireless.md](skills/wireless.md) | Blacklist, VLAN, AC management |
| Advanced | [advanced.md](skills/advanced.md) | FTP, HTTP, Samba, SNMPD |
| Batch | [batch.md](skills/batch.md) | Multi-command workflows |

## Command Shape

```
ikuai-cli <resource> <action> [args] [flags]
```

- Read: `ikuai-cli monitor system --format json`
- List: `ikuai-cli network dhcp list --page 1 --page-size 50`
- Write: `ikuai-cli network dhcp create --name "Office" --interface lan1 --addr-pool 192.168.1.100-200`
- Write (escape hatch): `ikuai-cli network dns set --data '{"key":"value"}'`
- REPL: `ikuai-cli` (interactive shell)

## Error Handling

```json
{"code": 3007, "message": "Invalid token"}
```

| Code | Meaning | Fix |
|------|---------|-----|
| 3007 | Invalid/expired token | `auth set-token <TOKEN>` |
| 1008 | Session expired | `auth set-token <TOKEN>` |
| 10014 | Permission denied | Check user permissions |
