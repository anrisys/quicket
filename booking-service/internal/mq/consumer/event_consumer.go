package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	eventsnapshot "quicket/booking-service/internal/event_snapshot"
	"quicket/booking-service/pkg/mq/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

const (
	exchangeName = "events.exchange"
	queueName = "booking-service.events.changes"
)

type EventConsumer struct {
	rabbitConsumer *rabbitmq.Consumer
	logger zerolog.Logger
	evSrv eventsnapshot.Service
}

func NewEventConsumer(consumer *rabbitmq.Consumer, logger zerolog.Logger, evSrv eventsnapshot.Service) *EventConsumer {
	return &EventConsumer{
		rabbitConsumer: consumer,
		logger: logger,
		evSrv: evSrv,
	}
}

func (c *EventConsumer) Start(ctx context.Context) error {
	if err := c.rabbitConsumer.DeclareExchange(exchangeName, "topic"); err != nil {
		return fmt.Errorf("failed to declare event exchange: %w", err)
	}

	queueConfig := rabbitmq.DefaultQueueConfig(queueName).WithDLQ("events.dlx")

	queue, err := c.rabbitConsumer.DeclareQueue(queueConfig)
	if err != nil {
		return fmt.Errorf("failed to declare events queue: %w", err)
	}

	routingKeys := []string{
        "event.created",
        "event.updated",
        "event.deleted",
        "event.seats.updated",
    }

	for _, routingKey := range routingKeys {
        if err := c.rabbitConsumer.BindQueue("events.exchange", queue.Name, routingKey); err != nil {
            return fmt.Errorf("failed to bind queue with routing key %s: %w", routingKey, err)
        }
    }

    c.logger.Info().
        Str("queue", queue.Name).
        Strs("routing_keys", routingKeys).
        Msg("Event consumer setup complete")

    // Start consuming messages
    return c.rabbitConsumer.StartConsuming(ctx, queue.Name, c.handleMessage)
}

func (c *EventConsumer) handleMessage(msg amqp.Delivery) {
	log := c.logger.With().
			Str("routing_key", msg.RoutingKey).
			Str("message_id", msg.MessageId).
			Logger()
	
	// Acknowledge message when done (or nack on error)
    defer func() {
        if err := recover(); err != nil {
            log.Error().Interface("error", err).Msg("Panic during message processing")
            msg.Nack(false, false) // Discard message
            return
        }
    }()

	    // Route message based on routing key
    switch msg.RoutingKey {
        case "event.created":
            if err := c.handleEventCreated(msg); err != nil {
                log.Error().Err(err).Msg("Failed to handle event creation")
                msg.Nack(false, true) // Requeue for retry
                return
            }

        case "event.updated":
            if err := c.handleEventUpdated(msg); err != nil {
                log.Error().Err(err).Msg("Failed to handle event update")
                msg.Nack(false, true)
                return
            }

        case "event.deleted":
            if err := c.handleEventDeleted(msg); err != nil {
                log.Error().Err(err).Msg("Failed to handle event deletion")
                msg.Nack(false, true)
                return
            }

        case "event.seats.updated":
            if err := c.handleSeatsUpdated(msg); err != nil {
                log.Error().Err(err).Msg("Failed to handle seats update")
                msg.Nack(false, true)
                return
            }

        default:
            log.Warn().Msg("Unknown routing key, acknowledging and ignoring")
    }

    // Acknowledge successful processing
    msg.Ack(false)
}

