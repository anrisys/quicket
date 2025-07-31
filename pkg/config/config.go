package config

import (
	"log"

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

	return config, nil
}
