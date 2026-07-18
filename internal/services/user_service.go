package services

import (
	"context"
	"strings"
	"task-management-api/internal/apperrors"
	userdto "task-management-api/internal/dto/user"
	"task-management-api/internal/mappers"
	"task-management-api/internal/repositories"
)

type UserService interface {
	GetProfile(ctx context.Context, id uint) (*userdto.UserResponse, error)
	UpdateProfile(ctx context.Context, id uint, req userdto.UpdateUserRequest) (*userdto.UserResponse, error)
	DeleteUser(ctx context.Context, id uint) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(
	repo repositories.UserRepository,
) UserService {
	return &userService{
		userRepo: repo,
	}
}

func (s *userService) GetProfile(
	ctx context.Context,
	id uint,
) (*userdto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	res := mappers.ToUserResponse(*user)

	return &res, nil
}

func (s *userService) UpdateProfile(
	ctx context.Context,
	id uint,
	req userdto.UpdateUserRequest,
) (*userdto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		user.Name = strings.TrimSpace(*req.Name)
	}

	if req.Email != nil {
		user.Email = strings.TrimSpace(*req.Email)
	}

	if strings.TrimSpace(user.Name) == "" || strings.TrimSpace(user.Email) == "" {
		return nil, apperrors.ErrBadRequest
	}

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	res := mappers.ToUserResponse(*user)

	return &res, nil
}

func (s *userService) DeleteUser(
	ctx context.Context,
	id uint,
) error {
	err := s.userRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
