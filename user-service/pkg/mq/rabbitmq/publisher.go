package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

type Publisher struct {
	Channel *amqp.Channel
}

func NewPublisher(client *Client) (*Publisher, error) {
	ch, err := client.GetChannel()
	if err != nil {
		return nil, err
	}
	return &Publisher{Channel: ch}, nil
}

// DeclareExchange ensures the exchange exists before publishing.
func (p *Publisher) DeclareExchange(name, kind string) error {
	return p.Channel.ExchangeDeclare(
		name,   // exchange name
		kind,   // type: direct, topic, fanout
		true,   // durable
		false,  // auto-deleted
		false,  // internal
		false,  // no-wait
		nil,    // arguments
	)
}

// Publish sends a message to an exchange with a routing key.
func (p *Publisher) Publish(exchange, routingKey string, body []byte) error {
	return p.Channel.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

// Close closes the publisher channel.
func (p *Publisher) Close() error {
	return p.Channel.Close()
}