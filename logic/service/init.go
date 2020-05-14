package service

import (
	"time"

	"github.com/astaxie/beego/logs"
	"go.etcd.io/etcd/clientv3"
)

func initEtcd(appConf *LogicConf) (err error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{appConf.EtcdConfig.EtcdAddr},
		DialTimeout: time.Duration(appConf.EtcdConfig.EtcdTimeOut) * time.Second,
	})
	if err != nil {
		logs.Error("connect etcd failed, err:", err)
		return
	}

	logicContext.EtcdClient = cli
	return
}

func InitSecKillLogic(appConf *LogicConf) (err error) {
	err = initRedis(appConf)
	if err != nil {
		logs.Error("init redis failed, error: %v", err)
		return
	}

	err = initEtcd(appConf)
	if err != nil {
		logs.Error("init etcd failed, error: %v", err)
		return
	}

	err = loadSecConf(appConf)
	if err != nil {
		logs.Error("load product from etcd failed, error: %v", err)
		return
	}
	logicContext.logicConf = appConf
	logicContext.Read2HandleChan = make(chan *SecKillRequest, appConf.Read2HandleChanSize)
	logicContext.Handle2WriteChan = make(chan *SecKillResponse, appConf.Handle2WriteChanSize)
	return
}
