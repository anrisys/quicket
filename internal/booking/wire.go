package booking

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewGormRepository,
	NewService,
	NewHandler,
	wire.Bind(new(Repository), new(*GormRepository)),
	wire.Bind(new(Service), new(*service)),
)
