package update_user_by_id

import (
	"time"

	"github.com/google/uuid"
)

type Command struct {
	ID   uuid.UUID
	Data CommandData
}

type CommandData struct {
	Phone       *string    `json:"phone"`
	FullName    *string    `json:"full_name"`
	BirthDate   *time.Time `json:"birth_date"`
	IikoGuestID *uuid.UUID `json:"iiko_guest_id"`
	IsActive    *bool      `json:"is_active"`
	CreatedAt   *time.Time `json:"created_at"`
}
