package event

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewEventRepository,
	NewEventService, 
	NewEventHandler,
	wire.Bind(new(EventRepositoryInterface), new(*EventRepository)),
	wire.Bind(new(EventServiceInterface), new(*EventService)),
)