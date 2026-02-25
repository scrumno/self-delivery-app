package order_status_log

import (
	"time"

	"github.com/google/uuid"
	orderStatus "github.com/scrumno/scrumno-api/internal/orders/entity/order-status"
	"github.com/scrumno/scrumno-api/internal/users/entity/user"
)

// OrderStatusLog — аудит-лог смены статусов заказа.
// Позволяет восстановить хронологию и отлаживать зависания.
// ChangedBy = NULL — автоматическое изменение (воркер, webhook, таймаут).
type OrderStatusLog struct {
	ID        uuid.UUID                `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	OrderID   uuid.UUID                `gorm:"type:uuid;not null;index:idx_order_status_log"   json:"order_id"`
	OldStatus *orderStatus.OrderStatus `gorm:"type:varchar(20)"                                json:"old_status,omitempty"`
	NewStatus orderStatus.OrderStatus  `gorm:"type:varchar(20);not null"                       json:"new_status"`
	// NULL = автоматическое изменение (воркер, webhook). UUID сотрудника если ручное.
	ChangedBy *uuid.UUID `gorm:"type:uuid" json:"changed_by,omitempty"`
	Reason    *string    `gorm:"type:text" json:"reason,omitempty"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`

	Staff *user.User `gorm:"foreignKey:ChangedBy" json:"staff,omitempty"`
}

func (OrderStatusLog) TableName() string { return "order_status_logs" }
