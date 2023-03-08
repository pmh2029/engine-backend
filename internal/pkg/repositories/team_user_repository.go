package repositories

import (
	"engine/internal/pkg/domains/interfaces"
	"engine/internal/pkg/domains/models/entities"

	"gorm.io/gorm"
)

type TeamUserRepository struct {
	DBConn *gorm.DB
}

func NewTeamUserRepository(dbConn *gorm.DB) interfaces.TeamUserRepository {
	return &TeamUserRepository{
		DBConn: dbConn,
	}
}

func (tur *TeamUserRepository) CreateTeamUser(teamUser entities.TeamUser) (entities.TeamUser, error) {
	result := tur.DBConn.Create(&teamUser)

	return teamUser, result.Error
}
