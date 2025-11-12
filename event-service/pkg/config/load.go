package config

import (
	"github.com/spf13/viper"
)

func Load() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AddConfigPath("../..")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var mysqlConfig MySQLConfig
	if err := viper.Unmarshal(&mysqlConfig); err != nil {
		return nil, err
	}

	var logConfig LogConfig
	if err := viper.Unmarshal(&logConfig); err != nil {
		return nil, err
	}

	var redisConfig RedisConfig
	if err := viper.Unmarshal(&redisConfig); err != nil {
		return nil, err
	}

	var jwtConfig JWTConfig
	if err := viper.Unmarshal(&jwtConfig); err != nil {
		return nil, err
	}

	var serverConfig ServerConfig
	if err := viper.Unmarshal(&serverConfig); err != nil {
		return nil, err
	}

	var clientsConfig ClientServices
	if err := viper.Unmarshal(&clientsConfig); err != nil {
		return nil, err
	}

	var rabbitMQConfig RabbitMQConfig
	if err := viper.Unmarshal(&rabbitMQConfig); err != nil {
		return nil, err
	}

	cfg := &Config{
		MySQL:  &mysqlConfig,
		Log:    &logConfig,
		Redis:  &redisConfig,
		JWT:    &jwtConfig,
		Server: &serverConfig,
		Clients: &clientsConfig,
		RabbitMQ: &rabbitMQConfig,
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func validateConfig(config *Config) error {
	if err := config.MySQL.Validate(); err != nil {
		return err
	}
	if err := config.Redis.Validate(); err != nil {
		return err
	}
	if err := config.JWT.Validate(); err != nil {
		return err
	}
	if err := config.Server.Validate(); err != nil {
		return err
	}
	if err := config.Clients.Validate(); err != nil {
		return err
	}
	if err := config.RabbitMQ.Validate(); err != nil {
		return err
	}
	return nil
}