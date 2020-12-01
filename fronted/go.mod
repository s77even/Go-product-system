module imooc-product/fronted

go 1.15

require (
	github.com/kataras/iris/v12 v12.1.8
	imooc-product/common v0.0.0-00010101000000-000000000000
	imooc-product/fronted/middlerware v0.0.0-00010101000000-000000000000
	imooc-product/fronted/web/controllers v0.0.0-00010101000000-000000000000
	imooc-product/repositories v0.0.0-00010101000000-000000000000
	imooc-product/services v0.0.0-00010101000000-000000000000
)

replace imooc-product/services => D:\GOPRO\src\imooc-product\services
replace imooc-product/rabbitmq => D:\GOPRO\src\imooc-product\rabbitmq
replace imooc-product/common => D:\GOPRO\src\imooc-product\common

replace imooc-product/datamodels => D:\GOPRO\src\imooc-product\datamodels

replace imooc-product/repositories => D:\GOPRO\src\imooc-product\repositories

replace imooc-product/fronted/web/controllers => D:\GOPRO\src\imooc-product\fronted\web\controllers

replace imooc-product/encrypt => D:\GOPRO\src\imooc-product\encrypt

replace imooc-product/fronted/middlerware => D:\GOPRO\src\imooc-product\fronted\middlerware
