package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"engine/internal/pkg/domains/interfaces"
	"engine/internal/pkg/domains/models/dtos"
	"engine/internal/pkg/repositories"
	"engine/internal/pkg/usecases"
	"engine/pkg/shared/middleware"
	"engine/pkg/shared/utils"
)

type TeamHandler struct {
	UserUsecase interfaces.UserUsecase
	TeamUsecase interfaces.TeamUsecase
}

func NewTeamHandler(dbConn *gorm.DB) *TeamHandler {
	userRepo := repositories.NewUserRepository(dbConn)
	teamRepo := repositories.NewTeamRepository(dbConn)
	teamUserRepo := repositories.NewTeamUserRepository(dbConn)
	userUsecase := usecases.NewUserUsecase(userRepo)
	teamUsecase := usecases.NewTeamUsecase(teamRepo, teamUserRepo)
	return &TeamHandler{
		UserUsecase: userUsecase,
		TeamUsecase: teamUsecase,
	}
}

func (th *TeamHandler) CreateTeam(c *gin.Context) {
	req := dtos.CreateTeamRequest{}
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

	email, err := middleware.GetUserInfoFromToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.BaseResponse{
			Status: "failed",
			Error: &dtos.ErrorResponse{
				ErrorMessage: err.Error(),
			},
		})
		return
	}

	user, err := th.UserUsecase.TakeUserByConditions(map[string]interface{}{
		"email": email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.BaseResponse{
			Status: "failed",
			Error: &dtos.ErrorResponse{
				ErrorMessage: err.Error(),
			},
		})
		return
	}

	team, err := th.TeamUsecase.CreateTeam(req, user)
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
			"team_info": utils.ConvertTeamEntityToTeamResponse(team),
		},
	})
}
