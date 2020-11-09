package datamodels

// Order 订单数据结构体
type Order struct {
	ID          int64 `sql:"ID"`
	UserId      int64 `sql:"userID"`
	ProductId   int64 `sql:"productID"`
	OrderStatus int64 `sql:"orderStatus"`
}


const (
	OrderWait    = iota
	OrderSuccess //成功
	OrderFailed  //失败
)