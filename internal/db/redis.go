package db

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{Addr: cfg.Redis.Addr})
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}
	return client, nil
}
