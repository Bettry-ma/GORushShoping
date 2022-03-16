package main

import (
	"GORushShoping/common"
	"GORushShoping/datamodels"
	"GORushShoping/encrypt"
	"GORushShoping/rabbitmq"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

/*
Redis
单机或者主从方式 --> 单机QPS 8W左右
Redis Cluster集群方式 --> QPS能千万以上
Redis分布式方式使用 --> QPS能千万以上
*/

//设置集群地址,最好是内网IP,访问速度会更快
var hostArray = []string{"127.0.0.1", "127.0.0.1"}

var localHost = "127.0.0.1"

// GetOneIP 数量控制接口服务器内网IP,或者getOne的SLB内网IP
var GetOneIP = "127.0.0.1"

var GetOnePort = "8088"

var port = "8083"

var hashConsistent *common.Consistent

var rabbitMqValidate *rabbitmq.RabbitMQ

// AccessControl 用来存放控制信息
type AccessControl struct {
	// 用来存放用户想要存放的信息
	sourcesArray map[int]interface{}
	sync.RWMutex
}

//创建全局变量
var accessControl = &AccessControl{sourcesArray: make(map[int]interface{})}

// GetNewRecord 获取指定的数据
func (a *AccessControl) GetNewRecord(uid int) interface{} {
	a.RWMutex.RLock()
	defer a.RWMutex.RUnlock()
	data := a.sourcesArray[uid]
	return data
}

// SetNewRecord 设置记录
func (a *AccessControl) SetNewRecord(uid int) {
	a.RWMutex.Lock()
	a.sourcesArray[uid] = "2"
	a.RWMutex.Unlock()
}

func (a *AccessControl) GetDistributedRight(req *http.Request) bool {
	//获取用户UID
	/*uid, err := req.Cookie("uid")
	if err != nil {
		fmt.Println(err)
		return false
	}*/
	//采用一致性hash算法,根据用户ID,判断获取具体机器
	/*hostRequest, err := hashConsistent.Get(uid.Value)
	if err != nil {
		fmt.Println(err)
		return false
	}*/
	hostRequest := localHost
	//判断是否为本机
	//fmt.Println("判断是否为本机: ", hostRequest)
	//fmt.Println("uid.value:", uid.Value)
	if hostRequest == localHost {
		//执行本机数据读取和校验
		return a.GetDataFromMap("2") //"uid.Value
	} else {
		//不是本机充当代理访问数据返回结果
		return a.GetDataFromOtherMap(hostRequest, req)
	}
}

// GetDataFromMap 获取本机map,并且处理业务逻辑,返回的结果类型为bool类型
func (a *AccessControl) GetDataFromMap(uid string) (isOk bool) {
	/*uidInt, err := strconv.Atoi(uid)
	if err != nil {
		fmt.Println("failed to strconv for uidInt ", err)
		return false
	}
	data := a.GetNewRecord(uidInt)
	//执行逻辑判断
	if data != nil {
		return true
	}*/
	return true
}

// GetDataFromOtherMap 从其他节点处理结果
func (a *AccessControl) GetDataFromOtherMap(host string, request *http.Request) bool {
	hostUrl := "http://" + host + ":" + port + "/checkRight"
	response, body, err := GetUrl(hostUrl, request)
	if err != nil {
		return false
	}
	//判断状态
	if response.StatusCode == 200 {
		if string(body) == "true" {
			return true
		} else {
			return false
		}
	}
	return false
	/*//获取uid
	uidPre, err := request.Cookie("uid")
	if err != nil {
		fmt.Println(err)
		return false
	}
	//获取sign
	uidSign, err := request.Cookie("sign")
	if err != nil {
		fmt.Println(err)
		return false
	}
	//模拟接口访问
	client := &http.Client{}
	req, err := http.NewRequest("Get", "http://"+host+":"+port+"/check", nil)
	if err != nil {
		fmt.Println(err)
		return false
	}
	//手动指定,排查多余cookies
	cookieUid := &http.Cookie{Name: "uid", Value: uidPre.Value, Path: "/"}
	cookieSign := &http.Cookie{Name: "sign", Value: uidSign.Value, Path: "/"}
	//添加Cookie到模拟的请求中
	req.AddCookie(cookieUid)
	req.AddCookie(cookieSign)

	//获取返回结果
	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}

	//判断状态
	if response.StatusCode == 200 {
		if string(body) == "true" {
			return true
		} else {
			return false
		}
	}
	return false*/
}

func GetUrl(hostUrl string, request *http.Request) (response *http.Response, body []byte, err error) {
	//获取Uid
	/*uidPre, err := request.Cookie("uid")
	if err != nil {
		return
	}
	//获取sign
	uidSign, err := request.Cookie("sign")
	if err != nil {
		return
	}*/

	//模拟接口访问
	client := &http.Client{}
	req, err := http.NewRequest("GET", hostUrl, nil)
	if err != nil {
		return
	}
	//手动指定,排查多余cookies
	/*cookieUid := &http.Cookie{Name: "uid", Value: uidPre.Value, Path: "/"}
	cookieSign := &http.Cookie{Name: "sign", Value: uidSign.Value, Path: "/"}
	//添加Cookie到模拟的请求中
	req.AddCookie(cookieUid)
	req.AddCookie(cookieSign)*/

	//获取返回结果
	response, err = client.Do(req)
	defer response.Body.Close()
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(response.Body)
	return
}

func CheckRight(w http.ResponseWriter, r *http.Request) {
	right := accessControl.GetDistributedRight(r)
	if !right {
		w.Write([]byte("false"))
		return
	}
	w.Write([]byte("true"))
	return
}

// Check 执行正常业务逻辑
func Check(w http.ResponseWriter, r *http.Request) {
	//执行正常业务逻辑
	//fmt.Println("执行check!")
	/*queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || len(queryForm["productID"]) <= 0 {
		w.Write([]byte("false,productID get error"))
		return
	}
	productString := queryForm["productID"][0]*/
	//获取用户Cookie
	/*userCookie, err := r.Cookie("uid")
	if err != nil {
		w.Write([]byte("false,uid get error"))
		return
	}
	fmt.Println("productID:", productString, "userID:", userCookie.Value)*/
	//fmt.Println("productID:", productString)
	//fmt.Println("productID:", 1)
	//1.分布式权限验证
	right := accessControl.GetDistributedRight(r)
	if right == false {
		w.Write([]byte("false"))
		return
	}
	//2.获取数量控制权限,防止秒杀出现超卖现象
	hostUrl := "http://" + GetOneIP + ":" + GetOnePort + "/getOne?productID=1" //todo 设计可配置接口
	responseValidate, validateBody, err := GetUrl(hostUrl, r)
	if err != nil {
		w.Write([]byte("false"))
		return
	}
	//判断数量控制接口请求状态
	if responseValidate.StatusCode == 200 {
		if string(validateBody) == "true" {
			//整合下单
			//1.获取商品ID
			/*productID, err := strconv.ParseInt(productString, 10, 64)
			if err != nil {
				w.Write([]byte("false"))
				return
			}*/
			//2.获取用户ID
			/*userID, err := strconv.ParseInt(userCookie.Value, 10, 64)
			if err != nil {
				w.Write([]byte("false"))
				return
			}*/
			//3.创建消息体
			message := datamodels.Message{1, 2}
			//类型转换
			byteMessage, err := json.Marshal(message)
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			//4.生产消息
			err = rabbitMqValidate.PublishSimple(string(byteMessage))
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			w.Write([]byte("true"))
			return
		}
	}
	w.Write([]byte("false"))
	return
}

// Auth 统一验证拦截器,每个接口都需要提前验证
func Auth(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("执行验证")
	//添加基于Cookie的权限验证
	err := checkUserInfo(r)
	if err != nil {
		return err
	}
	return nil
}

//身份校验函数
func checkUserInfo(r *http.Request) error {
	//获取UID,Cookie
	uidCookie, err := r.Cookie("uid")
	if err != nil {
		return errors.New("用户UID Cookie获取失败")
	}
	//获取用户加密串
	signCookie, err := r.Cookie("sign")
	if err != nil {
		return errors.New("用户加密串 Cookie获取失败")
	}
	//对加密串进行解密
	signByte, err := encrypt.DePwdCode(signCookie.Value)
	if err != nil {
		return errors.New("加密串已被篡改")
	}
	//结果比对
	check := checkInfo(uidCookie.Value, string(signByte))
	if check {
		return nil
	}
	return errors.New("身份校验失败")
}

func checkInfo(uidCookie, sign string) bool {
	return uidCookie == sign
}

func main() {
	//采用负载均衡器设置
	//采用一致性哈希算法
	hashConsistent = common.NewConsistent()
	//采用一致性hash算法,添加节点
	for _, v := range hostArray {
		hashConsistent.Add(v)
	}

	//获取本机向外暴露的IP地址
	localIP, errIp := common.GetInternIp()
	if errIp != nil {
		fmt.Println(errIp)
	}
	fmt.Println(localIP)
	rabbitMqValidate = rabbitmq.NewRabbitMQSimple("product")
	defer rabbitMqValidate.Destroy()

	//设置静态文件目录
	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir("./fronted/web/htmlProductShow"))))
	//设置资源目录
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("./fronted/web/public"))))

	//1.过滤器
	//filter := common.NewFilter()
	//注册拦截器
	//filter.RegisterFilterUri("/check", Auth)
	//filter.RegisterFilterUri("/checkRight", Auth)
	//2.启动服务
	//http.HandleFunc("/check", filter.Handle(Check))
	http.HandleFunc("/check", Check)
	http.HandleFunc("/checkRight", CheckRight)
	//启动服务
	err := http.ListenAndServe(":8083", nil)
	if err != nil {
		fmt.Println(err)
	}
}
