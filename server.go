package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type server struct {
	db_stockstore   *sql.DB
	unix_time_start time.Time
}

func StartServer() error {
	svr := &server{
		unix_time_start: time.Date(1970, 1, 1, 9, 0, 0, 0, time.Local),
	}
	err := svr.openDB()
	if err != nil {
		return err
	}

	http.HandleFunc("/", handler)
	http.HandleFunc("/api/stock/store", svr.store_handler)

	log.Println("start listen and serve")
	return http.ListenAndServe(":8080", nil)
}

func (s *server) openDB() error {
	log.Println("open stock_store db")
	db_stockstore, err := sql.Open("mysql", "pi:raspberry@/mysql")
	if err != nil {
		return err
	}
	s.db_stockstore = db_stockstore
	return nil
}

func (s *server) closeDB() {
	if s.db_stockstore != nil {
		s.db_stockstore.Close()
	}
}

func (s *server) makeDateString(tkey string) (string, error) {
	v, err := strconv.ParseUint(tkey, 10, 64)
	if err != nil {
		return "", err
	}
	t := s.unix_time_start.Add(time.Duration(v) * time.Nanosecond)
	return t.Format("2006-01-02 15:04:05.000000"), nil
}

func setStatus(w http.ResponseWriter, code int, msg string) {
	if code != http.StatusOK {
		log.Println(msg)
	}
	w.WriteHeader(code)
	fmt.Fprintf(w, msg)
}

func handler(w http.ResponseWriter, r *http.Request) {
	setStatus(w, http.StatusOK, "Hello, World")
}

/**
 * URL: api/store_by_brand
 * パラメータ
 *   tkey          : 時刻キー(unix時刻)
 *   brand_code    : 銘柄コード
 *   price         : 株価
 *   volume        : 出来高
 *   usd_jpy       : ドル円
 *   average_nikkei: 日経平均
 *   nikkei_futures: 日経225先物
 */
func (s *server) store_handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// パラメータ取得
	tkey := query.Get("tkey") // 時刻キー(unix時刻)
	if len(tkey) <= 0 {
		setStatus(w, http.StatusBadRequest, "FAIL: missed tkey")
		return
	}
	tkeyDate, err := s.makeDateString(tkey)
	if err != nil {
		setStatus(w, http.StatusBadRequest, fmt.Sprintf("FAIL: invalid format tkey(%v)", err))
		return
	}

	bid := query.Get("brand_code") // 銘柄コード
	if len(bid) <= 0 {
		setStatus(w, http.StatusBadRequest, "FAIL: missed brand_code")
		return
	}
	price := query.Get("price") // 株価
	if len(price) <= 0 {
		setStatus(w, http.StatusBadRequest, "FAIL: missed price")
		return
	}
	volume := query.Get("volume") // 出来高
	if len(volume) <= 0 {
		setStatus(w, http.StatusBadRequest, "FAIL: missed volume")
		return
	}
	usdJpy := query.Get("usd_jpy") // ドル円
	if len(usdJpy) <= 0 {
		setStatus(w, http.StatusBadRequest, "FAIL: missed usd_jpy")
		return
	}
	averageNikkei := query.Get("average_nikkei") // 日経平均
	if len(averageNikkei) <= 0 {
		setStatus(w, http.StatusBadRequest, "FAIL: missed average_nikkei")
		return
	}
	nikkeiFutures := query.Get("nikkei_futures") // ドル円
	if len(nikkeiFutures) <= 0 {
		setStatus(w, http.StatusBadRequest, "FAIL: missed nikkei_futures")
		return
	}

	ins, err := s.db_stockstore.Prepare("INSERT INTO stock_data(tkey,brand_code,price,volume,usd_jpy,average_nikkei,nikkei_futures) VALUES(?,?,?,?,?,?,?)")
	if err != nil {
		setStatus(w, http.StatusBadRequest, fmt.Sprintf("FAIL: failed to insert to db s: %v", err))
		return
	}
	_, err = ins.Exec(tkeyDate, bid, price, volume, usdJpy, averageNikkei, nikkeiFutures)
	if err != nil {
		setStatus(w, http.StatusBadRequest, fmt.Sprintf("FAIL: failed to insert to db: %v", err))
		return
	}
	setStatus(w, http.StatusOK, "OK")
}
