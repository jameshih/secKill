package service

import (
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	ProductStatusNormal       = 0
	ProductStatusSaleOut      = 1
	ProductStatusForceSaleOut = 2
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

type AccesLimitConf struct {
	IPSecAccessLimit   int
	UserSecAccessLimit int
	IPMinAccessLimit   int
	UserMinAccessLimit int
}

type SecKillConf struct {
	RedisBlacklistConf           RedisConf
	RedisProxy2LayerConf         RedisConf
	RedisLayer2ProxyConf         RedisConf
	EtcdConf                     EtcdConf
	LogPath                      string
	LogLevel                     string
	ProductInfoMap               map[int]*ProductInfoConf
	RwSecKillProductLock         sync.RWMutex
	CookieSecretKey              string
	RefererWhiteList             []string
	IPBlacklistMap               map[string]bool
	IDBlacklistMap               map[int]bool
	BlacklistRedisPool           *redis.Pool
	Proxy2LayerRedisPool         *redis.Pool
	Layer2ProxyRedisPool         *redis.Pool
	AccessLimitConf              AccesLimitConf
	RequestLimitMgr              *RequestLimitMgr
	RWBlacklistLock              sync.RWMutex
	WriteProxy2LayerGoroutineNum int
	ReadProxy2LayerGoroutineNum  int
	ReadLayer2ProxyGoroutineNum  int
	SecKillRequestChan           chan *SecKillRequest
	SecKillRequestChanSize       int

	UserConnMap     map[string]chan *SecKillResult
	UserConnMapLock sync.Mutex
}

type ProductInfoConf struct {
	ProductID int
	StartTime int64
	EndTime   int64
	Status    int
	Total     int
	Left      int
}

type SecKillResult struct {
	ProductID int
	UserID    int
	Code      int
	Token     string
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
	CloseNotify   <-chan bool

	ResultChan chan *SecKillResult
}

type Limit struct {
	secLimit TimeLimit
	minLimit TimeLimit
}

type TimeLimit interface {
	Count(nowTime int64) (curCount int)
	Check(nowTime int64) int
}
