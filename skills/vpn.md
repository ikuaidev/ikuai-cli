---
name: ikuai-vpn
description: iKuai VPN — PPTP, L2TP, OpenVPN, IKEv2, IPSec, WireGuard server config, client CRUD, tunnels and peers.
---

# VPN

## PPTP / L2TP

```bash
# 服务配置
ikuai-cli vpn pptp get
ikuai-cli vpn pptp set --enabled yes --server-ip "10.0.0.1" --server-port 1723 --addr-pool "10.0.0.2-10.0.0.254" --dns1 "114.114.114.114" --dns2 "119.29.29.29" --open-mppe 2 --mtu 1400 --mru 1400
# 全量更新，需传所有字段

# 客户端 CRUD
ikuai-cli vpn pptp clients
ikuai-cli vpn pptp client-create --name pptp_office --server "vpn.example.com" --username user1 --password "123456" --interface auto
ikuai-cli vpn pptp client-update <ID> --server "vpn2.example.com"
ikuai-cli vpn pptp kick <ID>
```

L2TP server set 同理：`--enabled`, `--server-ip`, `--server-port`, `--addr-pool`, `--dns1`, `--dns2`, `--mtu`, `--mru`, `--ipsec-secret`, `--leftid`, `--rightid`, `--force-ipsec`（全量更新）。
L2TP 客户端 name 需以 `l2tp` 开头。

OpenVPN/IKEv2 server set 字段多（含证书），建议用 `--data` 全量传 JSON。

## OpenVPN

```bash
ikuai-cli vpn openvpn get
ikuai-cli vpn openvpn clients
ikuai-cli vpn openvpn client-create --name ovpn_office --remote-addr "vpn.example.com" --username user1 --password "123456" --ca "<CA证书>" --interface auto
ikuai-cli vpn openvpn kick <ID>
```

`--ca` 必须传真实 CA 证书内容。

## IKEv2

```bash
ikuai-cli vpn ikev2 get
ikuai-cli vpn ikev2 clients
ikuai-cli vpn ikev2 client-create --name iked_office --remote-addr "vpn.example.com" --interface auto --left-id "localid" --username user1 --password "123456"
ikuai-cli vpn ikev2 kick <ID>
```

name 需以 `iked` 开头，authby=mschapv2 时 username 必填。

## IPSec

```bash
ikuai-cli vpn ipsec clients
ikuai-cli vpn ipsec client-create --name ipsec_site --remote-addr "10.0.0.1" --interface wan1 --left-subnet "192.168.1.0/24" --right-subnet "192.168.2.0/24" --secret "psk123"
# defaults: keyexchange=ikev2, authby=secret, dpdaction=none, ikelifetime=3, lifetime=1
ikuai-cli vpn ipsec kick <ID>
```

## WireGuard

```bash
# 隧道
ikuai-cli vpn wireguard list
ikuai-cli vpn wireguard create --name wg_site --address "10.9.0.1/24" --private-key "<base64>" --public-key "<base64>"
# defaults: interface=auto, port=5000, mtu=1420
ikuai-cli vpn wireguard get <ID>
ikuai-cli vpn wireguard toggle <ID> --enabled no
ikuai-cli vpn wireguard delete <ID>

# 对端
ikuai-cli vpn wireguard peers <TUNNEL_ID>
ikuai-cli vpn wireguard peer-create <TUNNEL_ID> --public-key "<base64>" --allow-ips "10.9.0.2/32" --interface wg0
ikuai-cli vpn wireguard peer-delete <TUNNEL_ID> <PEER_ID>
```
