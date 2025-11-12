package di

import (
	"github.com/anrisys/quicket/event-service/internal"
	"github.com/anrisys/quicket/event-service/pkg/config"
)

type App struct {
	Config  *config.Config
	Handler *internal.EventHandler
}