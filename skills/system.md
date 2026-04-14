---
name: ikuai-system
description: iKuai system config — hostname, NTP, reboot schedules, remote access, VRRP, ALG, kernel, CPU freq, web password reset.
---

# System

## 基础配置

```bash
ikuai-cli system get
ikuai-cli system set --hostname "my-router" --language 1 --time-zone 8
ikuai-cli system ntp-sync
```

## 定时重启

```bash
ikuai-cli system schedules list
ikuai-cli system schedules create --name "weekly_reboot" --strategy week --cycle-time "7" --time "03:00"
# defaults: event=reboot
ikuai-cli system schedules toggle <ID> --enabled no
ikuai-cli system schedules delete <ID>
```

## 远程访问

```bash
ikuai-cli system remote-access get
ikuai-cli system remote-access set --ssh "1" --ssh-port 22 --wan-web "1" --http-port 80 --https-port 443
```

## VRRP（高可用）

```bash
ikuai-cli system vrrp get
ikuai-cli system vrrp set --type "1" --priority "150" --gateway "192.168.1.1" --enabled yes
ikuai-cli system vrrp start
ikuai-cli system vrrp stop
```

## ALG

```bash
ikuai-cli system alg get
ikuai-cli system alg set --ftp "1" --sip "1" --tftp "1"
```

## 内核参数

```bash
ikuai-cli system kernel get
ikuai-cli system kernel set --bbr "1" --established-timeout "3600"
```

## CPU 调频

```bash
ikuai-cli system cpufreq get
ikuai-cli system cpufreq set --mode "performance" --turbo "1"
ikuai-cli system cpufreq mode-set --mode "powersave"
```

## Web 密码重置

```bash
ikuai-cli system web-passwd reset --ssh-user root --ssh-password "pass" --yes
# 需要 SSH 访问，--save 保存 SSH 凭证
```
