package repositories

import (
	"engine/internal/pkg/domains/interfaces"
	"engine/internal/pkg/domains/models/entities"

	"gorm.io/gorm"
)

type TeamRepository struct {
	DBConn *gorm.DB
}

func NewTeamRepository(dbConn *gorm.DB) interfaces.TeamRepository {
	return &TeamRepository{
		DBConn: dbConn,
	}
}

func (tr *TeamRepository) CreateTeam(team entities.Team) (entities.Team, error) {
	result := tr.DBConn.Create(&team)

	return team, result.Error
}

func (tr *TeamRepository) TakeTeamByConditions(conditions map[string]interface{}) (entities.Team, error) {
	team := entities.Team{}
	result := tr.DBConn.Where(conditions).Take(&team)

	return team, result.Error
}

func (tr *TeamRepository) GetTeamMemberList(team entities.Team) ([]entities.User, error) {
	result := tr.DBConn.Preload("Users").Where("id = ?", team.ID).Find(&team)

	return team.Users, result.Error
}
