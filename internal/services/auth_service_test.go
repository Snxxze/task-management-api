package services_test

import (
	"context"
	"testing"

	"task-management-api/internal/apperrors"
	authdto "task-management-api/internal/dto/auth"
	"task-management-api/internal/models"
	"task-management-api/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register(t *testing.T) {
	ctx := context.Background()
	jwtSecret := "test-secret-key"

	tests := []struct {
		name          string
		req           authdto.RegisterRequest
		mockSetup     func(userRepo *MockUserRepository)
		expectedError error
		checkResult   func(t *testing.T, res *authdto.RegisterResponse)
	}{
		{
			name: "Happy Path: Register user successfully",
			req: authdto.RegisterRequest{
				Name:     "  Test User  ",
				Email:    "Test@Example.com",
				Password: "password123",
			},
			mockSetup: func(userRepo *MockUserRepository) {
				userRepo.On("FindByEmail", ctx, "test@example.com").Return(nil, apperrors.ErrNotFound)
				userRepo.On("Create", ctx, mock.AnythingOfType("*models.User")).Return(nil)
			},
			expectedError: nil,
			checkResult: func(t *testing.T, res *authdto.RegisterResponse) {
				assert.Equal(t, "Test User", res.User.Name)
				assert.Equal(t, "test@example.com", res.User.Email)
				assert.NotEmpty(t, res.AccessToken)
			},
		},
		{
			name: "Error Path: Email already registered returns Conflict",
			req: authdto.RegisterRequest{
				Name:     "Test User",
				Email:    "existing@example.com",
				Password: "password123",
			},
			mockSetup: func(userRepo *MockUserRepository) {
				existing := &models.User{Email: "existing@example.com"}
				existing.ID = 1
				userRepo.On("FindByEmail", ctx, "existing@example.com").Return(existing, nil)
			},
			expectedError: apperrors.ErrConflict,
			checkResult:   nil,
		},
		{
			name: "Error Path: Password too short returns BadRequest",
			req: authdto.RegisterRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "123",
			},
			mockSetup:     func(userRepo *MockUserRepository) {},
			expectedError: apperrors.ErrBadRequest,
			checkResult:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(MockUserRepository)
			tt.mockSetup(userRepo)

			service := services.NewAuthService(userRepo, jwtSecret)
			require.NotNil(t, service)

			res, err := service.Register(ctx, tt.req)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, res)
				if tt.checkResult != nil {
					tt.checkResult(t, res)
				}
			}
			userRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	ctx := context.Background()
	jwtSecret := "test-secret-key"
	plainPassword := "correctpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)

	tests := []struct {
		name          string
		req           authdto.LoginRequest
		mockSetup     func(userRepo *MockUserRepository)
		expectedError error
		checkResult   func(t *testing.T, res *authdto.LoginResponse)
	}{
		{
			name: "Happy Path: Login successfully with valid credentials",
			req: authdto.LoginRequest{
				Email:    "user@example.com",
				Password: plainPassword,
			},
			mockSetup: func(userRepo *MockUserRepository) {
				user := &models.User{
					Name:         "Valid User",
					Email:        "user@example.com",
					PasswordHash: string(hashedPassword),
				}
				user.ID = 5
				userRepo.On("FindByEmail", ctx, "user@example.com").Return(user, nil)
			},
			expectedError: nil,
			checkResult: func(t *testing.T, res *authdto.LoginResponse) {
				assert.Equal(t, uint(5), res.User.ID)
				assert.NotEmpty(t, res.AccessToken)
			},
		},
		{
			name: "Error Path: Wrong password returns InvalidCredentials",
			req: authdto.LoginRequest{
				Email:    "user@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func(userRepo *MockUserRepository) {
				user := &models.User{
					Email:        "user@example.com",
					PasswordHash: string(hashedPassword),
				}
				user.ID = 5
				userRepo.On("FindByEmail", ctx, "user@example.com").Return(user, nil)
			},
			expectedError: apperrors.ErrInvalidCredentials,
			checkResult:   nil,
		},
		{
			name: "Error Path: User not found returns InvalidCredentials",
			req: authdto.LoginRequest{
				Email:    "unknown@example.com",
				Password: plainPassword,
			},
			mockSetup: func(userRepo *MockUserRepository) {
				userRepo.On("FindByEmail", ctx, "unknown@example.com").Return(nil, apperrors.ErrNotFound)
			},
			expectedError: apperrors.ErrInvalidCredentials,
			checkResult:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(MockUserRepository)
			tt.mockSetup(userRepo)

			service := services.NewAuthService(userRepo, jwtSecret)
			require.NotNil(t, service)

			res, err := service.Login(ctx, tt.req)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, res)
				if tt.checkResult != nil {
					tt.checkResult(t, res)
				}
			}
			userRepo.AssertExpectations(t)
		})
	}
}
