#!/bin/sh

# 自動機能の停止
sudo systemctl disable data_store_server.service

# サービスの停止
sudo systemctl stop data_store_server.service

# サービスファイルの削除
sudo rm /etc/systemd/system/data_store_server.service

# ワーキングディレクトリの削除
sudo rm -rf /etc/data_store_server