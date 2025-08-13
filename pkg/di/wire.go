//go:build wireinject
// +build wireinject

package di

import (
	"github.com/anrisys/quicket/internal/booking"
	"github.com/anrisys/quicket/internal/event"
	"github.com/anrisys/quicket/internal/user"
	"github.com/anrisys/quicket/pkg/config"
	"github.com/google/wire"
)

func InitializeApp() (*App, error) {
	wire.Build(
		CoreSet,
		user.ProviderSet,
		booking.ProviderSet,
		event.ProviderSet,
		wire.Struct(new(App), "*"),
	)
	return &App{}, nil
}

type App struct {
	Config 		*config.AppConfig
	UserHandler *user.UserHandler
	BookingHandler *booking.Handler
	EventHandler *event.EventHandler
}

