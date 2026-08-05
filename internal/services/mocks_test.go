package services_test

import (
	"context"
	"task-management-api/internal/models"

	"github.com/stretchr/testify/mock"
)

// MockUserRepository mocks repositories.UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	if user.ID == 0 {
		user.ID = 1 // Simulate primary key generation
	}
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	args := m.Called(ctx, id)
	if res := args.Get(0); res != nil {
		return res.(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if res := args.Get(0); res != nil {
		return res.(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
