package payment

import (
	"quicket-payment-service/internal/domain"
	"time"

	"gorm.io/datatypes"
)

type PaymentSession struct {
	ID               string               `gorm:"column:ID"`
	BookingID        string               `gorm:"column:booking_id"`
	Amount           int64                `gorm:"column:amount"`
	Currency         string               `gorm:"column:currency"`
	Status           domain.PaymentStatus `gorm:"column:status"`
	ExpiredAt        time.Time            `gorm:"column:expired_at"`
	FailureReason    *string              `gorm:"column:failure_reason"`
	GatewayReference *string              `gorm:"column:gateway_reference"`
}

// Append only table
type PaymentLogs struct {
	ID               int64          `gorm:"column:ID"`
	PaymentSessionID string         `gorm:"column:payment_session_id"`
	Direction        string         `gorm:"direction"`   // inbound (webhook received) or outbound (refund request)
	HttpStatus       *int           `gorm:"http_status"` // for outbound calls
	Payload          datatypes.JSON `gorm:"column:payload"`
	CreatedAt        time.Time
}
