package create_user

import (
	"github.com/google/uuid"
	"time"
)

type Command struct {
	Phone    string  `json:"phone"`
	FullName *string `json:"full_name,omitempty"`
}

type UserDTO struct {
	ID        uuid.UUID `json:"id"`
	Phone     string    `json:"phone"`
	FullName  *string   `json:"full_name,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}
