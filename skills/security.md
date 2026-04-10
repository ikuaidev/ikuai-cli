---
name: ikuai-security
description: iKuai security rules — ACL, MAC filtering, L7 app rules, URL filtering, domain blacklist, peerconn, terminals.
---

# Security

## ACL（访问控制）

```bash
ikuai-cli security acl list
ikuai-cli security acl get <ID>
ikuai-cli security acl create --name "block_ssh" --action drop --protocol tcp --dst-port "22" --enabled yes
ikuai-cli security acl toggle <ID> --enabled no
ikuai-cli security acl delete <ID>
```

`--src-addr`, `--dst-addr`, `--src-port`, `--dst-port` 支持逗号分隔（addrFields）。

## MAC 过滤

```bash
ikuai-cli security mac get-mode
ikuai-cli security mac set-mode --acl-mac 0
# --data '{"acl_mac":0}' 也可用
ikuai-cli security mac list
ikuai-cli security mac create --name "allow1" --mac "00:11:22:33:44:55" --enabled yes
ikuai-cli security mac toggle <ID> --enabled no
ikuai-cli security mac delete <ID>
```

## L7 应用层规则

```bash
ikuai-cli security l7 list
ikuai-cli security l7 get <ID>
ikuai-cli security l7 create --name "block_p2p" --action drop --app-proto "BT,eMule" --enabled yes
ikuai-cli security l7 toggle <ID> --enabled no
ikuai-cli security l7 delete <ID>
```

## URL 过滤

```bash
# 黑名单
ikuai-cli security url black list
ikuai-cli security url black create --name "block_ads" --mode 0 --domain "ads.example.com" --enabled yes
ikuai-cli security url black delete <ID>

# 关键词
ikuai-cli security url keywords list
ikuai-cli security url keywords create --name "kw1" --mode exact --src-url "example.com" --ori-keyword "bad" --rep-keyword "good" --enabled yes

# 重定向
ikuai-cli security url redirect list
ikuai-cli security url redirect create --name "redir1" --mode exact --src-url "old.com" --dst-url "new.com" --enabled yes

# 替换
ikuai-cli security url replace list
ikuai-cli security url replace create --name "rep1" --mode exact --src-url "example.com" --param-keyword "track" --rep-keyword "" --enabled yes
```

## 域名黑名单

```bash
ikuai-cli security domain-blacklist list
ikuai-cli security domain-blacklist create --name "blocked" --domain-group "evil.com" --enabled yes
ikuai-cli security domain-blacklist delete <ID>
```

## 连接数限制（Peerconn）

```bash
ikuai-cli security peerconn list
ikuai-cli security peerconn create --name "limit1" --limits 500 --protocol tcp --src-addr "192.168.9.0/24" --enabled yes
```

## 终端标注（Terminals）

```bash
ikuai-cli security terminals list
ikuai-cli security terminals create --name "printer" --mac "AA:BB:CC:DD:EE:FF"
```

## 高级安全配置

```bash
ikuai-cli security advanced-get
ikuai-cli security advanced-set --data '{"noping_lan":0,"noping_wan":1,"notracert":0,"hijack_ping":0,"invalid":0,"dos_lan":1,"dos_lan_num":500,"tcp_mss":1,"tcp_mss_num":1400}'
# 全量更新，需传所有必填字段；也支持 --noping-wan 等 flags 覆盖
ikuai-cli security secondary-route-get
ikuai-cli security secondary-route-set --data '{"nol2rt":0,"nol2rt_ip":{"custom":[],"object":[]},"ttl_num":1,"time":""}'
# 全量更新；也支持 --nol2rt / --ttl-num flags
```
