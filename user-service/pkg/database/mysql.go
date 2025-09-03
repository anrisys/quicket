package database

import (
	"fmt"

	"github.com/anrisys/quicket/user-service/pkg/config"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectMySQL(cfg *config.AppConfig) (*gorm.DB, error) {
	cfgDB := cfg.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
		cfgDB.DBUser,
		cfgDB.DBPassword,
		cfgDB.DBHost,
		cfgDB.DBPort,
		cfgDB.DBName,
	)

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