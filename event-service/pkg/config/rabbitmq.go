package config

import (
	"errors"
	"fmt"
)

type RabbitMQConfig struct {
    Host     string `mapstructure:"RABBITMQ_HOST" default:"localhost"`
    Port     string `mapstructure:"RABBITMQ_PORT" default:"5672"`
    User     string `mapstructure:"RABBITMQ_USER" default:"guest"`
    Password string `mapstructure:"RABBITMQ_PASSWORD" default:"guest"`
    VHost    string `mapstructure:"RABBITMQ_VHOST" default:"/"`
}

func (r *RabbitMQConfig) Validate() error {
    if r.Host == "" {
        return errors.New("rabbitmq host has not been set")
    }
    if r.Port == "" {
        return errors.New("rabbitmq port has not been set")
    }
    return nil
}

func (r *RabbitMQConfig) URL() string {
    return fmt.Sprintf("amqp://%s:%s@%s:%s%s", 
        r.User, r.Password, r.Host, r.Port, r.VHost)
}