package dto

type CreatePaymentRequest struct {
	Amount    float32 `json:"amount" binding:"required"`
	BookingID string  `json:"booking_id" binding:"required"`
	Status    string  `json:"status" binding:"required,payStatus"`
}