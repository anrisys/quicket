package dto

type SimpleBookingDTO struct {
	ID     uint
	UserID uint
	Amount float32
}

type CreatePaymentRequest struct {
	Amount    float32 `json:"amount" binding:"required"`
	BookingID string  `json:"booking_id" binding:"required"`
	Status    string  `json:"status" binding:"required,payStatus"`
}

type SimulateBookingPayment struct {
	BookingID uint
	UserID    uint
	Amount    float32
}