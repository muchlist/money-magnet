package cache

import (
	"github.com/muchlist/moneymagnet/cfg"
	"github.com/redis/go-redis/v9"
)

func InitRedis(cfg *cfg.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.RedisURL,
		Password: cfg.Redis.RedisPass,
		DB:       cfg.Redis.RedisDB,
	})
}
