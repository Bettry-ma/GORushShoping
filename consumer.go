package main

import (
	"GORushShoping/common"
	"GORushShoping/rabbitmq"
	"GORushShoping/repositories"
	"GORushShoping/services"
	"fmt"
)

func main() {
	db, err := common.NewMysqlConn()
	if err != nil {
		fmt.Println(err)
	}
	//创建product数据库操作实例
	product := repositories.NewProductManager("product", db)
	//创建productService
	productService := services.NewProductService(product)
	//创建order数据库实例
	order := repositories.NewOrderManager("order", db)
	//创建orderService
	orderService := services.NewOrderService(order)
	//创建RabbitMQ消费实例
	rabbitmqConsumeSimple := rabbitmq.NewRabbitMQSimple("product")
	rabbitmqConsumeSimple.ConsumeSimple(orderService, productService)
}
