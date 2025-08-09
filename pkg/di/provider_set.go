package di

import (
	"github.com/anrisys/quicket/pkg/config"
	"github.com/anrisys/quicket/pkg/config/logger"
	"github.com/anrisys/quicket/pkg/database"
	"github.com/anrisys/quicket/pkg/security"
	"github.com/anrisys/quicket/pkg/token"
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
	)
	CoreSet = wire.NewSet(
		ConfigSet,
		DatabaseSet, 
		LoggerSet,
		SecuritySet,
		TokenSet,
		wire.Bind(new(security.AccountSecurityInterface), new(*security.AccountSecurity)),
	)
)