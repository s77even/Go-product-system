package main

import (
	"fmt"
	"imooc-product/common"
	"imooc-product/rabbitmq"
	"imooc-product/repositories"
	"imooc-product/services"
)

func main(){
	db, err := common.NewMysqlConn()
	if err != nil {
		fmt.Println(err)
	}
	product := repositories.NewProductManagerRepostory("product",db)
	productService := services.NewProductService(product)

	order := repositories.NewOrderManagerRepository("orders",db)
	orderService := services.NewOrderService(order)

	rabbitmqConsumerSimple := rabbitmq.NewSimpleRabbitMQ("imoocProduct")
	rabbitmqConsumerSimple.ConsumeSimple(orderService,productService)
}
