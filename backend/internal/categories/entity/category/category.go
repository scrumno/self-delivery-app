package category

import (
	"time"

	"github.com/google/uuid"
	"github.com/scrumno/scrumno-api/internal/organization/entity/venue"
)

// Category — категория меню точки.
type Category struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"                  json:"id"`
	VenueID uuid.UUID `gorm:"type:uuid;not null;index:idx_categories_venue_visible,priority:1" json:"venue_id"`
	// UUID категории в POS-системе. NULL = создана вручную в нашей админке.
	ExternalID *uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_categories_external_unique,priority:2"  json:"external_id,omitempty"`
	// Для вложенных категорий. NULL = верхний уровень.
	ParentID  *uuid.UUID `gorm:"type:uuid"                                                        json:"parent_id,omitempty"`
	SortOrder int        `gorm:"default:0"                                                        json:"sort_order"`
	IsVisible bool       `gorm:"default:true;index:idx_categories_venue_visible,priority:2"       json:"is_visible"`
	DeletedAt *time.Time `gorm:"index"                                                            json:"deleted_at,omitempty"`

	Venue  venue.Venue `gorm:"foreignKey:VenueID"    json:"venue,omitempty"`
	Parent *Category   `gorm:"foreignKey:ParentID"   json:"parent,omitempty"`
}
