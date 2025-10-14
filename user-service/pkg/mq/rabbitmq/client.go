package rabbitmq

import (
	"fmt"

	"github.com/anrisys/quicket/user-service/pkg/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type Client struct {
	conn *amqp.Connection
	logger zerolog.Logger
}

func NewClient(config *config.Config, logger zerolog.Logger) (*Client, error) {
	conn, err := amqp.Dial(config.RabbitMQConfig.URL())
	if err != nil {
		logger.Error().Err(err).Msg("failed to connect to RabbitMQ")
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	logger.Info().Msg("RabbitMQ connected successfully")

	return &Client{
		conn: conn,
		logger: logger,
	}, nil
}

func (c *Client) GetChannel() (*amqp.Channel, error) {
	ch, err := c.conn.Channel()
	    if err != nil {
        return nil, fmt.Errorf("failed to create channel: %w", err)
    }
    return ch, nil
}