package usecases

import (
	"engine/internal/pkg/domains/interfaces"
	"engine/internal/pkg/domains/models/entities"
)

type UserUsecase struct {
	UserRepo interfaces.UserRepository
}

func NewUserUsecase(ur interfaces.UserRepository) interfaces.UserUsecase {
	return &UserUsecase{
		UserRepo: ur,
	}
}

func (uu *UserUsecase) TakeUserByConditions(conditions map[string]interface{}) (entities.User, error) {
	user, err := uu.UserRepo.TakeUserByConditions(conditions)

	return user, err
}

func (uu *UserUsecase) FindUserByConditions(conditions map[string]interface{}) ([]entities.User, error) {
	users, err := uu.UserRepo.FindUserByConditions(conditions)

	return users, err
}
