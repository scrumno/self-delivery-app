package order_item_modifier

import "github.com/google/uuid"

type OrderItemModifier struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	OrderItemID  uuid.UUID  `gorm:"type:uuid;not null;index:idx_oim_order_item" json:"order_item_id"`
	ModifierID   *uuid.UUID `gorm:"type:uuid" json:"modifier_id,omitempty"`
	NameSnapshot string     `gorm:"not null" json:"name_snapshot"`
	ExtraPrice   float64    `gorm:"type:numeric(12,2);default:0;not null" json:"extra_price"`
}
