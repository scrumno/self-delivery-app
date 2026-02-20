package dto

import (
	"time"

	"github.com/google/uuid"
)

type UpdateData struct {
	Phone       string
	FullName    *string
	BirthDate   *time.Time
	IikoGuestID *uuid.UUID
	IsActive    bool
	CreatedAt   time.Time
}
