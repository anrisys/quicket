package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

/*
|--------------------------------------------------------------------------
| MySQL
|--------------------------------------------------------------------------
*/

func NewMySQL(cfg *Config) (*gorm.DB, error) {
	dsn := cfg.MySQL.DSN()

	gormLogger := logger.Default.LogMode(logger.Info)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect mysql: %w", err)
	}

	return db, nil
}

/*
|--------------------------------------------------------------------------
| Redis
|--------------------------------------------------------------------------
*/

func NewRedis(cfg *Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect redis: %w", err)
	}

	return rdb, nil
}
