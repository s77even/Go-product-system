module imooc-product/services

go 1.15

replace imooc-product/datamodels => D:\GOPRO\src\imooc-product\datamodels

require (
	golang.org/x/crypto v0.0.0-20201016220609-9e8e0b390897 // indirect
	imooc-product/datamodels v0.0.0-00010101000000-000000000000
	imooc-product/repositories v0.0.0-00010101000000-000000000000
)

replace imooc-product/repositories => D:\GOPRO\src\imooc-product\repositories

replace imooc-product/common => D:\GOPRO\src\imooc-product\common
