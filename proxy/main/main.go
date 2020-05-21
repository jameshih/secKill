package main

import (
	"github.com/astaxie/beego"
	_ "github.com/jameshih/secKill/proxy/router"
)

func main() {
	err := initConfig()
	if err != nil {
		panic(err)
	}

	err = initSec()
	if err != nil {
		panic(err)
	}

	beego.Run()
}
