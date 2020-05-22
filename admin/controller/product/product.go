package product

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/jameshih/secKill/admin/model"
)

type ProductController struct {
	beego.Controller
}

func (p *ProductController) ListProduct() {
	productModel := model.NewProductModel()
	productList, err := productModel.GetProductList()
	if err != nil {
		logs.Warn("get product list failed, error: %v", err)
		return
	}

	p.Data["productList"] = productList
	p.Layout = "layout/layout.html"
	p.TplName = "product/list.html"
}

func (p *ProductController) NewProduct() {

	p.Layout = "layout/layout.html"
	p.TplName = "product/form.html"
}
