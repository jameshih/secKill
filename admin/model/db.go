package model

import (
	"github.com/astaxie/beego/logs"
	"github.com/jmoiron/sqlx"
	"go.etcd.io/etcd/clientv3"
)

var (
	Db         *sqlx.DB
	EtcdClient *clientv3.Client
)

func Init(db *sqlx.DB, etcdClient *clientv3.Client) (err error) {
	Db = db
	EtcdClient = etcdClient
	logs.Debug("init all succ")
	return
}
