package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"

	"engine/config"
	"engine/internal/pkg/domains/interfaces"
	"engine/internal/pkg/domains/models/dtos"
	"engine/internal/pkg/domains/models/entities"
	"engine/internal/pkg/repositories"
	"engine/internal/pkg/usecases"
	"engine/pkg/shared/auth"
	"engine/pkg/shared/constants"
	"engine/pkg/shared/utils"
)

type AuthHandler struct {
	AuthUsecase interfaces.AuthUsecase
}

func NewAuthHandler(dbConn *gorm.DB) *AuthHandler {
	authRepo := repositories.NewUserRepository(dbConn)
	authUsecase := usecases.NewAuthUsecase(authRepo)
	return &AuthHandler{
		AuthUsecase: authUsecase,
	}
}

func (ah *AuthHandler) SignUp(c *gin.Context) {
	req := dtos.CreateUserRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.BaseResponse{
			Status: "failed",
			Error: &dtos.ErrorResponse{
				ErrorMessage: err.Error(),
			},
		})
		return
	}

	user, err := ah.AuthUsecase.SignUp(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Status: "failed",
			Error: &dtos.ErrorResponse{
				ErrorMessage: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, dtos.BaseResponse{
		Status: "success",
		Data: gin.H{
			"message":   "send verify email success",
			"user_info": utils.ConvertUserEntityToUserResponse(user),
		},
	})
}

func (ah *AuthHandler) SignIn(c *gin.Context) {
	req := dtos.SignInRequest{}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.BaseResponse{
			Status: "failed",
			Error: &dtos.ErrorResponse{
				ErrorMessage: err.Error(),
			},
		})
		return
	}

	user, token, err := ah.AuthUsecase.SignIn(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Status: "failed",
			Error: &dtos.ErrorResponse{
				ErrorMessage: err.Error(),
			},
		})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusBadRequest, dtos.BaseResponse{
			Status: "failed",
			Error: &dtos.ErrorResponse{
				ErrorMessage: "user is not active",
			},
		})
		return
	}

	c.JSON(http.StatusOK, dtos.BaseResponse{
		Status: "success",
		Data: gin.H{
			"access_token": token,
			"user_info":    utils.ConvertUserEntityToUserResponse(user),
		},
	})
}

func (ah *AuthHandler) SignInWithGoogle(c *gin.Context) {
	// Create oauthState cookie
	oauthState := utils.GenerateStateOauthCookie(c)
	/*
		AuthCodeURL receive state that is a token to protect the user
		from CSRF attacks. You must always provide a non-empty string
		and validate that it matches the the state query parameter
		on your redirect callback.
	*/
	url := config.AppConfig.GoogleLoginConfig.AuthCodeURL(oauthState)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (ah *AuthHandler) Redirect(c *gin.Context) {
	// get oauth state from cookie for this user
	oauthState, _ := c.Request.Cookie("oauthstate")
	state := c.Request.FormValue("state")
	code := c.Request.FormValue("code")
	c.Header("content-type", "application/json")

	// ERROR : Invalid OAuth State
	if state != oauthState.Value {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Status: "failed",
			Error: &dtos.ErrorResponse{
				ErrorMessage: "invalid oauth state",
			},
		})
		return
	}

	// Exchange Auth Code for Tokens
	oauthToken, err := config.AppConfig.GoogleLoginConfig.Exchange(context.Background(), code)

	// ERROR : Auth Code Exchange Failed
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Status: "failed",
			Error: &dtos.ErrorResponse{
				ErrorMessage: "falied code exchange: " + err.Error(),
			},
		})
		return
	}

	// Fetch User Data from google server
	response, err := http.Get(constants.OauthGoogleUrlAPI + oauthToken.AccessToken)
	// ERROR : Unable to get user data from google
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Status: "failed",
			Error: &dtos.ErrorResponse{
				ErrorMessage: "failed getting user info: " + err.Error(),
			},
		})
		return
	}

	// Parse user data JSON Object
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Status: "failed",
			Error: &dtos.ErrorResponse{
				ErrorMessage: "failed read response: " + err.Error(),
			},
		})
		return
	}

	err = json.Unmarshal(contents, &config.GoogleUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Status: "failed",
			Error: &dtos.ErrorResponse{
				ErrorMessage: "failed unmarshal response: " + err.Error(),
			},
		})
		return
	}

	user, err := ah.AuthUsecase.TakeByConditions(map[string]interface{}{
		"email": config.GoogleUser.Email,
	})

	if err != nil && err == gorm.ErrRecordNotFound {
		req := dtos.CreateUserRequest{
			Username: config.GoogleUser.Name,
			Email:    config.GoogleUser.Email,
			Password: config.GoogleUser.ID,
			IsActive: true,
		}

		user, err = ah.AuthUsecase.SignUp(req)
		if err != nil {
			c.JSON(http.StatusBadRequest, dtos.BaseResponse{
				Status: "failed",
				Error: &dtos.ErrorResponse{
					ErrorMessage: err.Error(),
				},
			})
			return
		}

		jwtToken, err := auth.GenerateHS256JWT(map[string]interface{}{
			"email":    config.GoogleUser.Email,
			"username": config.GoogleUser.Name,
			"exp":      time.Now().Add(time.Hour * 72).Unix(),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
				Status: "failed",
				Error: &dtos.ErrorResponse{
					ErrorMessage: "failed while generate token: " + err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusOK, dtos.BaseResponse{
			Status: "success",
			Data: gin.H{
				"access_token": jwtToken,
				"user_info":    utils.ConvertUserEntityToUserResponse(user),
			},
		})
	} else if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Status: "failed",
			Error: &dtos.ErrorResponse{
				ErrorMessage: err.Error(),
			},
		})
		return
	} else {
		jwtToken, err := auth.GenerateHS256JWT(map[string]interface{}{
			"email":    config.GoogleUser.Email,
			"username": config.GoogleUser.Name,
			"exp":      time.Now().Add(time.Hour * 72).Unix(),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
				Status: "failed",
				Error: &dtos.ErrorResponse{
					ErrorMessage: "failed while generate token: " + err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusOK, dtos.BaseResponse{
			Status: "success",
			Data: gin.H{
				"access_token": jwtToken,
				"user_info":    utils.ConvertUserEntityToUserResponse(user),
			},
		})
	}
}

func (ah *AuthHandler) ForgotPassword(c *gin.Context) {
	req := dtos.ForgotPasswordRequest{}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.BaseResponse{
			Status: "failed",
			Error: &dtos.ErrorResponse{
				ErrorMessage: err.Error(),
			},
		})
		return
	}

	err = ah.AuthUsecase.SendMailForgotPassword(req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Status: "failed",
			Error: &dtos.ErrorResponse{
				ErrorMessage: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, dtos.BaseResponse{
		Status: "success",
		Data:   gin.H{"message": "send success"},
	})
}

