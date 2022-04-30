package repositories

import (
	"GORushShoping/common"
	"GORushShoping/datamodels"
	"database/sql"
	"fmt"
	"strconv"
)

type IOrderRepository interface {
	Conn() (err error)
	Insert(order *datamodels.Order) (id int64, err error)
	Delete(orderId int64) (status bool)
	Update(order *datamodels.Order) (err error)
	SelectById(orderId int64) (orderObj *datamodels.Order, err error)
	SelectAll() (orderObjs []*datamodels.Order, err error)
	SelectAllWithInfo() (infos map[int]map[string]string, err error) //查询某订单的信息
}

type OrderManager struct {
	table     string  //对应的order表
	mysqlConn *sql.DB //数据库连接对象
}

func NewOrderManager(table string, db *sql.DB) IOrderRepository {
	return &OrderManager{table: table, mysqlConn: db}
}

func (o *OrderManager) Conn() (err error) {
	if o.mysqlConn == nil {
		conn, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		o.mysqlConn = conn
	}
	if o.table == "" {
		o.table = "order"
	}
	return nil
}

func (o *OrderManager) Insert(order *datamodels.Order) (id int64, err error) {
	if err = o.Conn(); err != nil {
		return
	}
	//sqlStr := `insert order set ID = ?,UserId = ?,ProductId = ?,OrderStatus = ?`
	//sqlStr := "insert " + o.table + " set userID = ?,productID = ?,orderStatus = ?"
	sqlStr := "insert  `order` set userID = ?,productID = ?,orderStatus = ?,productNum= ?"
	stmt, err := o.mysqlConn.Prepare(sqlStr)
	if err != nil {
		return 0, err
	}
	exec, err := stmt.Exec(order.UserId, order.ProductId, order.OrderStatus, order.ProductNum)
	if err != nil {
		return 0, err
	}
	return exec.LastInsertId()
}

func (o *OrderManager) Delete(orderId int64) (status bool) {
	if err := o.Conn(); err != nil {
		return false
	}
	sqlStr := "delete from `order` where ID = ?"
	stmt, err := o.mysqlConn.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("failed to %s for prepare %s\n", sqlStr, err)
		return false
	}
	_, err = stmt.Exec(orderId)
	if err != nil {
		return false
	}
	fmt.Println("delete from product where id = ", orderId, " success")
	return true
}

func (o *OrderManager) Update(order *datamodels.Order) (err error) {
	if err = o.Conn(); err != nil {
		return err
	}
	sqlStr := "update " + o.table + " set userID = ?,productID = ?,orderStatus = ? where ID = " + strconv.FormatInt(order.ID, 10)
	stmt, err := o.mysqlConn.Prepare(sqlStr)
	if err != nil {
		fmt.Println("prepare for sql ", sqlStr, " failed", err)
		return err
	}
	exec, err := stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)
	affected, err := exec.RowsAffected()
	if err != nil {
		fmt.Println("get affected failed", err)
		return err
	}
	fmt.Println("update success,affected ", affected)
	return nil
}

func (o *OrderManager) SelectById(orderId int64) (orderObj *datamodels.Order, err error) {
	if err = o.Conn(); err != nil {
		return &datamodels.Order{}, err
	}
	sqlStr := "select * from `order` where ID = " + strconv.FormatInt(orderId, 10)
	row, err := o.mysqlConn.Query(sqlStr)
	if err != nil {
		return &datamodels.Order{}, err
	}
	defer row.Close()
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Order{}, nil
	}
	orderObj = &datamodels.Order{}
	common.DataToStructByTagSql(result, orderObj)
	return orderObj, nil
}

func (o *OrderManager) SelectAll() (orderObjs []*datamodels.Order, err error) {
	if err = o.Conn(); err != nil {
		return nil, err
	}
	sqlStr := "select *from `order`"
	//fmt.Println("执行到此处了,查询全部")
	rows, err := o.mysqlConn.Query(sqlStr)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	results := common.GetResultRows(rows)
	if len(results) == 0 {
		return nil, nil
	}
	for _, v := range results {
		orderObj := &datamodels.Order{}
		common.DataToStructByTagSql(v, orderObj)
		orderObjs = append(orderObjs, orderObj)
	}
	return orderObjs, nil
}

func (o *OrderManager) SelectAllWithInfo() (infos map[int]map[string]string, err error) {
	if err = o.Conn(); err != nil {
		return nil, err
	}
	sqlStr := "select o.ID,p.productName,o.orderStatus from `order` as o left join product as p on o.productID=p.ID" //连接查询
	rows, err := o.mysqlConn.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	return common.GetResultRows(rows), nil
}
