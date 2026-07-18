package services

import (
	"context"
	authdto "task-management-api/internal/dto/auth"
)

type AuthService interface {
	Register(ctx context.Context, req authdto.RegisterRequest) (*authdto.RegisterResponse, error)
	Login(ctx context.Context, req authdto.LoginRequest) (*authdto.LoginResponse, error)
}