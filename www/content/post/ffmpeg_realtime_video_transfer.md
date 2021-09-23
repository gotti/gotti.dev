---
title: "ffmpegだけでデスクトップキャプチャを他のPCに転送"
date: "2021-09-23"
tags: ["ffmpeg", "linux"]
---

# Motivation

discordとかやる端末と開発機/ゲーム機が違うとき、画面共有できませんよね？ネットワーク経由で転送して開発機の動画を見せましょう

# 前提条件

- 両方の端末にffmpegがインストールされていること
- デスクトップキャプチャ側はx11が動いていること

# こまんど

デスクトップキャプチャ&送信元
```bash
ffmpeg -f x11grab -video_size 3840_2160 -i :0.0+0,0 -vf "scale=1920:-1" -framerate 30 -preset superfast -tune zerolatency -f mpegts 'srt://${target_host}:30000'
```

画面を表示する方
```bash
ffplay -analyzeduration 1 -fflags -nobuffer -probesize 32 -i "srt://0.0.0.0:30000?mode=listener"
```

おわり

# 参考

https://stackoverflow.com/questions/42953616/very-low-latency-streaminig-with-ffmpeg-using-a-webcam
