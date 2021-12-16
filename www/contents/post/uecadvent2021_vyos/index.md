---
title: "VyOSでソフトウェアルータ"
date: "2021-12-13"
tags: ["infra", "linux"]
---

## はじめに

これはUEC 2 Advent Calendar 2021の12/14の記事です。前の方は12/13にpizacさんが投稿している[外付けSSDにUbuntuをインストールした話](https://pizac.hatenablog.com/entry/2021/12/13/145722)です。

USBメモリにLinux入れとくの便利ですよね。わたしの推しのArch Linuxは[wikiにUSBメモリにインストールする方法が載ってたりします](https://wiki.archlinux.jp/index.php/%E3%83%AA%E3%83%A0%E3%83%BC%E3%83%90%E3%83%96%E3%83%AB%E3%83%A1%E3%83%87%E3%82%A3%E3%82%A2%E3%81%AB_Arch_Linux_%E3%82%92%E3%82%A4%E3%83%B3%E3%82%B9%E3%83%88%E3%83%BC%E3%83%AB)。わたしもPCが死んだときの復旧に重宝しています。

## 自己紹介

uec20のごっちです。12/13の23:30ごろに急に記事書きてえなとなり12/14のadvent calendar枠を確保して頑張って書いています。

ところでみなさんも一度は家庭用ルータではできない複雑なルーティングをしたりVPNを張りまくったりしたくなりますよね。

ルータ用OSとして有名なものとして、家庭用ルータに焼けるOpenWrtなどもあります。

しかし組み込み機器向けなOSなので(わたしは)かなりつらくなって嫌になりました。そこでVyOSです。

## 特徴

VyOSはその名の通りOSで、x86_64のパソコンにインストールして使います。debianベースらしいです。

### 利点

- 家庭用ルータのCPUとは比べものにならないx86_64の性能で殴れます
- VPNにいろいろ対応しています。WireGuardやOpenVPN、L2TPなどなど
- ansibleから叩くことでデプロイを自動化できます
- コマンド書式が綺麗

### 欠点

- x86_64機用意するのめんどい
- 用途によってはtoo muchなことも多い

## やっていく

ある用途にOpenVPNのサーバが欲しくなったのと強いソフトウェアルータを触ってみたかったなどがありVyOSを使っていくことにしました。

インストールですが、VyOSはLTSが有償という謎なライセンス体系です。snapshotとnightly buildsは無償です。snapshotをダウンロードしましょう。

https://vyos.net/get/

インストール手順は解説しません。公式ドキュメントを見てください。

https://docs.vyos.io/en/equuleus/installation/index.html

これはQuick Start Guideです。コマンドが良い感じなのがわかりますね。コマンドで設定するタイプのルータを触ったことがある人は雰囲気で触れるはず。

https://docs.vyos.io/en/equuleus/quick-start.html

ユーザーを追加してファイアウォールなどを設定したらansibleもやっていきましょう。

### memo

ところで私はNTTPC Indigoの350円/月なVPSサーバにVyOSをインストールしました。このVPSはUbuntuとCentOSとWindowsしか動かせないのですが、黒魔術を使うとVyOSを起動できます。

https://scrapbox.io/otofune/NTTPC_Indigo_%E3%81%A7_OS_%E3%82%92%E5%85%A5%E3%82%8C%E6%9B%BF%E3%81%88%E3%82%8B

ubuntu20だからかsnapshotで不適切なのを選んだからか、いくつかの組み合せではカーネルパニックが起きて失敗しました。これをやる際はいろいろな組み合わせを試してみてください。ちなみにVyOSはqcow2イメージも配っています、優しい。

## ansible

また実際のコードは置きません。ドキュメントはこれです。

https://docs.vyos.io/en/equuleus/automation/vyos-ansible.html

コマンド列挙するだけなんですが自動化できるのは良いことです。私はGitHubに適用するコマンドを列挙してGitHub Actionsでansibleを叩いてデプロイするようにしています。

