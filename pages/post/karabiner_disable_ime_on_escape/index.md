---
title: "Karabiner-Elementsで、Vimのノーマルモードに戻ると同時にIMEをオフにする"
date: "2025-04-07"
tags: ["MacOS", "Vim", ""Karabiner-Elements", "AquaSKK", "SKK"]
---

こんにちは、ごっちと言います。
最近省電力でパフォーマンスの良いパソコンを探していて、M4 MacBook Airを購入し、破産しました。
今までArch Linuxのデスクトップ環境にいたので慣れていませんが、Homebrewは便利ですね。
AquaSKKとKarabiner-Elementsを使って、LinuxでやっていたSKKの設定を再現できたのでメモしておきます。

## SKKとVimの組み合わせ

Vimでノーマルモードに戻るにはEscapeキーを押します。ノーマルモードではIMEが無効になっている必要があります。
一方、かな英数でIMEを切り替えていると、ノーマルモードでIMEオンの状態になり、何も操作できない状態になりがちです。
そのため、ノーマルモードへの遷移とIMEオフを(ほぼ)同時にやりたくなります。
LinuxのFcitxではEscapeをIMEオフに割り当てることで実現していました。MacOSではKarabiner-Elementsを使う必要がありそうです。

## Karabiner-Elementsの設定

以下の設定を読み込ませてください。
お好みでEscapeキーを押しやすいCaps Lockとかに割り当てても良いでしょう。

```json
{
    "description": "日本語入力モードでEscapeを英数に変更する",
    "manipulators": [
        {
            "conditions": [
                {
                    "input_sources": [{ "language": "ja" }],
                    "type": "input_source_if"
                }
            ],
            "from": {
                "key_code": "escape",
                "modifiers": { "optional": ["any"] }
            },
            "to": [{ "key_code": "japanese_eisuu" }],
            "type": "basic"
        }
    ]
}
```

これで、日本語入力モードで一度目にEscapeを押すと英数に切り替わります。
二度目に押すと普通のEscapeになり、ノーマルモードに戻ります。
VimでEscapeを連打すればノーマルモードに戻るようにできました。ハッピー!
