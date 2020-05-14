package service

import (
	"fmt"
	"sync"

	"github.com/astaxie/beego/logs"
)

type RequestLimitMgr struct {
	UserLimitMap map[int]*Limit
	IPLimitMap   map[string]*Limit
	lock         sync.Mutex
}

func antiSpam(req *SecKillRequest) (err error) {
	_, ok := secKillConf.IDBlacklistMap[req.UserID]
	if ok {
		err = fmt.Errorf("invalid request")
		logs.Error("userID[%v] is blocked", req.UserID)
		return
	}
	_, ok = secKillConf.IPBlacklistMap[req.ClientAddr]
	if ok {
		err = fmt.Errorf("invalid request")
		logs.Error("userID[%v] ip[%v] is blocked", req.UserID, req.ClientAddr)
		return
	}
	secKillConf.RequestLimitMgr.lock.Lock()
	limit, ok := secKillConf.RequestLimitMgr.UserLimitMap[req.UserID]
	if !ok {
		limit = &Limit{
			secLimit: &SecLimit{},
			minLimit: &MinLimit{},
		}
		secKillConf.RequestLimitMgr.UserLimitMap[req.UserID] = limit
	}

	secIDCount := limit.secLimit.Count(req.AccessTime.Unix())
	minIDCount := limit.secLimit.Count(req.AccessTime.Unix())

	limit, ok = secKillConf.RequestLimitMgr.IPLimitMap[req.ClientAddr]
	if !ok {
		limit = &Limit{
			secLimit: &SecLimit{},
			minLimit: &MinLimit{},
		}
		secKillConf.RequestLimitMgr.IPLimitMap[req.ClientAddr] = limit
	}
	secIPCount := limit.secLimit.Count(req.AccessTime.Unix())
	minIPCount := limit.secLimit.Count(req.AccessTime.Unix())

	secKillConf.RequestLimitMgr.lock.Unlock()

	if secIPCount > secKillConf.AccessLimitConf.IPSecAccessLimit {
		err = fmt.Errorf("invalid request")
		return
	}

	if minIPCount > secKillConf.AccessLimitConf.IPMinAccessLimit {
		err = fmt.Errorf("invalid request")
		return
	}

	if secIDCount > secKillConf.AccessLimitConf.UserSecAccessLimit {
		err = fmt.Errorf("invalid request")
		return
	}

	if minIDCount > secKillConf.AccessLimitConf.UserMinAccessLimit {
		err = fmt.Errorf("invalid request")
		return
	}
	return
}
