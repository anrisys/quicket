//go:build wireinject
// +build wireinject

package di

import "github.com/google/wire"

func InitializeApp() (*App, error) {
	wire.Build(AppProviderSet)
	return &App{}, nil
}