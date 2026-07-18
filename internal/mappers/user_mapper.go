package mappers

import (
	userdto "task-management-api/internal/dto/user"
	"task-management-api/internal/models"
)

func ToUserResponse(user models.User) userdto.UserResponse {
	return userdto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}
