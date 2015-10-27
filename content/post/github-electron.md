+++
date = "2015-10-28T01:06:00+09:00"
description = "electron便利"
title = "electronでGitHubアプリを作ってみた"

+++

ソースはこちら。

* [hatajoe/github-electron](https://github.com/hatajoe/github-electron)
<br /><br />

使い方はまず、プロジェクトをクローンしてから

```
$ npm install
```

次に、以下のような `~/.github-electron` ファイルを作成。

```
{
    "token": "GitHubのPersonal Access Token"
}
```

Personal Access Tokenは `notifications` にチェックして作成。  
そしてリポジトリ内の `GitHub-darwin-x64/GitHub.app` を実行することで使用可能。

ちょっと前に以下の記事を見ていけるかもってことで作ってみた。

* [ElectronでChatworkをデスクトップアプリ化 (Webview + badge) - Qiita](http://qiita.com/mottox2/items/7a1373f23ba02245d0e0)
<br /><br />

仕事でGitHubとSlackを連携してたりするんだけど、例えば、プルリクエストをレビューする時なんかにGitHubでメンション飛ばしてもSlackに通知はされるがSlackのメンションにはならなくて、結局なかなか気づくことが出来なくて困ってた。  
このアプリを使うと、Watchしてるリポジトリの更新や自分に対するメンションがあった場合にバッヂとして通知を受けることが出来る。

<img src="/blog/images/github-electron-1.png" class="image" alt="sample">

やっていることは単純で、3秒に1回GitHub APIを叩いて自分への通知があればバッヂを付けるって感じ。
今のところ、webviewの戻るや進むに対応してなくて使いづらい。

electronは簡単にこういったことが出来て便利だなーと思った。  
良かったら使ってみて下さい。
