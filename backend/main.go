package main

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"imooc-product/backend/web/controllers"
	"imooc-product/common"
	"imooc-product/repositories"
	"imooc-product/services"
	"log"
)

func main() {
	//创建iris 实例
	app := iris.New()
	// 设置错误模式 在mvc下提示错误
	app.Logger().SetLevel("debug")
	//注册末班
	template := iris.HTML("./web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)
	//设置模板目标
	app.HandleDir("/assets", "./web/assets") //替代了staticweb
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
	productRepository := repositories.NewProductManagerRepostory("product",db)
	productService := services.NewProductService(productRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	product.Register(ctx,productService)
	product.Handle(new(controllers.ProductController))

	orderRepository := repositories.NewOrderManagerRepository("orders", db)
	orderService := services.NewOrderService(orderRepository)
	orderParty := app.Party("/order")
	order := mvc.New(orderParty)
	order.Register(
		ctx,
		orderService,
	)
	order.Handle(new(controllers.OrderController))

	// 启动服务
	app.Run(
		iris.Addr("localhost:8080"),
		iris.WithoutServerError(iris.ErrServerClosed), // 忽略iris框架的错误
		iris.WithOptimizations,
		)

}
