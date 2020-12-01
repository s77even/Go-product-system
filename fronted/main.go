package main

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"imooc-product/common"
	"imooc-product/fronted/middlerware"
	"imooc-product/fronted/web/controllers"
	"imooc-product/rabbitmq"
	"imooc-product/repositories"
	"imooc-product/services"
	"log"
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
	app.HandleDir("public", "./fronted/web/public") //替代了staticweb
	//访问生成的html静态文件
	app.HandleDir("html","./fronted/web/htmlProductShow")
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
	userRepository := repositories.NewUserManagerRepository("user",db)
	userService := services.NewUserService(userRepository)
	user := mvc.New(app.Party("/user"))
	user.Register(userService,ctx)
	user.Handle(new(controllers.UserController))

	rabbitmq := rabbitmq.NewSimpleRabbitMQ("imoocProduct")


	productRepository := repositories.NewProductManagerRepostory("product",db)
	productService := services.NewProductService(productRepository)
	order := repositories.NewOrderManagerRepository("orders" , db)
	orderService := services.NewOrderService(order)
	proProduct:=app.Party("/product")
	product := mvc.New(proProduct)
	proProduct.Use(middlerware.AuthConProduct)
	product.Register(productService,orderService,ctx,rabbitmq)
	product.Handle(new(controllers.ProductController))
	// 启动服务
	app.Run(
		iris.Addr("localhost:8082"),
		iris.WithoutServerError(iris.ErrServerClosed), // 忽略iris框架的错误
		iris.WithOptimizations,
	)
}
