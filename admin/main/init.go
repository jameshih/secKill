package main

import (
	"fmt"

	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jameshih/secKill/admin/model"
	"github.com/jmoiron/sqlx"
)

var (
	Db *sqlx.DB
)

func initDB() (err error) {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", AppConf.mysqlConf.Username, AppConf.mysqlConf.Password, AppConf.mysqlConf.Host, AppConf.mysqlConf.Port, AppConf.mysqlConf.Database)
	database, err := sqlx.Open("mysql", dns)
	if err != nil {
		logs.Error("failed to connect to mysql, err: %v", err)
		return
	}
	Db = database
	logs.Debug("init MysqlDB succ")
	return
}

func initAll() (err error) {
	err = initConfig()
	if err != nil {
		logs.Warn("init config failed, error: %v", err)
		return
	}
	err = initDB()
	if err != nil {
		logs.Warn("init db failed, error: %v", err)
		return
	}

	err = model.Init(Db)
	if err != nil {
		logs.Warn("init model failed, error: %v", err)
		return
	}
	return
}
