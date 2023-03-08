package interfaces

import "engine/internal/pkg/domains/models/entities"

type UserUsecase interface {
	FindUserByConditions(conditions map[string]interface{}) ([]entities.User, error)
	TakeUserByConditions(conditions map[string]interface{}) (entities.User, error)
}

type UserRepository interface {
	CreateUser(user entities.User) (entities.User, error)
	FindUserByConditions(conditions map[string]interface{}) ([]entities.User, error)
	TakeUserByConditions(conditions map[string]interface{}) (entities.User, error)
	UpdateUser(user entities.User, data map[string]interface{}) error
}
