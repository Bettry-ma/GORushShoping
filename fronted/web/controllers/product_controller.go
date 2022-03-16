package controllers

import (
	"GORushShoping/datamodels"
	"GORushShoping/rabbitmq"
	"GORushShoping/services"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"os"
	"path/filepath"
	"strconv"
	"text/template"
)

type ProductController struct {
	Ctx            iris.Context
	ProductService services.IProductService
	OrderService   services.IOrderService
	RabbitMQ       *rabbitmq.RabbitMQ
	Sessions       *sessions.Session
}

var (
	//生成HTML保存目录
	htmlOutPath = "./fronted/web/htmlProductShow/"
	//静态文件模板目录
	templatePath = "./fronted/web/views/template/"
)

func (p *ProductController) GetGenerateHtml() {
	productIDString := p.Ctx.URLParam("productID")
	productID, err := strconv.Atoi(productIDString)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	//1.获取模板
	contentTmp, err := template.ParseFiles(filepath.Join(templatePath, "product.html"))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	//2.获取HTML生成路径并指定生成文件的名称
	fileName := filepath.Join(htmlOutPath, "htmlProduct.html")
	//3.获取模板的渲染数据
	product, err := p.ProductService.GetProductById(int64(productID))
	if err != nil {
		fmt.Println("获取模板数据错误")
		p.Ctx.Application().Logger().Debug(err)
	}
	//4.生成静态文件
	generateStaticHtml(p.Ctx, contentTmp, fileName, product)
}

//生成静态文件
func generateStaticHtml(ctx iris.Context, tem *template.Template, fileName string, product *datamodels.Product) {
	//1.判断静态文件是否存在
	if exist(fileName) {
		err := os.Remove(fileName)
		if err != nil {
			ctx.Application().Logger().Debug(err)
		}
	}
	//2.生成静态文件
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		ctx.Application().Logger().Debug(err)
	}
	defer file.Close()
	err = tem.Execute(file, &product)
	if err != nil {
		ctx.Application().Logger().Debug(err)
	}
}

//判断文件是否存在
func exist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

func (p *ProductController) GetDetail() mvc.View {
	product, err := p.ProductService.GetProductById(1)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/view.html",
		Data: iris.Map{
			"product": product,
		},
	}
}

/*
基础功能架构缺陷:
	1.web服务器压力大,安全验证成本高
	2.实时读取数据对数据库会造成较大压力
	3.在高并发下无法保障数据一致性
*/

func (p *ProductController) GetOrder() []byte {
	/*productIDString := p.Ctx.URLParam("productID")
	uidString := p.Ctx.GetCookie("uid")
	productID, err := strconv.ParseInt(productIDString, 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	userID, err := strconv.ParseInt(uidString, 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}*/
	//创建消息体
	//message := datamodels.NewMessage(productID, userID)
	message := datamodels.NewMessage(1, 2)
	//类型转换
	byteMessage, err := json.Marshal(message)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	err = p.RabbitMQ.PublishSimple(string(byteMessage))
	if err != nil {
		fmt.Println(err)
	}
	return []byte("true")
	/*//判断商品数量是否满足需求
	showMessage := "抢购失败! "
	var orderID int64
	if product.ProductNum > 0 {
		//扣除商品数量
		product.ProductNum -= 1
		err := p.ProductService.UpdateProduct(product)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		}
		//创建订单
		uid, err := strconv.Atoi(uidString)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		}
		order := &datamodels.Order{
			UserId:      int64(uid),
			ProductId:   int64(productID),
			OrderStatus: datamodels.OrderSuccess,
		}
		orderID, err = p.OrderService.InsertOrder(order)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		} else {
			showMessage = "抢购成功! "
		}
	}
	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/result.html",
		Data: iris.Map{
			"showMessage": showMessage,
			"orderID":     orderID,
		},
	}*/
}

/*
分布式架构设计:
思想:筛选有效流量,异步处理数据
	1.静态资源放置CDN处理
	2.抢购数据的真实流量放到SLB(流量负载均衡器)中
	3.SLB后增加流量拦截,提供分布式安全验证
	4.秒杀数量控制,防止超卖,增加性能
	5.web服务
	6.RabbitMQ防止爆库
	7.MySQL
*/
