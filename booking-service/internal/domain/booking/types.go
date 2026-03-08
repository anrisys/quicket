package booking

type BookingStatus string
type PaymentMethod string
type Channel string

const (
	// Booking Status
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusSuccess   BookingStatus = "success"
	BookingStatusFailed    BookingStatus = "failed"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusExpired   BookingStatus = "expired"

	// Payment Method
	PaymentMethodCreditCard   PaymentMethod = "credit_card"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodEWallet      PaymentMethod = "e_wallet"
	PaymentMethodCash         PaymentMethod = "cash"

	// Channel
	ChannelWeb     Channel = "web"
	ChannelMobile  Channel = "mobile"
	ChannelPartner Channel = "partner"
)
