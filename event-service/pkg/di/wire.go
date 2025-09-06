//go:build wireinject
// +build wireinject

package di

import "github.com/google/wire"

func InitializeEventServiceApp() (*EventServiceApp, error) {
	wire.Build(EventAppProviderSet)
	return &EventServiceApp{}, nil
}