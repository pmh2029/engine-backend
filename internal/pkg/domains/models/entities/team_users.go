package entities

// TeamUsersTableName TableName
var TeamUsersTableName = "team_users"

type TeamUser struct {
	BaseEntity
	TeamID uint 
	UserID uint 
}

// TableName func
func (i *TeamUser) TableName() string {
	return TeamUsersTableName
}
