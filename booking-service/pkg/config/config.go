package config

type Config struct {
	MySQL  *MySQLConfig
	Log    *LogConfig
	Redis  *RedisConfig
	JWT    *JWTConfig
	Server *ServerConfig
}