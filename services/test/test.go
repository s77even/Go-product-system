package main

import (
	"fmt"
	"imooc-product/common"
	"imooc-product/repositories"
	"imooc-product/services"
)

func main(){
	db, _ :=  common.NewMysqlConn()
	//productm := repositories.NewProductManager("product",db)
	//productArray,  _ := productm.SelectByKey(1)
	//fmt.Println(productArray)
	order := repositories.NewOrderManagerRepository("orders",db)
	oredeService := services.NewOrderService(order)
	result , err:= oredeService.GetAllOrder()
	if err != nil {
		fmt.Println(err)
	}
	for _ , v := range result{
		fmt.Println(v)
	}
}