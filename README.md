# data_store_server

## 環境構築

golangのインストール
```
$ wget https://dl.google.com/go/go1.14.4.linux-armv6l.tar.gz
$ sudo tar -C /usr/local -xzf go1.14.4.linux-armv6l.tar.gz
$ /usr/local/go/bin/go # 起動確認
$ vi ~/.bashrc
・・・
export PATH=$PATH:/usr/local/go/bin # 追加
```

mariadbのインストール & 初期設定
```
$ sudo apt-get install mariadb-server
$ sudo /usr/bin/mysql_secure_installation # 初期設定
$ sudo mariadb # 起動確認
```

ユーザ追加
```
$ sudo mariadb # 起動確認
> use mysql;
> drop user pi@localhost;
> flush privileges;
> create user pi@localhost identified by 'raspberry';
> grant all on *.* to pi@localhost;
```

テーブル作成
```
$ sudo mariadb # 起動確認
> use mysql;
> create table stock_data(tkey DATETIME(6), brand_code int, price int, volume int, usd_jpy int, average_nikkei int, nikkei_futures int, constraint stock_data_pk PRIMARY KEY (tkey, brand_code)); # テーブルの作成
> show create table stock_data; # テーブルの定義確認
> insert into stock_data values('2006-01-02 15:04:05.000001', 2, 3, 4, 5, 6, 7); # データの挿入確認
> delete from stock_data where tkey='1' and brand_code='2'; # データの削除確認
```

## ビルド方法
```
go build
```

## 実行方法
```
$ sudo ./data_store_server
```