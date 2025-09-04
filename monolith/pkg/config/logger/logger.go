package logger

import (
	"os"
	"time"

	"github.com/anrisys/quicket/pkg/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)
func NewZerolog(cfg *config.AppConfig) zerolog.Logger {
	cfgLogger := cfg.Logging
	zerolog.TimeFieldFormat = time.RFC3339Nano

	logLevel, err := zerolog.ParseLevel(cfgLogger.Level)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(logLevel)

	if cfgLogger.Pretty {
		// Development purpose
		return log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	return log.Logger
}