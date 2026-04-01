package infrastructure

import (
	"context"
	"encoding/json"
	"quicket/booking-service/internal/booking"

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

func (r *RabbitMQPublisher) PublishBookingCreated(
	ctx context.Context,
	event booking.BookingCreatedEvent,
) error {
	return r.publish(ctx, "booking.created", event)
}

func (r *RabbitMQPublisher) PublishBookingCancelled(
	ctx context.Context,
	event booking.BookingCancelledEvent,
) error {
	return r.publish(ctx, "booking.cancelled", event)
}

func (r *RabbitMQPublisher) ReleaseEventSeats(ctx context.Context, payload *booking.BookingReleaseSeats) error {
	return r.publish(ctx, "booking.seats.release", payload)
}

func (r *RabbitMQPublisher) RefundRequired(ctx context.Context, payload booking.RefundRequiredEvent) error {
	return r.publish(ctx, "payment.refund.required", payload)
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
