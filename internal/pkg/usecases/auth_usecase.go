package usecases

import (
	"encoding/base64"
	"errors"
	"os"
	"time"

	"engine/internal/pkg/domains/interfaces"
	"engine/internal/pkg/domains/models/dtos"
	"engine/internal/pkg/domains/models/entities"
	"engine/pkg/shared/auth"
	"engine/pkg/shared/utils"
)

type AuthUsecase struct {
	UserRepo interfaces.UserRepository
}

func NewAuthUsecase(ur interfaces.UserRepository) interfaces.AuthUsecase {
	return &AuthUsecase{
		UserRepo: ur,
	}
}

func (au *AuthUsecase) SignUp(req dtos.CreateUserRequest) (entities.User, error) {
	user := entities.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		IsActive: req.IsActive,
	}

	hashPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return entities.User{}, err
	}

	user.Password = hashPassword
	user, err = au.UserRepo.CreateUser(user)
	if err != nil {
		return entities.User{}, err
	}

	encodedEmail := base64.StdEncoding.EncodeToString([]byte(user.Email))

	token, err := auth.GenerateHS256JWT(map[string]interface{}{
		"email": user.Email,
		"exp":   time.Now().Add(time.Minute * 15).Unix(),
	})
	if err != nil {
		return entities.User{}, err
	}

	templateData := utils.TemplateData{
		Path:    "pkg/shared/template/verify_email_template.html",
		To:      req.Email,
		Subject: "Confirm your email address",
		Url:     os.Getenv("BASE_URL") + "auth/verify_email/" + encodedEmail + "/" + token,
	}

	err = utils.SendTemplateEMail(templateData)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (au *AuthUsecase) SignIn(req dtos.SignInRequest) (entities.User, string, error) {
	user, err := au.UserRepo.TakeUserByConditions(map[string]interface{}{
		"email": req.Email,
	})
	if err != nil {
		return entities.User{}, "", err
	}

	checkPassword := utils.CheckHashPassword(req.Password, user.Password)
	if !checkPassword {
		return entities.User{}, "", errors.New("password is incorrect")
	}

	jwtToken, err := auth.GenerateHS256JWT(map[string]interface{}{
		"email": user.Email,
		"sub":   user.ID,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})
	if err != nil {
		return entities.User{}, "", errors.New("error while generating token")
	}

	return user, jwtToken, nil
}

func (au *AuthUsecase) SendMailForgotPassword(req dtos.ForgotPasswordRequest) error {
	user, err := au.UserRepo.TakeUserByConditions(map[string]interface{}{
		"email": req.Email,
	})
	if err != nil {
		return err
	}

	if !user.IsActive {
		return errors.New("user is not active")
	}

	encodedEmail := base64.StdEncoding.EncodeToString([]byte(user.Email))

	token, err := auth.GenerateHS256JWT(map[string]interface{}{
		"email": user.Email,
		"exp":   time.Now().Add(time.Minute * 15).Unix(),
	})
	if err != nil {
		return err
	}

	templateData := utils.TemplateData{
		Path:    "pkg/shared/template/forgot_password_template.html",
		To:      req.Email,
		Subject: "Forgot Password",
		Url:     os.Getenv("BASE_URL") + "auth/reset_password/" + encodedEmail + "/" + token,
	}

	err = utils.SendTemplateEMail(templateData)

	return err
}

func (au *AuthUsecase) ActiveUser(userID uint) error {
	user, err := au.UserRepo.TakeUserByConditions(map[string]interface{}{
		"id": userID,
	})
	if err != nil {
		return err
	}

	if user.IsActive {
		return errors.New("user already active")
	}

	err = au.UserRepo.UpdateUser(user, map[string]interface{}{
		"is_active": true,
	})
	return err
}

func (au *AuthUsecase) ResetPassword(userID uint, req dtos.ResetPasswordRequest) error {
	user, err := au.UserRepo.TakeUserByConditions(map[string]interface{}{
		"id": userID,
	})
	if err != nil {
		return err
	}

	if !user.IsActive {
		return errors.New("your account is not active")
	}

	newHashPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("error while hash password")
	}

	err = au.UserRepo.UpdateUser(user, map[string]interface{}{
		"password": newHashPassword,
	})

	return err
}
