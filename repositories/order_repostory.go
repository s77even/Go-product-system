package repositories

import (
	"database/sql"
	"imooc-product/common"
	"imooc-product/datamodels"
	"strconv"
)

type IOrderRepository interface {
	Conn() error
	Insert(order *datamodels.Order) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Order) error
	SelectByKey(int64) (*datamodels.Order, error)
	SelectAll() ([]*datamodels.Order, error)
	SelectAllWithInfo()(map[int]map[string]string,error)
}

type OrderManagerRepository struct {
	table string
	mysqlConn *sql.DB
}

//NewOrderManagerRepository 构造函数
func NewOrderManagerRepository (table string , db *sql.DB)IOrderRepository{
	return &OrderManagerRepository{table: table,mysqlConn: db}
}
//Conn ...
func (o *OrderManagerRepository)Conn()error{
	if o.mysqlConn == nil{
		mysql , err := common.NewMysqlConn()
		if err != nil{
			return err
		}
		o.mysqlConn = mysql
	}
	if o.table==""{
		o.table="orders"
	}
	return nil
}
//Insert 插入新订单
func (o *OrderManagerRepository) Insert(order *datamodels.Order)(orderID int64, err error) {
	if err = o.Conn();err!=nil{
		return
	}
	sql := "insert "+o.table+" set userID=?,productID=?,orderStatus=?"
	stem , errSql := o.mysqlConn.Prepare(sql)
	if errSql !=nil{
		return 0,errSql
	}
	result, errStem := stem.Exec(order.UserId,order.ProductId,order.OrderStatus)
	if errStem != nil{
		return 0,errStem
	}
	return result.LastInsertId()
}
//Delete ...
func (o *OrderManagerRepository)Delete(orderID int64)bool{
	if err := o.Conn();err!=nil{
		return false
	}
	sql := "delete from "+o.table+" where ID =?"
	stem , err := o.mysqlConn.Prepare(sql)
	if err != nil{
		return false
	}
	_ , err =stem.Exec(orderID)
	if err != nil {
		return false
	}
	return true
}
//Update ...
func (o *OrderManagerRepository)Update(order *datamodels.Order)error{
	if err:=o.Conn();err!=nil{
		return err
	}
	sql := "update "+o.table+" userID=?,productID=?,orderStatus=? where ID="+strconv.FormatInt(order.ID,10)
	stem , err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}
	_ ,err = stem.Exec(order.UserId, order.ProductId , order.OrderStatus)
	if err != nil {
		return err
	}
	return nil
}
//SelectByKey ..
func (o *OrderManagerRepository)SelectByKey(orderId int64)(order *datamodels.Order,err error){
	if err:=o.Conn();err!=nil{
		return &datamodels.Order{}, err
	}
	sql := "select * from "+ o.table + " where ID=" + strconv.FormatInt(orderId,10)
	row , errow := o.mysqlConn.Query(sql)
	if errow != nil {
		return &datamodels.Order{},err
	}
	result := common.GetResultRow(row)

	if len(result)==0{
		return &datamodels.Order{},nil
	}
	order = &datamodels.Order{}
	common.DataToStructByTagSql(result,order)
	return order,nil
}
//SelectAll ...
func (o *OrderManagerRepository)SelectAll() (orderArray []*datamodels.Order,err error){
	if err:=o.Conn();err!=nil{
		return nil ,err
	}
	sql := "select * from " + o.table
	rows , err := o.mysqlConn.Query(sql)
	if err != nil {
		return nil ,err
	}
	result := common.GetResultRows(rows)
	if len(result)==0{
		return nil, nil
	}
	for _ , v := range result{
		order :=&datamodels.Order{}
		common.DataToStructByTagSql(v,order)
		orderArray=append(orderArray, order)
	}
	return

}
//SelectAllWithInfo ...
func (o *OrderManagerRepository)SelectAllWithInfo()(orderMap map[int]map[string]string ,err error){
	if errConn := o.Conn(); errConn != nil {
		return nil, errConn
	}
	sql := "Select o.ID,p.productName,o.orderStatus From imooc.orders as o left join product as p on o.productID=p.id"
	rows, errRows := o.mysqlConn.Query(sql)
	if errRows != nil {
		return nil, errRows
	}
	return common.GetResultRows(rows), err
}
