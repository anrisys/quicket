package infrastructure

import (
	"context"
	"encoding/json"
	"quicket-payment-service/internal/payment"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublisher struct {
	channel  *amqp.Channel
	exchange string
}

func NewRabbitMQPublisher(conn *amqp.Connection, exchange string) (*RabbitMQPublisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQPublisher{
		channel:  ch,
		exchange: exchange,
	}, nil
}

func (r *RabbitMQPublisher) publish(
	ctx context.Context,
	routingKey string,
	payload any,
) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return r.channel.PublishWithContext(
		ctx,
		r.exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (r *RabbitMQPublisher) PaymentInitialized(ctx context.Context, payload payment.PaymentSession) error {
	return r.publish(ctx, "payment.created", payload)
}
