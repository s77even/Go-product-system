package main

import (
	"fmt"
	"imooc-product/common"
	"imooc-product/repositories"
)

func main(){
	db, _ :=  common.NewMysqlConn()
	productm := repositories.NewProductManager("product",db)
	productArray,  _ := productm.SelectAll()
	for i , v := range productArray{
		fmt.Println(i,v)
	}
}