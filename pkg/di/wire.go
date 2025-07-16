//go:build wireinject
// +build wireinject

package di

import (
	"github.com/anrisys/quicket/internal/user"
	"github.com/anrisys/quicket/pkg/config"
	"github.com/google/wire"
)

func InitializeApp() (*App, error) {
	wire.Build(
		CoreSet,
		user.ProviderSet,
		wire.Struct(new(App), "*"),
	)
	return &App{}, nil
}

type App struct {
	Config 		*config.AppConfig
	UserHandler *user.UserHandler
}

