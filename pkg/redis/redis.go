package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/riad804/go_ecommerce_api/internals/config"
)

type RedisClient struct {
	*redis.Client
}

func NewRedisClient(cfg *config.Config) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	rc := &RedisClient{Client: client}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Test the connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return rc, nil
}

func (rc *RedisClient) Close() error {
	return rc.Client.Close()
}

func (rc *RedisClient) WithContext(ctx context.Context) *RedisClient {
	return &RedisClient{Client: rc.Client.WithContext(ctx)}
}
