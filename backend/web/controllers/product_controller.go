package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"imooc-product/common"
	"imooc-product/datamodels"
	"imooc-product/services"
	"strconv"
)

type ProductController struct {
	Ctx iris.Context
	ProductService services.IProductService
}

func (p *ProductController)GetAll() mvc.View{
	productArray , err := p.ProductService.GetAllProduct()
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	return mvc.View{
		Name: "product/view.html",
		Data: iris.Map{
			"productArray":productArray,
		},
	}
}

func (p *ProductController) PostUpdate() {
	product := &datamodels.Product{}
	p.Ctx.Request().ParseForm()
	dec := common.NewDecoder(&common.DecoderOptions{TagName: "imooc"})
	if err := dec.Decode(p.Ctx.Request().Form, product); err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	err := p.ProductService.UpdateProduct(product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	p.Ctx.Redirect("/product/all")
}

func (p *ProductController)GetAdd()mvc.View{
	return mvc.View{
		Name: "product/add.html",
	}
}

func (p *ProductController) PostAdd() {
	product := &datamodels.Product{}
	p.Ctx.Request().ParseForm()
	dec := common.NewDecoder(&common.DecoderOptions{TagName: "imooc"})
	if err := dec.Decode(p.Ctx.Request().Form, product); err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	_, err := p.ProductService.InsertProduct(product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	p.Ctx.Redirect("/product/all")
}

func (p * ProductController) GetManager()mvc.View{
	idString := p.Ctx.URLParam("id")
	id , err := strconv.ParseInt(idString,10,16)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	product, err := p.ProductService.GetProductByID(id)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	return mvc.View{
		Name:   "product/manager.html",
		Data:   iris.Map{
			"product":product,
		},
	}
}

func (p *ProductController)GetDelete(){
	idString := p.Ctx.URLParam("id")
	id , err := strconv.ParseInt(idString,10,16)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	isOk:=p.ProductService.DeleteProductByID(id)
	if isOk{
		p.Ctx.Application().Logger().Debug("delete product successful , id is "+idString)
	} else{
		p.Ctx.Application().Logger().Debug("delete product failed , id is "+idString)
	}
	p.Ctx.Redirect("/product/all")
}