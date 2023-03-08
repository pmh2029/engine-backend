package migrations

import (
	"gorm.io/gorm"

	"engine/internal/pkg/domains/models/entities"
)

func Migrate(dbConn *gorm.DB) error {
	err := dbConn.SetupJoinTable(entities.User{}, "Teams", entities.TeamUser{})
	if err != nil {
		return err
	}

	err = dbConn.AutoMigrate(
		entities.User{},
		entities.Team{},
	)

	return err
}
