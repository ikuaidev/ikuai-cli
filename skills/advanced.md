---
name: ikuai-advanced
description: iKuai advanced services — FTP, HTTP, Samba file servers and SNMPD config.
---

# Advanced Services

## FTP

```bash
ikuai-cli advanced ftp config-get
ikuai-cli advanced ftp config-set --open-ftp 1 --ftp-port 21 --ftp-access 1
ikuai-cli advanced ftp list
ikuai-cli advanced ftp create --name "ftpuser" --username "user1" --password "123456" --permission rw --home-dir "/sda1" --enabled yes
ikuai-cli advanced ftp toggle <ID> --enabled no
ikuai-cli advanced ftp delete <ID>
```

FTP/HTTP/Samba 用户 CRUD 需路由器挂载外接存储。

## HTTP

```bash
ikuai-cli advanced http list
ikuai-cli advanced http create --name "www" --port "8080" --home-dir "/sda1" --enabled yes
```

## Samba

```bash
ikuai-cli advanced samba config-get
ikuai-cli advanced samba config-set --enabled yes --workgroup WORKGROUP --wsdd2 1 --access 1
ikuai-cli advanced samba list
ikuai-cli advanced samba create --name "share1" --username "user1" --password "123456" --permission rw --enabled yes
```

## SNMPD

```bash
ikuai-cli advanced snmpd get
ikuai-cli advanced snmpd set --enabled yes --listen-port 161 --version 2 --community public
```

SNMPD set 是全量更新，需传所有 required 字段。
