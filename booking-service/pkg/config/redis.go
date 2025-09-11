package config

import (
	"errors"
	"fmt"
)

type RedisConfig struct {
	Host     string `mapstructure:"EVENT_REDIS_HOST"`
	Port     string `mapstructure:"EVENT_REDIS_PORT"`
	Password string `mapstructure:"EVENT_REDIS_PASSWORD"`
	DB       int    `mapstructure:"EVENT_REDIS_DB"`
}

func (r *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

func (r *RedisConfig) Validate() error {
	if r.Host == "" {
		return errors.New("redis host has not been set")
	}
	if r.Port == "" {
		return errors.New("redis port has not been set")
	}
	return nil
}