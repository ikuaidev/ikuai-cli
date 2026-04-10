---
name: ikuai-log
description: iKuai system logs — view and clear 9 types of logs (system, arp, auth, dhcp, pppoe, web, ddns, notice, wireless).
---

# Log

9 种日志，每种有 list + clear：

```bash
# 查看（支持 --page, --page-size, --filter, --human-time）
ikuai-cli log system --human-time
ikuai-cli log arp
ikuai-cli log auth
ikuai-cli log dhcp
ikuai-cli log pppoe
ikuai-cli log web
ikuai-cli log ddns
ikuai-cli log notice
ikuai-cli log wireless

# 清除（破坏性操作）
ikuai-cli log system-clear
ikuai-cli log dhcp-clear
# ... 同理其他 <type>-clear
```
