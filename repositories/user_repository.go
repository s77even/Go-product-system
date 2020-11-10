package repositories

import (
	"database/sql"
	"errors"
	"imooc-product/common"
	"imooc-product/datamodels"
	"strconv"
)

type IUserRepository interface {
	Conn()error
	Select(userName string)(user *datamodels.User,err error)
	Insert(user *datamodels.User)(userId int64, err error)
}

type UserManagerRepository struct {
	table string
	mysqlConn *sql.DB
}

func NewUserManagerRepository(table string , db *sql.DB)IUserRepository{
	return &UserManagerRepository{
		table: table,
		mysqlConn: db,
	}
}

func (u *UserManagerRepository)Conn()error  {
	if u.mysqlConn == nil {
		newConn, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		u.mysqlConn	= newConn
	}
	if u.table==""{
		u.table = "user"
	}
	return nil
}

func (u *UserManagerRepository)Select(userName string)(user *datamodels.User,err error){
	if userName==""{
		return &datamodels.User{},errors.New("userName should not be empty")
	}
	if err := u.Conn(); err!=nil{
		return &datamodels.User{},err
	}
	sql := "select * from " + u.table+" where userName=?"
	rows , errRows:= u.mysqlConn.Query(sql , userName)
	if errRows!=nil{
		return &datamodels.User{},errRows
	}
	defer rows.Close()
	result := common.GetResultRows(rows)
	if len(result)==0{
		return &datamodels.User{},errors.New("user not exist")
	}
	user = new(datamodels.User)
	common.DataToStructByTagSql(result[0],user)
	return user,nil
}

func (u *UserManagerRepository)Insert(user *datamodels.User)(userId int64, err error){
	if err = u.Conn(); err!=nil{
		return
	}
	sql := "insert "+u.table+" set nickName=?,userName=?,password=?"
	stem, errStem := u.mysqlConn.Prepare(sql)
	if errStem != nil{
		return userId, errStem
	}
	defer stem.Close()
	result, errResult := stem.Exec(user.NickName,user.UserName,user.HashPassword)
	if errResult != nil{
		return userId, errResult
	}
	return result.LastInsertId()
}

func (u *UserManagerRepository)SelectByID(userID int64)(user *datamodels.User,err error){
	if err = u.Conn(); err!=nil{
		return
	}
	sql := "select * from "+u.table+" where ID="+strconv.FormatInt(userID,10)
	row, errRow := u.mysqlConn.Query(sql)
	if errRow!=nil{
		return &datamodels.User{},errRow
	}
	defer row.Close()
	result := common.GetResultRow(row)
	if len(result)==0{
		return &datamodels.User{},errors.New("user not exist")
	}
	user = &datamodels.User{}
	common.DataToStructByTagSql(result,user)
	return user,nil
}