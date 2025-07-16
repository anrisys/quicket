//go:build wireinject
// +build wireinject

package di

import (
	"github.com/anrisys/quicket/internal/user"
	"github.com/google/wire"
)

func InitializeApp() (*App, err) {
	wire.Build(
		CoreSet,
		user.ProviderSet,
		wire.Struct(new(App), "*")
	)
	return &App{}, nil
}

type App struct {
	UserHandler *user.Handler
}

