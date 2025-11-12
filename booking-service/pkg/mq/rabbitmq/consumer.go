package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type Consumer struct {
	Channel *amqp.Channel
	logger zerolog.Logger
}

func NewConsumer(client *Client, logger zerolog.Logger) (*Consumer, error) {
	ch, err := client.GetChannel()
	if err != nil {
		return nil, err
	}
	return &Consumer{Channel: ch, logger: logger}, nil
}

func (c *Consumer) DeclareExchange(name, kind string) error {
	return c.Channel.ExchangeDeclare(
		name,   // exchange name
		kind,   // type: direct, topic, fanout
		true,   // durable
		false,  // auto-deleted
		false,  // internal
		false,  // no-wait
		nil,    // arguments
	)
}

type QueueConfig struct {
    Name       string
    Durable    bool
    AutoDelete bool
    Exclusive  bool
    NoWait     bool
    Args       amqp.Table
}

func DefaultQueueConfig(name string) QueueConfig {
    return QueueConfig{
        Name:       name,
        Durable:    true,
        AutoDelete: false,
        Exclusive:  false,
        NoWait:     false,
        Args:       nil,
    }
}

func (q QueueConfig) WithDLQ(deadLetterExchange string) QueueConfig {
    if q.Args == nil {
        q.Args = make(amqp.Table)
    }
    q.Args["x-dead-letter-exchange"] = deadLetterExchange
    return q
}

func (c *Consumer) DeclareQueue(config QueueConfig) (amqp.Queue, error) {
    return c.Channel.QueueDeclare(
        config.Name,
        config.Durable,
        config.AutoDelete,
        config.Exclusive,
        config.NoWait,
        config.Args,
    )
}

// BindQueue binds the queue to an exchange with a routing key.
func (c *Consumer) BindQueue(exchange, queue, routingKey string) error {
	return c.Channel.QueueBind(
		queue,      // queue
		routingKey, // routing key
		exchange,   // exchange
		false,      // no-wait
		nil,        // args
	)
}

func (c *Consumer) StartConsuming(ctx context.Context, queueName string, handler func(amqp.Delivery)) error {
    // Set up consumer
    messages, err := c.Channel.Consume(
        queueName,
        "",    // consumer tag
        false, // auto-ack
        false, // exclusive
        false, // no-local
        false, // no-wait
        nil,
    )
    if err != nil {
        c.Channel.Close() // Close Channel if setup fails
        return err
    }
    
    // Start message processing in goroutine
    go func() {
        for msg := range messages {
            handler(msg)
        }
    }()
    
    return nil
}


// Close closes the consumer Channel.
func (c *Consumer) Close() error {
	c.logger.Info().Msg("Closing consumer channel")
	return c.Channel.Close()
}