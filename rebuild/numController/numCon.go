package numController

import (
	"GORushShoping/common"
	"database/sql"
	"log"
	"sync"
)

var db *sql.DB

// NumController 数量管理结构体
type NumController struct { //数量控制,因为map是线程不安全的,所以需要加锁
	Left map[int]int //productID 数量
	Lock sync.Mutex
}

// NewNumController 新建一个数量控制器
func NewNumController() *NumController {
	return &NumController{
		Left: make(map[int]int),
	}
}

// InitNum 初始化数量
func (n *NumController) InitNum() {
	// 获取数据库连接对象
	var err error
	db, err = common.NewMysqlConn()
	if err != nil {
		log.Println("数据库连接失败,location: numCon No:31")
		panic(err)
	}
	// 查询数据库中的数据
	sqlStr := "select ID,productNum from product"
	rows, err := db.Query(sqlStr)
	defer rows.Close()
	for rows.Next() {
		var (
			productID  int
			productNum int
		)
		err := rows.Scan(&productID, &productNum)
		if err != nil {
			log.Println("数据库查询失败,location: numCon No:45")
			panic(err)
		}
		//将数据放入map中
		n.Left[productID] = productNum
	}
}
