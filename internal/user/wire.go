package user

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewUserRepository,
	NewUserService,
	NewUserHandler,
	wire.Bind(new(UserRepositoryInterface), new(*UserRepository)),
	wire.Bind(new(UserServiceInterface), new(*UserService)),
)