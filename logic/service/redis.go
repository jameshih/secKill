package service

import (
	"encoding/json"
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
	for {
		conn := logicContext.Proxy2LayerRedisPool.Get()
		for {
			data, err := redis.String(conn.Do("BLPOP", "q", 0))
			if err != nil {
				logs.Error("pop from queue failed, error: %v", err)
				break
			}
			logs.Debug("pop from queue, data:%s", data)

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
	_, err = redis.String(conn.Do("RPUSH", "layer2proxy_queue", string(data)))
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

	return
}
