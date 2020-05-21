package router

import (
	"github.com/astaxie/beego"
	"github.com/jameshih/secKill/admin/controller/product"
)

func init() {
	beego.Router("/product/list", &product.ProductController{}, "*:ListProduct")
}
