package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `grom:"primaryKey" json:"id"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Name      string         `json:"name"`
	Username  string         `gorm:"uniqueIndex;not null" json:"username"`
	Password  *string        `json:"password"`
	Status    bool           `gorm:"default:true" json:"status"`
	CreatedAt time.Time      `gorm:"column:created_at;" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;" json:"deleted_at"`
}
