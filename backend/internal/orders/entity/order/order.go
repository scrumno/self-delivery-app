package order

import (
	"time"

	"github.com/google/uuid"
	orderItem "github.com/scrumno/scrumno-api/internal/orders/entity/order-item"
	"github.com/scrumno/scrumno-api/internal/orders/entity/payment"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	StatusNew            OrderStatus = "NEW"
	StatusWaitingPayment OrderStatus = "WAITING_PAYMENT"
	StatusCooking        OrderStatus = "COOKING"
	StatusReady          OrderStatus = "READY"
	StatusCompleted      OrderStatus = "COMPLETED"
	StatusCancelled      OrderStatus = "CANCELLED"
)

type OrderType string

const (
	OrderTypeTakeaway OrderType = "takeaway"
	OrderTypeDineIn   OrderType = "dine_in"
	OrderTypeDelivery OrderType = "delivery"
)

type PosSyncStatus string

const (
	SyncNotSent PosSyncStatus = "NOT_SENT"
	SyncSent    PosSyncStatus = "SENT"
	SyncError   PosSyncStatus = "ERROR"
)

type Order struct {
	ID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	VenueID uuid.UUID `gorm:"type:uuid;not null;index:idx_orders_venue_status_date" json:"venue_id"`
	UserID  uuid.UUID `gorm:"type:uuid;not null;index:idx_orders_user_date" json:"user_id"`

	Status      OrderStatus `gorm:"type:varchar;default:'NEW';index:idx_orders_venue_status_date" json:"status"`
	OrderType   OrderType   `gorm:"type:varchar;default:'takeaway'" json:"order_type"`
	TotalAmount float64     `gorm:"type:numeric(12,2);not null" json:"total_amount"`
	ScheduledAt *time.Time  `json:"scheduled_at,omitempty"`

	Comment       *string `gorm:"type:text" json:"comment,omitempty"`
	CutleryNeeded bool    `gorm:"default:false" json:"cutlery_needed"`

	// Логика отмены
	CancelledReason *string    `gorm:"type:text" json:"cancelled_reason,omitempty"`
	RefundID        *string    `json:"refund_id,omitempty"`
	RefundedAt      *time.Time `json:"refunded_at,omitempty"`

	// Платеж (может быть null до оплаты)
	PaymentID *uuid.UUID       `gorm:"type:uuid;index:idx_orders_payment" json:"payment_id,omitempty"`
	Payment   *payment.Payment `gorm:"foreignKey:PaymentID" json:"payment,omitempty"`

	// iiko
	PosOrderID    *uuid.UUID    `gorm:"type:uuid" json:"pos_order_id,omitempty"`
	PosSyncStatus PosSyncStatus `gorm:"type:varchar;default:'NOT_SENT';index:idx_orders_sync" json:"pos_sync_status"`
	OrderSource   string        `gorm:"default:'PWA'" json:"order_source"`

	// Промокод
	PromocodeID    *uuid.UUID `gorm:"type:uuid" json:"promocode_id,omitempty"`
	DiscountAmount float64    `gorm:"type:numeric(12,2);default:0" json:"discount_amount"`

	CreatedAt time.Time      `gorm:"default:now();index:idx_orders_venue_status_date;index:idx_orders_user_date" json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-"`

	Items []orderItem.OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}
