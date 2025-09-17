package di

import (
	"quicket/booking-service/internal"
	"quicket/booking-service/internal/mq/consumer"
	"quicket/booking-service/internal/mq/producer"
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
	// SnapshotSet = wire.NewSet(
	// 	internal.NewEvSnapshotRepo,
	// 	wire.Bind(new(internal.EventSnapshotRepository), new(*internal.EvSnapshotRepo)),
	// )
	RabbitMQSet = wire.NewSet(
		rabbitmq.SetUpProviderSet,
		producer.NewEventProducer,
		consumer.NewEventConsumer,
	)
	AppProviderSet = wire.NewSet(
		ConfigSet,
		RabbitMQSet,
		clients.ClientServices,
		internal.InternalProviderSet,
		wire.Struct(new(App), "*"),
	)
)