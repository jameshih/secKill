package main

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jameshih/secKill/admin/model"
	"github.com/jmoiron/sqlx"
	"go.etcd.io/etcd/clientv3"
)

func initDB() (db *sqlx.DB, err error) {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", AppConf.mysqlConf.Username, AppConf.mysqlConf.Password, AppConf.mysqlConf.Host, AppConf.mysqlConf.Port, AppConf.mysqlConf.Database)
	database, err := sqlx.Open("mysql", dns)
	if err != nil {
		logs.Error("failed to connect to mysql, err: %v", err)
		return
	}
	db = database
	logs.Debug("init MysqlDB succ")
	return
}

func initEtcd() (etcdClient *clientv3.Client, err error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{AppConf.etcdConf.Addr},
		DialTimeout: time.Duration(AppConf.etcdConf.EtcdTimeout) * time.Second,
	})
	if err != nil {
		logs.Error("connect etcd failed, err:", err)
		return
	}

	etcdClient = cli
	logs.Debug("init Etcd succ")
	return
}

func initAll() (err error) {
	err = initConfig()
	if err != nil {
		logs.Warn("init config failed, error: %v", err)
		return
	}
	db, err := initDB()
	if err != nil {
		logs.Warn("init db failed, error: %v", err)
		return
	}

	etcdClient, err := initEtcd()
	if err != nil {
		logs.Warn("init etcd failed, error: %v", err)
		return
	}

	err = model.Init(db, etcdClient, AppConf.etcdConf.EtcdSecProductKey)
	if err != nil {
		logs.Warn("init model failed, error: %v", err)
		return
	}
	return
}
