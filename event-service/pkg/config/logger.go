package config

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func NewZerolog(cfg *AppConfig) zerolog.Logger {
	cfgLogger := cfg.Logging
	zerolog.TimeFieldFormat = time.RFC3339Nano

	logLevel, err := zerolog.ParseLevel(cfgLogger.Level)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(logLevel)

	if cfgLogger.Pretty {
		return log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	return log.Logger
}