package product

import (
	"fmt"

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

func (p *ProductController) CreateProduct() {
	var err error
	defer func() {
		if err != nil {
			p.Data["Error"] = err.Error()
			p.Layout = "layout/layout.html"
			p.TplName = "layout/error.html"
		}
	}()

	productName := p.GetString("name")
	if len(productName) == 0 {
		err = fmt.Errorf("invalid product name")
		logs.Warn(err)
		return
	}

	productQTY, err := p.GetInt("qty")
	if err != nil {
		err = fmt.Errorf("invalid product quantity, error: %v", err)
		logs.Warn(err)
		return
	}

	productStatus, err := p.GetInt("status")
	if err != nil {
		err = fmt.Errorf("invalid product status, error: %v", err)
		logs.Warn(err)
		return
	}

	productModel := model.NewProductModel()
	product := model.Product{
		ProductName: productName,
		Total:       productQTY,
		Status:      productStatus,
	}
	err = productModel.CreateProduct(&product)
	if err != nil {
		err = fmt.Errorf("failed to submit, error: %v", err)
		logs.Warn(err)
		return
	}

	logs.Debug("product name[%s], product qty[%d], product status[%d]", product.ProductName, product.Total, product.Status)
	p.Redirect("/product/list", 302)
}
