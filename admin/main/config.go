package main

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type MysqlConfig struct {
	Username string
	Password string
	Port     int
	Database string
	Host     string
}

type Config struct {
	mysqlConf MysqlConfig
}

var (
	AppConf Config
)

func initConfig() (err error) {
	username := beego.AppConfig.String("mysql_user_name")
	if len(username) == 0 {
		err = fmt.Errorf("load config of mysql_user_name failed, is null")
		logs.Error(err)
		return
	}
	AppConf.mysqlConf.Username = username

	password := beego.AppConfig.String("mysql_password")
	if len(password) == 0 {
		err = fmt.Errorf("load config of mysql_password failed, is null")
		logs.Error(err)
		return
	}
	AppConf.mysqlConf.Password = password

	host := beego.AppConfig.String("mysql_host")
	if len(host) == 0 {
		err = fmt.Errorf("load config of mysql_host failed, is null")
		logs.Error(err)
		return
	}
	AppConf.mysqlConf.Host = host

	db := beego.AppConfig.String("mysql_db")
	if len(db) == 0 {
		err = fmt.Errorf("load config of mysql_db failed, is null")
		logs.Error(err)
		return
	}
	AppConf.mysqlConf.Database = db

	port, err := beego.AppConfig.Int("mysql_port")
	if err != nil {
		err = fmt.Errorf("load config of mysql_port failed, is null")
		logs.Error(err)
		return
	}
	AppConf.mysqlConf.Port = port
	return
}
