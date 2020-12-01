module imooc-product/rabbitmq

go 1.15

require (
	github.com/kataras/golog v0.1.5
	github.com/streadway/amqp v1.0.0
	golang.org/x/sys v0.0.0-20201126233918-771906719818 // indirect
	imooc-product/datamodels v0.0.0-00010101000000-000000000000
	imooc-product/services v0.0.0-00010101000000-000000000000
)

replace imooc-product/services => D:\GOPRO\src\imooc-product\services

replace imooc-product/datamodels => D:\GOPRO\src\imooc-product\datamodels

replace imooc-product/repositories => D:\GOPRO\src\imooc-product\repositories

replace imooc-product/common => D:\GOPRO\src\imooc-product\common
