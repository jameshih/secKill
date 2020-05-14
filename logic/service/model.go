package service

import (
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	"go.etcd.io/etcd/clientv3"
)

var (
	logicContext = &LogicContext{}
)

type RedisConf struct {
	RedisAddr        string
	RedisMaxIdle     int
	RedisMaxActive   int
	RedisIdleTimeout int
}

type EtcdConf struct {
	EtcdAddr          string
	EtcdTimeOut       int
	EtcdSecKeyPrefix  string
	EtcdSecProductKey string
}

type ProductInfoConf struct {
	ProductID int
	StartTime int64
	EndTime   int64
	Status    int
	Total     int
	Left      int
}

type LogicConf struct {
	Proxy2LayerRedis RedisConf
	Layer2ProxyRedis RedisConf
	EtcdConfig       EtcdConf
	LogPath          string
	LogLevel         string

	WriteGoroutineNum      int
	ReadGoroutineNum       int
	HandleUserGoroutineNum int
	Read2HandleChanSize    int
	Handle2WriteChanSize   int
	MaxRequestWaitTimeout  int
	ProductInfoMap         map[int]*ProductInfoConf

	SendToWriteChanTimeout  int
	SendToHandleChanTimeout int
}

type LogicContext struct {
	Proxy2LayerRedisPool *redis.Pool
	Layer2ProxyRedisPool *redis.Pool
	EtcdClient           *clientv3.Client
	RwSecKillProductLock sync.RWMutex
	logicConf            *LogicConf
	waitgroup            sync.WaitGroup
	Read2HandleChan      chan *SecKillRequest
	Handle2WriteChan     chan *SecKillResponse
}

type SecKillRequest struct {
	ProductID     int
	Source        string
	AuthCode      string
	SecTime       string
	Nounce        string
	UserID        int
	UserAuthSig   string
	AccessTime    time.Time
	ClientAddr    string
	ClientReferer string
	// CloseNotify   <-chan bool

	// ResultChan chan *SecKillResult
}

type SecKillResponse struct {
	ProductID int
	UserID    int
	Token     string
	Code      int
}
