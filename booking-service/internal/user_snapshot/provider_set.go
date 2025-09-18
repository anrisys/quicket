package usersnapshot

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewRepo,
	NewSrv,
	wire.Bind(new(Repository), new(*repo)),
	wire.Bind(new(Service), new(*srv)),
)