package venue

import (
	"time"

	"github.com/google/uuid"
	billingStatus "github.com/scrumno/scrumno-api/internal/organization/entity/billing-status"
	"gorm.io/datatypes"
)

// Venue — торговая точка. Центральная сущность всей системы.
// Почти каждая таблица ссылается на venue_id.
type Venue struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index"                        json:"organization_id"`
	Name           string    `gorm:"not null"                                        json:"name"`
	// Slug используется в URL: brand-kazan.foodapp.ru
	// Только [a-z0-9-]. Менять осторожно — старые QR-коды перестанут работать.
	Slug      string   `gorm:"uniqueIndex;not null" json:"slug"`
	Address   *string  `gorm:"type:text"            json:"address,omitempty"`
	Latitude  *float64 `gorm:"type:numeric(10,7)"   json:"latitude,omitempty"`
	Longitude *float64 `gorm:"type:numeric(10,7)"   json:"longitude,omitempty"`

	// is_active = false → точка закрыта навсегда (PWA → 404).
	// is_emergency_stop = true → временная пауза (PWA → 503, 1 клик в админке).
	IsActive          bool    `gorm:"default:true"                 json:"is_active"`
	IsEmergencyStop   bool    `gorm:"default:false"                json:"is_emergency_stop"`
	MinOrderAmount    float64 `gorm:"type:numeric(12,2);default:0" json:"min_order_amount"`
	AvgCookingMinutes int     `gorm:"default:20"                   json:"avg_cooking_minutes"`

	// Расписание работы (оверрайд над POS).
	// Формат: {"mon":{"open":"09:00","close":"22:00"}, ...}. Null = берём из POS.
	WorkHoursJSON datatypes.JSON `gorm:"type:jsonb" json:"work_hours_json,omitempty"`

	// billing_status — кэш из organization_billing для быстрой проверки на каждый запрос.
	// Источник правды — organization_billing. Синхронизируется cron-задачей каждый час
	// и при каждом изменении статуса в organization_billing.
	BillingStatus billingStatus.BillingStatus `gorm:"type:varchar(20);default:'ACTIVE'" json:"billing_status"`
	PaidUntil     *time.Time                  `                                         json:"paid_until,omitempty"`

	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt *time.Time `                       json:"updated_at,omitempty"`
	DeletedAt *time.Time `gorm:"index"          json:"deleted_at,omitempty"`
}
