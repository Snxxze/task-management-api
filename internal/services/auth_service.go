package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"task-management-api/internal/apperrors"
	authdto "task-management-api/internal/dto/auth"
	"task-management-api/internal/mappers"
	"task-management-api/internal/models"
	"task-management-api/internal/repositories"
	"task-management-api/internal/util"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req authdto.RegisterRequest) (*authdto.RegisterResponse, error)
	Login(ctx context.Context, req authdto.LoginRequest) (*authdto.LoginResponse, error)
}

type authService struct {
	userRepo  repositories.UserRepository
	jwtSecret string
}

func NewAuthService(
	userRepo repositories.UserRepository,
	jwtSecret string,
) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) Register(
	ctx context.Context,
	req authdto.RegisterRequest,
) (*authdto.RegisterResponse, error) {
	name := strings.TrimSpace(req.Name)
	email := strings.ToLower(strings.TrimSpace(req.Email))
	password := req.Password

	if name == "" || email == "" || len(password) < 6 {
		return nil, apperrors.ErrBadRequest
	}

	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, apperrors.ErrConflict
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := models.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	if err := s.userRepo.Create(ctx, &user); err != nil {
		return nil, err
	}

	token, err := util.GenerateToken(user.ID, s.jwtSecret, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &authdto.RegisterResponse{
		User:        mappers.ToUserResponse(user),
		AccessToken: token,
	}, nil
}

func (s *authService) Login(
	ctx context.Context,
	req authdto.LoginRequest,
) (*authdto.LoginResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	password := req.Password

	if email == "" || password == "" {
		return nil, apperrors.ErrBadRequest
	}

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrInvalidCredentials
		}

		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	token, err := util.GenerateToken(user.ID, s.jwtSecret, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &authdto.LoginResponse{
		User:        mappers.ToUserResponse(*user),
		AccessToken: token,
	}, nil
}