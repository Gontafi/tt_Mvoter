package db

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"tt/config"
)

func ConnectRedisClient(ctx context.Context, cfg *config.Config) (*redis.Client, error) {
	url := fmt.Sprintf("redis://:%s@%s:%s/%s?protocol=%s",
		cfg.RedisPassword, cfg.RedisHost, cfg.RedisPort, cfg.RedisDB, cfg.RedisProtocol)
	opts, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opts)

	status := client.Ping(ctx)
	if status.Err() != nil {
		return nil, status.Err()
	}

	return client, nil
}
