package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	usersnapshot "quicket/booking-service/internal/user_snapshot"
	"quicket/booking-service/pkg/mq/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type UserConsumer struct {
	rabbitConsumer *rabbitmq.Consumer
	logger zerolog.Logger
	usrSrv usersnapshot.Service
}

func NewUserConsumer(rabbitConsumer *rabbitmq.Consumer, logger zerolog.Logger, usrSrv usersnapshot.Service) *UserConsumer {
	return &UserConsumer{
		rabbitConsumer: rabbitConsumer,
		logger: logger,
		usrSrv: usrSrv,
	}
}

const (
	userExchangeName = "user.exchange"
	userQueueName = "booking-service.users.changes"
)

func (u *UserConsumer) Start(ctx context.Context) error {
	if err := u.rabbitConsumer.DeclareExchange(userExchangeName, "topic"); err != nil {
		return fmt.Errorf("failed to declare user exchange: %w", err)
	}

	queueConfig := rabbitmq.DefaultQueueConfig(userExchangeName).WithDLQ("users.dlx")

	queue, err := u.rabbitConsumer.DeclareQueue(queueConfig)
	if err != nil {
		return fmt.Errorf("failed to declare users queue: %w", err)
	}

	routingKeys := []string{
		"users.created",
		"users.deleted",
	}

	for _, routingKey := range routingKeys {
		if err := u.rabbitConsumer.BindQueue(userExchangeName, queue.Name, routingKey); err != nil {
			return fmt.Errorf("failed to bind queue with routing key %s: %w", routingKey, err)
		}
	}

	u.logger.Info().
		Str("queue", userQueueName).
		Strs("routing_keys", routingKeys).
		Msg("user consumer setup complete")

	
	return u.rabbitConsumer.StartConsuming(ctx, queue.Name, u.handleMessage)
}

func (u *UserConsumer) handleMessage(msg amqp.Delivery) {
	log := u.logger.With().
		Str("routing_key", msg.RoutingKey).
		Str("message_id", msg.MessageId).
		Logger()

	defer func() {
		if err := recover(); err != nil {
			log.Error().Interface("error", err).Msg("panic during message processing")
			msg.Nack(false, false)
			return
		}
	}()

	switch msg.RoutingKey {
	case "user.created": 
		if err := u.handleCreatedMessage(msg); err != nil {
			log.Error().Err(err).Msg("failed to handle user creation")
			msg.Nack(false, true)
			return 
		}
	
	case "user.deletion": 
		if err := u.handleDeletedMessage(msg); err != nil {
			log.Error().Err(err).Msg("failed to handle user deletion")
			msg.Nack(false, true)
			return
		}
	
	default: 
		log.Warn().Msg("unknown routing key, acknowleging and ignoring")
	}
	
	msg.Ack(false)
}

func (u *UserConsumer) handleCreatedMessage(msg amqp.Delivery) error {
	var userMsg UserCreatedMessage	
	if err := json.Unmarshal(msg.Body, &userMsg); err != nil {
		return fmt.Errorf("failed to unmarshal user message: %w", err)
	}

	log := u.logger.With().
		Uint("user_id", userMsg.ID).
		Str("public_id", userMsg.PublicID).
		Logger()

	log.Info().Msg("Processing user creation message")

	if err := u.usrSrv.CreateUserSnapshot(context.Background(), userMsg.ID, userMsg.PublicID); err != nil {
		return fmt.Errorf("failed to create user snpashot: %w", err)
	}

	log.Info().Msg("user snapshot created successfully")

	return nil
}

func (u *UserConsumer) handleDeletedMessage(msg amqp.Delivery) error {
	var userMsg UserDeletedMessage
	if err := json.Unmarshal(msg.Body, &userMsg); err != nil {
		return fmt.Errorf("failed to unmarshal user message: %w", err)
	}

	log := u.logger.With().
		Uint("user_id", userMsg.ID).
		Logger()

	log.Info().Msg("Processing user deletion message")

	if err := u.usrSrv.DeleteUserSnapshot(context.Background(), userMsg.ID); err != nil {
		return fmt.Errorf("failed to delete user snapshot: %w", err)
	}

	log.Info().Msg("user snapshot deleted successfully")

	return nil
}