---
name: ikuai-wireless
description: iKuai wireless — blacklist/whitelist rules, VLAN rules, AC management.
---

# Wireless

## 黑白名单

```bash
ikuai-cli wireless blacklist list
ikuai-cli wireless blacklist create --name "block1" --mac "00:11:22:33:44:55" --enabled yes
# defaults: mode=0(黑名单), lssid=ALL, lap=ALL, week=1234567, time=00:00-23:59
ikuai-cli wireless blacklist toggle <ID> --enabled no
ikuai-cli wireless blacklist delete <ID>
```

## 无线 VLAN

```bash
ikuai-cli wireless vlan list
ikuai-cli wireless vlan create --name "iot_vlan" --vlan-id 100 --mac "00:11:22:33:44:55" --enabled yes
ikuai-cli wireless vlan toggle <ID> --enabled no
ikuai-cli wireless vlan delete <ID>
```

## AC 管理

```bash
ikuai-cli wireless ac get
ikuai-cli wireless ac start
ikuai-cli wireless ac stop
# ac set 为预留接口（PUT），AC 开关用 start/stop
ikuai-cli wireless ac ap-list
ikuai-cli wireless ac ap-get <ID>
ikuai-cli wireless ac ap-update <ID> --ssid1 "NewSSID" --enc1 wpa2 --key1 "12345678"
# 不常用字段仍可用 --data 传
```
