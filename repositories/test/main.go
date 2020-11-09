package main

import (
	"fmt"
	"imooc-product/common"
	"imooc-product/repositories"
)

func main(){
	db, _ :=  common.NewMysqlConn()
	//productm := repositories.NewProductManager("product",db)
	//productArray,  _ := productm.SelectByKey(1)
	//fmt.Println(productArray)
	order := repositories.NewOrderManagerRepository("orders",db)
	orderArray , _ := order.SelectAll()

	for i , v := range orderArray{
		fmt.Println(i,v)
	}
}