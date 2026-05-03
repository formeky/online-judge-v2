package database

import (
	"gorm.io/gorm"

	"online-judge/internal/model"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Problem{},
		&model.TestCase{},
		&model.Submission{},
	)
}
