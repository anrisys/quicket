package config

import (
	"errors"
	"fmt"
)

type MySQLConfig struct {
	Host     string `mapstructure:"USER_MYSQL_HOST" default:"localhost"`
	Port     string `mapstructure:"USER_MYSQL_PORT" default:"3306"`
	User     string `mapstructure:"USER_MYSQL_USER" default:"root"`
	Password string `mapstructure:"USER_MYSQL_PASSWORD" default:""`
	Database string `mapstructure:"USER_MYSQL_DATABASE"`
	Charset  string `mapstructure:"USER_MYSQL_CHARSET" default:"utf8mb4"`
}

func (m *MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True",
		m.User, m.Password, m.Host, m.Port, m.Database, m.Charset)
}

func (m *MySQLConfig) Validate() error {
	if m.Database == "" {
		return errors.New("mysql database has not been set")
	}
	return nil
}