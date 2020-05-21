package main

import (
	"fmt"
	"log"

	"github.com/gchaincl/dotsql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/swithek/dotsqlx"
)

var (
	db *sqlx.DB
)

func init() {
	database, err := sqlx.Open("mysql", "root:toor@tcp(127.0.0.1:3306)/")
	if err != nil {
		fmt.Printf("failed to connect to mysql, err: %v", err)
		return
	}
	db = database
	return
}

func main() {
	dot, err := dotsql.LoadFromFile("./tables.sql")
	if err != nil {
		panic(err)
	}
	dotx := dotsqlx.Wrap(dot)

	_, err = dotx.Exec(db, "use-app-database")
	if err != nil {
		_, err = dotx.Exec(db, "create-app-database")
		if err != nil {
			panic(err)
		}
	}

	_, err = dotx.Exec(db, "create-product-table")
	if err != nil {
		panic(err)
	}
	log.Println("database creation succ")
}
