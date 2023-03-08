package interfaces

import "engine/internal/pkg/domains/models/entities"

type TeamUserRepository interface {
	CreateTeamUser(teamUser entities.TeamUser) (entities.TeamUser, error)
}
