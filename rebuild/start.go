package main

import (
	pb "GORushShoping/protobuf"
	"GORushShoping/rebuild/numController"
	"GORushShoping/rebuild/taskMgt"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

var ncl *numController.NumController
var task *taskMgt.Task
var bufferCache int
var left int

func GetGoBuy(w http.ResponseWriter, req *http.Request) {
	// 解析URL中的参数
	query, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		fmt.Println(err)
	}
	userIDStr := query["userID"][0]
	productIDStr := query["productID"][0]
	productNumStr := query["productNum"][0]
	userID, err := strconv.Atoi(userIDStr)
	productID, err := strconv.Atoi(productIDStr)
	productNum, err := strconv.Atoi(productNumStr)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(ncl.Left[productID])
	w.Write([]byte("true"))
	//fmt.Println(ncl.Left[1])
	//将性能消耗放在后台线程解决
	go func() {
		if ncl.Left[1] > 0 { //库存大于0则进行下单逻辑
			//if left > 0 { //库存大于0则进行下单逻辑
			ncl.Lock.Lock()
			ncl.Left[1] = ncl.Left[1] - 1 //dataRace情况,所以要加锁
			left--
			ncl.Lock.Unlock()

			//引入protobuf序列化
			or := &pb.Order{
				ProductID:  int64(productID),
				UserID:     int64(userID),
				ProductNum: int64(productNum),
			}
			/*or := &order.Order{
				ProductID:  1,
				UserID:     2,
				ProductNum: 1,
			}*/
			task.OrderChan <- or //异步处理,channel的一个数据可以相当于一个线程
			return
		} else {
			fmt.Println("库存不足,抢购失败")
		}
		return
	}()
	return
}

func main() {
	bufferCache = 500
	task = taskMgt.NewInstance(bufferCache)
	ncl = numController.NewNumController()
	ncl.InitNum()
	http.HandleFunc("/order", GetGoBuy)
	err := http.ListenAndServe(":8099", nil)
	fmt.Println("app is running...")
	if err != nil {
		log.Fatalln("err", err)
	}
}
