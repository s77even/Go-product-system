package datamodels

//Message 模拟一个简单的消息体
type Message struct {
	ProductID int64
	UserID    int64
}

func NewMassage(userID , productID int64) *Message {
	return &Message{
		ProductID: productID,
		UserID: userID,
	}
}
