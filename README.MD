# search-my-wiser

全文検索エンジンのサンプル。
wikipediaの日本語記事を登録して検索する。


### To start
```
./script/init.sh
```

# Required
- go
- mysql

### 登録
wikiのダンプから5件をインデクシングする。
```
cd ./wiser
./wiser -w ../files/jawiki-latest-pages-articles.xml -m 5
```

### 検索
"言語"という検索クエリを実行する。
```
cd ./wiser
wiser -q 言語

# 検索結果
docId: 5, title: 言語, score: 1.278822
docId: 3, title: アンパサンド, score: 0.043106
```
