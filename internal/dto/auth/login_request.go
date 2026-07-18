package auth

import "task-management-api/internal/dto/user"

type LoginRequest struct {
	AccessToken 		string	`json:"access_token"`
	User		user.UserResponse		`json:"user"`
}