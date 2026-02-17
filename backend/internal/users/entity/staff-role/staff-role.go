package staff_role

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleManager  UserRole = "manager"
	RoleCashier  UserRole = "cashier"
	RoleCustomer UserRole = "customer"
)

type StaffRole struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_staff_roles_user_venue" json:"user_id"`
	VenueID   uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_staff_roles_user_venue" json:"venue_id"`
	Role      UserRole  `gorm:"type:varchar;not null" json:"role"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"default:now()" json:"created_at"`
}
