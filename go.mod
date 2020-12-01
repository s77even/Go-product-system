module imooc-product

go 1.15

replace imooc-product/common => D:\GOPRO\src\imooc-product\common

require (
	imooc-product/common v0.0.0-00010101000000-000000000000
	imooc-product/encrypt v0.0.0-00010101000000-000000000000
	imooc-product/rabbitmq v0.0.0-00010101000000-000000000000
	imooc-product/repositories v0.0.0-00010101000000-000000000000
	imooc-product/services v0.0.0-00010101000000-000000000000
)

replace imooc-product/encrypt => D:\GOPRO\src\imooc-product\encrypt

replace imooc-product/repositories => D:\GOPRO\src\imooc-product\repositories

replace imooc-product/services => D:\GOPRO\src\imooc-product\services

replace imooc-product/rabbitmq => D:\GOPRO\src\imooc-product\rabbitmq

replace imooc-product/datamodels => D:\GOPRO\src\imooc-product\datamodels
