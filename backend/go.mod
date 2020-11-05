module imooc-product/backend

go 1.15

require (
	github.com/kataras/iris/v12 v12.1.8
	imooc-product/backend/web/controllers v0.0.0-00010101000000-000000000000
	imooc-product/common v0.0.0-00010101000000-000000000000
	imooc-product/repositories v0.0.0-00010101000000-000000000000
	imooc-product/services v0.0.0-00010101000000-000000000000
)

replace imooc-product/services => D:\GOPRO\src\imooc-product\services

replace imooc-product/common => D:\GOPRO\src\imooc-product\common

replace imooc-product/datamodels => D:\GOPRO\src\imooc-product\datamodels

replace imooc-product/repositories => D:\GOPRO\src\imooc-product\repositories

replace imooc-product/backend/web/controllers => D:\GOPRO\src\imooc-product\backend\web\controllers