// handleEventCreated processes event creation messages
func (c *EventConsumer) handleEventCreated(msg amqp.Delivery) error {
    var eventMsg EventCreatedMessage
    if err := json.Unmarshal(msg.Body, &eventMsg); err != nil {
        return fmt.Errorf("failed to unmarshal event message: %w", err)
    }

    log := c.logger.With().
        Uint("event_id", eventMsg.ID).
        Str("public_id", eventMsg.PublicID).
        Logger()

    log.Info().Msg("Processing event creation message")

    // Create event snapshot in local database
    eventSnapshot := eventsnapshot.EventSnapshot{
        ID:             eventMsg.ID,
        PublicID:       eventMsg.PublicID,
        Title:          eventMsg.Title,
        StartDate:      eventMsg.StartDate,
        EndDate:        eventMsg.EndDate,
        AvailableSeats: eventMsg.AvailableSeats,
        Version:        eventMsg.Version,
        UpdatedAt:    	eventMsg.CreatedAt,
    }

    if err := c.evSrv.CreateSnapshot(context.Background(), &eventSnapshot); err != nil {
        return fmt.Errorf("failed to create event snapshot: %w", err)
    }

    log.Info().
        Uint64("available_seats", eventMsg.AvailableSeats).
        Msg("Event snapshot created successfully")

    return nil
}

// handleEventUpdated handles event update messages
func (c *EventConsumer) handleEventUpdated(msg amqp.Delivery) error {
    var eventMsg EventUpdatedMessage
    if err := json.Unmarshal(msg.Body, &eventMsg); err != nil {
        return fmt.Errorf("failed to unmarshal event message: %w", err)
    }

    log := c.logger.With().
        Uint("event_id", eventMsg.ID).
        Str("public_id", eventMsg.PublicID).
        Logger()

    log.Info().Msg("Processing event update message")

    eventSnapshot := eventsnapshot.EventSnapshot{
        ID:             eventMsg.ID,
        PublicID:       eventMsg.PublicID,
        Title:          eventMsg.Title,
        StartDate:      eventMsg.StartDate,
        EndDate:        eventMsg.EndDate,
        AvailableSeats: eventMsg.AvailableSeats,
        Version:        eventMsg.Version,
        UpdatedAt:    	eventMsg.CreatedAt,
    }

    if err := c.evSrv.UpdateSnapshot(context.Background(), &eventSnapshot); err != nil {
        return fmt.Errorf("failed to update event snapshot: %w", err)
    }

    log.Info().
        Uint64("available_seats", eventMsg.AvailableSeats).
        Msg("Event snapshot updated successfully")

    return nil
}

// handleEventDeleted handles event deletion messages
func (c *EventConsumer) handleEventDeleted(msg amqp.Delivery) error {
    var deleteMsg struct {
        EventID uint `json:"event_id"`
    }

    if err := json.Unmarshal(msg.Body, &deleteMsg); err != nil {
        return fmt.Errorf("failed to unmarshal delete message: %w", err)
    }

    if err := c.evSrv.DeleteSnapshot(context.Background(), deleteMsg.EventID); err != nil {
        return fmt.Errorf("failed to delete event snapshot: %w", err)
    }

    c.logger.Info().
        Uint("event_id", deleteMsg.EventID).
        Msg("Event snapshot deleted successfully")

    return nil
}

// handleSeatsUpdated handles seat availability updates
func (c *EventConsumer) handleSeatsUpdated(msg amqp.Delivery) error {
    var seatsMsg struct {
        EventID        uint `json:"event_id"`
        AvailableSeats int    `json:"available_seats"`
        Version        int    `json:"version"`
    }

    if err := json.Unmarshal(msg.Body, &seatsMsg); err != nil {
        return fmt.Errorf("failed to unmarshal seats message: %w", err)
    }

    if err := c.evSrv.UpdateSeatsSnapshot(
        context.Background(),
        seatsMsg.EventID,
        seatsMsg.AvailableSeats,
        seatsMsg.Version,
    ); err != nil {
        return fmt.Errorf("failed to update event seats: %w", err)
    }

    c.logger.Info().
        Uint("event_id", seatsMsg.EventID).
        Int("available_seats", seatsMsg.AvailableSeats).
        Int("version", seatsMsg.Version).
        Msg("Event seats updated successfully")

    return nil
}

// Stop gracefully stops the consumer
func (c *EventConsumer) Stop() {
    c.logger.Info().Msg("Event consumer stopping")
}