func (ah *AuthHandler) VerifyResetPasswordLink(c *gin.Context) {
	_, ok := ah.VerifyParam(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Status: "failed",
			Error: &dtos.ErrorResponse{
				ErrorMessage: "verify param error",
			},
		})
	}

	c.JSON(http.StatusOK, dtos.BaseResponse{
		Status: "success",
		Data:   gin.H{"message": "valid link"},
	})
}

func (ah *AuthHandler) VerifyEmailAddress(c *gin.Context) {
	user, ok := ah.VerifyParam(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Status: "failed",
			Error: &dtos.ErrorResponse{
				ErrorMessage: "verify param error",
			},
		})
	}

	err := ah.AuthUsecase.ActiveUser(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Status: "failed",
			Error: &dtos.ErrorResponse{
				ErrorMessage: "failed to verify email address",
			},
		})
	}

	c.JSON(http.StatusOK, dtos.BaseResponse{
		Status: "success",
		Data: gin.H{
			"message": "active user success",
		},
	})
}

func (ah *AuthHandler) PatchResetPassword(c *gin.Context) {
	user, ok := ah.VerifyParam(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Status: "failed",
			Error: &dtos.ErrorResponse{
				ErrorMessage: "verify param error",
			},
		})
	}

	req := dtos.ResetPasswordRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.BaseResponse{
			Status: "failed",
			Error: &dtos.ErrorResponse{
				ErrorMessage: err.Error(),
			},
		})
		return
	}

	err = ah.AuthUsecase.ResetPassword(user.ID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Status: "failed",
			Error: &dtos.ErrorResponse{
				ErrorMessage: "failed to reset password",
			},
		})
	}

	c.JSON(http.StatusOK, dtos.BaseResponse{
		Status: "success",
		Data: gin.H{
			"message": "reset password success",
		},
	})
}

func (ah *AuthHandler) VerifyParam(c *gin.Context) (entities.User, bool) {
	email := c.Param("email")
	if email == "" {
		return entities.User{}, false
	}

	decodedEmail, err := base64.StdEncoding.DecodeString(email)
	if err != nil {
		return entities.User{}, false
	}

	token := c.Param("token")
	if token == "" {
		return entities.User{}, false
	}

	decodedToken, err := auth.DecodeJWT(token)
	if err != nil {
		return entities.User{}, false
	}

	user, err := ah.AuthUsecase.TakeByConditions(map[string]interface{}{
		"email": decodedToken.Claims.(jwt.MapClaims)["email"],
	})
	if err != nil {
		return entities.User{}, false
	}

	if string(decodedEmail) != user.Email {
		return entities.User{}, false
	}

	verifyToken := auth.VerifyJWT(token)
	if !verifyToken {
		return entities.User{}, false
	}

	return user, true
}
