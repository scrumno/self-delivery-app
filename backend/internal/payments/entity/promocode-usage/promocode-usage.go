package promocode_usage

import (
	"time"

	"github.com/google/uuid"
	"github.com/scrumno/scrumno-api/internal/orders/entity/order"
	"github.com/scrumno/scrumno-api/internal/users/entity/user"
)

// PromocodeUsage — лог использования промокода.
// Уникальный индекс (promocode_id, user_id) ограничивает промокод одним
// использованием на пользователя. Если промокод многоразовый — убрать этот индекс.
// Никогда не удалять записи — финансовая история.
type PromocodeUsage struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"                  json:"id"`
	PromocodeID uuid.UUID `gorm:"type:uuid;not null;index:idx_promo_usage_code"                    json:"promocode_id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_promo_usage_user_code,priority:1" json:"user_id"`
	// Один заказ может использовать только один промокод (и наоборот).
	OrderID   uuid.UUID   `gorm:"type:uuid;not null;uniqueIndex;uniqueIndex:idx_promo_usage_user_code,priority:2" json:"order_id"`
	CreatedAt time.Time   `gorm:"autoCreateTime" json:"created_at"`
	User      user.User   `gorm:"foreignKey:UserID"      json:"user,omitempty"`
	Order     order.Order `gorm:"foreignKey:OrderID"     json:"order,omitempty"`
}

func (PromocodeUsage) TableName() string { return "promocode_usages" }
