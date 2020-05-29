package service

import (
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
)

func InitService(serviceConf *SecKillConf) (err error) {
	secKillConf = serviceConf
	err = loadBlacklist()
	if err != nil {
		logs.Error("load blacklist failed, error: %v", err)
		return
	}

	err = initProxy2LayerRedis()
	if err != nil {
		logs.Error("load proxy2layer redis pool failed, error: %v", err)
		return
	}

	secKillConf.RequestLimitMgr = &RequestLimitMgr{
		UserLimitMap: make(map[int]*Limit, 1000),
		IPLimitMap:   make(map[string]*Limit, 10000),
	}

	err = initLayer2ProxyRedis()
	if err != nil {
		logs.Error("load layer2proxy redis pool failed, error: %v", err)
		return
	}

	secKillConf.SecKillRequestChan = make(chan *SecKillRequest, secKillConf.SecKillRequestChanSize)
	secKillConf.UserConnMap = make(map[string]chan *SecKillResult, 10000)

	initRedisProcessFunc(secKillConf)
	logs.Debug("init service succ, config: %v", secKillConf)
	return
}

func initRedisProcessFunc(serviceConf *SecKillConf) {
	for i := 0; i < secKillConf.WriteProxy2LayerGoroutineNum; i++ {
		go WriteHandle()
	}
	for i := 0; i < secKillConf.ReadLayer2ProxyGoroutineNum; i++ {
		go ReadHandle()
	}
	return
}

func initProxy2LayerRedis() (err error) {
	secKillConf.Proxy2LayerRedisPool = &redis.Pool{
		MaxIdle:     secKillConf.RedisProxy2LayerConf.RedisMaxIdle,
		MaxActive:   secKillConf.RedisProxy2LayerConf.RedisMaxActive,
		IdleTimeout: time.Duration(secKillConf.RedisProxy2LayerConf.RedisIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", secKillConf.RedisProxy2LayerConf.RedisAddr)
		},
	}
	conn := secKillConf.Proxy2LayerRedisPool.Get()
	defer conn.Close()
	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed, error: %v", err)
		return
	}
	return
}

func initLayer2ProxyRedis() (err error) {
	secKillConf.Layer2ProxyRedisPool = &redis.Pool{
		MaxIdle:     secKillConf.RedisLayer2ProxyConf.RedisMaxIdle,
		MaxActive:   secKillConf.RedisLayer2ProxyConf.RedisMaxActive,
		IdleTimeout: time.Duration(secKillConf.RedisLayer2ProxyConf.RedisIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", secKillConf.RedisLayer2ProxyConf.RedisAddr)
		},
	}
	conn := secKillConf.Layer2ProxyRedisPool.Get()
	defer conn.Close()
	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed, error: %v", err)
		return
	}
	return
}

func initBlackListRedis() (err error) {
	secKillConf.BlacklistRedisPool = &redis.Pool{
		MaxIdle:     secKillConf.RedisBlacklistConf.RedisMaxIdle,
		MaxActive:   secKillConf.RedisBlacklistConf.RedisMaxActive,
		IdleTimeout: time.Duration(secKillConf.RedisBlacklistConf.RedisIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", secKillConf.RedisBlacklistConf.RedisAddr)
		},
	}
	conn := secKillConf.BlacklistRedisPool.Get()
	defer conn.Close()
	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed, error: %v", err)
		return
	}
	return
}

func loadBlacklist() (err error) {
	secKillConf.IDBlacklistMap = make(map[int]bool, 1000)
	secKillConf.IPBlacklistMap = make(map[string]bool, 1000)
	err = initBlackListRedis()
	if err != nil {
		logs.Error("init blacklist redis failed, error: %v", err)
		return
	}

	conn := secKillConf.BlacklistRedisPool.Get()
	defer conn.Close()

	reply, err := conn.Do("HGETALL", "idblacklist")
	idList, err := redis.Strings(reply, err)
	if err != nil {
		logs.Warn("hget all failed, error: %v", err)
		return
	}
	for _, v := range idList {
		id, err := strconv.Atoi(v)
		if err != nil {
			logs.Warn("invalid user id[%v]", id)
			continue
		}
		secKillConf.IDBlacklistMap[id] = true
	}

	reply, err = conn.Do("HGETALL", "ipblacklist")
	ipList, err := redis.Strings(reply, err)
	if err != nil {
		logs.Warn("hget all failed, error: %v", err)
		return
	}
	for _, v := range ipList {
		secKillConf.IPBlacklistMap[v] = true
	}

	// Todo fix go routine leak
	// syncIPBlackList()
	// syncIDBlackList()

	return
}

func syncIPBlackList() {
	var ipList []string
	lastTime := time.Now().Unix()
	for {
		conn := secKillConf.BlacklistRedisPool.Get()
		defer conn.Close()
		reply, err := conn.Do("BLPOP", "blackiplist", time.Second)
		ip, err := redis.String(reply, err)
		if err != nil {
			continue
		}

		curTime := time.Now().Unix()
		ipList = append(ipList, ip)

		if len(ipList) > 100 || curTime-lastTime > 5 {
			secKillConf.RWBlacklistLock.Lock()
			for _, v := range ipList {
				secKillConf.IPBlacklistMap[v] = true
			}
			secKillConf.RWBlacklistLock.Unlock()
			lastTime = curTime
			logs.Info("sync ip list from redis succ, ip[%v]", ipList)
		}
	}
}

func syncIDBlackList() {
	// var idList []int
	// lastTime := time.Now().Unix()
	for {
		conn := secKillConf.BlacklistRedisPool.Get()
		defer conn.Close()
		reply, err := conn.Do("BLPOP", "blackidlist", time.Second)
		id, err := redis.Int(reply, err)
		if err != nil {
			continue
		}
		secKillConf.RWBlacklistLock.Lock()
		secKillConf.IDBlacklistMap[id] = true
		secKillConf.RWBlacklistLock.Unlock()
		logs.Info("sync id list from redis succ, ip[%v]", id)

		// Todo fix go routine leak
		// curTime := time.Now().Unix()
		// idList = append(idList, id)

		// if len(idList) > 100 || curTime-lastTime > 5 {
		// 	secKillConf.RWBlacklistLock.Lock()
		// 	for _, v := range idList {
		// 		secKillConf.IDBlacklistMap[v] = true
		// 	}
		// 	secKillConf.RWBlacklistLock.Unlock()
		// 	lastTime = curTime
		// 	logs.Info("sync id list from redis succ, ip[%v]", idList)
		// }
	}
}
