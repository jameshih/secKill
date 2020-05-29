package service

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
)

func initRedisPool(redisConf RedisConf) (pool *redis.Pool, err error) {
	pool = &redis.Pool{
		MaxIdle:     redisConf.RedisMaxIdle,
		MaxActive:   redisConf.RedisMaxActive,
		IdleTimeout: time.Duration(redisConf.RedisIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisConf.RedisAddr)
		},
	}
	conn := pool.Get()
	defer conn.Close()
	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed, error: %v", err)
		return
	}
	return
}

func initRedis(appConf *LogicConf) (err error) {
	logicContext.Proxy2LayerRedisPool, err = initRedisPool(appConf.Proxy2LayerRedis)
	if err != nil {
		logs.Error("init proxy2layer redis pool failed, error: %v", err)
		return
	}

	logicContext.Layer2ProxyRedisPool, err = initRedisPool(appConf.Layer2ProxyRedis)
	if err != nil {
		logs.Error("init layer2proxy redis pool failed, error: %v", err)
		return
	}
	return
}

func RunProcess() (err error) {
	for i := 0; i < logicContext.logicConf.ReadGoroutineNum; i++ {
		logicContext.waitgroup.Add(1)
		go HandleRead()
	}
	for i := 0; i < logicContext.logicConf.WriteGoroutineNum; i++ {
		logicContext.waitgroup.Add(1)
		go HandleWrite()
	}

	for i := 0; i < logicContext.logicConf.HandleUserGoroutineNum; i++ {
		logicContext.waitgroup.Add(1)
		go HandleUser()
	}

	logs.Debug("all process goroutine started")
	logicContext.waitgroup.Wait()
	logs.Debug("goroutine waitgroup exited")
	return
}

func HandleRead() {
	logs.Debug("handle read running")
	for {
		conn := logicContext.Proxy2LayerRedisPool.Get()
		for {
			ret, err := conn.Do("BLPOP", logicContext.logicConf.Proxy2LayerRedis.RedisQueueName, 0)
			if err != nil {
				break
			}
			tmp, ok := ret.([]interface{})
			if !ok || len(tmp) != 2 {
				logs.Error("pop from queue failed, error: %v", err)
				continue
			}
			data, ok := tmp[1].([]byte)
			if !ok {
				logs.Error("pop from queue failed, error: %v", err)
				continue
			}

			logs.Debug("pop from queue, data: %s", string(data))

			var req SecKillRequest
			json.Unmarshal([]byte(data), &req)
			if err != nil {
				logs.Error("unmarshal to seckillrequest failed, error: %v", err)
				continue
			}
			now := time.Now().Unix()
			if now-req.AccessTime.Unix() >= int64(logicContext.logicConf.MaxRequestWaitTimeout) {
				logs.Warn("req[%v] is expire", req)
				continue
			}

			timer := time.NewTicker(time.Millisecond * time.Duration(logicContext.logicConf.SendToHandleChanTimeout))
			select {
			case logicContext.Read2HandleChan <- &req:
			case <-timer.C:
				logs.Warn("send to handle chan timeout, req: %v", req)
				break

			}
		}
		conn.Close()
	}
}

func HandleWrite() {
	logs.Debug("handle write running")
	for res := range logicContext.Handle2WriteChan {
		err := sendToRedis(res)
		if err != nil {
			logs.Error("send to redis failed, error: %v, res: %v", err, res)
			continue
		}
	}
}

func sendToRedis(res *SecKillResponse) (err error) {
	data, err := json.Marshal(res)
	if err != nil {
		logs.Error("marshal failed, error: %v", err)
		return
	}
	conn := logicContext.Layer2ProxyRedisPool.Get()
	_, err = conn.Do("RPUSH", logicContext.logicConf.Layer2ProxyRedis.RedisQueueName, string(data))
	if err != nil {
		logs.Warn("RPUSH to redis failed, error: %v", err)
		return
	}
	return
}

func HandleUser() {
	logs.Debug("handle user running")
	for req := range logicContext.Read2HandleChan {
		logs.Debug("begin process request: %v", req)
		res, err := HandleSecKill(req)
		if err != nil {
			logs.Warn("process request %v failed, error: %v", err)
			res = &SecKillResponse{
				Code: ErrServiceBusy,
			}
		}

		//request timeout
		timer := time.NewTicker(time.Millisecond * time.Duration(logicContext.logicConf.SendToWriteChanTimeout))
		select {
		case logicContext.Handle2WriteChan <- res:
		case <-timer.C:
			logs.Warn("send to response chan timeout, res: %v", res)
			break

		}
	}
	return
}

func HandleSecKill(req *SecKillRequest) (res *SecKillResponse, err error) {
	logicContext.RwSecKillProductLock.RLock()
	defer logicContext.RwSecKillProductLock.RUnlock()

	res = &SecKillResponse{}
	res.UserID = req.UserID
	res.ProductID = req.ProductID
	product, ok := logicContext.logicConf.ProductInfoMap[req.ProductID]
	if !ok {
		logs.Error("cannot find product: %v", req.ProductID)
		res.Code = ErrNotFoundProduct
		return
	}
	if product.Status == ProductStatusSoldout {
		res.Code = ErrSoldOut
		return
	}

	now := time.Now().Unix()
	alreadySoldCount := product.SecLimit.Check(now)
	if alreadySoldCount >= product.SoldMaxLimit {
		res.Code = ErrRetry
		return
	}

	logicContext.UserBuyHistoryMapLock.Lock()
	defer logicContext.UserBuyHistoryMapLock.Unlock()
	userHistory, ok := logicContext.UserBuyHistoryMap[req.UserID]
	if !ok {
		userHistory = &UserBuyHistory{
			History: make(map[int]int, 16),
		}
		logicContext.UserBuyHistoryMap[req.UserID] = userHistory
	}

	historyCount := userHistory.GetProductBuyCount(req.ProductID)
	if historyCount >= product.OnePersonBuyLimit {
		res.Code = ErrAlreadyBuy
		return
	}

	curSold := logicContext.productCountMgr.Count(req.ProductID)
	if curSold >= product.Total {
		res.Code = ErrSoldOut
		product.Status = ProductStatusSoldout
		return
	}

	curRate := rand.Float64()
	logs.Debug("curRate: %v, product: %v", curRate, product.BuyRate)
	// Todo fix curRate
	if curRate > 0.8 {
		// if curRate > product.BuyRate {
		res.Code = ErrRetry
		return
	}

	userHistory.Add(req.ProductID, 1)
	logicContext.productCountMgr.Add(req.ProductID, 1)

	// generate token (userID +productID + currentTime + tokenSecret)

	res.Code = ErrSecKillSucc
	tokenData := fmt.Sprintf("userID=%d&productID=%d&timestamp=%d&secret=%s",
		req.UserID, req.ProductID, now, logicContext.logicConf.TokenSecret)
	res.Token = fmt.Sprintf("%x", md5.Sum([]byte(tokenData)))
	res.TokenTime = now //time must be the same as in Token or else md5 hash will be different
	return
}
