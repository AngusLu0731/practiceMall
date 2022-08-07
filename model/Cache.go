package model

import (
	"github.com/chenyahui/gin-cache/persist"
	"github.com/go-redis/redis/v8"
	"practiceMall/config"
)

var (
	RedisStore *persist.RedisStore
)

func init() {
	RedisStore = persist.NewRedisStore(redis.NewClient(&redis.Options{
		Network:  "tcp",
		Password: config.Conf.RedisPwd,
		Addr:     config.Conf.RedisConn,
		DB:       0,
	}))
}
