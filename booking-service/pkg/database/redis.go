package database

import (
	"context"
	"encoding/json"
	"fmt"
	"quicket/booking-service/pkg/config"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	client *redis.Client
}

func NewRedisClient(cfg *config.Config) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr(),
		Password: cfg.Redis.Password,
		DB: cfg.Redis.DB,
	})

	return &Client{ client: rdb }
}

const BookingKey = "booking"

func (c *Client) Connect(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) Get(ctx context.Context, key string, dest any) error {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return fmt.Errorf("redis get failed: %w", err)
	}

	return json.Unmarshal([]byte(val), dest)
}

func (c *Client) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return c.client.Set(ctx, key, data, ttl).Err()
}