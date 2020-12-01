package datamodels

type Product struct {
	//商品id
	ID int64 `json:"id" sql:"id" imooc:"ID"`
	//商品名称
	ProductName string `json:"ProductName" sql:"productName" imooc:"ProductName"`
	//商品数量
	ProductNum int64 `json:"ProductNum" sql:"productNum" imooc:"ProductNum"`
	//商品图片
	ProductImage string `json:"ProductImage" sql:"productImage" imooc:"ProductImage"`
	//商品地址
	ProductUrl string `json:"ProductUrl" sql:"productUrl" imooc:"ProductUrl"`
}
