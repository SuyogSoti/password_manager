package storage

import (
	"time"

	"gorm.io/gorm"
)

type Password struct {
	UserEmail      string `gorm:"primaryKey"`
	Site           string `gorm:"primaryKey"`
	SiteUserName   string `gorm:"primaryKey"`
	HashedPassword string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}
