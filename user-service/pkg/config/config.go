package config

type Config struct {
	Bcrypt *BcryptConfig
	Server *ServerConfig
	Mysql *MySQLConfig
	Log *LogConfig
	JWT *JWTConfig
}