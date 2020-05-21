package product

import (
	"github.com/astaxie/beego"
)

type ProductController struct {
	beego.Controller
}

func (p *ProductController) ListProduct() {
	// productModel := model.NewProductModel()
	p.Layout = "layout/layout.html"
	p.TplName = "product/list.html"
}
