package booking

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
		NewRepo,
		Newsrv,
		NewHandler,
		wire.Bind(new(RepositoryInterface), new(*repo)),
		wire.Bind(new(ServiceInterface), new(*srv)),
)