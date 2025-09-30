package config

import "errors"

type BcryptConfig struct {
	BcryptCost int `mapstructure:"bcrypt_cost"`
}

func (b *BcryptConfig) Validate() error {
	if b.BcryptCost == 0 {
		return errors.New("bcrypt cost has not been set yet")
	}
	return nil
}