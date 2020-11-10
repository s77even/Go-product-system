package controllers

import (
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"imooc-product/common"
	"imooc-product/datamodels"
	"imooc-product/services"
	"strconv"
)

type UserController struct {
	Ctx iris.Context
	Service services.IUserService
	Session *sessions.Session
}
func(c *UserController)GetRegister()mvc.View{
	return mvc.View{
		Name: "user/register.html",
	}
}

func(c *UserController)PostRegister(){
	var (
		nickName =c.Ctx.FormValue("nickName")
		userName =c.Ctx.FormValue("userName")
		password =c.Ctx.FormValue("password")
	)
	user := &datamodels.User{
		UserName: userName,
		NickName: nickName,
		HashPassword: password,
	}
	_ ,err := c.Service.AddUser(user)
	if err != nil{
		c.Ctx.Redirect("/user/error")
	}
	c.Ctx.Redirect("/user/login")
	return
}
func (c *UserController)GetLogin() mvc.View {
return mvc.View{
	Name: "user/login.html",
}
}
func (c *UserController)PostLogin() mvc.Response{
	var(
		userName =c.Ctx.FormValue("userName")
		password =c.Ctx.FormValue("password")
		)
	user , isOk := c.Service.IsPwdSuccess(userName,password)
	if !isOk{
		golog.Error("密码错误")
		return mvc.Response{
			Path: "/user/login",
		}
	}

	common.GlobalCookie(c.Ctx,"uid",strconv.FormatInt(user.ID,10),0)
	c.Session.Set("userID",strconv.FormatInt(user.ID,10))
	return mvc.Response{
		Path: "/product",
	}
}