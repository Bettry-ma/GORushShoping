package middleware

import "github.com/kataras/iris/v12"

//验证用户是否登录,暂时在商品详情页

func AuthConProduct(ctx iris.Context) {
	uid := ctx.GetCookie("uid")
	if uid == "" {
		ctx.Application().Logger().Debug("You must log in before purchasing")
		ctx.Redirect("/user/login")
		return
	}
	ctx.Application().Logger().Debug("logged in")
	ctx.Next()
}
