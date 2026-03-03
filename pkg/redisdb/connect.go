package redisdb

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisConnectionConfig struct {
	Addr     string //"localhost:6379",
	Password string //"" no password set
	DB       int    // 0 use default DB
}

func NewRedisConnection(cfg *RedisConnectionConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := pingRedis(rdb); err != nil {
		return nil, err
	}

	return rdb, nil
}

func pingRedis(rdb *redis.Client) error {
	_, err := rdb.Ping(context.Background()).Result()
	return err
}
