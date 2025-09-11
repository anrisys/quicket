package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port string `mapstructure:"EVENT_SERVICE_PORT"`
}

type RedisServer struct {
	Host string `mapstructure:"EVENT_REDIS_HOST"`
	Port string `mapstructure:"EVENT_REDIS_PORT"`
	Password string `mapstructure:"EVENT_REDIS_PASSWORD"`
	DB int `mapstructure:"EVENT_REDIS_DB"`
}

type DBConfig struct {
	DBHost     string `mapstructure:"EVENT_SERVICE_DB_HOST"`
	DBPort     string `mapstructure:"EVENT_SERVICE_DB_PORT"`
	DBUser     string `mapstructure:"EVENT_SERVICE_DB_USER"`
	DBPassword string `mapstructure:"EVENT_SERVICE_DB_PASSWORD"`
	DBName     string `mapstructure:"EVENT_SERVICE_DB_NAME"`
}

type LogConfig struct {
	Level  string `mapstructure:"log_level"`
	Pretty bool   `mapstructure:"log_pretty"`
}

type SecurityConfig struct {
	BcryptCost int           `mapstructure:"bcrypt_cost"`
	JWTSecret  string        `mapstructure:"jwt_secret"`
	JWTIssuer  string        `mapstructure:"jwt_issuer"`
	JWTExpiry  time.Duration `mapstructure:"jwt_expiry"`
}

type AppConfig struct {
	UserServiceURL 	string `mapstructure:"USER_SERVICE_URL"`
	Server   		ServerConfig
	Logging  		LogConfig
	Database 		DBConfig       `mapstructure:",squash"`
	Security 		SecurityConfig `mapstructure:",squash"`
	Redis 			RedisServer
}

func DefaultConfig() *AppConfig {
	return &AppConfig{
		Server:   ServerConfig{},
		Logging:  LogConfig{Level: "debug", Pretty: true},
		Security: SecurityConfig{BcryptCost: 14},
		Database: DBConfig{},
		Redis: RedisServer{},
	}
}

func Load() (*AppConfig, error) {
	config := DefaultConfig()

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

	checkDatabaseConfig(config)

	checkSecurityConfig(config)

	checkClientServices(config)

	checkRedisServer(config)

	return config, nil
}

func checkDatabaseConfig(config *AppConfig) {
	if config.Database == (DBConfig{}) {
		log.Fatal("Database has not been set yet")
	}
}

func checkSecurityConfig(config *AppConfig) {
	if config.Security.BcryptCost == 0 {
		log.Fatal("BCRYPT Cose has not been set yet")
	}

	if config.Security.JWTSecret == "" {
		log.Fatal("JWT Secret has not been set yet")
	}

	if config.Security.JWTExpiry == 0 {
		log.Fatal("JWT Expiry has not been set or is invalid")
	}
}

func checkClientServices(config *AppConfig) {
	if config.UserServiceURL == "" {
		log.Fatal("USER CLIENT URL has not been set yet")
	}
}

func checkRedisServer(config *AppConfig) {
	if config.Redis.Host == "" {
		log.Fatal("Redis host has not been set yet")
	}

	if config.Redis.Port == "" {
		log.Fatal("Redis password has not been set or is invalid")
	}
}