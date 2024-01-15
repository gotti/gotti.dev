---
title: "VyOSでDS-LiteとDHCPv6-PD"
date: "2024-01-15"
tags: ["VyOS", "Linux", "Network"]
---

NTTの光クロス(10G回線)を契約した．ISPはASAHIネット．NTTの黒いルータが送られてきたが，せっかくなのでVyOSで接続することにした．

## DS-Liteの終端アドレス(AFTR)を確認する

DS-Liteはipip6トンネルで接続できる．トンネルを張るための対向終端アドレスを調べる．

NTT東のDNSサーバに問い合わせる．

```
$ dig TXT 4over6.info @2404:1a8:7f01:b::3

;; ANSWER SECTION:
4over6.info.		3600	IN	TXT	"v=v6mig-1 url=https://prod.v6mig.v6connect.net/cpe/v1/config t=b"
```

urlの部分にオプションを付けてcurlで問い合わせる．[オプションはここを参照](https://github.com/v6pc/v6mig-prov/blob/1.1/spec.md)．一応コピペで動く．

```
$ curl "https://prod.v6mig.v6connect.net/cpe/v1/config?vendorid=acde48-v6pc_swg_hgw&product=V6MIG-ROUTER&version=0_00&capability=dslite"
{"ttl":86400,"token":"<redacted>","service_name":"v6 コネクト","enabler_name":"v6 コネクト","dslite":{"aftr":"dslite.v6connect.net"},"order":["dslite"]}
```

aftrの部分が終端アドレスとなる．digで問い合わせると，AFTRのアドレスがわかる．

```
$ dig AAAA dslite.v6connect.net @2404:1a8:7f01:b::3
;; ANSWER SECTION:
dslite.v6connect.net.	3600	IN	AAAA	2001:<redacted>
```

## VyOSでDHCPv6-PD

使っているVyOSのバージョンは以下．[DHCPv6-PDまわりがバグってるので注意](https://forum.vyos.io/t/how-to-find-installed-vyatta-release/683)．

```
Version:          VyOS 1.3.0-rc6
```

DHCPv6-PDでIPv6アドレスを取得する．
64bitと指定してるのに設定によっては56bitが降ってくることがある．これもバグ？
[DUIDについてはこれを参照のこと](https://blog.ytn86.net/2020/02/edgerouter-dhcp-pd-ntteast-flets/)．

```
set interfaces ethernet eth0 dhcpv6-options duid '00:03:00:01:<MAC Address of eth0>'
set interfaces ethernet eth0 dhcpv6-options pd 0 interface eth1 sla-id '1'
set interfaces ethernet eth0 dhcpv6-options pd 0 length '64'
#set interfaces ethernet eth0 hw-id '<MAC Address of eth0>'←コレ
```

うまくいくと，56bitか64bitのプレフィックスが降ってくる．

```
$ ip a
eth1
    inet6 <redacted>/56 scope global 
       valid_lft forever preferred_lft forever
```

/56の場合，適当に/64を切り出してRAでLAN側に配布する．/64の場合は`::/64`でもいけるのかな？

```
set service router-advert interface eth1 default-preference 'high'
set service router-advert interface eth1 prefix <redacted>::/64
```

## VyOSでDS-Lite

最初に手に入れたAFTRのアドレスへipip6トンネルを張る．NATは向こうがやってくれるので書かなくてもよい．

```
set interfaces tunnel tun0 address '192.168.1.1/32' # ルータのIPv4アドレスを書いておかないと，ルータ自身がインターネットに問い合わせられない．
set interfaces tunnel tun0 encapsulation 'ipip6'
set interfaces tunnel tun0 multicast 'disable'
set interfaces tunnel tun0 remote '<AFTR Address>' # AFTRのアドレスに置換
set interfaces tunnel tun0 source-address '<Router Address>' #ルータのIPv6アドレス．eth1に割り当てられてるIPでよい．
set protocols static interface-route 0.0.0.0/0 next-hop-interface tun0 # デフォルトルートをDS-Liteへ向ける
```

最後に，[KAMEが踊っているか確認する](https://www.kame.net/)．

## 課題

スピードテストを回してみるとIPv6だと4Gbps，IPv4だと3Gbpsぐらいは出た．
いっぽうVyOSのCPUが1つ使用率100%になる．VLANかルーティングのどっちが重いんだろうか．
今のインフラではルータがproxmoxというハイパーバイザ上で動いている．そのうえproxmox内でLinuxブリッジを乱用している．
SR-IOVやVLANオフロードも試してみたい．今後の課題とする．

## おわりに

DS-Liteは一つのIPv4アドレスを多くのユーザで共有するため，セッション数の制限が厳しい．
光ネクストでよく使われているPPPoEでは一つのグローバルIPを占有できていたが，[光クロスではPPPoEを提供するISPはほとんど無い](https://flets.com/cross/pppoe/isp.html)．
しばらく運用して困ったら，ZOOT NATIVEの固定IPプランを契約しようと思う．この固定IPプランもipip6トンネルを利用して接続するとのことである．
