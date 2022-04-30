package datamodels

//简单的订单消息体

type Message struct {
	ProductID int64
	UserID    int64
}

// NewMessage 构造函数
func NewMessage(productID int64, userID int64) *Message {
	return &Message{productID, userID}
}
