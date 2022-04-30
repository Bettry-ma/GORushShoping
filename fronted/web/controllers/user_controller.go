package controllers

import (
	"GORushShoping/common"
	"GORushShoping/datamodels"
	"GORushShoping/encrypt"
	"GORushShoping/services"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"strconv"
)

type UserController struct {
	Ctx         iris.Context
	UserService services.IUserService
	Session     *sessions.Session
}

func (u *UserController) GetRegister() mvc.View {
	return mvc.View{
		Name: "user/register.html",
	}
}

// PostRegister  用户注册
func (u *UserController) PostRegister() {
	var (
		nickName = u.Ctx.FormValue("nickName")
		userName = u.Ctx.FormValue("userName")
		passWord = u.Ctx.FormValue("password")
	)
	user := &datamodels.User{
		UserName:     userName,
		NickName:     nickName,
		HashPassword: passWord,
	}
	_, err := u.UserService.AddUser(user)
	if err != nil {
		u.Ctx.Redirect("error")
		fmt.Println("注册用户失败", err)
		return
	}
	u.Ctx.Redirect("login")
	return
}

// GetLogin 用户登录界面
func (u *UserController) GetLogin() mvc.View {
	return mvc.View{
		Name: "user/login.html",
	}
}

func (u *UserController) PostLogin() mvc.Response {
	//获取用户提交的表单信息
	var (
		userName = u.Ctx.FormValue("userName")
		password = u.Ctx.FormValue("password")
	)
	//验证账号密码是否正确
	user, ok := u.UserService.IsPwdSuccess(userName, password)
	if !ok {
		return mvc.Response{
			Path: "login",
		}
	}
	//写入用户ID到Cookie中
	common.GlobalCookie(u.Ctx, "uid", strconv.FormatInt(user.ID, 10))
	uidByte := []byte(strconv.FormatInt(user.ID, 10))
	uidString, err := encrypt.EnPwdCode(uidByte)
	if err != nil {
		u.Ctx.Application().Logger().Debug(err)
	}
	//写入用户浏览器
	common.GlobalCookie(u.Ctx, "sign", uidString)
	//u.Session.Set("userID", strconv.FormatInt(user.ID, 10))
	return mvc.Response{
		Path: "/product/detail",
	}
}

/*
Cookie和Session的区别:
	Cookie:采用的是在客户端保存状态和数据的方案,Cookie数据存放在客户的浏览器上,单个Cookie保存的数据不能超过4KB
	Session:采用的是在服务器端保存状态和数据的方案,Session数据存放在服务器上,当访问量增多,会比较占用大量资源
*/
/*
什么是分布式?
	分布式系统是相对于集中式系统说的概念;
	集中式系统:所有的应用程序和组件放在同台机器上运行;
	分布式定义:分布式系统是若干独立计算机的集合,这些计算机对于用户来说就像是单个系统

分布式更直观的感受
	1.分布式系统是由多台机器组成
	2.分布式系统是一个整体,从外部感觉不到多台机器的存在;
*/
