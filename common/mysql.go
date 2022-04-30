package common

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// NewMysqlConn 创建mysql连接
func NewMysqlConn() (db *sql.DB, err error) {
	//conf, err := ini.Load("config.ini")
	//if err != nil {
	//	return nil, fmt.Errorf("fail to read file: %v", err)
	//}
	//mysql := conf.Section("MySql")
	//host := mysql.Key("host").String()
	//port := mysql.Key("port").String()
	//user := mysql.Key("user").String()
	//password := mysql.Key("password").String()
	//database := mysql.Key("db").String()
	//db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", user, password, host, port, database))
	db, err = sql.Open("mysql", "root:62748226@tcp(192.168.1.102:3306)/gorush?charset=utf8")
	return
}

// GetResultRow 获取返回值,获取一条
func GetResultRow(rows *sql.Rows) map[string]string {
	columns, err := rows.Columns()
	if err != nil {
		fmt.Println("rows get column failed,", err)
		panic("rows get column failed")
	}
	scanArgs := make([]interface{}, len(columns))
	//values := make([]interface{}, len(columns))
	values := make([][]byte, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	record := make(map[string]string)
	for rows.Next() {
		//将行数据保存到record字典
		err := rows.Scan(scanArgs...)
		if err != nil {
			fmt.Println("rows scan failed", err)
		}
		for i, value := range values {
			if value != nil {
				//record[columns[i]] = string(value.([]byte))
				record[columns[i]] = string(value)
			}
		}
	}
	return record
}

//获取所有
func GetResultRows(rows *sql.Rows) map[int]map[string]string {
	//返回所有列
	columns, err := rows.Columns()
	if err != nil {
		fmt.Println("rows get column failed,", err)
		panic("rows get column failed")
	}
	vals := make([][]byte, len(columns))       //在这里表示一行所有列的值,用[]byte表示
	scans := make([]interface{}, len(columns)) //这里用scans引用vals,把数据填充到[]byte里
	for k := range vals {
		scans[k] = &vals[k]
	}
	var i int
	i = 0
	result := make(map[int]map[string]string)
	for rows.Next() {
		//填充数据
		rows.Scan(scans...)
		//每行数据
		row := make(map[string]string)
		for k, v := range vals { //把vals中的数复制到row中
			key := columns[k]
			row[key] = string(v) //这里吧[]byte数据转成string
		}
		result[i] = row
		i++
	}
	return result
}
