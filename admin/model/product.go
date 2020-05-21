package model

import (
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type ProductModel struct {
	Db *sqlx.DB
}

type Product struct {
	ProductID   int    `db:"id"`
	ProductName string `db:"name"`
	Total       int    `db:"total"`
	Status      int    `db:"status"`
}

func NewProductModel(db *sqlx.DB) (productModel *ProductModel) {
	productModel = &ProductModel{
		Db: db}
	return
}

func (p *ProductModel) GetProductList() (productList []*Product, err error) {
	sql := "SELECT id, name, total, status FROM  produt"
	err = p.Db.Select(productList, sql)
	if err != nil {
		logs.Warn("SELECT * FROM mysql failed, error: %v, sql: %v", err, sql)
		return
	}
	return
}
