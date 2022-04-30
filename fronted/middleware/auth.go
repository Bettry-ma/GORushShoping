package middleware

import "github.com/kataras/iris/v12"

// AuthConProduct 验证用户是否登录,暂时在商品详情页面使用
func AuthConProduct(ctx iris.Context) {
	uid := ctx.GetCookie("uid") //获取用户id
	if uid == "" {              //如果没有登录,则uid为空,跳转到登录页面
		ctx.Application().Logger().Debug("You must log in before purchasing")
		ctx.Redirect("/user/login") //跳转到登录页面
		return
	}
	ctx.Application().Logger().Debug("logged in") //登录日志
	ctx.Next()                                    //继续执行
}
