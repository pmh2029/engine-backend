package interfaces

import (
	"engine/internal/pkg/domains/models/dtos"
	"engine/internal/pkg/domains/models/entities"
)

type AuthUsecase interface {
	FindByConditions(conditions map[string]interface{}) ([]entities.User, error)
	TakeByConditions(conditions map[string]interface{}) (entities.User, error)
	SignUp(req dtos.CreateUserRequest) (entities.User, error)
	SignIn(req dtos.SignInRequest) (entities.User, string, error)
	SendMailForgotPassword(req dtos.ForgotPasswordRequest) error
	ActiveUser(userID uint) error
	ResetPassword(userId uint, req dtos.ResetPasswordRequest) error
}
