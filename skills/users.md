---
name: ikuai-users
description: iKuai user management — auth accounts CRUD, online sessions, kick, auth packages CRUD.
---

# Users

## 在线用户

```bash
ikuai-cli users online
ikuai-cli users kick <SESSION_ID>
```

## 认证账号

```bash
ikuai-cli users accounts list
ikuai-cli users accounts get <ID>
ikuai-cli users accounts create --username "guest1" --password "123456"
# defaults: ppptype=any, upload=0, download=0, share=1, expires=0
ikuai-cli users accounts delete <ID>
```

可选 flags：`--ppptype`, `--upload`, `--download`, `--share`

## 套餐管理

```bash
ikuai-cli users packages list
ikuai-cli users packages get <ID>
ikuai-cli users packages create --name "月卡" --time "1m" --price 100 --up-speed 500 --down-speed 1000
ikuai-cli users packages delete <ID>
```

`--time` 格式：`30d`(天), `1m`(月), `24h`(小时)
