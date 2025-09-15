package clients

import "github.com/google/wire"

var (
	ClientServices = wire.NewSet(
		NewUserServiceClient,
	)
)