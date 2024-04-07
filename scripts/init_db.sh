#!/bin/bash

# カレントディレクトリを取得
CURRENT=$(cd $(dirname $0); pwd)
# 一つ上のディレクトリを取得
CURRENT_PATH=$(dirname "$CURRENT")


source "${CURRENT_PATH}/wiser/.env"

# MySQLの接続情報
MYSQL_HOST="$DB_HOST"
MYSQL_PORT="$DB_PORT"
MYSQL_USER="$DB_USER"
MYSQL_PASSWORD="$DB_PASSWORD"
DATABASE_NAME="$DB_NAME"

# MySQLに接続してseedデータを投入する関数
function seed_database() {
    local sql_file="$1"

    # MySQLに接続してSQLファイルを実行
    mysql -h "$MYSQL_HOST" -P "$MYSQL_PORT" -u "$MYSQL_USER" -p"$MYSQL_PASSWORD" "$DATABASE_NAME" < "$sql_file"
    #mysql -h "$MYSQL_HOST" -P "$MYSQL_PORT" -u "$MYSQL_USER" "$DATABASE_NAME" < "$sql_file"
}

# MySQLを起動
mysql.server start "--defaults-extra-file=${CURRENT_PATH}/db/my.cnf"

# MySQLが起動するまで待つ
while true; do
    mysql_status=$(mysql.server status | grep -c "SUCCESS")
    if [[ "$mysql_status" -eq 1 ]]; then
        echo "MySQL started successfully."
        break
    else
        sleep 1
    fi
done

# MySQLのパスワードを入力して環境変数に保存
#read -s -p "Enter MySQL password: " mysql_password
#export MYSQL_PASSWORD="$mysql_password"

mysql -h "$MYSQL_HOST" -P "$MYSQL_PORT" -u "$MYSQL_USER" -p"$MYSQL_PASSWORD" -e "CREATE DATABASE IF NOT EXISTS $DATABASE_NAME"

# seedデータを投入
seed_database "${CURRENT_PATH}/db/sql/seed.sql"

echo "Seed data inserted successfully."
