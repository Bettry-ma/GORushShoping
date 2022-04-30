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
	"net/url"
	"strconv"
	"sync"
)

/*
商品数量控制接口
从数据库中获取相应商品的数量
*/

var db *sql.DB
var left int //商品的库存,在初始化时就绪
var rab *rabbitmq.RabbitMQ
var Here = make(chan bool, 1000) //全局通道变量,用于传输订单消费状态

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
	////加锁
	//mutex.Lock()
	//defer mutex.Unlock()
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
			Here <- true
		}()
	} else {
		Here <- false
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
	if msg {
		w.Write([]byte("true"))
	} else {
		w.Write([]byte("false"))
	}
	//msg := <-Here
	//w.Write([]byte(msg))
	return
}

func GetBuy(w http.ResponseWriter, req *http.Request) {
	query, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		fmt.Println(err)
	}
	userID := query["userID"][0]
	productID := query["productID"][0]
	price := query["price"][0]
	num := query["num"][0]
	numInt, err := strconv.Atoi(num)
	time := query["time"][0]
	userIDInt, err := strconv.Atoi(userID)
	fmt.Println(numInt, userIDInt)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(userID, productID, price, num, time)

	// 库存减去 订单数量
	var left int
	err = db.QueryRow("select productNum from product where ID = ?", productID).Scan(&left)
	//todo 调用数量控制接口
	//todo if left>0 创建订单

	orderInsert := &datamodels.Order{
		//填充数据
	}
	fmt.Println(orderInsert)
	sqlStr := ""
	result, err := db.Exec(sqlStr)
	affected, err := result.RowsAffected()
	if affected > 0 {

	}
	//todo 通道传输剩余商品数量
	//todo else
	//失败逻辑
}

func GetSettleAccount(w http.ResponseWriter, req *http.Request) {
	query, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		fmt.Println(err)
	}
	userID := query["userID"][0]
	userIDInt, err := strconv.Atoi(userID)
	fmt.Println(userIDInt)
	if err != nil {
		fmt.Println(err)
	}
	sqlStr := "select o.productNum,p.productPriceNow " +
		" from 'order' o,product p" +
		"where o.productId = p.ID and o.userId = " + userID + " and o.orderStatus = 0"
	var productNum int
	var productPriceNow int
	account := 0
	rows, err := db.Query(sqlStr)
	if rows.Next() {
		err := rows.Scan(&productNum, &productPriceNow)
		if err != nil {
			fmt.Println(err)
		}
		account += productPriceNow * productNum
	}
	w.Write([]byte("订单付款金额" + strconv.Itoa(account)))
	//account -->
}

func main() {
	defer rab.Destroy()
	http.HandleFunc("/getOne", GetProduct)
	http.HandleFunc("/buy", GetBuy)
	err := http.ListenAndServe(":8088", nil)
	fmt.Println("app is running...")
	if err != nil {
		log.Fatalln("err", err)
	}
}
