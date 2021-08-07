---
title: "KVMとヤマハルータでタグVLAN"
date: "2020-12-06"
---

こんにちは、ごっちと言います。UEC20で工研に入っていたりします。以下の画像のような圧力を掛けられ、夜中に泣きながらクソ記事を錬成しています。(今日付が変わって12/6になったあたりです。はやぶさ2のカプセル再突入、楽しみですね!)
![上下からの強い圧力](./advent2020.jpg)


これは[UEC Advent Calendar 2020]("https://adventar.org/calendars/5070")の6日目の記事です。
前: [スマプリは黄色が一番かわいい - ぺんぎんさんのおうち]("https://ykm11.hatenablog.com/entry/2020/12/05/UECAC2020") ykm11さんの多倍長整数の記事です。
後: [独選！明日使えないc/c++豆知識集 - Exte’s blog]("https://exte.hateblo.jp/entry/2020/12/07/000000") Exte Externalさんのc/c++の記事です。

この記事ではKVMのゲストをホストの属するネットワークから切り離す話をします。初めて触ったのとまともに検証していないので一切を保証しません。切り離せてないぞ!などあれば[@0xgotti]("https://twitter.com/intent/user?user_id=3721840992") にまさかりを投げてください。

# 目的

KVMでゲストがホストにアクセスしてほしくないとき、あると思います。そんなときに役立つのがタグVLAN([IEEE 802.1Q]("https://ja.wikipedia.org/wiki/IEEE_802.1Q"))です。

# しくみ

タグVLANはL2スイッチなどでパケットのイーサネットフレームにVLAN用の情報を追加したり読みとったりすることで、物理的には1つのネットワークに属するものを論理的に分割することができます。今回の目的はKVMのゲストOSをホストのネットワークから切り離すことなので、ホストからタグ付けをやります。そしてルータでタグが付いたものを別のサブネットに置くようにしましょう。最後にルータでサブネット間の通信を遮断すれば隔離できるはずです。


# やったこと
環境
- ホストOS: ArchLinux (2020年12月ぐらいのやつ)
- ゲストOS: CentOS8 (一応書いていますがなんでもいいと思います)
- ルータ  : YAMAHA NVR500

どうしてかホストはArchLinuxです、CentOS8をホストにしろよという声が聞こえてきますがぐうの音も出ません。私もそう思います。

やったこと
netctlでVLANあたりを弄れます。サンプルを見てみましょう。/etc/netctl/examples/vlan-dhcpにあります。(DHCPを使うなら)
```config
Description='Virtual LAN 55 on interface eth0 using DHCP'
Interface=eth0.55
Connection=vlan
# The variable name is plural, but needs precisely one interface
BindsToInterfaces=eth0
VLANID=55
IP=dhcp
```

/etc/netctl/vlan-dhcpに持ってきて以下のように編集します。xxはVLANのIDで好きな番号にしてください。同じIDだと同じVLANに属することになります。なおホストのNICはeno1です。

```config
Description='Virtual LAN xx on interface eno1 using DHCP'
Interface=eno1.xx
Connection=vlan
# The variable name is plural, but needs precisely one interface
BindsToInterfaces=eno1
VLANID=xx
IP=dhcp
MACAddress=好きなMACアドレス
```
これでvlanしてくれる気がします。次にKVMから繋げられるように、eno1.xxとKVMとを繋げるブリッジを作りましょう。ブリッジのサンプルは/etc/netctl/examples/bridgeにあります。見てみましょう。
```config
Description="Example Bridge connection"
Interface=br0
Connection=bridge
BindsToInterfaces=(eth0 eth1 tap0)
## Use the MAC address of eth1 (may also be set to a literal MAC address)
#MACAddress=eth1
IP=dhcp
## Ignore (R)STP and immediately activate the bridge
#SkipForwardingDelay=yes
```
/etc/netctl/bridgeにコピーし編集しました。以下のようになりました。
```config
Description="Bridge connection"
Interface=br0
Connection=bridge
BindsToInterfaces=(eno1.xx)
IP=dhcp
MACAddress=好きなMACアドレス2
```
これで
```
インターネット <----> NVR500 <----> eno1 <----> eno1.xx <----> br0 <----> KVMの何か
```
のようになる....はずです。
次はルータ側の設定をやりましょう。YAMAHAが[タグVLANの設定例を公開]("https://network.yamaha.com/setting/switch_swx/simple_smart/switch_swx-command/tag_vlan")してくれています。コピペして弄りましょう。以下のようになりました。

xxはVLANのIDです。
```config
vlan lan1/1 802.1q vid=xx name=VLANxx
ip lan1/1 address 192.168.y.1/24               #VLANに割り当てる範囲です。好きな範囲にしてください。
ip lan1/1 secure filter in 1040 1050
ip lan1/1 secure filter out 1041 1050
ip filter 1040 reject * 192.168.z.0/24         #zはホストなどが属するようにしてください。ゲストからホストのネットワークへの接続を拒否するようにします。
ip filter 1041 reject 192.168.z.0/24 *         #zは同上。ホストのネットワークからゲストのネットワークへの接続を拒否します。
ip filter 1050 pass * * * * *                  #その他は通します。
dhcp scope xx 192.168.y.2-192.168.y.191/24     #DHCPでアドレスを割り当てます
```
ルータの設定はおしまいです。あとは実行すればKVMの全ての通信にVLANのタグが付けられ、ルータにより別のサブネットに置かれます。インターフェースやKVMを立ち上げてみましょう。
次のコマンドを打つと設定したようにブリッジなどが立ち上がります。
```bash
sudo netctl start vlan-dhcp
sudo netctl start bridge
```

次のコマンドでブリッジを通してKVMを起動できます。
```bash
qemu-system-x86_64 -drive file=file,format=qcow2 --enable-kvm -m 2G -net nic -net bridge,br=br0
```
おわりです、手順抜けてるかもしれませんが頑張って補完すれば多分できます。再度ですが隔離できてないぞ!などあれば連絡をくれると嬉しいです。
