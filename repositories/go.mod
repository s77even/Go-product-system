module imooc-product/repositories

go 1.15

replace imooc-product/common => D:\GOPRO\src\imooc-product\common

replace imooc-product/datamodels => D:\GOPRO\src\imooc-product\datamodels

require (
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	imooc-product/common v0.0.0-00010101000000-000000000000
	imooc-product/datamodels v0.0.0-00010101000000-000000000000
)
