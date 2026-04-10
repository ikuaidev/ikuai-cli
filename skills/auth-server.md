---
name: ikuai-auth-server
description: iKuai web auth server — get/set portal authentication config.
---

# Auth Server

```bash
ikuai-cli auth-server get
ikuai-cli auth-server get --columns enabled,interface,idle_time,max_time
ikuai-cli auth-server set --data '{"enabled":"no","max_time":0,"idle_time":60,...}' --user-auth 1
```

全量更新（141 字段），建议 get → 修改 JSON → set。支持 flags 覆盖 `--data` 中的值。

常用 flags：`--enabled`, `--max-time`, `--idle-time`, `--user-auth`, `--coupon-auth`, `--phone-auth`,
`--static-pwd`, `--nopasswd`, `--weixin`, `--interface`, `--passwd`, `--whitelist`, `--whitelist-https`,
`--whiteip`, `--noauth-mac`, `--radius-ip`, `--radius-key`, `--https-redirect`
