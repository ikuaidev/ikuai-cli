---
name: ikuai-network
description: iKuai network config — DNS, DHCP, VLAN, NAT/DNAT, WAN, LAN, PPPoE, DMZ, DNS proxy.
---

# Network

## DNS

```bash
ikuai-cli network dns get
ikuai-cli network dns set --dns1 223.5.5.5 --dns2 119.29.29.29
ikuai-cli network dns stats

# DNS 代理
ikuai-cli network dns proxy list
ikuai-cli network dns proxy create --domain "example.com" --dns-addr "8.8.8.8" --parse-type ipv4 --comment "office dns"
ikuai-cli network dns proxy update <ID> --domain "example.com" --dns-addr "1.1.1.1" --parse-type ipv4 --enabled yes
ikuai-cli network dns proxy delete <ID>
```

## DHCP

```bash
ikuai-cli network dhcp list
ikuai-cli network dhcp get <ID>
ikuai-cli network dhcp create --name "Office" --interface lan1 --phy-ifnames "all" --addr-pool "192.168.1.100-192.168.1.200" --gateway "192.168.1.1" --netmask "255.255.255.0" --lease 120 --dns1 "223.5.5.5"
ikuai-cli network dhcp update <ID> --dns1 "8.8.8.8"
ikuai-cli network dhcp toggle <ID> --enabled no
ikuai-cli network dhcp delete <ID>
ikuai-cli network dhcp clients
ikuai-cli network dhcp restart / start / stop

# 静态绑定
ikuai-cli network dhcp static list
ikuai-cli network dhcp static create --name "Printer" --ip "192.168.1.50" --mac "AA:BB:CC:DD:EE:FF" --interface lan1 --gateway "192.168.1.1"
ikuai-cli network dhcp static update <ID> --dns1 "223.5.5.5" --comment "printer"
ikuai-cli network dhcp static toggle <ID> --enabled no
ikuai-cli network dhcp static delete <ID>

# 接入控制
ikuai-cli network dhcp access-mode get
ikuai-cli network dhcp access-mode set --mode 0
ikuai-cli network dhcp access-rule list
ikuai-cli network dhcp access-rule create --name "allow-printer" --mac "AA:BB:CC:DD:EE:FF" --comment "printer"
ikuai-cli network dhcp access-rule delete <ID>

# DHCPv6
ikuai-cli network dhcp6 clients
ikuai-cli network dhcp6 access-mode get
ikuai-cli network dhcp6 access-mode set --mode 0
ikuai-cli network dhcp6 access-rule list
ikuai-cli network dhcp6 access-rule create --name "allow-v6" --mac "AA:BB:CC:DD:EE:FF" --enabled yes
ikuai-cli network dhcp6 access-rule update <ID> --name "allow-v6-new" --mac "AA:BB:CC:DD:EE:FF" --enabled yes
ikuai-cli network dhcp6 access-rule toggle <ID> --enabled no
ikuai-cli network dhcp6 access-rule delete <ID>
```

## WAN / LAN / 接口

```bash
ikuai-cli network wan
ikuai-cli network wan-vlan
ikuai-cli network lan
ikuai-cli network physical
```

## VLAN

```bash
ikuai-cli network vlan list
ikuai-cli network vlan create --name "IoT" --vlan-id 100 --interface lan1 --ip "10.0.100.1" --netmask "255.255.255.0"
ikuai-cli network vlan update <ID> --comment "iot vlan"
ikuai-cli network vlan toggle <ID> --enabled no
ikuai-cli network vlan delete <ID>
```

## NAT / DNAT / DMZ

```bash
ikuai-cli network nat list
ikuai-cli network nat create --name "allow_office" --action filter --in-interface any --out-interface any --src-addr "192.168.1.0/24" --protocol any
ikuai-cli network nat update <ID> --comment "office nat"
ikuai-cli network nat toggle <ID> --enabled no
ikuai-cli network nat delete <ID>

ikuai-cli network dnat list
ikuai-cli network dnat create --name "web_forward" --lan-addr "192.168.1.10" --lan-port 80 --wan-port 8080 --protocol tcp --interface all
ikuai-cli network dnat update <ID> --comment "web forward"
ikuai-cli network dnat toggle <ID> --enabled no
ikuai-cli network dnat delete <ID>

ikuai-cli network dmz list
ikuai-cli network dmz create --name "safe_dmz_test" --interface "203.0.113.254" --lan-addr "192.168.1.250" --protocol tcp --excl-port "80,443,18440" --enabled no
ikuai-cli network dmz update <ID> --name "safe_dmz_test" --interface "203.0.113.254" --lan-addr "192.168.1.250" --protocol tcp --excl-port "80,443,18440" --enabled no
ikuai-cli network dmz toggle <ID> --enabled no
ikuai-cli network dmz delete <ID>
```

## PPPoE

```bash
ikuai-cli network pppoe get
ikuai-cli network pppoe set --comment "maintenance" --mtu 1480 --mru 1480 --lcp-echo-interval 10 --lcp-echo-failure 3 --maxconnect 0
# set 会先读取当前配置并回填完整 PUT body；优先使用语义化 flags。--data 仅作为复杂字段 escape hatch。
```
