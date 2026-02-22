package order

import (
	"time"

	"github.com/google/uuid"
	orderItem "github.com/scrumno/scrumno-api/internal/orders/entity/order-item"
	orderStatus "github.com/scrumno/scrumno-api/internal/orders/entity/order-status"
	orderStatusLog "github.com/scrumno/scrumno-api/internal/orders/entity/order-status-log"
	orderType "github.com/scrumno/scrumno-api/internal/orders/entity/order-type"
	"github.com/scrumno/scrumno-api/internal/organization/entity/venue"
	"github.com/scrumno/scrumno-api/internal/payments/entity/payment"
	"github.com/scrumno/scrumno-api/internal/payments/entity/promocode"
	posSyncStatus "github.com/scrumno/scrumno-api/internal/pos/entity/pos-sync-status"
	"github.com/scrumno/scrumno-api/internal/users/entity/user"
)

// Order — заказ клиента.
// Статус менять только через сервисные методы — иначе пуши и логи не сработают.
// Заказы НИКОГДА не удалять физически — финансовая история.
type Order struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"                   json:"id"`
	VenueID uuid.UUID `gorm:"type:uuid;not null;index:idx_orders_venue_status_date,priority:1" json:"venue_id"`
	UserID  uuid.UUID `gorm:"type:uuid;not null;index:idx_orders_user_date,priority:1"         json:"user_id"`

	Status    orderStatus.OrderStatus `gorm:"type:varchar(20);default:'NEW';index:idx_orders_venue_status_date,priority:2" json:"status"`
	OrderType orderType.OrderType     `gorm:"type:varchar(20);default:'takeaway'"                                           json:"order_type"`
	// Итог после скидки. Фиксируется при создании заказа.
	TotalAmount float64 `gorm:"type:numeric(12,2);not null" json:"total_amount"`
	// NULL = "как можно скорее". Шаг 15 мин, текущий или следующий день.
	ScheduledAt *time.Time `json:"scheduled_at,omitempty"`

	Comment       *string `gorm:"type:text"     json:"comment,omitempty"`
	CutleryNeeded bool    `gorm:"default:false" json:"cutlery_needed"`

	CancelledReason *string    `gorm:"type:text" json:"cancelled_reason,omitempty"`
	RefundID        *string    `gorm:"type:text" json:"refund_id,omitempty"`
	RefundedAt      *time.Time `                 json:"refunded_at,omitempty"`

	// payment_id — денормализованный кэш для быстрого JOIN.
	// Источник правды — payments.order_id.
	// Порядок: создать order → создать payment → обновить payment_id транзакционно.
	PaymentID *uuid.UUID `gorm:"type:uuid;index:idx_orders_payment" json:"payment_id,omitempty"`

	PosOrderID    *uuid.UUID                  `gorm:"type:uuid"                                                json:"pos_order_id,omitempty"`
	PosSyncStatus posSyncStatus.PosSyncStatus `gorm:"type:varchar(20);default:'NOT_SENT';index:idx_orders_sync" json:"pos_sync_status"`

	// Передаётся в POS как orderSource для маркетинговой аналитики.
	OrderSource string `gorm:"default:'PWA'" json:"order_source"`

	PromocodeID *uuid.UUID `gorm:"type:uuid" json:"promocode_id,omitempty"`
	// Хранится для аналитики. total_amount уже включает скидку.
	DiscountAmount float64 `gorm:"type:numeric(12,2);default:0" json:"discount_amount"`

	CreatedAt time.Time  `gorm:"autoCreateTime;index:idx_orders_venue_status_date,priority:3;index:idx_orders_user_date,priority:2" json:"created_at"`
	DeletedAt *time.Time `gorm:"index"                                                                                              json:"deleted_at,omitempty"`

	Venue venue.Venue `gorm:"foreignKey:VenueID"               json:"venue,omitempty"`
	User  user.User   `gorm:"foreignKey:UserID"                json:"user,omitempty"`
	// Payment ищем через payments.order_id (не через orders.payment_id).
	// orders.payment_id — денормализованный кэш, не FK в смысле GORM association.
	Payment    *payment.Payment                `gorm:"foreignKey:OrderID;references:ID" json:"payment,omitempty"`
	Promocode  *promocode.Promocode            `gorm:"foreignKey:PromocodeID"           json:"promocode,omitempty"`
	Items      []orderItem.OrderItem           `gorm:"foreignKey:OrderID"               json:"items,omitempty"`
	StatusLogs []orderStatusLog.OrderStatusLog `gorm:"foreignKey:OrderID"               json:"status_logs,omitempty"`
}
