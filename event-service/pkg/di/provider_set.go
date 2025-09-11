package di

import (
	"github.com/anrisys/quicket/event-service/internal"
	"github.com/anrisys/quicket/event-service/pkg/config"
	"github.com/anrisys/quicket/event-service/pkg/database"
	"github.com/anrisys/quicket/event-service/pkg/redis"
	"github.com/google/wire"
)

var (
	ConfigSet = wire.NewSet(
		config.Load,
		database.ConnectMySQL,
		config.NewZerolog,
		config.LoadRedisConfig,
		redis.NewClient,
	)
	EventAppProviderSet = wire.NewSet(
		ConfigSet,
		internal.NewEventRepository,
		internal.NewUserServiceClient,
		internal.NewEventService,
		internal.NewEventHandler,
		wire.Bind(new(internal.UserReader), new(*internal.UserServiceClient)),
		wire.Bind(new(internal.EventRepositoryInterface), new(*internal.EventRepository)),
		wire.Bind(new(internal.EventServiceInterface), new(*internal.EventService)),
		wire.Struct(new(EventServiceApp), "*"),
	)
)