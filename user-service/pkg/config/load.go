package config

import "github.com/spf13/viper"

func Load() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AddConfigPath("../..")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var serverConfig ServerConfig
	if err := viper.Unmarshal(&serverConfig); err != nil {
		return nil, err
	}

	var mysqlConfig MySQLConfig
	if err := viper.Unmarshal(&mysqlConfig); err != nil {
		return nil, err
	}

	var logConfig LogConfig
	if err := viper.Unmarshal(&logConfig); err != nil {
		return nil, err
	}

	var jwtConfig JWTConfig
	if err := viper.Unmarshal(&jwtConfig); err != nil {
		return nil, err
	}

	var bcryptConfig BcryptConfig
	if err := viper.Unmarshal(&bcryptConfig); err != nil {
		return nil, err
	}

	var rabbitMQConfig RabbitMQConfig
	if err := viper.Unmarshal(&rabbitMQConfig); err != nil {
		return nil, err
	}

	cfg := &Config{
		Bcrypt: &bcryptConfig,
		Server: &serverConfig,
		Mysql: &mysqlConfig,
		Log: &logConfig,
		JWT: &jwtConfig,
		RabbitMQConfig: &rabbitMQConfig,
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func validateConfig(config *Config) error {
	if err := config.Server.Validate(); err != nil {
		return err
	}
	if err := config.Mysql.Validate(); err != nil {
		return err
	}
	if err := config.JWT.Validate(); err != nil {
		return err
	}
	if err := config.Bcrypt.Validate(); err != nil {
		return err
	}
	if err := config.RabbitMQConfig.Validate(); err != nil {
		return err
	}
	return nil
}