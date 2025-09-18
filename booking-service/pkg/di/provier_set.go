package di

import (
	"quicket/booking-service/internal/booking"
	eventsnapshot "quicket/booking-service/internal/event_snapshot"
	"quicket/booking-service/internal/mq/consumer"
	"quicket/booking-service/internal/mq/producer"
	usersnapshot "quicket/booking-service/internal/user_snapshot"
	"quicket/booking-service/pkg/clients"
	"quicket/booking-service/pkg/config"
	"quicket/booking-service/pkg/database"
	"quicket/booking-service/pkg/mq/rabbitmq"

	"github.com/google/wire"
)

var (
	ConfigSet = wire.NewSet(
		config.Load,
		config.NewZerolog,
		database.ConnectMySQL,
	)
	SnapshotSet = wire.NewSet(
		eventsnapshot.ProviderSet,
		usersnapshot.ProviderSet,
	)
	RabbitMQSet = wire.NewSet(
		rabbitmq.SetUpProviderSet,
		producer.NewEventProducer,
		consumer.NewEventConsumer,
	)
	AppProviderSet = wire.NewSet(
		ConfigSet,
		RabbitMQSet,
		clients.ClientServices,
		SnapshotSet,
		booking.ProviderSet,
		wire.Struct(new(App), "*"),
	)
)