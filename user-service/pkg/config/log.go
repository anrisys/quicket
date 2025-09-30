package config

type LogConfig struct {
	Level  string `mapstructure:"log_level" default:"debug"`
	Pretty bool   `mapstructure:"log_pretty" default:"true"`
}