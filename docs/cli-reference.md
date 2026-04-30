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
ikuai-cli auth set-url https://192.168.1.1
ikuai-cli auth set-token <your-api-token>
ikuai-cli auth status
ikuai-cli auth clear
```

## Monitor

```bash
ikuai-cli monitor system
ikuai-cli monitor cpu
ikuai-cli monitor cpu --time-range day --aggregate max
ikuai-cli monitor memory
ikuai-cli monitor disk
ikuai-cli monitor temp
ikuai-cli monitor terminals
ikuai-cli monitor connections
ikuai-cli monitor network-load
ikuai-cli monitor downstream
ikuai-cli monitor interfaces
ikuai-cli monitor interfaces-traffic
ikuai-cli monitor interfaces-config
ikuai-cli monitor interfaces-physical
ikuai-cli monitor interfaces-traffic-v6
ikuai-cli monitor clients-online
ikuai-cli monitor clients-offline
ikuai-cli monitor clients-ip6-online
ikuai-cli monitor clients-ip6-offline
ikuai-cli monitor traffic-summary
ikuai-cli monitor traffic-load --ip 192.168.1.100 --mac 08:9b:4b:01:7e:7c
ikuai-cli monitor client-protocols --ip 192.168.1.100 --mac 08:9b:4b:01:7e:7c
ikuai-cli monitor client-protocols-history --ip 192.168.1.100 --mac 08:9b:4b:01:7e:7c
ikuai-cli monitor client-app-protocols --ip 192.168.1.100 --mac 08:9b:4b:01:7e:7c
ikuai-cli monitor protocols
ikuai-cli monitor protocols-history
ikuai-cli monitor app-traffic-summary
ikuai-cli monitor app-protocols-load
ikuai-cli monitor app-protocols-history --appids 2580003,2580004 --starttime 1773215100 --stoptime 1773218700
ikuai-cli monitor app-protocols-terminals --appid 2580003
ikuai-cli monitor wireless-stats
ikuai-cli monitor wireless-score
ikuai-cli monitor wireless-traffic
ikuai-cli monitor ssid-clients
ikuai-cli monitor channel-clients
ikuai-cli monitor cameras
ikuai-cli monitor flow-shunting
ikuai-cli monitor switch
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
ikuai-cli users accounts list
ikuai-cli users accounts create --data '{"username":"guest","password":"guest123"}'
ikuai-cli users online
ikuai-cli users kick --data '{"id":1}'
ikuai-cli users packages list
```

## System

```bash
ikuai-cli system get
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
ikuai-cli security acl create --name "BlockSSH" --action drop --protocol tcp --direction forward --dst-port 22 --priority 30 --enabled no
ikuai-cli security mac get-mode
ikuai-cli security mac set-mode --acl-mac 0
ikuai-cli security url black list
ikuai-cli security url keywords create --name "KW1" --mode exact --src-url "example.com" --ori-keyword "bad" --rep-keyword "good" --hit-rate 1 --priority 10 --enabled no
ikuai-cli security domain-blacklist list
ikuai-cli security l7 list
ikuai-cli security advanced-get
ikuai-cli security secondary-route-get
```

## VPN

```bash
ikuai-cli vpn pptp get
ikuai-cli vpn pptp clients
ikuai-cli vpn pptp client-create --name pptpoffice --server vpn.example.com --username user1 --password 123456 --interface auto --enabled no
ikuai-cli vpn pptp client-get <ID>
ikuai-cli vpn pptp client-update <ID> --comment updated
ikuai-cli vpn pptp client-toggle <ID> --enabled no
ikuai-cli vpn pptp client-delete <ID> --yes
ikuai-cli vpn l2tp get
ikuai-cli vpn l2tp clients
ikuai-cli vpn openvpn get
ikuai-cli vpn openvpn client-create --name ovpnoffice --remote-addr vpn.example.com --username user1 --password 123456 --ca "<CA证书>" --interface auto --enabled no
ikuai-cli vpn ikev2 get
ikuai-cli vpn ikev2 client-create --name ikedoffice --remote-addr vpn.example.com --interface auto --left-id localid --username user1 --password 123456 --enabled no
ikuai-cli vpn ipsec clients
ikuai-cli vpn ipsec client-create --name ipsecsite --remote-addr 10.0.0.1 --interface wan1 --left-subnet 192.168.1.0/24 --right-subnet 192.168.2.0/24 --secret psk123 --enabled no
ikuai-cli vpn ipsec client-get <ID>
ikuai-cli vpn ipsec client-update <ID> --comment updated
ikuai-cli vpn ipsec client-toggle <ID> --enabled no
ikuai-cli vpn ipsec client-delete <ID> --yes
ikuai-cli vpn wireguard list
ikuai-cli vpn wireguard create --name wgsite --address 10.9.0.1/24 --interface auto --private-key "<base64>" --public-key "<base64>" --enabled no
ikuai-cli vpn wireguard update <ID> --interface auto --port 5001
ikuai-cli vpn wireguard peers <ID>
ikuai-cli vpn wireguard peer-create <ID> --public-key "<base64>" --allow-ips 10.9.0.2/32 --interface wgsite --enabled no
ikuai-cli vpn wireguard peer-get <ID> <PEER_ID>
ikuai-cli vpn wireguard peer-update <ID> <PEER_ID> --interface wgsite --comment updated
ikuai-cli vpn wireguard peer-toggle <ID> <PEER_ID> --enabled no
ikuai-cli vpn wireguard peer-delete <ID> <PEER_ID> --yes
```

## Routing

```bash
ikuai-cli routing static list
ikuai-cli routing static get <ID>
ikuai-cli routing static create --name to_lan2 --dst-addr 10.10.10.0 --gateway 10.66.0.1 --netmask 255.255.255.0 --interface wan1
ikuai-cli routing static update <ID> --name to_lan2_bak --priority 10
ikuai-cli routing static toggle <ID> --enabled no
ikuai-cli routing static delete <ID> --yes
ikuai-cli routing stream domain list
ikuai-cli routing stream domain get <ID>
ikuai-cli routing stream domain create --name baidu --domain www.baidu.com,baidu.com --interface wan2
ikuai-cli routing stream domain update <ID> --name baidu_wan1 --interface wan1
ikuai-cli routing stream domain toggle <ID> --enabled no
ikuai-cli routing stream domain delete <ID> --yes
ikuai-cli routing stream five-tuple list
ikuai-cli routing stream five-tuple get <ID>
ikuai-cli routing stream five-tuple create --name web_wan2 --protocol tcp --dst-port 80,443 --interface wan2
ikuai-cli routing stream five-tuple update <ID> --dst-port 8080,8443
ikuai-cli routing stream five-tuple toggle <ID> --enabled no
ikuai-cli routing stream five-tuple delete <ID> --yes
ikuai-cli routing stream l7 list
ikuai-cli routing stream l7 get <ID>
ikuai-cli routing stream l7 create --name dns_wan2 --app-proto DNS --interface wan2
ikuai-cli routing stream l7 update <ID> --name dns_wan1 --interface wan1
ikuai-cli routing stream l7 toggle <ID> --enabled no
ikuai-cli routing stream l7 delete <ID> --yes
ikuai-cli routing stream load-balance list
ikuai-cli routing stream load-balance get <ID>
ikuai-cli routing stream load-balance create --name lb_wan1 --interface wan1 --mode 0 --weight 1 --isp-name all
ikuai-cli routing stream load-balance update <ID> --weight 2
ikuai-cli routing stream load-balance toggle <ID> --enabled no
ikuai-cli routing stream load-balance delete <ID> --yes
ikuai-cli routing stream updown list
ikuai-cli routing stream updown get <ID>
ikuai-cli routing stream updown create --name split --upiface wan1 --downiface wan2
ikuai-cli routing stream updown update <ID> --protocol tcp --dst-port 80,443
ikuai-cli routing stream updown toggle <ID> --enabled no
ikuai-cli routing stream updown delete <ID> --yes
```

## QoS

```bash
ikuai-cli qos ip list
ikuai-cli qos ip create --name limit100 --ip-addr 192.168.1.100 --interface wan1 --upload 1000 --download 1000
ikuai-cli qos ip update <ID> --name limit100u --upload 1200 --download 1300
ikuai-cli qos ip toggle <ID> --enabled no
ikuai-cli qos ip delete <ID> --yes
ikuai-cli qos mac list
ikuai-cli qos mac create --name limitmac --mac-addr 00:11:22:33:44:55 --interface wan1 --upload 500 --download 500
ikuai-cli qos mac delete <ID> --yes
```

## Log

```bash
ikuai-cli log system list --page 1
ikuai-cli log system list --human-time
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
# FTP
ikuai-cli advanced ftp config-get
ikuai-cli advanced ftp config-set --open-ftp 1 --ftp-port 21 --ftp-access 1
ikuai-cli advanced ftp list
ikuai-cli advanced ftp create --username "user1" --password "pass" --permission rw --home-dir "/sda1"
ikuai-cli advanced ftp toggle <ID> --enabled no
ikuai-cli advanced ftp delete <ID>

