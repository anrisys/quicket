package dto

type PaymentDTO struct {
	PublicID  string
	Amount    float32
	Status    string
	BookingID string
}