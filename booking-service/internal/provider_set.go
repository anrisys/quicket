package internal

import (
	"github.com/google/wire"
)

var InternalProviderSet = wire.NewSet(
		NewRepo,
		Newsrv,
		NewHandler,
		wire.Bind(new(RepositoryInterface), new(*repo)),
		wire.Bind(new(ServiceInterface), new(*srv)),
)