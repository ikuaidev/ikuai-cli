---
name: ikuai-log
description: iKuai system logs — view and clear 9 types of logs (system, arp, auth, dhcp, pppoe, web, ddns, notice, wireless).
---

# Log

9 种日志，每种有 list + delete：

```bash
# 查看（支持 --page, --page-size, --filter, --key, --pattern, --human-time）
ikuai-cli log system list --format json
ikuai-cli log arp list --format json
ikuai-cli log auth list --format json
ikuai-cli log dhcp list --format json
ikuai-cli log pppoe list --format json
ikuai-cli log web list --format json
ikuai-cli log ddns list --format json
ikuai-cli log notice list --format json
ikuai-cli log wireless list --format json

# 清除（破坏性操作）
ikuai-cli log system delete --yes --format json
ikuai-cli log dhcp delete --yes --format json
# ... 同理其他 <type> delete
```
