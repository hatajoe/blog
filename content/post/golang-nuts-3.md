+++
date = "2015-10-08T01:34:11+09:00"
description = "今のところGoでGoのプラグインは作れない"
title = "c-sharedでruntime error 結末"

+++

[前回](http://hatajoe.github.io/blog/post/golang-nuts-1/)までで、Darwin/arm64環境ではGoでGo shared libraryはロード出来ないことがわかった。  
そのため、ちょっと環境を替えてLinux/amd64環境にて検証をすすめることにした。  

そもそもやりたいこととしては汎用なプロキシサーバーを書くことで、機能自体はplaggableに拡張出来るものを目指したかった。  
そうすることで、コアな部分はシンプルに保ちつつ機能拡張はコアを侵食することなく行えるだろうと。  

この文脈で言うところのプロキシとは、HTTPリクエストをフィルターし、然るべきタイミングでフックされ処理が行われるものを想定していた。  
そのため、http.HandlerFuncをプロキシコアとプラグインとの間のインターフェースに採用すると良さそうだった。  
処理をhttp.HandlerFuncでラップするテクニックとしては以下がとてもわかりやすく参考になると思う。  

* [The http.HandlerFunc wrapper technique in #golang — Medium](https://medium.com/@matryer/the-http-handlerfunc-wrapper-technique-in-golang-c60bf76e6124)<br /><br />

そんなわけで、http.HandlerFuncをGoとGo shared libraryでやり取り出来ると良さそうだった。
そうして以下のコードを書いた。

```
package main

import (
	"C"
	"log"
	"net/http"
	"time"
)

//export Profile
func Profile(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b := time.Now()
		defer log.Printf("elapsed: %v\n", time.Now().Sub(b))
		fn(w, r)
	}
}

func init() {
	log.Println("loaded")
}

func main() {
}
```

すると、コンパイル時に以下のエラー。

```
 Go type not supported in export: http.HandlerFunc
```

shared libraryではhttp.HandlerFunc型のexportをサポートしていないらしい。
他に方法は無いかgolang-nutsで[聞いてみた](https://groups.google.com/forum/#!topic/golang-nuts/cHiFHNyKXRw)ところ、  
どうやら現在のところ、Go shared libraryをGoで読み込むこと自体がunsupportedらしいことがわかった。  

まぁでもにわかには納得がいかなくて、ポインタを渡してデリファレンスすればいけんじゃね？と思ったが、  
実際にやってみると、渡す前と渡した後でアドレスが違っていたりしてよくわからなかった。（ランタイムを詳しく読めばわかるのかな・・・？）  

結論としては、取り敢えず今のところGoでGoのプラグインを書くことは出来そうにないようだ。残念。

P.S  
この内容、来週頭の第1回関西golang勉強会で話しようと思う。ネタバレになっちゃうけど。

