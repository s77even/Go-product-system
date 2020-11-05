package repositories

import (
	"database/sql"
	"fmt"
	"imooc-product/common"
	"imooc-product/datamodels"
	"strconv"
)

//先开发接口
//实现定义的接口

type Iproduct interface {
	//连接数据
	Conn() error
	Insert(*datamodels.Product) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Product) error
	SelectByKey(int64) (*datamodels.Product, error)
	SelectAll() ([]*datamodels.Product, error)
}
type  ProductManager struct {
	table string
	mysqlConn *sql.DB
}
//NewProductManager 构造函数 并且能够检查是否实现接口
func NewProductManager (table string , db *sql.DB)Iproduct{
	return &ProductManager{table: table,mysqlConn: db}
}

// Conn 初始化连接
func (p *ProductManager)Conn()(err error){
	if p.mysqlConn == nil{
		mysql , err := common.NewMysqlConn()
		if err != nil{
			return err
		}
		p.mysqlConn = mysql
	}
	if p.table==""{
		p.table="product"
	}
	return
}
// Insert 插入...
func (p *ProductManager)Insert(product *datamodels.Product) (productId int64,err error){
	if err = p.Conn(); err != nil{
		return
	}
	sql := "INSERT product SET productName=?,productNum=?,productImage=?,productUrl=?"
	stem , errSql := p.mysqlConn.Prepare(sql)
	if errSql !=nil{
		return 0,errSql
	}
	result, errStem := stem.Exec(product.ProductName,product.ProductNum,product.ProductImage,product.ProductUrl)
	if errStem != nil{
		return 0,errStem
	}
	return result.LastInsertId()
}
//Delete 删除...
func (p *ProductManager)Delete(productId int64) bool{
	if err := p.Conn(); err != nil{
		return false
	}//判断连接是否可用
	sql := "delete from product where ID=?"
	stem , err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return false
	}
	_ , err =stem.Exec(productId)
	if err != nil {
		return false
	}
	return true
}
//Update 更新...
func (p *ProductManager)Update(product *datamodels.Product)(err error){
	if err = p.Conn(); err != nil{
	return err
	}//判断连接是否可用
	sql := "Updata product SET productName=?,productNum=?,productImage=?,productUrl=? where ID="+strconv.FormatInt(product.ID,10)
	stem , err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}
	_ , err =stem.Exec(product.ProductName,product.ProductNum,product.ProductImage,product.ProductUrl)
	if err != nil {
		return err
	}
	return nil
}
//SelectByKey 根据id查询记录
func (p *ProductManager)SelectByKey(productID int64) (productResult *datamodels.Product,err error){
	if err=p.Conn();err!=nil{
		return &datamodels.Product{} , err
	}
	sql := "select * from "+ p.table + " where ID= " + strconv.FormatInt(productID,10)
	row ,errRow :=  p.mysqlConn.Query(sql)

	if errRow != nil {
		return &datamodels.Product{},errRow
	}
	result := common.GetResultRow(row)

	if len(result)==0{
		return &datamodels.Product{},nil
	}
	fmt.Println(result)
	productResult = &datamodels.Product{}             // 指针要分配内存
	common.DataToStructByTagSql(result,productResult) //数据映射
	return productResult,nil
}
//SelectAll 查询全表
func (p *ProductManager)SelectAll() (productArray []*datamodels.Product,errProduct error){
	if err:=p.Conn();err!=nil{
		return nil ,err
	}
	sql := "select * from " + p.table
	rows , err := p.mysqlConn.Query(sql)
	if err != nil {
		return nil ,err
	}
	result := common.GetResultRows(rows)
	if len(result)==0{
		return nil, nil
	}
	for _ , v := range result{
		product :=&datamodels.Product{}
		common.DataToStructByTagSql(v,product)
		productArray=append(productArray, product)
	}
	return productArray,nil
}
