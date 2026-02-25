package organization_billing

import (
	"time"

	"github.com/google/uuid"
	billingStatus "github.com/scrumno/scrumno-api/internal/organization/entity/billing-status"
	subscriptionPlan "github.com/scrumno/scrumno-api/internal/organization/entity/subscription-plan"
)

type OrganizationBilling struct {
	ID             uuid.UUID                   `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	OrganizationID uuid.UUID                   `gorm:"type:uuid;not null;index"                        json:"organization_id"`
	PlanID         uuid.UUID                   `gorm:"type:uuid;not null"                              json:"plan_id"`
	BillingStatus  billingStatus.BillingStatus `gorm:"type:varchar(20);default:'ACTIVE'"               json:"billing_status"`
	// При оплате: paid_until = MAX(paid_until, now()) + 1 month.
	// Не затирать будущую дату при досрочной оплате.
	PaidUntil          *time.Time `json:"paid_until,omitempty"`
	PaymentMethodToken *string    `gorm:"type:text" json:"payment_method_token,omitempty"`
	LastInvoiceAt      *time.Time `json:"last_invoice_at,omitempty"`
	CreatedAt          time.Time  `gorm:"autoCreateTime" json:"created_at"`

	Plan subscriptionPlan.SubscriptionPlan `gorm:"foreignKey:PlanID"         json:"plan,omitempty"`
}
