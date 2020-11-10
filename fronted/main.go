package main

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"imooc-product/common"
	"imooc-product/fronted/web/controllers"
	"imooc-product/repositories"
	"imooc-product/services"
	"log"
	"time"
)

func main(){
	//创建iris 实例
	app := iris.New()
	// 设置错误模式 在mvc下提示错误
	app.Logger().SetLevel("debug")
	//注册模板
	template := iris.HTML("./web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)
	//设置模板目标
	app.HandleDir("/public", "./web/public") //替代了staticweb
	//异常跳转
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})
	db, err := common.NewMysqlConn()
	if err != nil {
		log.Fatal(err)
	}
	ctx ,cancel := context.WithCancel(context.Background())
	defer cancel()
	//注册控制器
	sess := sessions.New(sessions.Config{
		Cookie: "helloworld",
		Expires: 60*time.Second,
	})
	userRepository := repositories.NewUserManagerRepository("user",db)
	userService := services.NewUserService(userRepository)
	user := mvc.New(app.Party("/user"))
	user.Register(userService,ctx,sess.Start)
	user.Handle(new(controllers.UserController))

	// 启动服务
	app.Run(
		iris.Addr("localhost:8082"),
		iris.WithoutServerError(iris.ErrServerClosed), // 忽略iris框架的错误
		iris.WithOptimizations,
	)
}
