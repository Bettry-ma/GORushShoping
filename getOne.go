package main

import (
	"GORushShoping/common"
	"GORushShoping/datamodels"
	"GORushShoping/rabbitmq"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

/*
商品数量控制接口
从数据库中获取相应商品的数量
*/

var db *sql.DB
var left int //商品的库存,在初始化时就绪
var rab *rabbitmq.RabbitMQ
var Here = make(chan string, 100) //全局通道变量,用于传输订单消费状态

func init() {
	rab = rabbitmq.NewRabbitMQSimple("product")
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

// 数量控制对象
type numControl struct {
	table      string //表名
	productNum string //需要数量控制的商品名称的数量字段
	productID  int    //商品ID
	//left       int    //该商品的库存
}

func NewProduct(table, productNum string) *numControl {
	return &numControl{table: table, productNum: productNum}
}

// 内部方法,获取数量
/*func getNum() {
	var err error
	db, err = common.NewMysqlConn()
	sqlStr := "select " + n.productNum + " from " + n.table + " where ID =" + strconv.Itoa(n.productID)
	rows, err := db.Query(sqlStr)
	if err != nil {
		log.Fatalln("numControl get num error")
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&n.left)
		if err != nil {
			log.Println("query product Num error")
		}
	}
}*/

//互斥锁
var mutex sync.Mutex

// GetOneProduct 获取秒杀商品
func GetOneProduct() {
	//加锁
	mutex.Lock()
	defer mutex.Unlock()
	if left > 0 {
		go func() {
			message := datamodels.Message{1, 2}
			//类型转换
			byteMessage, err := json.Marshal(message)
			if err != nil {
				log.Fatalln(err)
			}
			//4.生产消息
			err = rab.PublishSimple(string(byteMessage))
			if err != nil {
				log.Fatalln(err)
			}
			left--
			fmt.Println("left: ", left)
			Here <- "success"
		}()
	}
}

func GetProduct(w http.ResponseWriter, req *http.Request) {
	/*//1.获取URL中的商品ID
	query, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil || len(query["productID"]) <= 0 {
		log.Fatalln("couldn't get ProductID,location:getOne")
	}
	productID := query["productID"][0]
	//2.新建 numControl 实例对象
	nc := NewProduct("product", "productNum")
	nc.productID, err = strconv.Atoi(productID)
	if err != nil {
		log.Fatalln("failed to set nc productID")
	}
	//fmt.Println("nc.ProductID", nc.productID)
	if nc.GetOneProduct() {
		w.Write([]byte("true"))
		return
	}
	w.Write([]byte("false"))*/
	GetOneProduct()
	msg := <-Here
	w.Write([]byte(msg))
	return
}

func main() {
	defer rab.Destroy()
	http.HandleFunc("/getOne", GetProduct)
	err := http.ListenAndServe(":8088", nil)
	fmt.Println("app is running...")
	if err != nil {
		log.Fatalln("err", err)
	}
}
