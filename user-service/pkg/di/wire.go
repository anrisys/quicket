//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
)

func InitializeUserServiceApp() (*UserServiceApp, error) {
	wire.Build(UserAppProviderSet)
	return &UserServiceApp{}, nil
}
