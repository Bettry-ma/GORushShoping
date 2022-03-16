package controllers

import (
	"GORushShoping/common"
	"GORushShoping/datamodels"
	"GORushShoping/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"strconv"
)

type ProductController struct {
	Ctx            iris.Context             //请求上下文
	ProductService services.IProductService //接口类型
}

//定义动作

// GetAll 用Get方法展示全部商品
func (p *ProductController) GetAll() mvc.View { //这里的Get是iris框架的方法标识,说明使用Get方式,真正名称为All
	productArray, _ := p.ProductService.GetAllProduct()
	return mvc.View{
		Name: "product/view.html", //渲染的模板
		Data: iris.Map{
			"productArray": productArray, //前是模板的名称,后是该方法的变量名称
		},
	}
}

// PostUpdate 修改商品
func (p *ProductController) PostUpdate() {
	product := &datamodels.Product{}
	p.Ctx.Request().ParseForm()
	dec := common.NewDecoder(&common.DecoderOptions{TagName: "form"}) //调用库来对表单进行处理
	if err := dec.Decode(p.Ctx.Request().Form, product); err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	err := p.ProductService.UpdateProduct(product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	p.Ctx.Redirect("/product/all") //this is GetAll method
}

// GetAdd 添加商品HTML页面跳转
func (p *ProductController) GetAdd() mvc.View {
	return mvc.View{
		Name: "product/add.html",
	}
}

// PostAdd 将添加的商品信息表单提交
func (p *ProductController) PostAdd() {
	product := &datamodels.Product{}
	p.Ctx.Request().ParseForm()
	dec := common.NewDecoder(&common.DecoderOptions{TagName: "form"}) //调用库来对表单进行处理
	if err := dec.Decode(p.Ctx.Request().Form, product); err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	_, err := p.ProductService.InsertProduct(product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	p.Ctx.Redirect("/product/all") //this is GetAll method
}

// GetManager 修改商品属性页面
func (p *ProductController) GetManager() mvc.View {
	idString := p.Ctx.URLParam("ID")              //通过URL获取到该行的ID
	id, err := strconv.ParseInt(idString, 10, 64) //将字符串id转换为十进制id,数据类型为int64
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	//fmt.Println("已经运行到此处")
	product, err := p.ProductService.GetProductById(id)
	if err != nil {
		//fmt.Println("这里打印的是错误信息")
		p.Ctx.Application().Logger().Debug(err)
	}
	return mvc.View{
		Name: "product/manager.html",
		Data: iris.Map{
			"product": product,
		},
	}
}

// GetDelete 删除商品
func (p *ProductController) GetDelete() {
	idString := p.Ctx.URLParam("ID")              //通过当前view页面获取到该行的ID
	id, err := strconv.ParseInt(idString, 10, 64) //将字符串id转换为十进制id,数据类型为int64
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	state := p.ProductService.DeleteProductById(id)
	if state {
		p.Ctx.Redirect("/product/all")
	} else {
		p.Ctx.Application().Logger().Debug("delete product failed! where id =", id)
	}
}
