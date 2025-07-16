package user

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewRepository,
	NewService,
	NewHandler,
	wire.Bind(new(RepositoryInterface), new(*repository)),
	wire.Bind(new(ServiceInterface), new(*service)),
)