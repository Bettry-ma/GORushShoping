package services

import (
	"GORushShoping/datamodels"
	"GORushShoping/repositories"
)

type IOrderService interface {
	GetOrderByID(id int64) (order *datamodels.Order, err error)
	DeleteOrderByID(id int64) (status bool)
	UpdateOrder(order *datamodels.Order) (err error)
	InsertOrder(order *datamodels.Order) (affected int64, err error)
	GetAllOrder() (orders []*datamodels.Order, err error)
	GetAllOrderInfo() (info map[int]map[string]string, err error)
	InsertOrderByMessage(message *datamodels.Message) (int64, error)
}

type OrderService struct {
	OrderRepository repositories.IOrderRepository
}

func NewOrderService(op repositories.IOrderRepository) IOrderService {
	return &OrderService{op}
}

func (o *OrderService) GetOrderByID(id int64) (order *datamodels.Order, err error) {
	return o.OrderRepository.SelectById(id)
}

func (o *OrderService) DeleteOrderByID(id int64) (status bool) {
	return o.OrderRepository.Delete(id)
}

func (o *OrderService) UpdateOrder(order *datamodels.Order) (err error) {
	return o.OrderRepository.Update(order)
}

func (o *OrderService) InsertOrder(order *datamodels.Order) (affected int64, err error) {
	return o.OrderRepository.Insert(order)
}

func (o *OrderService) GetAllOrder() (orders []*datamodels.Order, err error) {
	return o.OrderRepository.SelectAll()
}

func (o *OrderService) GetAllOrderInfo() (info map[int]map[string]string, err error) {
	return o.OrderRepository.SelectAllWithInfo()
}

// InsertOrderByMessage 根据消息创建订单
func (o *OrderService) InsertOrderByMessage(message *datamodels.Message) (orderID int64, err error) {
	order := &datamodels.Order{
		UserId:      message.UserID,
		ProductId:   message.ProductID,
		OrderStatus: datamodels.OrderSuccess,
	}
	return o.InsertOrder(order)
}
