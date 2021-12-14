---
title: "調布祭プラレールのラズパイのお話"
date: "2020-12-23"
tags: ["advent", "raspi"]
---

# 調布祭プラレールのラズパイのお話
この記事は[UEC koken Advent Calendar 2020](https://adventar.org/calendars/5692)の23日目の記事となります。

調布祭プラレール企画の関連記事は以下の通りです。

* 20: [概要、自動制御部分](https://qiita.com/takoyaki3/items/4a9984c4023ceae4791f)
* 21: [WebUI、ビデオ配信](https://stea.hatenablog.com/entry/2020/12/21/043801)
* 22: [WebAPIに関して](https://pfpfdev.hatenablog.com/entry/20201222/1608570444)
* 23: ラズパイ部分 ←この記事
* 24: ESP・列車回路
* 25: サーボモーターや配線周り

# 概要
初めまして、[ごっち](https://twitter.com/intent/user?user_id=3721840992)です。最近工研に入ったのでslackのチャンネルにかたっぱしから入っていたらプラレール企画のチャンネルに入ってしまい、気がついたらタスクが割り当てられてデスマ真っ只中にいました。(こうは書きましたがプラレール企画楽しかったです、ありがとうございます。)

さて、調布祭お疲れ様でした! 鉄研と合同で工研では11/21~11/23の調布祭3日間プラレールの遠隔操作と自動制御という企画をやっていました。今回の記事ではプラレール記事のうちラズパイ関連について書きたいと思います。

<blockquote class="twitter-tweet"><p lang="ja" dir="ltr">オタク!電車を脱線させて遊べるぞ!<br>工研🛠／鉄研🚎コラボ開催中！今ならプラレールの遠隔操作ができるよ！今すぐアクセス！ <a href="https://t.co/m5mFef1A7F">https://t.co/m5mFef1A7F</a> <a href="https://twitter.com/hashtag/%E8%AA%BF%E5%B8%83%E7%A5%AD?src=hash&amp;ref_src=twsrc%5Etfw">#調布祭</a> <a href="https://t.co/p77wEuqCH5">pic.twitter.com/p77wEuqCH5</a></p>&mdash; ゴッチャ (@_nil_a_) <a href="https://twitter.com/_nil_a_/status/1330719524610457601?ref_src=twsrc%5Etfw">November 23, 2020</a></blockquote> <script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>

これは準備日と1日目のデスマが終わりclusterで遊んでたときの様子

<blockquote class="twitter-tweet"><p lang="ja" dir="ltr">バーチャルデスマしに来た <a href="https://t.co/M5eyJUWdc5">pic.twitter.com/M5eyJUWdc5</a></p>&mdash; ゴッチャ (@_nil_a_) <a href="https://twitter.com/_nil_a_/status/1330735630779760640?ref_src=twsrc%5Etfw">November 23, 2020</a></blockquote> <script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>

# プラレールの全体像について
プラレールの遠隔操作企画ではポイントなどを遠隔操作することができます。1編成だけは速度の制御が可能(逆走も)で車載カメラによる映像配信も行なっていました。

私が担当したのはポイント遠隔操作用のラズパイと車載カメラのラズパイです。

# 遠隔操作のラズパイの話
遠隔操作部分はラズパイがWebAPIとのソケット通信を確立し、降ってくる命令に沿ってポイントを制御します。中心である調布駅にラズパイを設置しているので調布駅のサーボモーターはラズパイのGPIOから制御、ちょっと離れた駅(具体的には北野と笹塚)にはサーボモーターを繋いだESP32を設置しています。そのESP32のAPI(http)をラズパイから叩いています。また速度制御可能なスカイライナーにもESPが積まれており、それに対しても速度変更の命令を送信しています。

こんな感じの形式でw.deviceNum(ESP側のピン番号と対応)とb(onまたはoff)を設定する感じです。

```go
req,err := http.NewRequest(http.MethodGet,"http://"+w.address+"/"+b+w.deviceNum,nil)
```

pofさんのWebAPIが早々に完成していたのでコーディングはかなり簡単に進みました。ありがとうございます。(デバッグに1,2日溶かしましたが...) Go言語を使っていますがGPIOはPythonの方が扱いやすいのでPythonも併用しています。

# 配信用ラズパイの話
プラレールの遠隔操作で最も問題となるのはカメラ映像の遅延です。今回、駅ごとに様々なアングルから撮影できるようにカメラを多数設置し配信していました。(詳しくはすてあさんの記事を参照してください。)これをYoutube liveなどを介して配信しようとすると、遅延最小の設定でも3秒ほどの遅延が発生するようです。そこでWebRTCを用いてYoutubeなどの中継なしでユーザーに映像を直接配信しています。

定点カメラはWebRTC(SkyWayというサービスを使っています。)で配信していたのですが、車載カメラにはSkyWayが使えません。RaspberryPi zeroでは複数台にWebRTCで配信する負荷に耐えられないですし、そもそもCPUがARMv6であるためSkyWay Gatewayは起動すらできません。(RaspberryPi zero以外はARMv8なので動きます。負荷は知りません。)

今回、ラズパイに最適化されたWebRTCツールであるMomo(https://github.com/shiguredo/momo)を用いました。これはRaspberryPi zeroでも720p20fpsで配信可能な優れものです。なんとJetson nanoを使えば4k30fpsで配信できるそうです。内部では血の滲むような最適化がされているはずです...。さらにRaspberryPi zeroでは複数台へ配信する負荷に耐えられないので一度配信用パソコンだけに送信して、そのデスクトップキャプチャをSkyWayで配信するごり押しで解決しています。もっとスマートにやりたかった...(遺言)。

RaspberryPi zeroは最初は電池がもたないから駄目だろうと考えていたのですが、2日目にモバブを買って載せてみたところ意外と動きびっくりしていました。1本の18650でおそらく1日の展示には余裕で耐えられる気がします。なお配信のために購入したモバブは体積と重量を減らすため一瞬でバラバラにされていました。(昇圧回路と18650がセットで売られておりオトク!)

<blockquote class="twitter-tweet"><p lang="ja" dir="ltr">これはモバブだったもの、ショートが怖くて誰も持って帰りたがらない <a href="https://t.co/oYKYc2vhMU">pic.twitter.com/oYKYc2vhMU</a></p>&mdash; ゴッチャ (@_nil_a_) <a href="https://twitter.com/_nil_a_/status/1330790512719060994?ref_src=twsrc%5Etfw">November 23, 2020</a></blockquote> <script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script> 
なお車載カメラの配信は1日目はなし(遠隔操作系でデスマしていたので)、2日目はmomo+Youtube live、3日目はffmpeg+Youtube liveからmomo+SkyWayと次第に最適化？されています。なおそのあとラズパイが落ちて壊れました<br>
<blockquote class="twitter-tweet"><p lang="ja" dir="ltr">状況です、こいつは遠隔で速度を弄れて逆走などもできます、カメラ配信は現在準備中あとちょっと <a href="https://t.co/cq7hJWzULn">pic.twitter.com/cq7hJWzULn</a></p>&mdash; ゴッチャ (@_nil_a_) <a href="https://twitter.com/_nil_a_/status/1330362889262096386?ref_src=twsrc%5Etfw">November 22, 2020</a></blockquote> <script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>
<blockquote class="twitter-tweet"><p lang="ja" dir="ltr">車両映像を720p15fpsでwebrtc配信していたのですが，raspiがこわれました.....</p>&mdash; ゴッチャ (@_nil_a_) <a href="https://twitter.com/_nil_a_/status/1330730748366602240?ref_src=twsrc%5Etfw">November 23, 2020</a></blockquote> <script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script> 
<blockquote class="twitter-tweet"><p lang="ja" dir="ltr">寝てる間に車両前方の映像配信が最適化されてその後壊れてた</p>&mdash; る (@ruu_uec) <a href="https://twitter.com/ruu_uec/status/1330730533639262210?ref_src=twsrc%5Etfw">November 23, 2020</a></blockquote> <script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script> 

# 反省点や改善

* 早めに開発と統合テストをやりましょう。調布祭1日目にプラレールの配信ではなくデスマを配信することになります。(ごめんなさい....。)
* RaspberryPi zeroの電池もちが予想以上だったので実験してみないとわからないなと思いました。
* 360度動画にしたいねという意見が出ていました。ぼくもやりたい。(ツイート見つけられず)
* チャット機能など

<blockquote class="twitter-tweet"><p lang="ja" dir="ltr">ユーザー同士で会話するためのチャット機能とかあると良かったね，みたいな話は出ています</p>&mdash; ゴッチャ (@_nil_a_) <a href="https://twitter.com/_nil_a_/status/1330762988492521474?ref_src=twsrc%5Etfw">November 23, 2020</a></blockquote> <script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script> 
* WebRTC周りは触ったことがないのでSkyWayやmomoにおんぶにだっこでしたが、自分で書いてみたいですね。視聴者のパソコンを配信インフラに組み込んで中継サーバとし配信者の負荷を軽減する配信ネットワークの計画が進行中なのでご期待ください(本当に完成するんですか？)
<blockquote class="twitter-tweet"><p lang="ja" dir="ltr">FHD60fpsだけなら機材を選べば可能、同接1kは無限の予算でツリー状に中継させれば可能なので、無限の予算さえあれば遅延10ms以外いけそう</p>&mdash; ゴッチャ (@_nil_a_) <a href="https://twitter.com/_nil_a_/status/1330857851095379969?ref_src=twsrc%5Etfw">November 23, 2020</a></blockquote> <script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script> 
<blockquote class="twitter-tweet"><p lang="ja" dir="ltr">カップルは手を繋ぐけど僕はWebRTCでストリーム繋ぐのでデータ量で圧勝できる</p>&mdash; る (@ruu_uec) <a href="https://twitter.com/ruu_uec/status/1331959289196417026?ref_src=twsrc%5Etfw">November 26, 2020</a></blockquote> <script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script> 
<blockquote class="twitter-tweet"><p lang="ja" dir="ltr">WebRTC完全に理解した</p>&mdash; る (@ruu_uec) <a href="https://twitter.com/ruu_uec/status/1331950806766034945?ref_src=twsrc%5Etfw">November 26, 2020</a></blockquote> <script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script> 
