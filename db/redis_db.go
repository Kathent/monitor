package db

import (
	"paas/icsoc-monitor/config"

	"github.com/alecthomas/log4go"
	"github.com/go-redis/redis"
)

var client *redis.Client

func GetRedisClient() error {
	if client == nil {
		client = redis.NewClient(&redis.Options{
			Addr:     config.GetConf().RedisDb.Addr,
			Password: config.GetConf().RedisDb.Pass,
			DB:       config.GetConf().RedisDb.Database,
		})
	}

	_, err := client.Ping().Result()
	if err != nil {
		log4go.Error("fail to connect redis, error:", err)
		return err
	}

	log4go.Info("succeed connect to redis")
	return nil
}

func GetClient() *redis.Client {
	if client == nil {
		GetRedisClient()
	}
	return client
}

type RedisScanTriple struct {
	Cursor uint64
	Err    error
	Keys   []string
}

func (triple *RedisScanTriple) Reset() {
	triple.Cursor = 0
	triple.Keys = []string{}
	triple.Err = nil
}

func HGet(key, field string) (res string, err error) {
	return client.HGet(key, field).Result()
}
