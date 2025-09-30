package config

import "errors"

type ServerConfig struct {
	Port string `mapstructure:"SERVER_PORT"`
}

func (s *ServerConfig) Validate() error {
	if s.Port == "" {
		return errors.New("server port has not been set")
	}
	return nil
}