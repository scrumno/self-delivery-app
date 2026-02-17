package order_item

import (
	"github.com/google/uuid"
	orderItemModifier "github.com/scrumno/scrumno-api/internal/orders/entity/order-item-modifier"
)

type OrderItem struct {
	ID                  uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	OrderID             uuid.UUID `gorm:"type:uuid;not null;index:idx_order_items_order" json:"order_id"`
	ProductID           uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	Quantity            int       `gorm:"not null" json:"quantity"`
	UnitPrice           float64   `gorm:"type:numeric(12,2);not null" json:"unit_price"`
	ProductNameSnapshot string    `gorm:"not null" json:"product_name_snapshot"`

	Modifiers []orderItemModifier.OrderItemModifier `gorm:"foreignKey:OrderItemID" json:"modifiers,omitempty"`
}
