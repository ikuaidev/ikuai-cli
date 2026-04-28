---
name: ikuai-monitor
description: Monitor iKuai router — system status, CPU, memory, disk, temperature, online clients, interface traffic, wireless stats.
---

# Monitor

## Commands

All examples use `--format json` because this skill is primarily for agents and scripts. Use table output only for human inspection.

### System Overview
```bash
ikuai-cli monitor system --format json       # Uptime, load, firmware version, WAN IP
ikuai-cli monitor cpu --format json          # CPU load (default: last hour, avg)
ikuai-cli monitor cpu --format json --time-range day --aggregate max   # Last day, max values
ikuai-cli monitor memory --format json       # Memory usage history
ikuai-cli monitor disk --format json         # Disk usage history
ikuai-cli monitor temp --format json         # CPU temperature history
ikuai-cli monitor connections --format json  # Connection count history
ikuai-cli monitor terminals --format json    # Terminal count history
ikuai-cli monitor network-load --format json # Network load history
```

Load commands (cpu, memory, disk, temp, terminals, connections, network-load) support:
- `--time-range hour|day|week|month` — time range
- `--aggregate avg|max` — calculation method

### Traffic & Network
```bash
ikuai-cli monitor network-load --format json         # Overall network load + rate history
ikuai-cli monitor downstream --format json           # Downstream traffic detail
ikuai-cli monitor interfaces --format json           # WAN interface status (IP, gateway)
ikuai-cli monitor interfaces-traffic --format json   # Per-interface Bytes/s
ikuai-cli monitor interfaces-config --format json    # Interface config detail
ikuai-cli monitor interfaces-physical --format json  # Physical NIC info
ikuai-cli monitor interfaces-traffic-v6 --format json # IPv6 traffic
ikuai-cli monitor flow-shunting --format json        # Traffic shunting data
ikuai-cli monitor switch --format json               # Switch port monitoring
```

`downstream`, `cameras`, and `switch` are data-dependent list commands. Treat an empty `data: []` / `total: 0` response as environment-blocked, not PASS.

### Clients
```bash
ikuai-cli monitor clients-online --format json --page 1 --page-size 100    # Online IPv4 clients
ikuai-cli monitor clients-offline --format json                             # Historical offline
ikuai-cli monitor clients-ip6-online --format json                          # Online IPv6 clients
ikuai-cli monitor clients-ip6-offline --format json                         # Historical offline IPv6
ikuai-cli monitor traffic-summary --format json                             # Per-client traffic summary
```

Per-client commands (require `--ip` and `--mac`):
```bash
ikuai-cli monitor traffic-load --format json --ip 192.168.1.100 --mac 08:9b:4b:01:7e:7c
ikuai-cli monitor client-protocols --format json --ip 192.168.1.100 --mac 08:9b:4b:01:7e:7c
ikuai-cli monitor client-protocols-history --format json --ip 192.168.1.100 --mac 08:9b:4b:01:7e:7c
ikuai-cli monitor client-app-protocols --format json --ip 192.168.1.100 --mac 08:9b:4b:01:7e:7c
```

### Application & Protocol
```bash
ikuai-cli monitor protocols --format json            # Protocol distribution
ikuai-cli monitor protocols-history --format json    # Protocol history
ikuai-cli monitor app-traffic-summary --format json  # App-layer traffic summary (24h)
ikuai-cli monitor app-protocols-load --format json   # Current app protocol load
ikuai-cli monitor app-protocols-history --format json --appids 2580003,2580004  # App protocol rate history
ikuai-cli monitor app-protocols-terminals --format json --appid 2580003         # Terminals using an app
```

### Wireless
```bash
ikuai-cli monitor wireless-stats --format json    # Wireless site statistics
ikuai-cli monitor wireless-score --format json    # Quality score
ikuai-cli monitor wireless-traffic --format json  # Per-AP traffic
ikuai-cli monitor ssid-clients --format json      # SSID client distribution
ikuai-cli monitor channel-clients --format json   # Channel client distribution
ikuai-cli monitor cameras --format json           # IP camera list
```

## Common Workflows

### Health Check
```bash
ikuai-cli monitor system --format json
ikuai-cli monitor cpu --format json
ikuai-cli monitor memory --format json
ikuai-cli monitor disk --format json
ikuai-cli monitor connections --format json
```

### Traffic Investigation
```bash
# Step 1: Find top clients
ikuai-cli monitor clients-online --format json --page-size 200
# Step 2: Drill into a specific client
ikuai-cli monitor traffic-load --format json --ip 192.168.1.100 --mac 08:9b:4b:01:7e:7c
ikuai-cli monitor client-protocols --format json --ip 192.168.1.100 --mac 08:9b:4b:01:7e:7c
ikuai-cli monitor client-app-protocols --format json --ip 192.168.1.100 --mac 08:9b:4b:01:7e:7c
# Step 3: Check app-level traffic
ikuai-cli monitor app-traffic-summary --format json
ikuai-cli monitor app-protocols-terminals --format json --appid 2580003
```

### Regression Data Discovery
```bash
# Discover a real online client before per-client commands
ikuai-cli monitor clients-online --format json --page 1 --page-size 20

# Discover a real appid before app protocol detail commands
ikuai-cli monitor app-traffic-summary --format json --page 1 --page-size 20
```

For regression, per-client commands must use a real `ip_addr` and `mac` returned by `clients-online`. App protocol detail commands must use a real `appid` returned by `app-traffic-summary`.
