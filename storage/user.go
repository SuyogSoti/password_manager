package storage

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Email          string `gorm:"primaryKey"`
	HashedPassword string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}
