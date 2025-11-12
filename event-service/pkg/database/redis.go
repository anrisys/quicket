package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/anrisys/quicket/event-service/pkg/config"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(cfg *config.Config) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr(),
		Password: cfg.Redis.Password,
		DB: cfg.Redis.DB,
	})

	return &RedisClient{ client: rdb }
}

const EventKey = "event"

func (c *RedisClient) Connect(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

func (c *RedisClient) Close() error {
	return c.client.Close()
}

func (c *RedisClient) Get(ctx context.Context, key string, dest any) error {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return fmt.Errorf("redis get failed: %w", err)
	}

	return json.Unmarshal([]byte(val), dest)
}

func (c *RedisClient) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return c.client.Set(ctx, key, data, ttl).Err()
}