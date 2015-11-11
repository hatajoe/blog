+++
date = "2015-11-12T01:01:48+09:00"
description = "クロージャの話"
title = "gocql使って関係ないとこでハマった"

+++

Goでgojiみたいな薄いフレームワーク使ってWebサーバー書いてるとORMをどうするかって話があると思う。  
この記事ではCassandraを例に出すんだけど、CassandraにはBatchといって複数の操作をまとめて行える命令がある。
GoでCassandraと言えば[gocql/gocql](https://github.com/gocql/gocql)がすぐに見つかると思う。gocqlでBatchを実行するには以下のAPIを使う。  

> https://godoc.org/github.com/gocql/gocql#Session.ExecuteBatch
```
func (s *Session) ExecuteBatch(batch *Batch) error
```

引数で渡すBatchには以下のAPIを使ってステートメントとパラメータをバインドする。

> https://godoc.org/github.com/gocql/gocql#Batch.Bind
```
func (b *Batch) Bind(stmt string, bind func(q *QueryInfo) ([]interface{}, error))
```

ここで注目なのが第二引数で指定する `bind func(q *QueryInfo) ([]interface{}, error)` 。  
このbindは、Bindメソッドに渡すstmtに対応するパラメータを[]interface{}で返す関数ということ。  
そこで以下のようなBatchラッパーを書いた。

```
type BatchHandler func() (string, func(q *gocql.QueryInfo) ([]interface{}, error))
func (m *Cassandra) Batch(handlers ...BatchHandler) error {
	batch := gocql.NewBatch(gocql.UnloggedBatch)
	for _, handler := range handlers {
		// ステートメントとパラメータを取得してbatchにバインド
		stmt, bind := handler()
		batch.Bind(stmt, bind)
	}
	return m.session().ExecuteBatch(batch)
}
```

BatchHandlerは、Bindメソッドに渡すstmtとbindを戻り値で返すような関数型。  
このBatchHandlerをスライスにしてBatchに渡すことで、複数の命令をまとめて実行することが出来る。（果たして伝わるかどうか・・・）  
この設計に対して、ORMは以下のようにした。

```
type UserSlice []*User
func (s UserSlice) Save() (h []BatchHandler) {
	h = make([]BatchHandler, l)
	// UserSliceの各要素を保存
	for i, _ := range s {
		// BatchHandlerのスライスを作る
		h[i] = func() (string, func(q *gocql.QueryInfo) ([]interface{}, error)) {
			query := fmt.Sprintf("INSERT INTO %s (key, value) VALUES (?, ?)", users)
			// ここがクロージャになってる
			return query, func(q *gocql.QueryInfo) ([]interface{}, error) {
				b, err := json.Marshal(s[i])
				if err != nil {
					return nil, err
				}
				values := make([]interface{}, 2)
				values[0] = s[i].ID
				values[1] = b
				return values, nil
			}
		}
	}
	return
}
```

以下のようにして使う。

```
func main() {
	us := UserSlice{{ID:1, Name:"hatajoe"}}
	es := EntrySlice{{ID: 1, Title:"hoge"},{ID: 2, Title:"fuga"}}
	h := us.Save()
	h = append(h, es.Save())

	if err := cassandra.Batch(h); err != nil {
		panic(err)
	}
}
```

これにてUserSliceとEntrySliceは無事にCassandraに保存されたかに見えた。  
しかし、実際には、UserSliceは正しく保存されるが、EntrySliceはID 2のデータしか保存されない。なぜか。  

原因はクロージャの変数バインドにあった。  
Saveメソッドは、スライスをクロージャにバインドして、クロージャ実行時にs[i]からデータを取り出している。  
上記の例でEntrySliceの場合だと、このiはクロージャ実行時には常に1になる。  
それは、恐らくクロージャがiの参照を保存しているため。そら要素数２のスライスをforで回したら抜けた後のiは1になるわな。  

ということで、以下のようにループ中のiの状態を無名関数の即時実行で閉じ込めるのが正解だった。  

```
func (s UserSlice) Save() (h []cassandra.BatchHandler) {
	h = make([]BatchHandler, l)
	for i, _ := range s {
		// iを無名関数の引数idxとして閉じ込める
		h[i] = func(idx int) func() (string, func(q *gocql.QueryInfo) ([]interface{}, error)) {
			return func() (string, func(q *gocql.QueryInfo) ([]interface{}, error)) {
				query := fmt.Sprintf("INSERT INTO %s (key, value) VALUES (?, ?)", users)
				return query, func(q *gocql.QueryInfo) ([]interface{}, error) {
					b, err := json.Marshal(s[idx])
					if err != nil {
						return nil, err
					}
					values := make([]interface{}, 2)
					values[0] = s[idx].ID
					values[1] = b
					return values, nil
				}
			}
		}(i)
	}
	return
}
```

関数型というか普段からクロージャ使ってる人からしたらあるあるな内容なのかも。僕は結構ハマってツラかった。  
・・・やっぱこの話伝わらない気がする。。

