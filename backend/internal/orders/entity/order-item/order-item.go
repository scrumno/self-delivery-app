package order_item

import (
	"github.com/google/uuid"
	orderItemModifier "github.com/scrumno/scrumno-api/internal/orders/entity/order-item-modifier"
	"github.com/scrumno/scrumno-api/internal/products/entity/product"
)

// OrderItem — позиция в заказе (снапшот цены на момент оплаты).
type OrderItem struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null;index:idx_order_items_order"  json:"order_id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null"                              json:"product_id"`
	Quantity  int       `gorm:"not null"                                        json:"quantity"`
	// Снапшот цены на момент оплаты. Не меняется при изменении products.price.
	UnitPrice float64 `gorm:"type:numeric(12,2);not null" json:"unit_price"`
	// Снапшот названия — история не пустеет при удалении товара из меню.
	ProductNameSnapshot string `gorm:"not null" json:"product_name_snapshot"`

	Product   product.Product                       `gorm:"foreignKey:ProductID"   json:"product,omitempty"`
	Modifiers []orderItemModifier.OrderItemModifier `gorm:"foreignKey:OrderItemID" json:"modifiers,omitempty"`
}
