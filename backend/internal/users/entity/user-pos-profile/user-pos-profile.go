package user_pos_profile

import (
	"time"

	"github.com/google/uuid"
	posProvider "github.com/scrumno/scrumno-api/internal/pos/entity/pos-provider"
)

type UserPOSProfile struct {
	ID              uuid.UUID               `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"             json:"id"`
	UserID          uuid.UUID               `gorm:"type:uuid;not null;uniqueIndex:idx_user_pos,priority:1"      json:"user_id"`
	VenueID         uuid.UUID               `gorm:"type:uuid;not null;uniqueIndex:idx_user_pos,priority:2"      json:"venue_id"`
	Provider        posProvider.POSProvider `gorm:"type:varchar(20);not null;uniqueIndex:idx_user_pos,priority:3" json:"provider"`
	ExternalGuestID string                  `gorm:"not null"                                                    json:"external_guest_id"`
	CreatedAt       time.Time               `gorm:"autoCreateTime"                                              json:"created_at"`
}
