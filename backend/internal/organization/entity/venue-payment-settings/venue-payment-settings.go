package venue_payment_settings

import (
	"time"

	"github.com/google/uuid"
	paymentProvider "github.com/scrumno/scrumno-api/internal/payments/entity/payment-provider"
)

// VenuePaymentSettings — настройки платёжного провайдера для конкретной точки (1:1).
//
// credentials_encrypted содержит провайдер-специфичный JSON (шифровать через pgcrypto/Vault):
//
//	tinkoff:  {"terminal_key":"...", "password":"..."}
//	yookassa: {"shop_id":"...", "secret_key":"..."}
//	sbp:      {"merchant_id":"...", "api_key":"..."}
//	stripe:   {"publishable_key":"...", "secret_key":"..."}
//
// Ключ шифрования — в env. Никогда не логировать расшифрованный payload.
type VenuePaymentSettings struct {
	ID       uuid.UUID                       `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	VenueID  uuid.UUID                       `gorm:"type:uuid;not null;uniqueIndex"                  json:"venue_id"`
	Provider paymentProvider.PaymentProvider `gorm:"type:varchar(30);not null"                      json:"provider"`
	// Зашифрованный JSON с credentials провайдера.
	CredentialsEncrypted string `gorm:"type:text;not null" json:"credentials_encrypted"`
	// false = настройки сохранены но провайдер не используется (напр. при смене).
	IsActive  bool       `gorm:"default:true"   json:"is_active"`
	UpdatedAt *time.Time `                       json:"updated_at,omitempty"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

func (VenuePaymentSettings) TableName() string { return "venue_payment_settings" }
