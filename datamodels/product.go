package datamodels

type Product struct {
	ID              int64   `json:"id" sql:"ID" form:"ID"`                                        //商品ID
	ProductName     string  `json:"ProductName" sql:"productName" form:"ProductName"`             //商品名称
	ProductNum      int64   `json:"ProductNum" sql:"productNum" form:"ProductNum"`                //商品数量
	ProductImage    string  `json:"ProductImage" sql:"productImage" form:"ProductImage"`          //商品图片
	ProductUrl      string  `json:"ProductUrl" sql:"productUrl" form:"ProductUrl"`                //商品地址
	ProductPriceNow float64 `json:"ProductPriceNow" sql:"productPriceNow" form:"ProductPriceNow"` //商品价格(现)
	ProductPricePre float64 `json:"ProductPricePre" sql:"productPricePre" form:"ProductPricePre"` //商品价格(之前的)
}
