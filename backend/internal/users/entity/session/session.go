package session

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Session struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID     uuid.UUID      `gorm:"type:uuid;not null;index:idx_sessions_user" json:"user_id"`
	Token      string         `gorm:"type:text;not null;uniqueIndex:idx_sessions_token" json:"token"`
	DeviceInfo datatypes.JSON `gorm:"type:jsonb" json:"device_info"`
	ExpiresAt  time.Time      `gorm:"not null;index:idx_sessions_expires" json:"expires_at"`
	CreatedAt  time.Time      `gorm:"default:now()" json:"created_at"`
}
