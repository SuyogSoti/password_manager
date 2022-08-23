package storage

import (
	"time"
)

type User struct {
	Email          string `gorm:"primaryKey"`
	HashedPassword string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
