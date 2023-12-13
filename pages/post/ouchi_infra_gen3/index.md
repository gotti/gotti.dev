---
title: "新おうちインフラ(第三世代)"
date: "2020-12-23"
tags: ["Liunx", "Infra", "k8s"]
---

# 既存インフラ(第一、第二世代)の問題

既存インフラは第一世代の`camphor`と第二世代の`olive`からなっていました。

`camphor`のOSはArch Linuxで、NFSと録画鯖が昔ながらの構成で動いています。`olive`のOSはproxmoxで、仮想化された上でルータVMなどが走っています。

しばらく使っているとVLAN、ルータVM、仮想マシン、ストレージなどの管理が辛くなってしまいました。しかもproxmoxの上でPPPoEのルータが動いているので大きな単一障害店になっています。

このままのインフラだと壊れたときに再構築するのは大変面倒になるでしょう。そこでおうちインフラの全てをKubernetesを使って管理しようということになりました。Kubernetesで管理することで、kubesprayでクラスタを組んでGitHubからマニフェストを流すだけで同じ状態をすぐ再構成できるようにすることを目的としています。

物理マシンは以下の構成となりました。予算不足でノードは2台でHA構成ではないので単一障害点の問題はまだ残っています、許して。ノード間は10Gで繋いでいます。

- apricot (master)
  - CPU: 6C12T(Ivy Bridge E)
  - RAM: 32GB
  - SSD(System): 512GB SATA
  - SSD(Data): 1TB NVMe

- olive
  - CPU: 8C16T(Zen3)
  - RAM: 32GB
  - SSD(System): 512GB NVMe
  - SSD(Data): 1TB NVMe

# 要求

- 高速なローカルストレージ
  - 冗長性は担保しないが高速なもの。VMのephemeralなどに
- そこそこ高速な分散ストレージ
  - 冗長化されているが普通に使う分には問題ない速度を出せる
- LoadBalancerに割り当てるグローバルIPやローカルIPを管理できる
- VMを管理でき、必要に応じてグローバルIPを割り当てられる
- ルータを必要なだけ増やせる

# 使用するもの

Kubernetesを使えばこれを達成できます。OpenShiftなどを使うとより簡単にデプロイできると思いますが、趣味は逆張りしてなんぼなので全部自分で管理することにしました。Kubernetesに導入するコンポーネントは以下のようになりました。

- TopoLVM
ホストからLVMを切り出してPVを作る、rookの他高速なローカルストレージのために導入
- Rook
topolvmを束ねて分散ストレージであるceph
- KubeVirt && CDI && Cluster Network Addon Operator
VMを管理してPVCをVMのディスクとしてアタッチ、ホストのOpen vSwitchをVMにアタッチ
- FluxCD
Continuous Delivery
- MetalLB
BGPでグローバルIPとプライベートIPを広報、LoadBalancerのannotationで切り替えられる。
- Traefik
Webで何かしら公開するコンテナのために

# 構成

低速なストレージはとりあえず既存の`camphor`のNFSサーバに頼ることにしました。

クラスタはkubesprayを用いてベアメタルに構築しています。


