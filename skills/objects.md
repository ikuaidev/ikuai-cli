---
name: ikuai-objects
description: iKuai network objects — IP, IPv6, MAC, port, protocol, domain, time object groups for rule references.
---

# Objects

7 种对象组：ip, ip6, mac, port, proto, domain, time

每种支持 CRUD + refs 查询：

```bash
# List
ikuai-cli objects ip list
ikuai-cli objects mac list

# Create（--value 逗号分隔，自动转为 group_value 数组）
ikuai-cli objects ip create --name "servers" --value "192.168.1.10,192.168.1.11"
ikuai-cli objects mac create --name "printers" --value "AA:BB:CC:DD:EE:FF"
ikuai-cli objects port create --name "web_ports" --value "80,443,8080"
ikuai-cli objects domain create --name "blocked" --value "ads.example.com,track.example.com"

# Update
ikuai-cli objects ip update <ID> --name "servers_v2" --value "192.168.1.10,192.168.1.12"

# Toggle / Delete
ikuai-cli objects ip toggle <ID> --enabled no
ikuai-cli objects ip delete <ID>

# 查看引用该对象的规则
ikuai-cli objects ip refs --group-name "servers"
```

time 对象较复杂，用 `--data` 创建。
