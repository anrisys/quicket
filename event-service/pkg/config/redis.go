package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Redis struct {
	Host     string `mapstructure:"EVENT_REDIS_HOST"`
	Port     string `mapstructure:"EVENT_REDIS_PORT"`
	Password string `mapstructure:"EVENT_REDIS_PASSWORD"`
	DB       int    `mapstructure:"EVENT_REDIS_DB"`
}

func LoadRedisConfig() (*Redis, error) {
	config := &Redis{
		Host: "localhost",
		Port: "6379",
		Password: "",
		DB: 0,
	}

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AddConfigPath("../..")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No config file found, using defaults")
		} else {
			return nil, err
		}
	}

	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	if err := checkRedisConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func checkRedisConfig(config *Redis) error {
	if config.Host == "" {
		return fmt.Errorf("redis host has not been set yet")
	}
	
	if config.Port == "" {
		return fmt.Errorf("redis port has not been set or is invalid")
	}
	
	return nil
}