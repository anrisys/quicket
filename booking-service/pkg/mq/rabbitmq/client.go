package rabbitmq

import (
	"fmt"
	"quicket/booking-service/pkg/config"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type Client struct {
    conn    *amqp.Connection
    logger  zerolog.Logger
}

func NewClient(config *config.RabbitMQConfig, logger zerolog.Logger) (*Client, error) {
    conn, err := amqp.Dial(config.URL())
    if err != nil {
		logger.Error().Err(err).Msg("failed to connect to RabbitMQ")
        return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
    }
    logger.Info().Msg("RabbitMQ connected successfully")
    
    return &Client{
        conn:    conn,
        logger:  logger,
    }, nil
}

func (c *Client) GetConnection() *amqp.Connection {
    return c.conn
}

func (c *Client) GetChannel() (*amqp.Channel, error) {
    channel, err := c.conn.Channel()
    if err != nil {
        return nil, fmt.Errorf("failed to create channel: %w", err)
    }
    return channel, nil
}

func (c *Client) Close() error {
    if err := c.conn.Close(); err != nil {
        c.logger.Warn().Err(err).Msg("Failed to close RabbitMQ connection")
        return err
    }
    
    c.logger.Info().Msg("RabbitMQ connection closed")
    return nil
}
/*
func (c *Client) DeclareExchange(name, kind string) error {
    return c.channel.ExchangeDeclare(
        name,
        kind,
        true,  // durable
        false, // auto-deleted
        false, // internal
        false, // no-wait
        nil,   // arguments
    )
}

func (c *Client) DeclareQueue(name string, args amqp.Table) (amqp.Queue, error) {
    return c.channel.QueueDeclare(
        name,
        true,  // durable
        false, // delete when unused
        false, // exclusive
        false, // no-wait
        args,  // arguments
    )
}
	*/