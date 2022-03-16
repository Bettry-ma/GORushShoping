package controllers

import (
	"GORushShoping/datamodels"
	"GORushShoping/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"strconv"
)

type OrderController struct {
	Ctx          iris.Context
	OrderService services.IOrderService
}

func (o *OrderController) GetAll() mvc.View {
	orderArray, err := o.OrderService.GetAllOrderInfo()
	if err != nil {
		o.Ctx.Application().Logger().Debug("订单查询失败 ", err)
	}
	return mvc.View{
		Name: "order/view.html",
		Data: iris.Map{
			"order": orderArray,
		},
	}
}

func (o *OrderController) GetMenu() mvc.View {
	orders, err := o.OrderService.GetAllOrder()
	if err != nil {
		o.Ctx.Application().Logger().Debug("获取订单信息失败 ", err)
	}
	return mvc.View{
		Name: "order/menu.html",
		Data: iris.Map{
			"orders": orders,
		},
	}
}

func (o *OrderController) GetSearch() mvc.View {
	idString := o.Ctx.URLParam("sid")
	var id int64
	if idString == "" {
		o.Ctx.Redirect("menu")
	} else {
		var err error
		id, err = strconv.ParseInt(idString, 10, 64)
		if err != nil {
			o.Ctx.Application().Logger().Debug(err)
		}
	}
	order, err := o.OrderService.GetOrderByID(id)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	nl := datamodels.Order{}
	if *order == nl {
		o.Ctx.Redirect("searchempty")
	}
	return mvc.View{
		Name: "order/search.html",
		Data: iris.Map{
			"order": order,
		},
	}
}

func (o *OrderController) GetSearchempty() mvc.View {
	return mvc.View{
		Name: "order/searchempty.html",
	}
}

func (o *OrderController) GetDelete() {
	idString := o.Ctx.URLParam("ID")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}
	status := o.OrderService.DeleteOrderByID(id)
	if status {
		o.Ctx.Redirect("menu")
	} else {
		o.Ctx.Application().Logger().Debug("delete product failed! where id =", id)
	}
}
