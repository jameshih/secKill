package main

import (
	"fmt"

	"github.com/astaxie/beego/logs"
	"github.com/jameshih/secKill/logic/service"
)

func main() {
	// load config file
	err := initConfig("ini", "./conf/logic.conf")
	if err != nil {
		err = fmt.Errorf("init config failed, error: %v", err)
		logs.Error(err)
		panic(err)
	}

	logs.Debug("init config succ, appConfig:%v", appConfig)

	// initialize logger
	err = initLogger()
	if err != nil {
		err = fmt.Errorf("init logger failed, error: %v", err)
		logs.Error(err)
		panic(err)
	}

	logs.Debug("init logger succ")

	// initialize logic
	err = service.InitSecKillLogic(appConfig)
	if err != nil {
		err = fmt.Errorf("init secKill logic failed, error: %v", err)
		logs.Error(err)
		panic(err)
	}

	logs.Debug("init logic succ")

	// run logic
	err = service.Run(appConfig)
	if err != nil {
		err = fmt.Errorf("serve logic failed, error: %v", err)
		logs.Error(err)
		panic(err)
	}

	logs.Debug("run start succ")
	logs.Info("logic service exited")
}
