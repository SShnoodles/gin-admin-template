package config

import (
	"context"
	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var CTX = context.Background()

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     AppConfig.Redis.Addr,
		Password: AppConfig.Redis.Password,
		DB:       AppConfig.Redis.Db,
	})
}
