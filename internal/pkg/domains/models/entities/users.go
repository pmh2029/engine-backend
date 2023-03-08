package entities

// UsersTableName TableName
var UsersTableName = "users"

type User struct {
	BaseEntity
	Username string  `gorm:"column:username;not null"`
	Email    string  `gorm:"column:email;not null;unique"`
	Password string  `gorm:"column:password;not null"`
	IsActive bool    `gorm:"column:is_active;default:false"`
	Teams    []*Team `gorm:"many2many:team_users;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName func
func (i *User) TableName() string {
	return UsersTableName
}
