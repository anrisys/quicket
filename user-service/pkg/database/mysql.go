package database

import (
	"github.com/anrisys/quicket/user-service/pkg/config"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectMySQL(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.Mysql.DSN()

	gormLogger := logger.New(
		&log.Logger,
		logger.Config{
			LogLevel: logger.Info,
		},
	)

	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
}