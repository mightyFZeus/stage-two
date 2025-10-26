package store

import (
	"github.com/mightyzeus/stage-two/internal/models"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Country{},
	)
}
