---
name: ikuai-auth-server
description: iKuai web auth server — get/set portal authentication config.
---

# Auth Server

```bash
ikuai-cli auth-server get --format json
ikuai-cli auth-server get --columns enabled,interface,idle_time,max_time --format json
ikuai-cli auth-server set --idle-time 60 --max-time 0 --format json
```

常用 flags：`--enabled`, `--max-time`, `--idle-time`, `--user-auth`, `--coupon-auth`, `--phone-auth`,
`--static-pwd`, `--nopasswd`, `--weixin`, `--interface`, `--passwd`, `--whitelist`, `--whitelist-https`,
`--whiteip`, `--noauth-mac`, `--radius-ip`, `--radius-key`, `--https-redirect`
