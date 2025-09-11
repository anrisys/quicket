package di

import (
	"quicket/booking-service/internal"
	"quicket/booking-service/pkg/config"
	"quicket/booking-service/pkg/database"

	"github.com/google/wire"
)

var (
	ConfigSet = wire.NewSet(
		config.Load,
		config.NewZerolog,
		database.ConnectMySQL,
		database.NewRedisClient,
	)
	AppProviderSet = wire.NewSet(
		ConfigSet,
		internal.InternalProviderSet,
		wire.Struct(new(App), "*"),
	)
)