package di

import (
	"github.com/anrisys/quicket/user-service/internal"
	"github.com/anrisys/quicket/user-service/pkg/config"
)

type UserServiceApp struct {
	Config  *config.Config
	Handler *internal.UserHandler
}