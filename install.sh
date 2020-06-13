#!/bin/sh

go build

# サービスの停止
sudo systemctl stop data_store_server.service

# ワーキングディレクトリの作成
sudo mkdir -p /etc/data_store_server

# 実行ファイルのコピー
sudo cp data_store_server /etc/data_store_server/

# サービスファイルをコピー
sudo cp data_store_server.service /etc/systemd/system/

# 自動機能の有効化
sudo systemctl disable data_store_server.service
sudo systemctl enable data_store_server.service

# serviceファイルの再読み込み
sudo systemctl daemon-reload

# サービスの開始
sudo systemctl start data_store_server.service