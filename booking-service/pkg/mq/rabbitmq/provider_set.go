package rabbitmq

import "github.com/google/wire"

var SetUpProviderSet = wire.NewSet(
	NewClient,
	NewPublisher,
	NewConsumer,
)