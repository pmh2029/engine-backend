package migrations

import (
	"fmt"

	"gorm.io/gorm"

	"engine/internal/pkg/domains/models/entities"
)

func Migrate(dbConn *gorm.DB) error {
	err := dbConn.SetupJoinTable(entities.User{}, "Teams", entities.TeamUser{})
	if err != nil {
		fmt.Println("err: ", err)
	}
	err = dbConn.SetupJoinTable(entities.Team{}, "Users", entities.TeamUser{})
	if err != nil {
		fmt.Println("err: ", err)
	}

	err = dbConn.AutoMigrate(
		entities.User{},
		entities.Team{},
		entities.TeamUser{},
	)
	if err != nil {
		fmt.Println("err: ", err)
	}
	return err
}
