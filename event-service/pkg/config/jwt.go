package config

import (
	"errors"
	"time"
)

type JWTConfig struct {
	JWTSecret string        `mapstructure:"jwt_secret"`
	JWTIssuer string        `mapstructure:"jwt_issuer"`
	JWTExpiry time.Duration `mapstructure:"jwt_expiry"`
}

func (j *JWTConfig) Validate() error {
	if j.JWTSecret == "" {
		return errors.New("JWT Secret has not been set yet")
	}

	if j.JWTIssuer== "" {
		return errors.New("JWT Issuer has not been set yet")
	}

	if j.JWTExpiry == 0 {
		return errors.New("JWT Expiry has not been set or is invalid")
	}

	return nil
}