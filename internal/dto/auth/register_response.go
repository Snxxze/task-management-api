package auth

import "task-management-api/internal/dto/user"

type RegisterResponse struct {
	User        	user.UserResponse 	`json:"user"`
	AccessToken 	string            	`json:"access_token"`
}