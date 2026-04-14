---
name: ikuai-routing
description: iKuai routing — static routes, traffic shunting (domain, five-tuple, L7, load-balance, updown).
---

# Routing

## 静态路由

```bash
ikuai-cli routing static list
ikuai-cli routing static create --name "to_lan2" --dst-addr "10.10.10.0" --gateway "10.66.0.1" --netmask "255.255.255.0" --interface wan1
ikuai-cli routing static toggle <ID> --enabled no
ikuai-cli routing static delete <ID>
```

## 分流规则（stream）

5 种分流，统一 CRUD 模式：

### 域名分流
```bash
ikuai-cli routing stream domain list
ikuai-cli routing stream domain create --name "baidu" --domain "www.baidu.com,baidu.com" --interface wan2
ikuai-cli routing stream domain toggle <ID> --enabled no
ikuai-cli routing stream domain delete <ID>
```

### 五元组分流
```bash
ikuai-cli routing stream five-tuple list
ikuai-cli routing stream five-tuple create --name "web_wan2" --protocol tcp --dst-port "80,443" --interface wan2
```

### L7 协议分流
```bash
ikuai-cli routing stream l7 list
ikuai-cli routing stream l7 create --name "dns_wan2" --app-proto "DNS" --interface wan2
```

### 负载均衡
```bash
ikuai-cli routing stream load-balance list
ikuai-cli routing stream load-balance create --name "lb_wan1" --interface wan1
```

### 上下行分离
```bash
ikuai-cli routing stream updown list
ikuai-cli routing stream updown create --name "split" --upiface wan1 --downiface wan2
```

所有 stream 子命令支持 `--src-addr`、`--dst-addr` 等 addrFields。
