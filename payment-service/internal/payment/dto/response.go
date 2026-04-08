package dto

type CreatePaymentResponse struct {
	GatewayID  string
	PaymentURL string
	Status     string
}
