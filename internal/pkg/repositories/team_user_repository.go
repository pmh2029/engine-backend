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

func (tur *TeamUserRepository) AddUserToTeam(team entities.Team, user entities.User) (entities.TeamUser, error) {
	teamUser := entities.TeamUser{
		TeamID: team.ID,
		UserID: user.ID,
	}

	result := tur.DBConn.Create(&teamUser)

	return teamUser, result.Error
}
