package consumer

import (
	"quicket/booking-service/pkg/mq/rabbitmq"

	"github.com/rs/zerolog"
)

type EventConsumer struct {
	consumer *rabbitmq.Consumer
	logger zerolog.Logger
}

func NewEventConsumer(consumer *rabbitmq.Consumer, logger zerolog.Logger) *EventConsumer {
	return &EventConsumer{consumer: consumer, logger: logger}
}
/*

func (ec *EventConsumer) UpdateAvailableSeats(eventID uint, seats uint) error {
	log := ec.logger.With().
		Str("consumer", "event_consumer").
		Uint("event_id", eventID).
		Uint("seats", seats).
		Logger()
	defer ec.rabbitConn.Close()

	ch, err := ec.rabbitConn.Conn.Channel()
	if err != nil {
		log.Error().Err(err).Msg("failed to open a channel")
		return fmt.Errorf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to declare queue")
		return fmt.Errorf("failed to declare queue: %v", err)
	}
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
	}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	return nil
}

func failOnError(err error, msg string) {
  if err != nil {
    log.Panicf("%s: %s", msg, err)
  }
}
  */