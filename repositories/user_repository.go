package repositories

import (
	"GORushShoping/common"
	"GORushShoping/datamodels"
	"database/sql"
	"errors"
	"strconv"
)

type IUserRepository interface {
	Conn() (err error)
	Select(userName string) (user *datamodels.User, err error)
	Insert(user *datamodels.User) (userId int64, err error)
}

type UserManager struct {
	table     string
	mysqlConn *sql.DB
}

func (u *UserManager) Conn() (err error) {
	if u.mysqlConn == nil {
		Conn, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		u.mysqlConn = Conn
	}
	if u.table == "" {
		u.table = "user"
	}
	return nil
}

func (u *UserManager) Select(userName string) (user *datamodels.User, err error) {
	if userName == "" {
		return &datamodels.User{}, errors.New("userName nil")
	}
	if err = u.Conn(); err != nil {
		return
	}
	sqlStr := "select * from " + u.table + " where userName = ?"
	row, err := u.mysqlConn.Query(sqlStr, userName)
	if err != nil {
		return &datamodels.User{}, err
	}
	defer row.Close()
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.User{}, errors.New("user not exist")
	}
	user = &datamodels.User{}
	common.DataToStructByTagSql(result, user)
	return user, nil
}

func (u *UserManager) Insert(user *datamodels.User) (userId int64, err error) {
	if err = u.Conn(); err != nil {
		return 0, err
	}
	sqlStr := "insert " + u.table + " set nickName = ?,userName = ?,passWord = ?"
	stmt, err := u.mysqlConn.Prepare(sqlStr)
	if err != nil {
		return 0, err
	}
	exec, err := stmt.Exec(user.NickName, user.UserName, user.HashPassword)
	if err != nil {
		return 0, err
	}
	return exec.LastInsertId()
}

func NewUserManager(table string, db *sql.DB) IUserRepository {
	return &UserManager{table, db}
}

func (u *UserManager) SelectByID(userId int64) (user *datamodels.User, err error) {
	if err = u.Conn(); err != nil {
		return &datamodels.User{}, err
	}
	sqlStr := "select *from " + u.table + " where ID =" + strconv.FormatInt(userId, 10)
	row, err := u.mysqlConn.Query(sqlStr)
	if err != nil {
		return &datamodels.User{}, err
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.User{}, nil
	}
	user = &datamodels.User{}
	common.DataToStructByTagSql(result, user)
	return user, nil
}
