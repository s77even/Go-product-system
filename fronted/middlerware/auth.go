package middlerware

import "github.com/kataras/iris/v12"

func AuthConProduct(ctx iris.Context){
	uid := ctx.GetCookie("uid")
	if uid == ""{
		ctx.Application().Logger().Debug("请先登录")
		ctx.Redirect("/user/login")
		return
	}
	ctx.Application().Logger().Debug("已经登录")
	ctx.Next()
}