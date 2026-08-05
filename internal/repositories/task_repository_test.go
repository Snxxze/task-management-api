package repositories_test

import (
	"context"
	"testing"

	"task-management-api/internal/apperrors"
	"task-management-api/internal/models"
	"task-management-api/internal/repositories"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskRepository_Integration(t *testing.T) {
	db, cleanup := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()

	userRepo := repositories.NewUserRepository(db)
	projectRepo := repositories.NewProjectRepository(db)
	taskRepo := repositories.NewTaskRepository(db)
	ctx := context.Background()

	userA := &models.User{Name: "User A", Email: "usera@example.com", PasswordHash: "secret"}
	require.NoError(t, userRepo.Create(ctx, userA))

	userB := &models.User{Name: "User B", Email: "userb@example.com", PasswordHash: "secret"}
	require.NoError(t, userRepo.Create(ctx, userB))

	projA := &models.Project{Name: "Project A", UserID: userA.ID}
	require.NoError(t, projectRepo.Create(ctx, projA))

	taskA := &models.Task{
		Title:     "Task A Title",
		Status:    models.Todo,
		Priority:  models.High,
		ProjectID: projA.ID,
	}
	require.NoError(t, taskRepo.Create(ctx, taskA))
	require.NotZero(t, taskA.ID)

	t.Run("User A can find Task A by ID", func(t *testing.T) {
		found, err := taskRepo.FindByID(ctx, taskA.ID, userA.ID)
		assert.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, "Task A Title", found.Title)
	})

	t.Run("User B CANNOT find Task A (IDOR Protection Verified)", func(t *testing.T) {
		found, err := taskRepo.FindByID(ctx, taskA.ID, userB.ID)
		assert.ErrorIs(t, err, apperrors.ErrNotFound)
		assert.Nil(t, found)
	})

	t.Run("User B CANNOT delete Task A (IDOR Guard)", func(t *testing.T) {
		err := taskRepo.Delete(ctx, taskA.ID, userB.ID)
		assert.ErrorIs(t, err, apperrors.ErrNotFound)
	})
}
