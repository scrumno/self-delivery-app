package user

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primary_key" json:"id"`
	Username  string         `gorm:"uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"not null" json:"-"`
	FullName  string         `json:"full_name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
