package di

import (
	"quicket/booking-service/internal"
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
		database.NewRedisClient,
		rabbitmq.ProviderSet,
	)
	AppProviderSet = wire.NewSet(
		ConfigSet,
		clients.ClientServices,
		internal.InternalProviderSet,
		wire.Struct(new(App), "*"),
	)
)