package promocode

import (
	"time"

	"github.com/google/uuid"
	"github.com/scrumno/scrumno-api/internal/organization/entity/venue"
	discountType "github.com/scrumno/scrumno-api/internal/payments/entity/discount-type"
)

type Promocode struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"                                                          json:"id"`
	VenueID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_promocodes_venue_code,priority:1;index:idx_promocodes_active,priority:1" json:"venue_id"`
	Code    string    `gorm:"not null;uniqueIndex:idx_promocodes_venue_code,priority:2"                                                json:"code"`

	DiscountType  discountType.DiscountType `gorm:"type:varchar(10);not null"    json:"discount_type"`
	DiscountValue float64                   `gorm:"type:numeric(12,2);not null"  json:"discount_value"`
	// 0 = без ограничений по сумме заказа.
	MinOrderAmount float64 `gorm:"type:numeric(12,2);default:0" json:"min_order_amount"`
	// NULL = безлимит. INCREMENT в транзакции с SELECT FOR UPDATE.
	MaxUses   *int `json:"max_uses,omitempty"`
	UsedCount int  `gorm:"default:0;not null" json:"used_count"`

	// NULL = бессрочно.
	ExpiresAt *time.Time `gorm:"index:idx_promocodes_active,priority:3"               json:"expires_at,omitempty"`
	IsActive  bool       `gorm:"default:true;index:idx_promocodes_active,priority:2" json:"is_active"`
	CreatedAt time.Time  `gorm:"autoCreateTime"                                      json:"created_at"`

	Venue venue.Venue `gorm:"foreignKey:VenueID"     json:"venue,omitempty"`
}
