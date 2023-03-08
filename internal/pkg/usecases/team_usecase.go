package usecases

import (
	"engine/internal/pkg/domains/interfaces"
	"engine/internal/pkg/domains/models/dtos"
	"engine/internal/pkg/domains/models/entities"
)

type TeamUsecase struct {
	TeamRepo     interfaces.TeamRepository
	TeamUserRepo interfaces.TeamUserRepository
}

func NewTeamUsecase(
	tr interfaces.TeamRepository,
	tur interfaces.TeamUserRepository,
) interfaces.TeamUsecase {
	return &TeamUsecase{
		TeamRepo:     tr,
		TeamUserRepo: tur,
	}
}

func (tu *TeamUsecase) CreateTeam(req dtos.CreateTeamRequest, user entities.User) (entities.Team, error) {
	team, err := tu.TeamRepo.CreateTeam(entities.Team{
		TeamName:   req.TeamName,
		TeamAvatar: req.TeamAvatar,
	})
	if err != nil {
		return entities.Team{}, err
	}

	_, err = tu.TeamUserRepo.CreateTeamUser(entities.TeamUser{
		TeamID: team.ID,
		UserID: user.ID,
	})
	if err != nil {
		return team, err
	}

	return team, err
}
