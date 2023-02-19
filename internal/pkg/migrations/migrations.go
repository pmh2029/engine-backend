package migrations

import (
	"gorm.io/gorm"

	"engine/internal/pkg/domains/models/entities"
)

func Migrate(dbConn *gorm.DB) error {
	err := dbConn.AutoMigrate(entities.User{})

	return err
}
