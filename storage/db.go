package storage

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDB(config postgres.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(config))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}
	// Migrate the schema
	// TODO(suyogsoti): get rid of this in the future
	if err := db.AutoMigrate(&User{}); err != nil {
		return nil, fmt.Errorf("failed to migrate db model user: %w", err)
	}
	// Migrate the schema
	// TODO(suyogsoti): get rid of this in the future
	if err := db.AutoMigrate(&Password{}); err != nil {
		return nil, fmt.Errorf("failed to migrate db model user: %w", err)
	}

	return db, nil
}
