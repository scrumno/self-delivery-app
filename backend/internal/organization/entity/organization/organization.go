package organization

import (
	"time"

	"github.com/google/uuid"
	organizationBilling "github.com/scrumno/scrumno-api/internal/organization/entity/organization-billing"
	subscriptionPlan "github.com/scrumno/scrumno-api/internal/organization/entity/subscription-plan"
	"github.com/scrumno/scrumno-api/internal/organization/entity/venue"
)

// Organization — юридическое лицо / владелец аккаунта.
// Не хранит операционных данных — только регистрационная информация.
type Organization struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name     string    `gorm:"not null"                                        json:"name"`
	INN      *string   `gorm:"type:varchar(12)"                                json:"inn,omitempty"`
	IsActive bool      `gorm:"default:true"                                    json:"is_active"`
	// PlanID — текущий тарифный план. При смене плана обновляется здесь
	// и создаётся новая запись в organization_billing (для истории).
	PlanID    uuid.UUID `gorm:"type:uuid;not null" json:"plan_id"`
	CreatedAt time.Time `gorm:"autoCreateTime"     json:"created_at"`

	Plan    subscriptionPlan.SubscriptionPlan         `gorm:"foreignKey:PlanID"         json:"plan,omitempty"`
	Venues  []venue.Venue                             `gorm:"foreignKey:OrganizationID" json:"venues,omitempty"`
	Billing []organizationBilling.OrganizationBilling `gorm:"foreignKey:OrganizationID" json:"billing,omitempty"`
}
