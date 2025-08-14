package di

import (
	"github.com/anrisys/quicket/internal/booking"
	"github.com/anrisys/quicket/internal/event"
	"github.com/anrisys/quicket/internal/user"
	"github.com/anrisys/quicket/pkg/config"
	"github.com/anrisys/quicket/pkg/config/logger"
	"github.com/anrisys/quicket/pkg/database"
	"github.com/anrisys/quicket/pkg/security"
	"github.com/anrisys/quicket/pkg/token"
	"github.com/anrisys/quicket/pkg/types"
	"github.com/google/wire"
)

var (
	ConfigSet = wire.NewSet(
		config.Load,
	)
	DatabaseSet = wire.NewSet(
		database.MySQLDB,
	)
	LoggerSet = wire.NewSet(
		logger.NewZerolog,
	)
	SecuritySet = wire.NewSet(
		security.NewAccountSecurity,
	)
	TokenSet = wire.NewSet(
		token.NewGenerator,
		wire.Bind(new(token.GeneratorInterface), new(*token.Generator)),
	)
	CoreSet = wire.NewSet(
		ConfigSet,
		DatabaseSet, 
		LoggerSet,
		SecuritySet,
		TokenSet,
		wire.Bind(new(security.AccountSecurityInterface), new(*security.AccountSecurity)),
	)
	AppProviderSet = wire.NewSet(
		CoreSet, 
		user.ProviderSet,
		event.ProviderSet,
		booking.ProviderSet,
		wire.Bind(new(types.UserReader), new(*user.UserService)),
		wire.Bind(new(types.EventReader), new(*event.EventService)),
		wire.Struct(new(App), "*"),
	)
)