package config

import (
	"strings"
	"github.com/jinzhu/configor"
	"os"
	log "github.com/alecthomas/log4go"
	"fmt"
	"paas/icsoc-monitor/im-flag"
)

var conf AppConfig

type AppConfig struct {
	//mongo配置
	MongoDb struct{
		Addr string
	}

	//Redis配置
	RedisDb struct{
		Addr string
		Pass string
		User string
		Database int
	}

	//rabbitMQ配置信息
	RabbitMQ struct {
		Addr string
	}

	//日志级别
	Logger struct {
		Level string
	}

	//http服务地址
	Http struct {
		Addr string
	}

	DebugConfig struct{
		Debug bool
	}
}

func Init(file string){
	conf = AppConfig{}
	if len(file) <= 0 || strings.TrimSpace(file) == "" {
		file = im_flag.ConfPath
	}

	err := configor.Load(&conf, file)
	if err != nil {
		log.Error("Failed to find configuration ", file)
		os.Exit(1)
	}

	log.Info(fmt.Sprintf("load config path %s", file))
	log.Info("redis addr:%s", conf.RedisDb.Addr)
	log.Info("mq addr:%s", conf.RabbitMQ.Addr)
	log.Info(fmt.Sprintf("debug config..%t", conf.DebugConfig.Debug))
	log.Info("mongodb addr:%s", conf.MongoDb.Addr)
	log.Info("logger lv:%s", conf.Logger.Level)
}

func GetConf() AppConfig{
	return conf
}
