package producer

import (
	"fmt"
	"quicket/booking-service/internal/mq"
	"quicket/booking-service/pkg/mq/rabbitmq"

	"github.com/rs/zerolog"
)

type EventProducer struct {
	publisher *rabbitmq.Publisher
	logger zerolog.Logger
}

func NewEventProducer(publisher *rabbitmq.Publisher) *EventProducer {
	return &EventProducer{publisher: publisher}
}

func (evp *EventProducer) PublishAvailableSeatsUpdate(eventID, seats uint) error {
	log := evp.logger.With().
		Str("producer", "event_producer").
		Str("action", "publish_available_seats_update").
		Uint("event_id", eventID).
		Uint("seats", seats).
		Logger()

	exchange := "booking.exchange"

	err := evp.publisher.DeclareExchange(exchange, "topic")
	if err != nil {
		log.Error().Err(err).Str("exchange", exchange).Msg("failed to declare exchange")
		return fmt.Errorf("%w: %v", mq.ErrFailedToDeclareExchange, err)
	}

	defer evp.publisher.Channel.Close()
	body := fmt.Appendf(nil, `{"event_id": %d, "available_seats": %d}`, eventID, seats)

	err = evp.publisher.Publish(exchange, "bookings.seats.updated", body)
	if err != nil {
		log.Error().Err(err).Str("exchange", exchange).Msg("failed to publish message")
		return err
	}

	log.Info().Msgf("Published seats updated: %s", string(body))
	return nil
}