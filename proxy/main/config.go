package main

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/jameshih/secKill/proxy/service"
)

var (
	secKillConf = &service.SecKillConf{
		ProductInfoMap: make(map[int]*service.ProductInfoConf, 1024),
	}
)

func initConfig() (err error) {

	secKillConf.LogPath = beego.AppConfig.String("log_path")
	secKillConf.LogLevel = beego.AppConfig.String("log_level")
	productKey := beego.AppConfig.String("etcd_seckill_product_key")
	if len(productKey) == 0 {
		err = fmt.Errorf("init config failed, read etcd_product_key error:%v", err)
		return
	}

	secKillConf.CookieSecretKey = beego.AppConfig.String("cookie_secretkey")
	secLimit, err := beego.AppConfig.Int("user_sec_access_limit")
	if err != nil {
		err = fmt.Errorf("init config failed, read user_sec_access_limit error:%v", err)
		return
	}

	referList := beego.AppConfig.String("refer_whitelist")
	if len(referList) > 0 {
		secKillConf.RefererWhiteList = strings.Split(referList, ",")
	}

	ipLimit, err := beego.AppConfig.Int("ip_sec_access_limit")
	if err != nil {
		err = fmt.Errorf("init config failed, read ip_sec_access_limit error:%v", err)
		return
	}

	minIdLimit, err := beego.AppConfig.Int("user_min_access_limit")
	if err != nil {
		err = fmt.Errorf("init config failed, read user_min_access_limit error:%v", err)
		return
	}

	minIpLimit, err := beego.AppConfig.Int("ip_min_access_limit")
	if err != nil {
		err = fmt.Errorf("init config failed, read ip_min_access_limit error:%v", err)
		return
	}

	secKillConf.AccessLimitConf.UserSecAccessLimit = secLimit
	secKillConf.AccessLimitConf.IPSecAccessLimit = ipLimit
	secKillConf.AccessLimitConf.UserMinAccessLimit = minIdLimit
	secKillConf.AccessLimitConf.IPMinAccessLimit = minIpLimit

	//Etcd
	etcdAddr := beego.AppConfig.String("etcd_addr")
	if len(etcdAddr) == 0 {
		err = fmt.Errorf("init config failed etcd[%s] config is null", etcdAddr)
		return
	}
	etcdTimeout, err := beego.AppConfig.Int("etcd_timeout")
	if err != nil {
		err = fmt.Errorf("init config failed, read etcd_timeout error:%v", err)
		return
	}

	secKillConf.EtcdConf.EtcdSecKeyPrefix = beego.AppConfig.String("etcd_seckill_key_prefix")
	if len(secKillConf.EtcdConf.EtcdSecKeyPrefix) == 0 {
		err = fmt.Errorf("init config failed, read etcd_seckill_key error:%v", err)
		return
	}

	secKillConf.EtcdConf.EtcdAddr = etcdAddr
	secKillConf.EtcdConf.EtcdTimeOut = etcdTimeout
	if strings.HasSuffix(secKillConf.EtcdConf.EtcdSecKeyPrefix, "/") == false {
		secKillConf.EtcdConf.EtcdSecKeyPrefix = secKillConf.EtcdConf.EtcdSecKeyPrefix + "/"
	}
	secKillConf.EtcdConf.EtcdSecProductKey = fmt.Sprintf("%s%s", secKillConf.EtcdConf.EtcdSecKeyPrefix, productKey)

	// Redis blacklist
	redisBlackAddr := beego.AppConfig.String("redis_blacklist_addr")
	if len(redisBlackAddr) == 0 {
		err = fmt.Errorf("init config failed redisBlack[%s] config is null", redisBlackAddr)
		return
	}
	redisMaxIdle, err := beego.AppConfig.Int("redis_blacklist_max_idle")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_blacklist_max_idle error:%v", err)
		return
	}

	redisMaxActive, err := beego.AppConfig.Int("redis_blacklist_max_active")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_blacklist_max_active error:%v", err)
		return
	}

	redisIdleTimeout, err := beego.AppConfig.Int("redis_blacklist_idle_timeout")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_blacklist_idle_timeout error:%v", err)
		return
	}

	secKillConf.RedisBlacklistConf.RedisAddr = redisBlackAddr
	secKillConf.RedisBlacklistConf.RedisMaxIdle = redisMaxIdle
	secKillConf.RedisBlacklistConf.RedisMaxActive = redisMaxActive
	secKillConf.RedisBlacklistConf.RedisIdleTimeout = redisIdleTimeout

	// Redis Proxy 2 Layer
	redisProxy2LayerAddr := beego.AppConfig.String("redis_proxy2layer_addr")
	if len(redisProxy2LayerAddr) == 0 {
		err = fmt.Errorf("init config failed, redis[%s]  config is null", redisProxy2LayerAddr)
		return
	}

	redisMaxIdle, err = beego.AppConfig.Int("redis_proxy2layer_max_idle")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_proxy2layer_max_idle error:%v", err)
		return
	}

	redisMaxActive, err = beego.AppConfig.Int("redis_proxy2layer_max_active")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_proxy2layer_max_active error:%v", err)
		return
	}

	redisIdleTimeout, err = beego.AppConfig.Int("redis_proxy2layer_idle_timeout")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_proxy2layer_idle_timeout error:%v", err)
		return
	}

	writeGoNums, err := beego.AppConfig.Int("write_proxy2layer_goroutine_num")
	if err != nil {
		err = fmt.Errorf("init config failed, read write_proxy2layer_goroutine_num error:%v", err)
		return
	}

	secKillConf.RedisProxy2LayerConf.RedisAddr = redisProxy2LayerAddr
	secKillConf.RedisProxy2LayerConf.RedisMaxIdle = redisMaxIdle
	secKillConf.RedisProxy2LayerConf.RedisMaxActive = redisMaxActive
	secKillConf.RedisProxy2LayerConf.RedisIdleTimeout = redisIdleTimeout
	secKillConf.WriteProxy2LayerGoroutineNum = writeGoNums

	// Redis Layer 2 Proxy
	redisLayer2ProxyAddr := beego.AppConfig.String("redis_layer2proxy_addr")
	if len(redisLayer2ProxyAddr) == 0 {
		err = fmt.Errorf("init config failed, redis[%s]  config is null", redisLayer2ProxyAddr)
		return
	}

	redisMaxIdle, err = beego.AppConfig.Int("redis_layer2proxy_max_idle")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_layer2proxy_max_idle error:%v", err)
		return
	}

	redisMaxActive, err = beego.AppConfig.Int("redis_layer2proxy_max_active")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_layer2proxy_max_active error:%v", err)
		return
	}

	redisIdleTimeout, err = beego.AppConfig.Int("redis_layer2proxy_idle_timeout")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_layer2proxy_idle_timeout error:%v", err)
		return
	}

	readGoNums, err := beego.AppConfig.Int("read_layer2proxy_goroutine_num")
	if err != nil {
		err = fmt.Errorf("init config failed, read read_layer2proxy_goroutine_num error:%v", err)
		return
	}

	secKillConf.RedisLayer2ProxyConf.RedisAddr = redisLayer2ProxyAddr
	secKillConf.RedisLayer2ProxyConf.RedisMaxIdle = redisMaxIdle
	secKillConf.RedisLayer2ProxyConf.RedisMaxActive = redisMaxActive
	secKillConf.RedisLayer2ProxyConf.RedisIdleTimeout = redisIdleTimeout
	secKillConf.ReadLayer2ProxyGoroutineNum = readGoNums

	return
}
