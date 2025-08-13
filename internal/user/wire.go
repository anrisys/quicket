package user

import (
	"github.com/anrisys/quicket/internal/booking"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewUserRepository,
	NewUserService,
	NewUserHandler,
	wire.Bind(new(UserRepositoryInterface), new(*UserRepository)),
	wire.Bind(new(UserServiceInterface), new(*UserService)),
	wire.Bind(new(booking.UserReader), new(*UserService)),
)