package interfaces

import (
	"engine/internal/pkg/domains/models/dtos"
	"engine/internal/pkg/domains/models/entities"
)

type TeamUsecase interface {
	CreateTeam(req dtos.CreateTeamRequest, user entities.User) (entities.Team, error)
}

type TeamRepository interface {
	TakeTeamByConditions(conditions map[string]interface{}) (entities.Team, error)
	CreateTeam(team entities.Team) (entities.Team, error)
}
