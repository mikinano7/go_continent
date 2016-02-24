# go_continent
Go言語で作成したシンプルな Windows OS 向けツイッタークライアントです。  
投稿パネルのみ表示、 Ctrl+Enter で投稿します。

# Usage
gitからソースコードをcloneします。
`> git clone https://github.com/mikinano7/go_continent.git`

go.envにツイッターの認証情報を書きます。（そのうちPINコード認証にします）
```
TWITTER_CONSUMER_KEY=1iyt7kOztWsYPQBVnF6APm910
TWITTER_CONSUMER_SECRET=I6SiognJI766oodNrqM3lorIDHg2H09UZeIukvGOC1U7Re6r5m
TWITTER_ACCESS_TOKEN=1707378048-Yt9BDenVGA1W7bwi67sPdKLMJx0OcbFvCrQQoqV
TWITTER_ACCESS_SECRET=EAQnfqDTn3Hmf8muaG4vr14YpYfqFQvaTjXC23gZAQe8f
```

goコンパイラでビルドします。（そのうち実行ファイルを公開します）
`> go build`

出てきた .exe ファイルを実行します。
`> ./go_continent.exe`