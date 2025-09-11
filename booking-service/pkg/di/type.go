package di

import (
	"quicket/booking-service/internal"
	"quicket/booking-service/pkg/config"
)

type App struct {
	Config *config.Config
	Handler *internal.Handler
}