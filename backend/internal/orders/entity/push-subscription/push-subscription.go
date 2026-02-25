package push_subscription

import (
	"time"

	"github.com/google/uuid"
)

// PushSubscription — браузерная Web Push подписка клиента.
// Ставить is_active = false если браузер вернул 410 Gone при отправке.
type PushSubscription struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"          json:"id"`
	UserID  uuid.UUID `gorm:"type:uuid;not null;index:idx_push_user_venue,priority:1"  json:"user_id"`
	VenueID uuid.UUID `gorm:"type:uuid;not null;index:idx_push_user_venue,priority:2"  json:"venue_id"`
	// URL push-сервиса браузера (Google FCM, Apple и др.).
	Endpoint string `gorm:"type:text;not null" json:"endpoint"`
	// Публичный ключ и auth-секрет для web-push шифрования payload.
	P256DH    string    `gorm:"type:text;not null" json:"p256dh"`
	Auth      string    `gorm:"type:text;not null" json:"auth"`
	IsActive  bool      `gorm:"default:true"       json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime"     json:"created_at"`
}
