package rabbitmq

import (
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

func (c *Consumer) DeclareQueue(queue string) error {
	_, err := c.Channel.QueueDeclare(
		queue, // name
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	return err
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

func (c *Consumer) StartConsuming(queueName string, handler func(amqp.Delivery)) error {
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
        defer c.Channel.Close() // Close Channel when consumer stops
        
        for msg := range messages {
            handler(msg)
        }
    }()
    
    return nil
}


// Close closes the consumer Channel.
func (c *Consumer) Close() error {
	return c.Channel.Close()
}