package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/pangjianan/qq_robot/conf"
)

var GlobalRedis *redis.Client

func Init(config *conf.Config) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})
	GlobalRedis = rdb
}
