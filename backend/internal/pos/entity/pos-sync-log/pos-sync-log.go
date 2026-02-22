package pos_sync_log

import (
	"time"

	"github.com/google/uuid"
	"github.com/scrumno/scrumno-api/internal/orders/entity/order"
	"github.com/scrumno/scrumno-api/internal/organization/entity/venue"
	posProvider "github.com/scrumno/scrumno-api/internal/pos/entity/pos-provider"
	"gorm.io/datatypes"
)

// PosSyncLog — лог одной попытки синхронизации заказа с POS-системой.
// TTL = 14 дней: pg_cron DELETE WHERE created_at < NOW() - INTERVAL '14 days'.
type PosSyncLog struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"                          json:"id"`
	VenueID uuid.UUID `gorm:"type:uuid;not null;index:idx_pos_logs_venue_date,priority:1"              json:"venue_id"`
	OrderID uuid.UUID `gorm:"type:uuid;not null;index:idx_pos_logs_order_date,priority:1"              json:"order_id"`
	// Провайдер POS-системы в момент синхронизации.
	// Позволяет различать логи iiko и rkeeper в одной таблице.
	Provider posProvider.POSProvider `gorm:"type:varchar(20);not null" json:"provider"`
	// Логировать ДО отправки. НЕ включать credentials в payload.
	RequestPayload datatypes.JSON `gorm:"type:jsonb" json:"request_payload,omitempty"`
	// NULL если нет ответа (таймаут). При ERROR — главный источник для диагностики.
	ResponsePayload datatypes.JSON `gorm:"type:jsonb" json:"response_payload,omitempty"`
	// 200 = успех → SENT. 500/503 = retry до 3 раз. 401 = невалидный ключ. NULL = таймаут.
	HTTPStatusCode *int      `json:"http_status_code,omitempty"`
	CreatedAt      time.Time `gorm:"autoCreateTime;index:idx_pos_logs_venue_date,priority:2;index:idx_pos_logs_order_date,priority:2" json:"created_at"`

	Venue venue.Venue `gorm:"foreignKey:VenueID" json:"venue,omitempty"`
	Order order.Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

func (PosSyncLog) TableName() string { return "pos_sync_logs" }
