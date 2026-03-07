package infrastructure

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	MySQL    MySQLConfig
	Log      LogConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Server   ServerConfig
	Clients  ClientServices
	RabbitMQ RabbitMQConfig
}

/*
|--------------------------------------------------------------------------
| MySQL
|--------------------------------------------------------------------------
*/

type MySQLConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Charset  string
}

func (m *MySQLConfig) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True",
		m.User, m.Password, m.Host, m.Port, m.Database, m.Charset,
	)
}

func (m *MySQLConfig) Validate() error {
	if m.Host == "" {
		return errors.New("mysql host is required")
	}
	if m.Database == "" {
		return errors.New("mysql database is required")
	}
	return nil
}

/*
|--------------------------------------------------------------------------
| Log
|--------------------------------------------------------------------------
*/

type LogConfig struct {
	Level  string
	Pretty bool
}

/*
|--------------------------------------------------------------------------
| Redis
|--------------------------------------------------------------------------
*/

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func (r *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

func (r *RedisConfig) Validate() error {
	if r.Host == "" {
		return errors.New("redis host is required")
	}
	if r.Port == "" {
		return errors.New("redis port is required")
	}
	return nil
}

/*
|--------------------------------------------------------------------------
| JWT
|--------------------------------------------------------------------------
*/

type JWTConfig struct {
	JWTSecret string
	JWTIssuer string
	JWTExpiry time.Duration
}

func (j *JWTConfig) Validate() error {
	if j.JWTSecret == "" {
		return errors.New("jwt secret is required")
	}
	if j.JWTIssuer == "" {
		return errors.New("jwt issuer is required")
	}
	if j.JWTExpiry == 0 {
		return errors.New("jwt expiry must be set")
	}
	return nil
}

/*
|--------------------------------------------------------------------------
| Server
|--------------------------------------------------------------------------
*/

type ServerConfig struct {
	Port string
}

func (s *ServerConfig) Validate() error {
	if s.Port == "" {
		return errors.New("server port is required")
	}
	return nil
}

/*
|--------------------------------------------------------------------------
| Client Services
|--------------------------------------------------------------------------
*/

type ClientServices struct {
	UserServiceURL  string
	EventServiceURL string
}

func (c *ClientServices) Validate() error {
	if c.UserServiceURL == "" {
		return errors.New("user service url is required")
	}
	if c.EventServiceURL == "" {
		return errors.New("event service url is required")
	}
	return nil
}

/*
|--------------------------------------------------------------------------
| RabbitMQ
|--------------------------------------------------------------------------
*/

type RabbitMQConfig struct {
	Host                   string
	Port                   string
	User                   string
	Password               string
	VHost                  string
	BookingServiceExchange string
}

func (r *RabbitMQConfig) URL() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%s%s",
		r.User, r.Password, r.Host, r.Port, r.VHost,
	)
}

func (r *RabbitMQConfig) Validate() error {
	if r.Host == "" {
		return errors.New("rabbitmq host is required")
	}
	if r.Port == "" {
		return errors.New("rabbitmq port is required")
	}
	return nil
}

/*
|--------------------------------------------------------------------------
| Load
|--------------------------------------------------------------------------
*/

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	cfg := &Config{
		MySQL: MySQLConfig{
			Host:     getEnv("BOOKING_MYSQL_HOST", "localhost"),
			Port:     getEnv("BOOKING_MYSQL_PORT", "3306"),
			User:     getEnv("BOOKING_MYSQL_USER", "root"),
			Password: os.Getenv("BOOKING_MYSQL_PASSWORD"),
			Database: getRequiredEnv("BOOKING_MYSQL_DB_NAME"),
			Charset:  getEnv("BOOKING_MYSQL_CHARSET", "utf8mb4"),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "debug"),
			Pretty: getBoolEnv("LOG_PRETTY", true),
		},
		Redis: RedisConfig{
			Host:     getRequiredEnv("REDIS_HOST"),
			Port:     getRequiredEnv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       getIntEnv("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			JWTSecret: getRequiredEnv("JWT_SECRET"),
			JWTIssuer: getRequiredEnv("JWT_ISSUER"),
			JWTExpiry: getDurationEnv("JWT_EXPIRY", time.Hour),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Clients: ClientServices{
			UserServiceURL: getRequiredEnv("USER_SERVICE_URL"),
		},
		RabbitMQ: RabbitMQConfig{
			Host:                   getEnv("RABBITMQ_HOST", "localhost"),
			Port:                   getEnv("RABBITMQ_PORT", "5672"),
			User:                   getEnv("RABBITMQ_USER", "guest"),
			Password:               getEnv("RABBITMQ_PASSWORD", "guest"),
			VHost:                  getEnv("RABBITMQ_VHOST", "/"),
			BookingServiceExchange: getRequiredEnv("RABBITMQ_BOOKING_EXCHANGE"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if err := c.MySQL.Validate(); err != nil {
		return err
	}
	if err := c.Redis.Validate(); err != nil {
		return err
	}
	if err := c.JWT.Validate(); err != nil {
		return err
	}
	if err := c.Server.Validate(); err != nil {
		return err
	}
	if err := c.Clients.Validate(); err != nil {
		return err
	}
	if err := c.RabbitMQ.Validate(); err != nil {
		return err
	}
	return nil
}

/*
|--------------------------------------------------------------------------
| Helpers
|--------------------------------------------------------------------------
*/

func getRequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic("missing required environment variable: " + key)
	}
	return value
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getIntEnv(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		panic("invalid integer for environment variable: " + key)
	}
	return value
}

func getBoolEnv(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		panic("invalid boolean for environment variable: " + key)
	}
	return value
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		panic("invalid duration for environment variable: " + key)
	}
	return value
}
