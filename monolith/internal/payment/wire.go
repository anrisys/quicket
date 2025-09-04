package payment

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewRepository,
	NewPaymentService,
	wire.Bind(new(Repository), new(*GormRepository)),
	wire.Bind(new(PaymentServiceInterface), new(*PaymentService)),
)