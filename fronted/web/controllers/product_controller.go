package controllers

import (
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"html/template"
	"imooc-product/datamodels"
	"imooc-product/services"
	"os"
	"path/filepath"
	"strconv"
)

type ProductController struct {
	Ctx iris.Context
	ProductService services.IProductService
	OrderService services.IOrderService
	Session *sessions.Session
}
var (
	htmlOutPath = "D:/GOPRO/src/imooc-product/fronted/web/htmlProductShow"
	templatePath = "D:/GOPRO/src/imooc-product/fronted/web/views/template"
)
func (p *ProductController) GetGenerateHtml() {
	productString := p.Ctx.URLParam("productID")
	productID , err :=strconv.Atoi(productString)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	//1. 获取模版文件地址
	contentTmp, err := template.ParseFiles(filepath.Join(templatePath, "product.html"))

	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
		return
	}
	//2. 获取html生成路径
	fileName := filepath.Join(htmlOutPath, "htmlProduct.html")
	//3. 获取模版渲染数据
	product, err := p.ProductService.GetProductByID(int64(productID))
	//fmt.Println(product.ProductNum, product.ProductUrl, productID)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
		return
	}
	//4. 生成静态文件
	generateStaticHtml(p.Ctx, contentTmp, fileName, product)
}

// 生成html静态文件
func generateStaticHtml(ctx iris.Context, template *template.Template, fileName string, product *datamodels.Product) {
	// 1.判断静态文件是否存在
	if exist(fileName) {
		err := os.Remove(fileName)
		if err != nil {
			ctx.Application().Logger().Error(err)
			return
		}
	}
	//2.生成静态文件
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		ctx.Application().Logger().Error(err)
		return
	}
	defer file.Close()
	//fmt.Println(file)
	err = template.Execute(file, &product)
	if err != nil {
		ctx.Application().Logger().Error(err)
	}
}
func exist(fileName string)bool{
	_ ,err := os.Stat(fileName)
	return err== nil || os.IsExist(err)
}

func (p *ProductController)GetDetail() mvc.View{
	product,err := p.ProductService.GetProductByID(1)
	if err != nil {
		p.Ctx.Application().Logger().Error(err)
	}
	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/view.html",
		Data: iris.Map{
			"product" : product,
		},
	}
}

func (p *ProductController)GetOrder()mvc.View	{
	productString := p.Ctx.URLParam("productID")
	userString := p.Ctx.GetCookie("uid")
	productID , err := strconv.Atoi(productString)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	product, err := p.ProductService.GetProductByID(int64(productID))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	var orderID int64
	showMessage := "抢购失败"
	if product.ProductNum >0 {
		//扣除商品数量
		product.ProductNum -= 1
		err := p.ProductService.UpdateProduct(product) //会出现超卖
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		}
		//
		userID ,err := strconv.Atoi(userString)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		}
		order := &datamodels.Order{
			UserId: int64(userID),
			ProductId: int64(productID),
			OrderStatus: datamodels.OrderSuccess,
		}
		orderID , err =p.OrderService.InsertOrder(order)//新建订单
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		}else {
			showMessage = "抢购成功！"
		}
	}
	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/result.html",
		Data: iris.Map{
			"orderID":orderID,
			"showMessage":showMessage,
		},
	}
}