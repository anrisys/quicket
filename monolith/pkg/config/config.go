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
	DBHost     string `mapstructure:"MONOLITH_MYSQL_HOST"`
    DBPort     string `mapstructure:"MONOLITH_MYSQL_PORT"`
    DBUser     string `mapstructure:"MONOLITH_MYSQL_USER"`
    DBPassword string `mapstructure:"MONOLITH_MYSQL_PASSWORD"`
    DBName     string `mapstructure:"MONOLITH_MYSQL_NAME"`
}

type LogConfig struct {
	Level  string `mapstructure:"log_level"`
	Pretty bool   `mapstructure:"log_pretty"`
}

type SecurityConfig struct {
	BcryptCost int `mapstructure:"bcrypt_cost"`
	JWTSecret string `mapstructure:"jwt_secret"`
	JWTIssuer string `mapstructure:"jwt_issuer"`
	JWTExpiry time.Duration `mapstructure:"jwt_expiry"`
}

type AppConfig struct {
	UserServiceURL string `mapstructure:"USER_SERVICE_URL"`
	Server   ServerConfig
	Logging  LogConfig
	Database DBConfig       `mapstructure:",squash"`
	Security SecurityConfig `mapstructure:",squash"`
}

func DefaultConfig() *AppConfig {
	return &AppConfig{
		Server:  ServerConfig{Port: "8080"},
		Logging: LogConfig{Level: "debug", Pretty: true},
		Security: SecurityConfig{BcryptCost: 14},
		Database: DBConfig{},
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