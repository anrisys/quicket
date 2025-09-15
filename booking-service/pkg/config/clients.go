package config

import "errors"

type ClientServices struct {
	UserServiceURL string `mapstructure:"USER_SERVICE_URL"`
}

func NewClientServices(UserServiceURL string) *ClientServices {
	return &ClientServices{
		UserServiceURL: UserServiceURL,
	}
}

func (cli *ClientServices) Validate() error {
	if cli.UserServiceURL == "" {
		return errors.New("user service url has not been set")
	}
	return nil
}