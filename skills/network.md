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
ikuai-cli network dns proxy create --domain "example.com" --dns-addr "8.8.8.8" --parse-type 0 --enabled yes
ikuai-cli network dns proxy delete <ID>
```

## DHCP

```bash
ikuai-cli network dhcp list
ikuai-cli network dhcp get <ID>
ikuai-cli network dhcp create --name "Office" --interface lan1 --addr-pool "192.168.1.100-200" --gateway "192.168.1.1" --netmask "255.255.255.0" --dns1 "223.5.5.5" --enabled yes
ikuai-cli network dhcp toggle <ID> --enabled no
ikuai-cli network dhcp delete <ID>
ikuai-cli network dhcp clients
ikuai-cli network dhcp restart / start / stop

# 静态绑定
ikuai-cli network dhcp static list
ikuai-cli network dhcp static create --name "Printer" --ip-addr "192.168.1.50" --mac "AA:BB:CC:DD:EE:FF" --interface lan1 --enabled yes
ikuai-cli network dhcp static toggle <ID> --enabled no
ikuai-cli network dhcp static delete <ID>

# 接入控制
ikuai-cli network dhcp access-mode get
ikuai-cli network dhcp access-mode set --mode 0
# --data '{"mode":0}' 也可用
ikuai-cli network dhcp access-rule list / create / delete
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
ikuai-cli network vlan create --name "IoT" --vlan-id 100 --interface lan1 --ip-addr "10.0.100.1" --netmask "255.255.255.0" --enabled yes
ikuai-cli network vlan toggle <ID> --enabled no
ikuai-cli network vlan delete <ID>
```

## NAT / DNAT / DMZ

```bash
ikuai-cli network nat list
ikuai-cli network nat create --name "forward_web" --action DNAT --in-interface wan1 --out-interface lan1 --src-addr "any" --dst-addr "192.168.1.10" --protocol tcp --enabled yes
ikuai-cli network nat toggle <ID> --enabled no
ikuai-cli network nat delete <ID>

ikuai-cli network dnat list
ikuai-cli network dmz list
```

## PPPoE

```bash
ikuai-cli network pppoe get
ikuai-cli network pppoe set --data '{"enabled":"no","server_ip":"10.1.1.1","addr_pool":"10.1.1.2-10.1.1.254","dns1":"114.114.114.114","dns2":"119.29.29.29","interface":"lan1","authmode":0,"mtu":1480,"mru":1480,...}'
# 全量更新（30+ 必填字段），建议 get → 修改 JSON → set；也支持 --enabled / --server-ip 等 flags 覆盖
```
