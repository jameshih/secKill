package service

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
)

func WriteHandle() {
	for {
		req := <-secKillConf.SecKillRequestChan
		conn := secKillConf.Proxy2LayerRedisPool.Get()
		data, err := json.Marshal(req)
		if err != nil {
			logs.Error("json.Marshal failed, error: %v req: %v", err, req)
			conn.Close()
			continue
		}
		_, err = conn.Do("LPUSH", "sec_queue", data)
		if err != nil {
			logs.Error("LPUSH failedd, error: %v, req: %v", err, req)
			conn.Close()
			continue
		}
		conn.Close()
	}
}

func ReadHandle() {
	for {
		conn := secKillConf.Proxy2LayerRedisPool.Get()
		reply, err := conn.Do("RPOP", "recv_queue")
		data, err := redis.String(reply, err)
		if err != nil {
			logs.Error("RPOP failed, error: %v", err)
			conn.Close()
			continue
		}
		var result SecKillResult

		err = json.Unmarshal([]byte(data), &result)
		if err != nil {
			logs.Error("json.Unmarshal failed, error: %v", err)
			conn.Close()
			continue
		}

		userKey := fmt.Sprintf("%d_%d", result.UserID, result.ProductID)

		secKillConf.UserConnMapLock.Lock()
		resultChan, ok := secKillConf.UserConnMap[userKey]
		secKillConf.UserConnMapLock.Unlock()

		if !ok {
			conn.Close()
			logs.Warn("user not found: %v", userKey)
			continue
		}

		resultChan <- &result
		conn.Close()
	}
}
