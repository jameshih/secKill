package model

import (
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
)

type ProductModel struct {
}

type Product struct {
	ProductID   int    `db:"id"`
	ProductName string `db:"name"`
	Total       int    `db:"total"`
	Status      int    `db:"status"`
}

func NewProductModel() (productModel *ProductModel) {
	productModel = &ProductModel{}
	return
}

func (p *ProductModel) GetProductList() (productList []*Product, err error) {
	sql := "SELECT id, name, total, status FROM  product"
	err = Db.Select(&productList, sql)
	if err != nil {
		logs.Warn("SELECT * FROM mysql failed, error: %v, sql: %v", err, sql)
		return
	}
	return
}

// check duplicate in database

func (p *ProductModel) CreateProduct(product *Product) (err error) {
	sql := "INSERT INTO product(name, total, status)values(?,?,?)"
	_, err = Db.Exec(sql, product.ProductName, product.Total, product.Status)
	if err != nil {
		logs.Warn("INSERT INTO product failed, error: %v, sql: %v", err, sql)
		return
	}
	return
}
