package product_modifier

import (
	"github.com/google/uuid"
)

// ProductModifier — добавка / модификатор к товару.
type ProductModifier struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ProductID  uuid.UUID  `gorm:"type:uuid;not null;index"                        json:"product_id"`
	ExternalID *uuid.UUID `gorm:"type:uuid"                                       json:"external_id,omitempty"`
	ExtraPrice float64    `gorm:"type:numeric(12,2);default:0"                    json:"extra_price"`
	// true = клиент не может добавить товар в корзину без выбора этого модификатора.
	IsRequired bool `gorm:"default:false" json:"is_required"`
	// 1 = чекбокс. >1 = можно выбрать несколько (напр. до 3 сиропов).
	MaxQuantity int `gorm:"default:1" json:"max_quantity"`
	// Одинаковый group_name = одна визуальная группа в UI. NULL = одиночный модификатор.
	GroupName *string `json:"group_name,omitempty"`
}
