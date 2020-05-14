package main

import (
	"fmt"

	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"github.com/jameshih/secKill/logic/service"
)

var (
	appConfig *service.LogicConf
)

func initConfig(confType, filename string) (err error) {
	conf, err := config.NewConfig(confType, filename)
	if err != nil {
		err = fmt.Errorf("new config failed, error: %v", err)
		return
	}

	// read log config
	appConfig = &service.LogicConf{}
	appConfig.LogLevel = conf.String("logs::log_level")
	appConfig.LogPath = conf.String("logs::log_path")

	if len(appConfig.LogLevel) == 0 {
		appConfig.LogLevel = "debug"
	}

	if len(appConfig.LogPath) == 0 {
		appConfig.LogPath = "./logs"
	}

	// read Proxy2LayerRedis config
	appConfig.Proxy2LayerRedis.RedisAddr = conf.String("redis::redis_proxy2layer_addr")
	if len(appConfig.Proxy2LayerRedis.RedisAddr) == 0 {
		err = fmt.Errorf("read redis::redis_proxy2layer_addr failed")
		logs.Error(err)
		return
	}

	appConfig.Proxy2LayerRedis.RedisMaxIdle, err = conf.Int("redis::redis_proxy2layer_max_idle")
	if err != nil {
		err = fmt.Errorf("read redis::redis_proxy2layer_max_idle failed, error: %v", err)
		logs.Error(err)
		return
	}

	appConfig.Proxy2LayerRedis.RedisMaxActive, err = conf.Int("redis::redis_proxy2layer_max_active")
	if err != nil {
		err = fmt.Errorf("read redis::redis_proxy2layer_max_active failed, error: %v", err)
		logs.Error(err)
		return
	}

	appConfig.Proxy2LayerRedis.RedisIdleTimeout, err = conf.Int("redis::redis_proxy2layer_idle_timeout")
	if err != nil {
		err = fmt.Errorf("read redis::redis_proxy2layer_idle_timeout failed, error: %v", err)
		logs.Error(err)
		return
	}

	// read Layer2ProxyRedis
	appConfig.Layer2ProxyRedis.RedisAddr = conf.String("redis::redis_layer2proxy_addr")
	if len(appConfig.Layer2ProxyRedis.RedisAddr) == 0 {
		err = fmt.Errorf("read redis::redis_layer2proxy_addr failed")
		logs.Error(err)
		return
	}

	appConfig.Layer2ProxyRedis.RedisMaxIdle, err = conf.Int("redis::redis_layer2proxy_max_idle")
	if err != nil {
		err = fmt.Errorf("read redis::redis_layer2proxy_addr failed, error: %v", err)
		logs.Error(err)
		return
	}

	appConfig.Layer2ProxyRedis.RedisMaxActive, err = conf.Int("redis::redis_layer2proxy_max_active")
	if err != nil {
		err = fmt.Errorf("read redis::redis_layer2proxy_max_active failed, error: %v", err)
		logs.Error(err)
		return
	}

	appConfig.Layer2ProxyRedis.RedisIdleTimeout, err = conf.Int("redis::redis_layer2proxy_idle_timeout")
	if err != nil {
		err = fmt.Errorf("read redis::redis_layer2proxy_idle_timeout failed, error: %v", err)
		logs.Error(err)
		return
	}

	// read Redis read/write

	appConfig.WriteGoroutineNum, err = conf.Int("service::write_proxy2layer_goroutine_num")
	if err != nil {
		err = fmt.Errorf("read service::write_proxy2layer_goroutine_num failed, error: %v", err)
		logs.Error(err)
		return
	}

	appConfig.ReadGoroutineNum, err = conf.Int("service::read_layer2proxy_goroutine_num")
	if err != nil {
		err = fmt.Errorf("read service::read_layer2proxy_goroutine_num failed, error: %v", err)
		logs.Error(err)
		return
	}

	appConfig.HandleUserGoroutineNum, err = conf.Int("service::handle_user_goroutine_num")
	if err != nil {
		err = fmt.Errorf("read service::handle_user_goroutine_num failed, error: %v", err)
		logs.Error(err)
		return
	}

	appConfig.Read2HandleChanSize, err = conf.Int("service::read2handle_chan_size")
	if err != nil {
		err = fmt.Errorf("read service::read2handle_chan_size failed, error: %v", err)
		logs.Error(err)
		return
	}

	appConfig.MaxRequestWaitTimeout, err = conf.Int("service::max_request_wait_timeout")
	if err != nil {
		err = fmt.Errorf("read service::max_request_wait_timeout failed, error: %v", err)
		logs.Error(err)
		return
	}

	appConfig.Handle2WriteChanSize, err = conf.Int("service::handle2write_chan_size")
	if err != nil {
		err = fmt.Errorf("read service::handle2write_chan_size failed, error: %v", err)
		logs.Error(err)
		return
	}

	appConfig.SendToHandleChanTimeout, err = conf.Int("service::send_to_handle_chan_timeout")
	if err != nil {
		err = fmt.Errorf("read service::send_to_handle_chan_timeout failed, error: %v", err)
		logs.Error(err)
		return
	}

	appConfig.SendToWriteChanTimeout, err = conf.Int("service::send_to_write_chan_timeout")
	if err != nil {
		err = fmt.Errorf("read service::send_to_write_chan_timeout failed, error: %v", err)
		logs.Error(err)
		return
	}

	return
}
