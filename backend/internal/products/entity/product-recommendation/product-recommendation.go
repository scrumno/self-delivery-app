package product_recommendation

import (
	"github.com/google/uuid"
	"github.com/scrumno/scrumno-api/internal/products/entity/product"
)

type ProductRecommendation struct {
	ProductID     uuid.UUID `gorm:"type:uuid;primaryKey;index:idx_recs_sort,priority:1" json:"product_id"`
	RecommendedID uuid.UUID `gorm:"type:uuid;primaryKey"                                json:"recommended_id"`
	SortOrder     int       `gorm:"default:0;index:idx_recs_sort,priority:2"           json:"sort_order"`
	
	Recommended product.Product `gorm:"foreignKey:RecommendedID" json:"recommended,omitempty"`
}
