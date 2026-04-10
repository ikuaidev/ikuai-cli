# CLI Reference

Full command reference for ikuai-cli. See [README](../README.md) for Quick Start.

## Command Shape

```text
ikuai-cli <resource> <action> [args] [flags]
```

## Global Flags

```bash
-f, --format table|json|yaml   # Output format (default: table)
    --raw                      # Full API envelope (debug); mutually exclusive with --format
    --dry-run                  # Preview API request without executing
    --human-time               # Convert timestamps to human-readable local time
```

## List Flags

Common flags for all `list` commands:

```bash
-p, --page INT          # Page number (default: 1)
    --page-size INT     # Items per page (default: 100)
-L, --limit INT         # Total items limit with auto-pagination
    --filter STRING     # Filter expression, e.g. "enabled==true"
-o, --order-by STRING   # Sort field
    --order asc|desc    # Sort direction
```

## Monitor Load Flags

Common flags for `cpu`, `memory`, `disk`, `temp`, `terminals`, `connections`, `network-load`:

```bash
--time-range hour|day|week|month   # Time window (default: hour)
--aggregate avg|max                # Aggregation method (default: avg)
```

---

## Auth

```bash
ikuai-cli auth set-url https://192.168.1.1   # Set router base URL
ikuai-cli auth set-token <your-api-token>    # Set API Bearer token
ikuai-cli auth status                        # Show session info
ikuai-cli auth status --format json          # Compact JSON
ikuai-cli auth clear                         # Clear host + token
```

## Monitor

```bash
ikuai-cli monitor system                     # CPU, memory, uptime, WAN IP
ikuai-cli monitor cpu                        # CPU load (default: last hour, avg)
ikuai-cli monitor cpu --time-range day --aggregate max
ikuai-cli monitor memory                     # Memory usage history
ikuai-cli monitor disk                       # Disk usage history
ikuai-cli monitor interfaces-traffic         # Per-interface traffic
ikuai-cli monitor clients-online             # Online IPv4 clients
ikuai-cli monitor client-protocols --ip 192.168.1.100 --mac 08:9b:4b:01:7e:7c
ikuai-cli monitor client-protocols-history --ip 192.168.1.100 --mac 08:9b:4b:01:7e:7c
ikuai-cli monitor client-app-protocols --ip 192.168.1.100 --mac 08:9b:4b:01:7e:7c
ikuai-cli monitor traffic-load --ip 192.168.1.100 --mac 08:9b:4b:01:7e:7c
ikuai-cli monitor app-protocols-history --appids 2580003,2580004
ikuai-cli monitor app-protocols-terminals --appid 2580003
```

## Network

```bash
# DNS
ikuai-cli network dns get
ikuai-cli network dns set --dns1 114.114.114.114 --dns2 8.8.8.8

# WAN / LAN
ikuai-cli network wan
ikuai-cli network lan

# DHCP
ikuai-cli network dhcp list --page 1 --page-size 50
ikuai-cli network dhcp create --name "Office" --interface lan1 --addr-pool 192.168.1.100-192.168.1.200
ikuai-cli network dhcp toggle 1 --enabled yes
ikuai-cli network dhcp static list
ikuai-cli network dhcp static create --name "Printer" --ip 192.168.1.50 --mac AA:BB:CC:DD:EE:FF
ikuai-cli network dhcp access-mode get

# NAT / DNAT
ikuai-cli network nat list --filter "enabled==true" --order asc --order-by id
ikuai-cli network nat create --name "Web" --in-interface wan1 --action DNAT
ikuai-cli network dnat create --name "SSH" --wan-port 2222 --lan-addr 192.168.1.10 --lan-port 22 --protocol tcp

# VLAN
ikuai-cli network vlan create --name "IoT" --vlan-id 100 --interface lan1
```

## Users

```bash
ikuai-cli users accounts list --format json
ikuai-cli users accounts create --data '{"username":"guest","password":"guest123"}'
ikuai-cli users online
ikuai-cli users kick --data '{"id":1}'
ikuai-cli users packages list
```

## System

```bash
ikuai-cli system get                         # System config
ikuai-cli system set --hostname "ikuai-gw"
ikuai-cli system schedules list
ikuai-cli system schedules create --name "NightReboot" --time "04:00"
ikuai-cli system remote-access get
ikuai-cli system remote-access set --data '{"enabled":true}'
ikuai-cli system vrrp get
ikuai-cli system alg get
ikuai-cli system kernel get
ikuai-cli system cpufreq get
ikuai-cli system web-passwd reset --ssh-user root --yes
```

## Security

```bash
ikuai-cli security acl list
ikuai-cli security acl create --name "BlockSSH" --action drop --protocol tcp --direction in --priority 100
ikuai-cli security mac get-mode
ikuai-cli security url list
ikuai-cli security domain-blacklist list
ikuai-cli security l7 list
```

## VPN

```bash
ikuai-cli vpn pptp get
ikuai-cli vpn pptp clients
ikuai-cli vpn l2tp get
ikuai-cli vpn l2tp clients
ikuai-cli vpn openvpn get
ikuai-cli vpn ikev2 get
ikuai-cli vpn ipsec clients
ikuai-cli vpn wireguard list
ikuai-cli vpn wireguard peers
```

## Routing

```bash
ikuai-cli routing static list
ikuai-cli routing stream domain list
ikuai-cli routing stream five-tuple list
ikuai-cli routing stream l7 list
```

## QoS

```bash
ikuai-cli qos ip list
ikuai-cli qos ip create --name "Limit" --ip-addr 192.168.1.100 --upload 10M --download 50M
ikuai-cli qos mac list
```

## Log

```bash
ikuai-cli log system list --page 1
ikuai-cli log system list --human-time       # Human-readable timestamps
ikuai-cli log arp list
ikuai-cli log auth list
ikuai-cli log dhcp list
ikuai-cli log web list
ikuai-cli log wireless list
```

## Wireless

```bash
ikuai-cli wireless ac get
ikuai-cli wireless blacklist list
ikuai-cli wireless vlan list
```

## Advanced

```bash
ikuai-cli advanced ftp config-get
ikuai-cli advanced http config-get
ikuai-cli advanced samba config-get
ikuai-cli advanced snmpd config-get
```

## Auth Server

```bash
ikuai-cli auth-server get
ikuai-cli auth-server set --data '{"enabled":true}'
```

## Objects

```bash
ikuai-cli objects ip list
ikuai-cli objects mac list
ikuai-cli objects port list
ikuai-cli objects domain list
ikuai-cli objects time list
```

## Other

```bash
ikuai-cli version                            # Build info
ikuai-cli version --format json
ikuai-cli completion bash                    # Shell completion
ikuai-cli completion zsh
ikuai-cli completion fish
ikuai-cli completion powershell
ikuai-cli repl                               # Interactive shell
```
