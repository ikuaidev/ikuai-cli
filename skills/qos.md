---
name: ikuai-qos
description: iKuai QoS bandwidth control — IP-based and MAC-based bandwidth limiting rules.
---

# QoS

## IP 限速

```bash
ikuai-cli qos ip list
ikuai-cli qos ip create --name "limit_100m" --ip-addr "192.168.9.0/24" --upload 100 --download 100 --interface wan1
ikuai-cli qos ip get <ID>
ikuai-cli qos ip toggle <ID> --enabled no
ikuai-cli qos ip delete <ID>
```

`--ip-addr` 支持逗号分隔多 IP（addrFields 模式）。

## MAC 限速

```bash
ikuai-cli qos mac list
ikuai-cli qos mac create --name "limit_mac" --mac-addr "00:11:22:33:44:55" --upload 50 --download 50 --interface wan1
ikuai-cli qos mac toggle <ID> --enabled no
ikuai-cli qos mac delete <ID>
```

`--mac-addr` 同样是 addrFields 模式。
