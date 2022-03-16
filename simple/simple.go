package main

import (
	"GORushShoping/common"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var db, _ = common.NewMysqlConn()
var left int
var mutex sync.Mutex

func init() {
	var err error
	db, err = common.NewMysqlConn()
	sqlStr := "select productNum from product where ID =1"
	rows, err := db.Query(sqlStr)
	if err != nil {
		log.Fatalln("numControl get num error")
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&left)
		if err != nil {
			log.Println("query product Num error")
		}
	}
}
func Simple(w http.ResponseWriter, req *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	sqlStr := "update product set productNum = productNum -1 where ID =1"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		log.Fatalln(err)
	}
	if left > 0 {
		_, err = stmt.Exec()
		left--
	}
	fmt.Println(left)
	w.Write([]byte("true"))
}

func main() {
	defer db.Close()
	http.HandleFunc("/simple", Simple)
	err := http.ListenAndServe(":8090", nil)
	fmt.Println("app is running...")
	if err != nil {
		log.Fatalln("err", err)
	}
}
