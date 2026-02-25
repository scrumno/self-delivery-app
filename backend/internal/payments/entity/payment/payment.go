package payment

import (
	"time"

	"github.com/google/uuid"
	paymentProvider "github.com/scrumno/scrumno-api/internal/payments/entity/payment-provider"
	paymentStatus "github.com/scrumno/scrumno-api/internal/payments/entity/payment-status"
)

type Payment struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"          json:"id"`
	OrderID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null;index:idx_payments_order"  json:"order_id"`
	// ID транзакции у провайдера. Используется для идентификации входящих webhook.
	ExternalID string                      `gorm:"not null;uniqueIndex:idx_payments_external_unique"  json:"external_id"`
	Status     paymentStatus.PaymentStatus `gorm:"type:varchar(20);default:'PENDING';index:idx_payments_status" json:"status"`
	Amount     float64                     `gorm:"type:numeric(12,2);not null"                        json:"amount"`
	// Провайдер нужен при возврате — выбираем правильный API.
	// Значения: tinkoff, yookassa, sbp, stripe.
	Provider paymentProvider.PaymentProvider `gorm:"type:varchar(30);not null" json:"provider"`
	// Момент фактического списания из данных webhook (не наше время).
	PaidAt    *time.Time `json:"paid_at,omitempty"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
}
