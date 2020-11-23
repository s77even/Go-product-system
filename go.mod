module imooc-product

go 1.15

replace imooc-product/common => D:\GOPRO\src\imooc-product\common

require (
	imooc-product/common v0.0.0-00010101000000-000000000000
	imooc-product/encrypt v0.0.0-00010101000000-000000000000 // indirect
)

replace imooc-product/encrypt => D:\GOPRO\src\imooc-product\encrypt
