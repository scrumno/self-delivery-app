package create_user

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserCommand struct {
	Phone string
}

type UserDTO struct {
	ID        uuid.UUID  `json:"id"`
	Phone     string     `json:"phone"`
	FullName  *string    `json:"full_name,omitempty"`
	BirthDate *time.Time `json:"birth_date"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
}
