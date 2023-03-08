package dtos

type CreateTeamRequest struct {
	TeamName   string `json:"team_name" binding:"required"`
	TeamAvatar string `json:"team_avatar"`
}

type TeamResponse struct {
	ID         uint   `json:"id"`
	TeamName   string `json:"team_name"`
	TeamAvatar string `json:"team_avatar"`
}
