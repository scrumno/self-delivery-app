package order_item_modifier

import (
	"github.com/google/uuid"
	productModifier "github.com/scrumno/scrumno-api/internal/products/entity/product-modifier"
)

// OrderItemModifier — модификатор в позиции заказа (снапшот на момент оплаты).
// Итог позиции: (unit_price + SUM(extra_price)) * quantity.
type OrderItemModifier struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	OrderItemID uuid.UUID `gorm:"type:uuid;not null;index:idx_oim_order_item"     json:"order_item_id"`
	// NULLABLE — модификатор может быть удалён из меню. Всегда LEFT JOIN.
	ModifierID   *uuid.UUID `gorm:"type:uuid"                             json:"modifier_id,omitempty"`
	NameSnapshot string     `gorm:"not null"                              json:"name_snapshot"`
	ExtraPrice   float64    `gorm:"type:numeric(12,2);not null;default:0" json:"extra_price"`

	Modifier *productModifier.ProductModifier `gorm:"foreignKey:ModifierID"  json:"modifier,omitempty"`
}

func (OrderItemModifier) TableName() string { return "order_item_modifiers" }
