package repositories

import (
	"GORushShoping/common"
	"GORushShoping/datamodels"
	"database/sql"
	"fmt"
	"strconv"
)

// IProduct 定义了商品的接口类型,开发对应的接口
type IProduct interface {
	Conn() (err error)                                            //数据库连接
	Insert(product *datamodels.Product) (id int64, err error)     //商品插入
	Delete(id int64) (state bool)                                 //商品删除
	Update(product *datamodels.Product) (err error)               //商品更新
	SelectById(id int64) (product *datamodels.Product, err error) //根据商品id查询商品信息
	SelectAll() (products []*datamodels.Product, err error)       //查询所有商品信息
	SubProductNum(productID int64) error
}

// ProductManager 将实现接口
type ProductManager struct { //面向对象思维,新建一个对象管理者的结构体
	table     string
	mysqlConn *sql.DB
}

func NewProductManager(table string, db *sql.DB) IProduct { //注意这里返回的是IProduct,我们要做一个实现接口
	return &ProductManager{table: table, mysqlConn: db} //如果方法没实现时会报错
}

// Conn 获取数据库连接
func (p *ProductManager) Conn() (err error) {
	if p.mysqlConn == nil { //如果数据库为空则为他创建连接
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		p.mysqlConn = mysql //如果无err则说明数据库创建成功,将数据库指针连接赋值给p
	}
	if p.table == "" {
		p.table = "product"
	}
	return nil
}

// Insert 插入商品
func (p *ProductManager) Insert(product *datamodels.Product) (id int64, err error) {
	if err = p.Conn(); err != nil { //如果获取到的链接
		return
	}
	sqlStr := "insert product set productName=?,productNum=?,productImage=?,productUrl=?,productPriceNow=?,productPricePre=?" //SQL语句
	stmt, err := p.mysqlConn.Prepare(sqlStr)
	if err != nil {
		return 0, err
	}
	exec, err := stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl, product.ProductPriceNow, product.ProductPricePre)
	if err != nil {
		return 0, err
	}
	return exec.LastInsertId()
}

// Delete 删除商品
func (p *ProductManager) Delete(id int64) (state bool) {
	if err := p.Conn(); err != nil {
		return false
	}
	sqlStr := "delete from product where id =?"
	stmt, err := p.mysqlConn.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("failed to %s for prepare %s\n", sqlStr, err)
		return false
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return false
	}
	fmt.Println("delete from product where id = ", id, " success")
	return true
}

// Update 更新商品
func (p *ProductManager) Update(product *datamodels.Product) (err error) {
	if err := p.Conn(); err != nil {
		return err
	}
	sqlStr := "update product set productName=?,productNum=?,productImage=?,productUrl=?,productPriceNow=?,productPricePre=? where id = " + strconv.FormatInt(product.ID, 10)
	stmt, err := p.mysqlConn.Prepare(sqlStr)
	if err != nil {
		fmt.Println("prepare for sql ", sqlStr, " failed", err)
		return err
	}
	exec, err := stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl, product.ProductPriceNow, product.ProductPricePre)
	affected, err := exec.RowsAffected()
	if err != nil {
		fmt.Println("get affected failed", err)
		return err
	}
	fmt.Println("update success,affected ", affected)
	return nil
}

// SelectById 根据商品Id查询商品
func (p *ProductManager) SelectById(id int64) (productResult *datamodels.Product, err error) {
	if err := p.Conn(); err != nil {
		return &datamodels.Product{}, err //当没有数据库连接时返回空结构
	}
	sqlStr := "select *from " + p.table + " where ID =" + strconv.FormatInt(id, 10)
	row, err := p.mysqlConn.Query(sqlStr)
	defer row.Close()
	if err != nil {
		return &datamodels.Product{}, err
	}
	result := common.GetResultRow(row)
	if len(result) == 0 { //如果是0的情况,则
		return &datamodels.Product{}, nil
	}
	productResult = &datamodels.Product{}
	common.DataToStructByTagSql(result, productResult)
	return productResult, nil
}

// SelectAll 获取所有商品
func (p *ProductManager) SelectAll() (products []*datamodels.Product, err error) {
	if err := p.Conn(); err != nil {
		return nil, err //当没有数据库连接时返回空结构
	}
	sqlStr := "select *from " + p.table
	rows, err := p.mysqlConn.Query(sqlStr)
	if err != nil {
		fmt.Println("failed to query sql ", sqlStr)
		return nil, err
	}
	defer rows.Close()
	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil, nil
	}
	for _, v := range result {
		product := &datamodels.Product{}
		common.DataToStructByTagSql(v, product)
		products = append(products, product)
	}
	return products, nil
}

// SubProductNum 商品数量抢购扣除
func (p *ProductManager) SubProductNum(productID int64) error {
	if err := p.Conn(); err != nil {
		return err
	}
	//sqlStr := "update " + p.table + " set productNum = productNum-1 where ID =" + strconv.FormatInt(productID, 10)
	sqlStr := "update product set productNum = productNum -1 where ID =1"
	stmt, err := p.mysqlConn.Prepare(sqlStr)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	return err
}
