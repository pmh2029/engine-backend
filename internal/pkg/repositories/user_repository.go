package repositories

import (
	"gorm.io/gorm"

	"engine/internal/pkg/domains/interfaces"
	"engine/internal/pkg/domains/models/entities"
)

type UserRepository struct {
	DBConn *gorm.DB
}

func NewUserRepository(dbConn *gorm.DB) interfaces.UserRepository {
	return &UserRepository{
		DBConn: dbConn,
	}
}

func (ur *UserRepository) CreateUser(user entities.User) (entities.User, error) {
	result := ur.DBConn.Create(&user)

	return user, result.Error
}

func (ur *UserRepository) FindByConditions(conditions map[string]interface{}) ([]entities.User, error) {
	users := []entities.User{}

	result := ur.DBConn.Where(conditions).Find(&users)
	return users, result.Error
}

func (ur *UserRepository) TakeByConditions(conditions map[string]interface{}) (entities.User, error) {
	user := entities.User{}
	result := ur.DBConn.Where(conditions).Take(&user)

	return user, result.Error
}

func (ur *UserRepository) UpdateUser(user entities.User, data map[string]interface{}) error {
	result := ur.DBConn.Model(&user).Where("id = ?", user.ID).Updates(data)

	return result.Error
}
