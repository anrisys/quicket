package di

import (
	"github.com/anrisys/quicket/pkg/config"
	"github.com/anrisys/quicket/pkg/config/logger"
	"github.com/anrisys/quicket/pkg/database"
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
	CoreSet = wire.NewSet(
		ConfigSet,
		DatabaseSet, 
		LoggerSet,
	)
)