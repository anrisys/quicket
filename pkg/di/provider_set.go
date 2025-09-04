package di

import (
	"github.com/anrisys/quicket/internal/booking"
	"github.com/anrisys/quicket/internal/event"
	"github.com/anrisys/quicket/internal/infrastructure"
	"github.com/anrisys/quicket/internal/payment"
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
		wire.Bind(new(security.AccountSecurityInterface), new(*security.AccountSecurity)),
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
	)
	UserServiceClientSet = wire.NewSet(
		infrastructure.NewUserServiceClient,
		wire.Bind(new(types.UserReader), new(*infrastructure.UserServiceClient)),
	)
	AppProviderSet = wire.NewSet(
		CoreSet, 
		event.ProviderSet,
		payment.ProviderSet,
		booking.ProviderSet,
		UserServiceClientSet,
		wire.Bind(new(types.EventReader), new(*event.EventService)),
		wire.Bind(new(types.SimulatePayment), new(*payment.PaymentService)),
		wire.Struct(new(App), "*"),
	)
)