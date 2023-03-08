package entities

// TeamUsersTableName TableName
var TeamUsersTableName = "team_users"

type TeamUser struct {
	BaseEntity
	TeamID uint `gorm:"primaryKey"`
	UserID uint `gorm:"primaryKey"`
}

// TableName func
func (i *TeamUser) TableName() string {
	return TeamUsersTableName
}
