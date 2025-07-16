package database

import (
	"fmt"

	"github.com/anrisys/quicket/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func MySQLDB(cfg *config.AppConfig) (*gorm.DB, error) {
	cfgDB := cfg.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfgDB.DBUser,
		cfgDB.DBPassword,
		cfgDB.DBHost,
		cfgDB.DBPort,
		cfgDB.DBName,
	)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}