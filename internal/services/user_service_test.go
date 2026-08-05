package services_test

import (
	"context"
	"testing"

	"task-management-api/internal/apperrors"
	userdto "task-management-api/internal/dto/user"
	"task-management-api/internal/models"
	"task-management-api/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserService_GetProfile(t *testing.T) {
	ctx := context.Background()
	userID := uint(1)

	t.Run("Happy Path: Get user profile successfully", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		existing := &models.User{
			Name:  "Jane Doe",
			Email: "jane@example.com",
		}
		existing.ID = userID

		userRepo.On("FindByID", ctx, userID).Return(existing, nil)

		service := services.NewUserService(userRepo)
		res, err := service.GetProfile(ctx, userID)

		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, "Jane Doe", res.Name)
		assert.Equal(t, "jane@example.com", res.Email)
		userRepo.AssertExpectations(t)
	})

	t.Run("Error Path: User not found returns NotFound", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		userRepo.On("FindByID", ctx, userID).Return(nil, apperrors.ErrNotFound)

		service := services.NewUserService(userRepo)
		res, err := service.GetProfile(ctx, userID)

		assert.ErrorIs(t, err, apperrors.ErrNotFound)
		assert.Nil(t, res)
		userRepo.AssertExpectations(t)
	})
}

func TestUserService_UpdateProfile(t *testing.T) {
	ctx := context.Background()
	userID := uint(1)
	newName := "Jane Updated"

	t.Run("Happy Path: Update profile successfully", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		existing := &models.User{
			Name:  "Jane Doe",
			Email: "jane@example.com",
		}
		existing.ID = userID

		userRepo.On("FindByID", ctx, userID).Return(existing, nil)
		userRepo.On("Update", ctx, existing).Return(nil)

		service := services.NewUserService(userRepo)
		res, err := service.UpdateProfile(ctx, userID, userdto.UpdateUserRequest{
			Name: &newName,
		})

		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, "Jane Updated", res.Name)
		userRepo.AssertExpectations(t)
	})
}
