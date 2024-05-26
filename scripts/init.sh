#!/bin/bash

# カレントディレクトリを取得
CURRENT=$(cd $(dirname $0); pwd)

# wikiのダンプファイルをダウンロードするスクリプト
# ダンプはサイズ大きいから気をつけてね
("${CURRENT}/dump_wiki.sh")

# dbを起動する
("${CURRENT}/init_db.sh")

# ビルド
(cd "${CURRENT}/../wiser" && go build -o wiser)
echo "Build completed successfully."

