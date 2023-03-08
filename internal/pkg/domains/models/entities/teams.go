package entities

// TeamsTableName TableName
var TeamsTableName = "teams"

type Team struct {
	BaseEntity
	TeamName   string `gorm:"column:team_name;not null"`
	TeamAvatar string `gorm:"column:team_avatar"`
	Users      []User `gorm:"many2many:team_users;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName func
func (i *Team) TableName() string {
	return TeamsTableName
}
