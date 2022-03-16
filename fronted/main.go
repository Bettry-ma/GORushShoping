package main

import (
	"GORushShoping/common"
	"GORushShoping/fronted/web/controllers"
	"GORushShoping/rabbitmq"
	"GORushShoping/repositories"
	"GORushShoping/services"
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func main() {
	//1.创建iris实例
	app := iris.New()
	//2.设置日志错误模式,在MVC模式下提示错误
	app.Logger().SetLevel("debug")
	//3.注册模板
	htmlEngine := iris.HTML("./fronted/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(htmlEngine)

	//4.设置模板目录
	app.HandleDir("/public", "./fronted/web/public")
	app.HandleDir("/html", "./fronted/web/htmlProductShow")
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
	/*sess := sessions.New(sessions.Config{
		Cookie:  "AdminCookie",
		Expires: 60 * time.Microsecond,
	})*/
	//连接数据库
	db, err := common.NewMysqlConn()
	if err != nil {
		app.Logger().Error(err)
	}
	//注册控制器
	//user
	userRepository := repositories.NewUserManager("user", db)
	userService := services.NewUserService(userRepository)
	UserParty := app.Party("user")
	user := mvc.New(UserParty)
	user.Register(userService, ctx)
	user.Handle(new(controllers.UserController))

	rabbitmq := rabbitmq.NewRabbitMQSimple("product")
	//product
	productRepository := repositories.NewProductManager("product", db)
	productService := services.NewProductService(productRepository)
	//order 嵌入Product中完成业务
	orderRepository := repositories.NewOrderManager("order", db)
	orderService := services.NewOrderService(orderRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	//productParty.Use(middleware.AuthConProduct) //验证用户是否登录
	product.Register(productService, orderService, ctx, rabbitmq)
	product.Handle(new(controllers.ProductController))

	err = app.Run(
		iris.Addr("localhost:8082"),                   //访问地址的端口号与backend不同
		iris.WithoutServerError(iris.ErrServerClosed), //忽略iris的关闭错误
		iris.WithOptimizations,
	)
	if err != nil {
		fmt.Println("failed to run app")
		return
	}
}
