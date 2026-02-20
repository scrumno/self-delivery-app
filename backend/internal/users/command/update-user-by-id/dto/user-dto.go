package dto

import (
	"time"

	"github.com/google/uuid"
	staffrole "github.com/scrumno/scrumno-api/internal/users/entity/staff-role"
)

type UserDTO struct {
	ID          uuid.UUID  `json:"id"`
	Phone       string     `json:"phone"`
	FullName    *string    `json:"full_name"`
	BirthDate   *time.Time `json:"birth_date"`
	IikoGuestID *uuid.UUID `json:"iiko_guest_id"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`

	// Связи
	StaffRoles []staffrole.StaffRole `json:"staff_roles,omitempty"`
}
