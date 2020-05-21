package main

import (
	"fmt"

	"github.com/astaxie/beego"
	_ "github.com/jameshih/secKill/admin/router"
)

func main() {
	err := initAll()
	if err != nil {
		panic(fmt.Sprintf("init failed, error: %v", err))
	}
	beego.Run()
}
