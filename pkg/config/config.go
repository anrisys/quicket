package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port string `mapstructure:"PORT"`
}

type DBConfig struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
}

type LogConfig struct {
	Level  string `mapstructure:"LOG_LEVEL"`
	Pretty bool   `mapstructure:"LOG_PRETTY"`
}

type SecurityConfig struct {
	BcryptCost int `mapstructure:"BCRYPT_COST"`
	JWTSecret string `mapstructure:"JWT_SECRET"`
	JWTIssuer string `mapstructure:"JWT_ISSUER"`
	JWTExpiry time.Duration `mapstructure:"JWT_EXPIRY"`
}

type AppConfig struct {
	Server   ServerConfig
	Logging  LogConfig
	Database DBConfig
	Security SecurityConfig
}

func DefaultConfig() *AppConfig {
	return &AppConfig{
		Server:  ServerConfig{Port: "8080"},
		Logging: LogConfig{Level: "debug", Pretty: true},
		Security: SecurityConfig{BcryptCost: 14},
	}
}

func Load() (*AppConfig, error) {
	config := DefaultConfig()

	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
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

	if config.Database == (DBConfig{}) {
		log.Fatal("Database has not been set yet")
	}

	if config.Security.JWTSecret == "" {
		log.Fatal("JWT Secret has not been set yet")
	}

	if config.Security.JWTExpiry == 0 {
		log.Fatal("JWT Expiry has not been set or is invalid")
	}

	return config, nil
}
