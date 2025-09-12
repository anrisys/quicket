package internal

import (
	"github.com/google/wire"
)

var InternalProviderSet = wire.NewSet(
		NewRepo,
		Newsrv,
		NewHandler,
		NewUsrReader,
		NewEvReader,
		wire.Bind(new(RepositoryInterface), new(*repo)),
		wire.Bind(new(UserReader), new (*UsrReader)),
		wire.Bind(new(EventReader), new(*EvReader)),
		wire.Bind(new(ServiceInterface), new(*srv)),
)