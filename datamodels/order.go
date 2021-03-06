package datamodels

type Order struct {
	ID          int64 `sql:"ID"`
	UserId      int64 `sql:"userID"`
	ProductId   int64 `sql:"productID"`
	OrderStatus int64 `sql:"orderStatus"`
	ProductNum  int64 `sql:"productNum"`
}

//状态值
const (
	OrderWait    = iota //0
	OrderSuccess        //1
	OrderFailed         //2
)
