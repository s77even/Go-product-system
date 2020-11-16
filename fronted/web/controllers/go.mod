module imooc-product/fronted/web/controllers

go 1.15

require (
	github.com/kataras/golog v0.0.18
	github.com/kataras/iris v0.0.2
	github.com/kataras/iris/v12 v12.1.8
	github.com/nats-io/nats-server/v2 v2.1.9 // indirect
	imooc-product/common v0.0.0-00010101000000-000000000000
	imooc-product/datamodels v0.0.0-00010101000000-000000000000
	imooc-product/encrypt v0.0.0-00010101000000-000000000000
	imooc-product/services v0.0.0-00010101000000-000000000000
)

replace imooc-product/services => D:\GOPRO\src\imooc-product\services

replace imooc-product/repositories => D:\GOPRO\src\imooc-product\repositories

replace imooc-product/common => D:\GOPRO\src\imooc-product\common

replace imooc-product/datamodels => D:\GOPRO\src\imooc-product\datamodels

replace imooc-product/encrypt => D:\GOPRO\src\imooc-product\encrypt
