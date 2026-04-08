package dto

import "time"

type InitPaymentRequest struct {
	BookingID string    `json:"booking_id" validate:"required,uuid"`
	Amount    int64     `json:"amount" validate:"required,gt=0"`
	Currency  string    `json:"currency" validate:"required"`
	ExpiredAt time.Time `json:"expired_at" validate:"required"`
}

type CreatePaymentRequest struct {
	ExternalID string
	Amount     int64
	Currency   string
	Expiry     time.Time
}

/*
type CreatePaymentRequest struct {
	ExternalID  string                `json:"external_id"`
	Amount      int                   `json:"amount"`
	Currency    string                `json:"currency"`
	Description string                `json:"description"`
	ExpiryDate  time.Time             `json:"expiry_date"`
	Customer    CreatePaymentCustomer `json:"customer"`
	Metadata    CreatePaymentMetadata `json:"booking_id"`
}

type CreatePaymentCustomer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreatePaymentMetadata struct {
	BookingID string `json:"booking_id"`
}
*/
