package router

import (
	"github.com/astaxie/beego"
	"github.com/jameshih/secKill/admin/controller/event"
	"github.com/jameshih/secKill/admin/controller/product"
)

func init() {
	beego.Router("/", &product.ProductController{}, "*:ListProduct")
	beego.Router("/product/list", &product.ProductController{}, "*:ListProduct")
	beego.Router("/product/new", &product.ProductController{}, "*:NewProduct")
	beego.Router("/product/create", &product.ProductController{}, "*:CreateProduct")
	beego.Router("/event/list", &event.EventController{}, "*:ListEvent")
	beego.Router("/event/new", &event.EventController{}, "*:NewEvent")
	beego.Router("/event/create", &event.EventController{}, "*:CreateEvent")
}
