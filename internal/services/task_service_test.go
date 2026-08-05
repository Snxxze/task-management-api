package services_test

import (
	"context"
	"testing"
	"time"

	"task-management-api/internal/apperrors"
	taskdto "task-management-api/internal/dto/task"
	"task-management-api/internal/models"
	"task-management-api/internal/repositories/filters"
	"task-management-api/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockTaskRepository mocks repositories.TaskRepository interface
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(ctx context.Context, task *models.Task) error {
	args := m.Called(ctx, task)
	if task.ID == 0 {
		task.ID = 100 // Simulate DB primary key generation
	}
	return args.Error(0)
}

func (m *MockTaskRepository) FindByID(ctx context.Context, id uint, userID uint) (*models.Task, error) {
	args := m.Called(ctx, id, userID)
	if res := args.Get(0); res != nil {
		return res.(*models.Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskRepository) Find(ctx context.Context, userID uint, filter filters.TaskFilter) ([]models.Task, error) {
	args := m.Called(ctx, userID, filter)
	if res := args.Get(0); res != nil {
		return res.([]models.Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskRepository) Update(ctx context.Context, task *models.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) Delete(ctx context.Context, id uint, userID uint) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

// MockProjectRepository mocks repositories.ProjectRepository interface
type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) Create(ctx context.Context, project *models.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *MockProjectRepository) FindByID(ctx context.Context, id uint, userID uint) (*models.Project, error) {
	args := m.Called(ctx, id, userID)
	if res := args.Get(0); res != nil {
		return res.(*models.Project), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProjectRepository) FindAllByUserID(ctx context.Context, userID uint) ([]models.Project, error) {
	args := m.Called(ctx, userID)
	if res := args.Get(0); res != nil {
		return res.([]models.Project), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProjectRepository) Update(ctx context.Context, project *models.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *MockProjectRepository) Delete(ctx context.Context, id uint, userID uint) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func TestTaskService_Create(t *testing.T) {
	ctx := context.Background()
	userID := uint(1)
	projectID := uint(10)

	tests := []struct {
		name          string
		req           taskdto.CreateTaskRequest
		mockSetup     func(taskRepo *MockTaskRepository, projRepo *MockProjectRepository)
		expectedError error
		checkResult   func(t *testing.T, res *taskdto.TaskResponse)
	}{
		{
			name: "Happy Path: Create task successfully",
			req: taskdto.CreateTaskRequest{
				Title:       "   Build Unit Tests   ",
				Description: "Write comprehensive unit tests",
				Priority:    "HIGH",
				ProjectID:   projectID,
			},
			mockSetup: func(taskRepo *MockTaskRepository, projRepo *MockProjectRepository) {
				// Arrange: User owns the project
				projRepo.On("FindByID", ctx, projectID, userID).Return(&models.Project{UserID: userID}, nil)
				// Arrange: Task creation succeeds
				taskRepo.On("Create", ctx, mock.AnythingOfType("*models.Task")).Return(nil)
			},
			expectedError: nil,
			checkResult: func(t *testing.T, res *taskdto.TaskResponse) {
				assert.Equal(t, "Build Unit Tests", res.Title) // Trimmed
				assert.Equal(t, string(models.High), res.Priority)
				assert.Equal(t, string(models.Todo), res.Status)
			},
		},
		{
			name: "Edge Case: Empty priority defaults to Medium",
			req: taskdto.CreateTaskRequest{
				Title:     "Default Priority Task",
				ProjectID: projectID,
			},
			mockSetup: func(taskRepo *MockTaskRepository, projRepo *MockProjectRepository) {
				projRepo.On("FindByID", ctx, projectID, userID).Return(&models.Project{UserID: userID}, nil)
				taskRepo.On("Create", ctx, mock.AnythingOfType("*models.Task")).Return(nil)
			},
			expectedError: nil,
			checkResult: func(t *testing.T, res *taskdto.TaskResponse) {
				assert.Equal(t, string(models.Medium), res.Priority)
			},
		},
		{
			name: "Error Path: Empty title returns BadRequest",
			req: taskdto.CreateTaskRequest{
				Title:     "   ",
				ProjectID: projectID,
			},
			mockSetup:     func(taskRepo *MockTaskRepository, projRepo *MockProjectRepository) {},
			expectedError: apperrors.ErrBadRequest,
			checkResult:   nil,
		},
		{
			name: "Error Path: Invalid priority returns BadRequest",
			req: taskdto.CreateTaskRequest{
				Title:     "Invalid Priority Task",
				Priority:  "SUPER_HIGH",
				ProjectID: projectID,
			},
			mockSetup: func(taskRepo *MockTaskRepository, projRepo *MockProjectRepository) {
				projRepo.On("FindByID", ctx, projectID, userID).Return(&models.Project{UserID: userID}, nil)
			},
			expectedError: apperrors.ErrBadRequest,
			checkResult:   nil,
		},
		{
			name: "Error Path: Project not found or not owned returns NotFound",
			req: taskdto.CreateTaskRequest{
				Title:     "Valid Task",
				ProjectID: projectID,
			},
			mockSetup: func(taskRepo *MockTaskRepository, projRepo *MockProjectRepository) {
				projRepo.On("FindByID", ctx, projectID, userID).Return(nil, apperrors.ErrNotFound)
			},
			expectedError: apperrors.ErrNotFound,
			checkResult:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			taskRepo := new(MockTaskRepository)
			projRepo := new(MockProjectRepository)
			tt.mockSetup(taskRepo, projRepo)

			service := services.NewTaskService(taskRepo, projRepo)
			require.NotNil(t, service)

			// Act
			res, err := service.Create(ctx, userID, tt.req)

			// Assert
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
			taskRepo.AssertExpectations(t)
			projRepo.AssertExpectations(t)
		})
	}
}

func TestTaskService_Update(t *testing.T) {
	ctx := context.Background()
	userID := uint(1)
	taskID := uint(5)

	doneStatus := "done"
	todoStatus := "todo"
	invalidStatus := "invalid_status"
	emptyTitle := "  "

	tests := []struct {
		name          string
		req           taskdto.UpdateTaskRequest
		mockSetup     func(taskRepo *MockTaskRepository, projRepo *MockProjectRepository)
		expectedError error
		checkResult   func(t *testing.T, res *taskdto.TaskResponse)
	}{
		{
			name: "Happy Path: Update status to Done sets CompletedAt",
			req: taskdto.UpdateTaskRequest{
				Status: &doneStatus,
			},
			mockSetup: func(taskRepo *MockTaskRepository, projRepo *MockProjectRepository) {
				existingTask := &models.Task{
					Status: models.Todo,
				}
				existingTask.ID = taskID

				taskRepo.On("FindByID", ctx, taskID, userID).Return(existingTask, nil)
				taskRepo.On("Update", ctx, mock.AnythingOfType("*models.Task")).Return(nil)
			},
			expectedError: nil,
			checkResult: func(t *testing.T, res *taskdto.TaskResponse) {
				assert.Equal(t, string(models.Done), res.Status)
				assert.NotNil(t, res.CompletedAt)
			},
		},
		{
			name: "Happy Path: Update status from Done back to Todo clears CompletedAt",
			req: taskdto.UpdateTaskRequest{
				Status: &todoStatus,
			},
			mockSetup: func(taskRepo *MockTaskRepository, projRepo *MockProjectRepository) {
				now := time.Now()
				existingTask := &models.Task{
					Status:      models.Done,
					CompletedAt: &now,
				}
				existingTask.ID = taskID

				taskRepo.On("FindByID", ctx, taskID, userID).Return(existingTask, nil)
				taskRepo.On("Update", ctx, mock.AnythingOfType("*models.Task")).Return(nil)
			},
			expectedError: nil,
			checkResult: func(t *testing.T, res *taskdto.TaskResponse) {
				assert.Equal(t, string(models.Todo), res.Status)
				assert.Nil(t, res.CompletedAt)
			},
		},
		{
			name: "Error Path: Task not found returns NotFound",
			req: taskdto.UpdateTaskRequest{
				Status: &doneStatus,
			},
			mockSetup: func(taskRepo *MockTaskRepository, projRepo *MockProjectRepository) {
				taskRepo.On("FindByID", ctx, taskID, userID).Return(nil, apperrors.ErrNotFound)
			},
			expectedError: apperrors.ErrNotFound,
			checkResult:   nil,
		},
		{
			name: "Error Path: Invalid status returns BadRequest",
			req: taskdto.UpdateTaskRequest{
				Status: &invalidStatus,
			},
			mockSetup: func(taskRepo *MockTaskRepository, projRepo *MockProjectRepository) {
				existingTask := &models.Task{Status: models.Todo}
				existingTask.ID = taskID
				taskRepo.On("FindByID", ctx, taskID, userID).Return(existingTask, nil)
			},
			expectedError: apperrors.ErrBadRequest,
			checkResult:   nil,
		},
		{
			name: "Error Path: Empty title returns BadRequest",
			req: taskdto.UpdateTaskRequest{
				Title: &emptyTitle,
			},
			mockSetup: func(taskRepo *MockTaskRepository, projRepo *MockProjectRepository) {
				existingTask := &models.Task{Title: "Existing Title"}
				existingTask.ID = taskID
				taskRepo.On("FindByID", ctx, taskID, userID).Return(existingTask, nil)
			},
			expectedError: apperrors.ErrBadRequest,
			checkResult:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			taskRepo := new(MockTaskRepository)
			projRepo := new(MockProjectRepository)
			tt.mockSetup(taskRepo, projRepo)

			service := services.NewTaskService(taskRepo, projRepo)

			// Act
			res, err := service.Update(ctx, taskID, userID, tt.req)

			// Assert
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
			taskRepo.AssertExpectations(t)
			projRepo.AssertExpectations(t)
		})
	}
}

func TestTaskService_Delete(t *testing.T) {
	ctx := context.Background()
	userID := uint(1)
	taskID := uint(5)

	t.Run("Happy Path: Delete task successfully", func(t *testing.T) {
		// Arrange
		taskRepo := new(MockTaskRepository)
		projRepo := new(MockProjectRepository)
		taskRepo.On("Delete", ctx, taskID, userID).Return(nil)

		service := services.NewTaskService(taskRepo, projRepo)

		// Act
		err := service.Delete(ctx, taskID, userID)

		// Assert
		assert.NoError(t, err)
		taskRepo.AssertExpectations(t)
	})

	t.Run("Error Path: Task to delete not found returns NotFound", func(t *testing.T) {
		// Arrange
		taskRepo := new(MockTaskRepository)
		projRepo := new(MockProjectRepository)
		taskRepo.On("Delete", ctx, taskID, userID).Return(apperrors.ErrNotFound)

		service := services.NewTaskService(taskRepo, projRepo)

		// Act
		err := service.Delete(ctx, taskID, userID)

		// Assert
		assert.ErrorIs(t, err, apperrors.ErrNotFound)
		taskRepo.AssertExpectations(t)
	})
}
