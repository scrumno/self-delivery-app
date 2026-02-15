package entity

import "time"

// UserToken — модель для хранения refresh-токенов (таблица создаётся миграцией).
type UserToken struct {
	UserID       string    `gorm:"primaryKey;not null" json:"user_id"`
	RefreshToken string    `gorm:"uniqueIndex;not null" json:"-"`
	ExpiresAt    time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

// TableName задаёт имя таблицы для GORM.
func (UserToken) TableName() string {
	return "user_tokens"
}
