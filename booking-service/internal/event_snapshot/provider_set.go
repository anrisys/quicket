package eventsnapshot

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewEvSnapshotRepo, 
	NewEvSnapshotSrv,
	wire.Bind(new(Repository), new(*EvSnapshotRepo)),
	wire.Bind(new(Service), new(*srv)),
)