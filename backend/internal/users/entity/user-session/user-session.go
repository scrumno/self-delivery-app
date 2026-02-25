package user_session

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type UserSession struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;not null;index:idx_sessions_user"      json:"user_id"`
	// Случайный непрозрачный токен (не JWT). Можно мгновенно отозвать.
	Token string `gorm:"uniqueIndex:idx_sessions_token;not null" json:"token"`
	// Формат: {"ua":"Mozilla...","platform":"Android","pwa":true}
	DeviceInfo datatypes.JSON `gorm:"type:jsonb"                          json:"device_info,omitempty"`
	ExpiresAt  time.Time      `gorm:"index:idx_sessions_expires;not null" json:"expires_at"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"                      json:"created_at"`
}

func (UserSession) TableName() string { return "user_sessions" }
