package interfaces

import "engine/internal/pkg/domains/models/entities"

type TeamUserRepository interface {
	CreateTeamUser(teamUser entities.TeamUser) (entities.TeamUser, error)
	AddUserToTeam(team entities.Team, user entities.User) (entities.TeamUser, error)
}
