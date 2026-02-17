package payment

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "PENDING"
	PaymentSuccess   PaymentStatus = "SUCCESS"
	PaymentFailed    PaymentStatus = "FAILED"
	PaymentRefunded  PaymentStatus = "REFUNDED"
	PaymentCancelled PaymentStatus = "CANCELLED"
)

type Payment struct {
	ID         uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	OrderID    uuid.UUID     `gorm:"type:uuid;unique;not null;index:idx_payments_order" json:"order_id"`
	ExternalID string        `gorm:"not null;uniqueIndex:idx_payments_external_unique" json:"external_id"`
	Status     PaymentStatus `gorm:"type:varchar;default:'PENDING';index:idx_payments_status" json:"status"`
	Amount     float64       `gorm:"type:numeric(12,2);not null" json:"amount"`
	Provider   *string       `json:"provider,omitempty"`
	PaidAt     *time.Time    `json:"paid_at,omitempty"`
	CreatedAt  time.Time     `gorm:"default:now()" json:"created_at"`
}
