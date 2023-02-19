package interfaces

import "engine/internal/pkg/domains/models/entities"

type UserRepository interface {
	CreateUser(user entities.User) (entities.User, error)
	FindByConditions(conditions map[string]interface{}) ([]entities.User, error)
	TakeByConditions(conditions map[string]interface{}) (entities.User, error)
	UpdateUser(user entities.User, data map[string]interface{}) error
}