# HTTP
ikuai-cli advanced http list
ikuai-cli advanced http create --name "www" --port 8080 --ssl 0 --autoindex 0 --download 0 --home-dir "/sda1"
ikuai-cli advanced http toggle <ID> --enabled no
ikuai-cli advanced http delete <ID>

# Samba
ikuai-cli advanced samba config-get
ikuai-cli advanced samba config-set --enabled yes --workgroup WORKGROUP
ikuai-cli advanced samba list
ikuai-cli advanced samba create --name "share1" --username "user1" --password "pass" --permission rw --guest yes --home-dir "/sda1"
ikuai-cli advanced samba toggle <ID> --enabled no
ikuai-cli advanced samba delete <ID>

# SNMPD
ikuai-cli advanced snmpd get
ikuai-cli advanced snmpd set --enabled yes --listen-port 161 --version 2 --community public
```

## Auth Server

```bash
ikuai-cli auth-server get
ikuai-cli auth-server set --data '{"enabled":true}'
```

## Objects

```bash
ikuai-cli objects ip list
ikuai-cli objects ip create --name servers --value 192.168.1.10,192.168.1.11
ikuai-cli objects ip get <ID>
ikuai-cli objects ip update <ID> --name servers_v2 --value 192.168.1.12
ikuai-cli objects ip refs --group-name servers_v2
ikuai-cli objects ip delete <ID> --yes
ikuai-cli objects ip6 list
ikuai-cli objects mac list
ikuai-cli objects port list
ikuai-cli objects proto list
ikuai-cli objects domain list
ikuai-cli objects time create --name office --type weekly --weekdays 12345 --start-time 09:00 --end-time 18:00
ikuai-cli objects time list
```

## Other

```bash
ikuai-cli version
ikuai-cli completion bash
ikuai-cli completion zsh
ikuai-cli completion fish
ikuai-cli completion powershell
ikuai-cli repl
```
