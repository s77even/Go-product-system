package controllers

import (
	"fmt"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"imooc-product/common"
	"imooc-product/datamodels"
	"imooc-product/encrypt"
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
	//获取用户提交的表单信息
	var(
		userName =c.Ctx.FormValue("userName")
		password =c.Ctx.FormValue("password")
		)
	//验证账号密码是否正确 从数据库中  （可改进未使用redis）
	user , isOk := c.Service.IsPwdSuccess(userName,password)
	if !isOk{
		golog.Error("密码错误")
		return mvc.Response{
			Path: "/user/login",
		}
	}
	// 写入用户ID到cookie
	common.GlobalCookie(c.Ctx,"uid",strconv.FormatInt(user.ID,10),0)
	uidByte := []byte(strconv.FormatInt(user.ID,10))
	uidString , err :=encrypt.EnPwdCode(uidByte)
	if err != nil {
		fmt.Println(err)
	}
	//c.Session.Set("userID",strconv.FormatInt(user.ID,10))
	common.GlobalCookie(c.Ctx,"sign",uidString,0)
	return mvc.Response{
		Path: "/product",
	}
}