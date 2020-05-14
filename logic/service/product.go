package service

import (
	"context"
	"encoding/json"

	"github.com/astaxie/beego/logs"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

func loadSecConf(appConf *LogicConf) (err error) {
	resp, err := logicContext.EtcdClient.Get(context.Background(), appConf.EtcdConfig.EtcdSecProductKey)
	if err != nil {
		logs.Error("get [%s] from etcd failed, err:%v", appConf.EtcdConfig.EtcdSecProductKey, err)
		return
	}

	var produtInfo []ProductInfoConf
	for k, v := range resp.Kvs {
		logs.Debug("key[%v] valud[%v]", k, v)
		err = json.Unmarshal(v.Value, &produtInfo)
		if err != nil {
			logs.Error("Unmarshal sec product info failed, err:%v", err)
			return
		}

		logs.Debug("sec info conf is [%v]", produtInfo)
	}

	updateSecProductInfo(appConf, produtInfo)
	initSecProductWatcher(appConf)
	return
}

func updateSecProductInfo(appConf *LogicConf, secProductInfo []ProductInfoConf) {

	var tmp map[int]*ProductInfoConf = make(map[int]*ProductInfoConf, 1024)
	for _, v := range secProductInfo {
		produtInfo := v
		tmp[v.ProductID] = &produtInfo
	}
	logicContext.RwSecKillProductLock.Lock()
	appConf.ProductInfoMap = tmp
	logicContext.RwSecKillProductLock.Unlock()
}

func initSecProductWatcher(appConf *LogicConf) {
	go watchSecProductKey(appConf)
}

func watchSecProductKey(appConf *LogicConf) {
	logs.Debug("begin watch key:%s", appConf.EtcdConfig.EtcdSecProductKey)
	for {
		rch := logicContext.EtcdClient.Watch(context.Background(), appConf.EtcdConfig.EtcdSecProductKey)
		var secProductInfo []ProductInfoConf
		var getConfSucc = true

		for wresp := range rch {
			for _, ev := range wresp.Events {
				if ev.Type == mvccpb.DELETE {
					logs.Warn("key[%s] 's config deleted", appConf.EtcdConfig.EtcdSecProductKey)
					continue
				}

				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == appConf.EtcdConfig.EtcdSecProductKey {
					err := json.Unmarshal(ev.Kv.Value, &secProductInfo)
					if err != nil {
						logs.Error("key [%s], Unmarshal[%s], err:%v ", err)
						getConfSucc = false
						continue
					}
				}
				logs.Debug("get config from etcd, %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}

			if getConfSucc {
				logs.Debug("get config from etcd succ, %v", secProductInfo)
				updateSecProductInfo(appConf, secProductInfo)
			}
		}
	}
}
