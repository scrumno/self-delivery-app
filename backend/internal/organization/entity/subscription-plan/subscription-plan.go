package subscription_plan

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type SubscriptionPlan struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name          string    `gorm:"not null"                                        json:"name"`
	PricePerMonth float64   `gorm:"type:numeric(12,2);not null"                     json:"price_per_month"`
	// Лимиты и флаги плана.
	// Формат: {"max_products":100,"max_venues":1,"push_marketing":false}
	FeaturesJSON datatypes.JSON `gorm:"type:jsonb"     json:"features_json,omitempty"`
	IsActive     bool           `gorm:"default:true"   json:"is_active"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
}
