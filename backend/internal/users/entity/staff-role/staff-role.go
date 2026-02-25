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
	ID      uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"           json:"id"`
	UserID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_staff_roles_user_venue" json:"user_id"`
	VenueID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_staff_roles_user_venue" json:"venue_id"`
	Role    UserRole  `gorm:"type:varchar(20);not null"                                 json:"role"`
	// false = сотрудник уволен. НЕ удалять запись.
	IsActive  bool      `gorm:"default:true"   json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
