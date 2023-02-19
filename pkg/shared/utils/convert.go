package utils

import (
	"engine/internal/pkg/domains/models/dtos"
	"engine/internal/pkg/domains/models/entities"
)

// convertUserEntityToUserResponse func
func ConvertUserEntityToUserResponse(user entities.User) dtos.UserResponse {
	return dtos.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		IsActive: user.IsActive,
	}
}
