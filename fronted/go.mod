module imooc-product/fronted

go 1.15

require (
	github.com/kataras/iris v0.0.2 // indirect
	github.com/kataras/iris/v12 v12.1.8
	github.com/nats-io/nats-server/v2 v2.1.9 // indirect
	imooc-product/common v0.0.0-00010101000000-000000000000
	imooc-product/fronted/web/controllers v0.0.0-00010101000000-000000000000
	imooc-product/repositories v0.0.0-00010101000000-000000000000
	imooc-product/services v0.0.0-00010101000000-000000000000
)

replace imooc-product/services => D:\GOPRO\src\imooc-product\services

replace imooc-product/common => D:\GOPRO\src\imooc-product\common

replace imooc-product/datamodels => D:\GOPRO\src\imooc-product\datamodels

replace imooc-product/repositories => D:\GOPRO\src\imooc-product\repositories

replace imooc-product/fronted/web/controllers => D:\GOPRO\src\imooc-product\fronted\web\controllers
