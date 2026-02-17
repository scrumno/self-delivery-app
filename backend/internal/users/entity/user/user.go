package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/scrumno/scrumno-api/internal/users/entity/session"
	staffrole "github.com/scrumno/scrumno-api/internal/users/entity/staff-role"
)

type User struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Phone       string     `gorm:"not null;uniqueIndex:idx_users_phone" json:"phone"`
	FullName    *string    `json:"full_name,omitempty"`
	BirthDate   *time.Time `gorm:"type:date" json:"birth_date,omitempty"` // Только дата
	IikoGuestID *uuid.UUID `gorm:"type:uuid" json:"iiko_guest_id,omitempty"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time  `gorm:"default:now()" json:"created_at"`

	// Связи
	StaffRoles []staffrole.StaffRole `gorm:"foreignKey:UserID" json:"staff_roles,omitempty"`
	Sessions   []session.Session     `gorm:"foreignKey:UserID" json:"-"`
}
