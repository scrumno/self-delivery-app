package product

import (
	"time"

	"github.com/google/uuid"
	"github.com/scrumno/scrumno-api/internal/categories/entity/category"
	"github.com/scrumno/scrumno-api/internal/organization/entity/venue"
	productModifier "github.com/scrumno/scrumno-api/internal/products/entity/product-modifier"
)

type Product struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"                        json:"id"`
	VenueID uuid.UUID `gorm:"type:uuid;not null;index:idx_products_pwa_main,priority:1"             json:"venue_id"`
	// ON DELETE RESTRICT — нельзя удалить категорию пока в ней есть товары.
	CategoryID uuid.UUID `gorm:"type:uuid;index:idx_products_pwa_main,priority:2"                    json:"category_id"`
	// UUID из POS-системы. Матчинг строго по UUID при синхронизации.
	// Нет в нашей БД → INSERT, есть → UPDATE, нет в POS → soft delete.
	ExternalID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_products_external_unique,priority:2" json:"external_id"`

	// Текущая цена из POS. Обновляется при синхронизации.
	// На цену в заказе не влияет — там снапшот в order_items.unit_price.
	Price       float64 `gorm:"type:numeric(12,2);not null"                              json:"price"`
	ImageURL    *string `gorm:"type:text"                                                json:"image_url,omitempty"`
	IsAvailable bool    `gorm:"default:true;index:idx_products_pwa_main,priority:3"      json:"is_available"`
	SortOrder   int     `gorm:"default:0;index:idx_products_pwa_main,priority:4"         json:"sort_order"`

	// Физические параметры и КБЖУ из POS.
	// Блок КБЖУ показываем только если calories IS NOT NULL.
	Weight   *int     `gorm:"type:integer"      json:"weight,omitempty"`
	Calories *float64 `gorm:"type:numeric(8,2)" json:"calories,omitempty"`
	Protein  *float64 `gorm:"type:numeric(8,2)" json:"protein,omitempty"`
	Fat      *float64 `gorm:"type:numeric(8,2)" json:"fat,omitempty"`
	Carbs    *float64 `gorm:"type:numeric(8,2)" json:"carbs,omitempty"`

	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	Venue     venue.Venue                       `gorm:"foreignKey:VenueID"    json:"venue,omitempty"`
	Category  category.Category                 `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Modifiers []productModifier.ProductModifier `gorm:"foreignKey:ProductID"  json:"modifiers,omitempty"`
}
