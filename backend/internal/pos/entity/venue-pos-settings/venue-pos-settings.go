package venue_pos_settings

import (
	"time"

	"github.com/google/uuid"
	posProvider "github.com/scrumno/scrumno-api/internal/pos/entity/pos-provider"
)

// VenuePOSSettings — настройки подключения POS-системы для точки (1:1).
// Абстрагирован от конкретной системы через поле provider.
//
// credentials_encrypted содержит провайдер-специфичный JSON (шифровать через pgcrypto/Vault):
//
//	iiko:    {"api_key":"...", "terminal_id":"uuid", "org_id":"uuid"}
//	syrve:   {"api_login":"...", "org_id":"uuid"}
//	rkeeper: {"ws_url":"...", "token":"..."}
//	manual:  {} — пустой объект, без интеграции
//
// Никогда не логировать расшифрованный payload.
type VenuePOSSettings struct {
	VenueID  uuid.UUID               `gorm:"type:uuid;primaryKey"      json:"venue_id"`
	Provider posProvider.POSProvider `gorm:"type:varchar(20);not null" json:"provider"`
	// Зашифрованный JSON с credentials POS-системы.
	CredentialsEncrypted string `gorm:"type:text;not null" json:"credentials_encrypted"`
	// Интервал полной синхронизации меню (минуты).
	SyncIntervalMinutes int `gorm:"default:1"  json:"sync_interval_minutes"`
	// Интервал опроса стоп-листа (секунды).
	StoplistPollSeconds int `gorm:"default:60" json:"stoplist_poll_seconds"`

	// Если LastMenuSyncAt > 5 мин назад — алерт в Sentry.
	// NULL = синхронизации ещё не было.
	LastMenuSyncAt     *time.Time `json:"last_menu_sync_at,omitempty"`
	LastStoplistSyncAt *time.Time `json:"last_stoplist_sync_at,omitempty"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty"`
}

func (VenuePOSSettings) TableName() string { return "venue_pos_settings" }
