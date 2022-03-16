package main

import (
	"GORushShoping/backend/web/controllers"
	"GORushShoping/common"
	"GORushShoping/repositories"
	"GORushShoping/services"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"golang.org/x/net/context"
)

func main() {
	//1.创建iris实例
	app := iris.New()
	//2.设置日志错误模式,在MVC模式下提示错误
	app.Logger().SetLevel("debug")
	//3.注册模板
	htmlEngine := iris.HTML("./backend/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(htmlEngine)

	//4.设置模板目录
	//app.StaticWeb("assets/","/backend/web/assets")
	app.HandleDir("/assets", "./backend/web/assets")
	//app.StaticContent("/assets", "pain/text", nil)
	//出现异常跳转到指定页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错"))
		ctx.ViewLayout("")
		err := ctx.View("shared/error.html")
		if err != nil {
			fmt.Println("error.html set failed")
			return
		}
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//连接数据库
	db, err := common.NewMysqlConn()
	if err != nil {
		app.Logger().Error(err)
	}
	//5.注册控制器
	//product
	productRepository := repositories.NewProductManager("product", db) // repository(携带table和db属性,且实现了IProduct方法 --> service ,
	productService := services.NewProductService(productRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	product.Register(ctx, productService)
	product.Handle(new(controllers.ProductController))

	//order
	orderRepository := repositories.NewOrderManager("order", db)
	orderService := services.NewOrderService(orderRepository)
	orderParty := app.Party("/order")
	order := mvc.New(orderParty)
	order.Register(ctx, orderService)
	order.Handle(new(controllers.OrderController))

	//this is Mac content for test
	//6.启动服务
	err = app.Run(
		iris.Addr("localhost:8080"),
		iris.WithoutServerError(iris.ErrServerClosed), //忽略iris的关闭错误
		iris.WithOptimizations,
	)
	if err != nil {
		fmt.Println("failed to run app")
		return
	}
}
