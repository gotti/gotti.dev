---
title: "自宅サーバ/Kubernetesを支える柔軟なネットワークを目指せ"
date: "2022-12-14"
tags: ["advent", "Linux", "Network", "Kubernetes"]
---

こんにちは、20のごっちです。みなさん自宅鯖を飼っていますか？もちろん1台はありますよね。

## 背景

自宅鯖を初めてやった人はまず最初に家のルータの設定を弄り、DMZや静的NATなどよくわからん設定を弄り、DDNSの設定をし、サーバにグローバルIPを割り当てていると思います。プロバイダごとの差異のアレコレやグローバルIPがかわったり面倒ですよね。

ところでサーバが2台以上になったり実家と下宿に別々にサーバ置いて連携させたいとき、IPアドレスの管理が面倒になりますよね。実家と下宿の間でVPNを張ってルーティングやればええやんという声もあると思いますが、それをやるには業務用のルータが必要です。実家に置いてあるルータなんてメンテしてられません。そんな皆様に。。。

## 提案

[meshover](https://github.com/gotti/meshover) を開発しています。これをセットアップするだけで`10.1.2.3`のようなアドレスがパソコンにランダムに割り当てられ、オーバーレイネットワークに参加されます。他のサーバにも同じようにインストールすれば、そのIPアドレスを使ってパソコン同士で通信できます。パソコンを置いてる場所やルータの設定なんて関係ありません。

この機能だけ見てると[tailscale](https://tailscale.com/)でいいじゃん!ってなると思います。meshoverはtailscaleと違い自宅サーバクラスタに特化していてパソコンがBGPを喋ります。そのためパソコンだけでなく、そのパソコンの上で動いているVMやコンテナも1つのネットワークに接続されます。

## 要求

- グローバルIPv6アドレスが割り当てられていること。
  - 普通のNTT環境でこれは満たされてます。
- グローバルIPを持ったサーバが1台はあること
  - コントローラを動かすために必要です。

## 手順

ここまでmeshoverについて利便性を推してきましたが使うのはある程度面倒です。トラシューにもある程度の知識が必要です。やりたい方は具体的な手順を載せてる[GitHubに上げたGetting Started](https://github.com/gotti/meshover/blob/main/docs/getting-started.md)を見てもらうとして、手順を簡単に説明します。

- meshoverコントローラを動かす。
  - 各パソコンで動いているエージェントのVPN公開鍵やIPアドレスなどを保存して全てのエージェントに共有します。
  - 自宅サーバのネットワークから独立させた方がいいかもしれません。私はGCPで動かしてます。
- meshoverエージェントをインストールする。
  - いくつか必要なツールがありますがバイナリを立ち上げるだけです。
  - ネットワークに参加させるパソコン全てで実施してください。

おわりです。簡単ですね(？)

ちなみに、これしかやらないならtailscaleの方が対応端末多いしIPアドレスの要件が緩いしでtailscaleの方が良いです。

## 追加手順

VMやコンテナを動かしましょう。ところでVMやコンテナの管理に何を使っていますか？派閥はあると思いますがmeshoverはKubernetesとの連携を目的として開発したのでKubernetesを使いましょう。

### ところでKubernetesってどう通信してるの

Kubernetes自身はほぼネットワークに関わらずCNI(Container Network Interface)に任せています。CNIにはいろいろな実装があり、CalicoとかFlannelとかCiliumが有名ではないでしょうか。CNIには次の最低限の2つの機能を持っているのがほとんどで、IPアドレスの割り当て、オーバーレイネットワーク構築などによるコンテナ間通信の確保です。さきほど挙げたCNIは全てBGPやVXLANを使って自身でオーバーレイネットワークを構築するなどの機能を持っています[^cni]。この機能についてはmeshoverが責任を持つためCNIには前者だけやってもらえばいいです。

前者のIPアドレス管理だけをやるCNIとして[Cilium(Native Routingモード)](https://docs.cilium.io/en/stable/concepts/networking/routing/#id2)や[Coil](https://github.com/cybozu-go/coil)などがあります[^calico]。Ciliumにはおまけ機能が充実していたりと楽しいので私はCiliumを選択しました。

これも詳細な手順はmeshoverのgithubに上げているので雑に説明します。

- kubernetesをkubeadmなどでインストールします。
  - controlplane ipなどはmeshover側のipにすること
- ciliumのtunnelを切ったりmasqueradeの設定をやります。
- loadbalancerを入れます。
  - 私は[PureLB](https://purelb.gitlab.io/docs/)を入れました。

## あとがき

このツールは

- Kubernetesを使っていて
- サーバが複数あって
- 複数リージョンに分散していて
- VMやコンテナをガンガン使っている

方向けのツールです。それ以外の方は以下のツールが自宅鯖やるのに楽なのでオススメです。

- [Cloudflare Tunnel](https://developers.cloudflare.com/cloudflare-one/connections/connect-apps/)
  - 内側からcloudflareにトンネルを掘ってサーバを世界に公開します。ルータの設定がいりません、すごいぜ。
- [tailscale](https://tailscale.com/)
  - 散々例に出したツールです。パソコン間で自動でVPNを張って学校でも自宅鯖にアクセスできます。ぼくも愛用しています。

## 注釈

[^cni]: meshoverはKubernetesクラスタ外とも通信したりIPv6でVPNを張ったりするため、これらでは代替できません。
[^calico]: calicoにも似たような機能があるらしいんですが未調査です。
