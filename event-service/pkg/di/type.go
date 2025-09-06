package di

import (
	"github.com/anrisys/quicket/event-service/internal"
	"github.com/anrisys/quicket/event-service/pkg/config"
)

type EventServiceApp struct {
	Config  *config.AppConfig
	Handler *internal.EventHandler
}