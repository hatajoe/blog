+++
date = "2015-09-30T01:30:23+09:00"
description = "Darwinは無理。"
title = "c-sharedでruntime error その後"

+++

前段は[こちら](http://hatajoe.github.io/blog/post/golang-nuts-1/)。  
簡単に言うと、darwin/amd64環境においてbuildmode=c-sharedでビルドしたGo Shared LibraryをGoで使おうとするとruntime errorになる。  

何故かわからなくて、golang-nutsで聞いてみたところ回答を頂いたので共有する。  
まぁ実際に見てもらった方が話は早くて。 

[How can I do to fix `runtime/cgo: could not obtain pthread_keys' on darwin/amd64 - Google グループ](https://groups.google.com/forum/#!topic/golang-nuts/Vy8r05reLyw)

要約すると、Darwin環境ではGoとGo Shared Libraryで同じpthread keyをアロケートすることが問題であるということらしい。
また、静的リンクについては本当に恐ろしいハックによって成り立っているらしく・・・恐らく、 `runtime/cgo/gcc_darwin_386.c` のマジックナンバーをインラインアセンブラでオフセットしている部分のことなんだろうと思うけど、僕には本当に意味がわからない。

取り敢えず、darwin/amd64でビルドしたGo Shared LibraryをGoで使うことは出来ないということが明らかになったのであった。  
一応、回答の全文を載せてこの件は終わりとする。

>
The problem is described in a comment in runtime/cgo/gcc_darwin_386.c.   
When linking statically on Darwin, the linker selects an offset for  
the single thread local variable used by the runtime.  The code in  
gcc_darwin_{386,amd64} grabs pthread keys to ensure that that offset  
is indeed allocated for Go.  This is a truly horrible hack, but it  
works well enough when linking statically.  
>
From your description, it sounds like you are trying to open a Go  
shared library in a Go program.  That is not going to work as both the  
Go program and the Go shared library are going to try to allocate the  
same pthread key.  That is why you are getting that error.  
>
The fix is going to be to have the Go linker emit whatever relocations  
are required to handle thread local variables in Darwin shared  
libraries.  Perhaps it does that already--I don't know.  Then the code  
in runtime/cgo need not run in this case.  
>
By the way, I want to be clear that -buildmode=c-shared is not  
intended to build a plugin library for a Go program.  Even if you fix  
this problem there may be future problems.  Those problems are bugs  
that should be fixed, but I want to caution you that you need to be  
prepared to run into problems.  
>
Ian 

thank you for your reply Ian! 

ということで、GoでGo Shared Libraryを使ったプラグイン機構を作る場合は、今のところdarwinを外さないといけないですね。 
