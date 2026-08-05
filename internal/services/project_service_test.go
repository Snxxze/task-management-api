package services_test

import (
	"context"
	"testing"

	"task-management-api/internal/apperrors"
	projectdto "task-management-api/internal/dto/project"
	"task-management-api/internal/models"
	"task-management-api/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProjectService_Create(t *testing.T) {
	ctx := context.Background()
	userID := uint(1)

	tests := []struct {
		name          string
		req           projectdto.CreateProjectRequest
		mockSetup     func(projRepo *MockProjectRepository)
		expectedError error
		checkResult   func(t *testing.T, res *projectdto.ProjectResponse)
	}{
		{
			name: "Happy Path: Create project successfully",
			req: projectdto.CreateProjectRequest{
				Name:        "  Backend Project  ",
				Description: "Build REST API",
				Color:       "#FF0000",
			},
			mockSetup: func(projRepo *MockProjectRepository) {
				projRepo.On("Create", ctx, mock.AnythingOfType("*models.Project")).Return(nil)
			},
			expectedError: nil,
			checkResult: func(t *testing.T, res *projectdto.ProjectResponse) {
				assert.Equal(t, "Backend Project", res.Name)
				assert.Equal(t, "Build REST API", res.Description)
				assert.Equal(t, "#FF0000", res.Color)
			},
		},
		{
			name: "Error Path: Empty name returns BadRequest",
			req: projectdto.CreateProjectRequest{
				Name: "   ",
			},
			mockSetup:     func(projRepo *MockProjectRepository) {},
			expectedError: apperrors.ErrBadRequest,
			checkResult:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projRepo := new(MockProjectRepository)
			tt.mockSetup(projRepo)

			service := services.NewProjectService(projRepo)
			require.NotNil(t, service)

			res, err := service.Create(ctx, userID, tt.req)

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
			projRepo.AssertExpectations(t)
		})
	}
}

func TestProjectService_FindByID(t *testing.T) {
	ctx := context.Background()
	userID := uint(1)
	projectID := uint(10)

	t.Run("Happy Path: Find project by ID successfully", func(t *testing.T) {
		projRepo := new(MockProjectRepository)
		existing := &models.Project{
			Name:   "My Project",
			UserID: userID,
		}
		existing.ID = projectID

		projRepo.On("FindByID", ctx, projectID, userID).Return(existing, nil)

		service := services.NewProjectService(projRepo)
		res, err := service.FindByID(ctx, projectID, userID)

		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, "My Project", res.Name)
		projRepo.AssertExpectations(t)
	})

	t.Run("Error Path: Project not found returns NotFound", func(t *testing.T) {
		projRepo := new(MockProjectRepository)
		projRepo.On("FindByID", ctx, projectID, userID).Return(nil, apperrors.ErrNotFound)

		service := services.NewProjectService(projRepo)
		res, err := service.FindByID(ctx, projectID, userID)

		assert.ErrorIs(t, err, apperrors.ErrNotFound)
		assert.Nil(t, res)
		projRepo.AssertExpectations(t)
	})
}
