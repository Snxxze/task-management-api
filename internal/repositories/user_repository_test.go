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

func TestUserRepository_Integration(t *testing.T) {
	db, cleanup := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanup()

	repo := repositories.NewUserRepository(db)
	ctx := context.Background()

	t.Run("Create user and FindByEmail / FindByID", func(t *testing.T) {
		user := &models.User{
			Name:         "Integration User",
			Email:        "integration@example.com",
			PasswordHash: "hashedsecret",
		}

		err := repo.Create(ctx, user)
		require.NoError(t, err)
		assert.NotZero(t, user.ID)

		foundEmail, err := repo.FindByEmail(ctx, "integration@example.com")
		assert.NoError(t, err)
		require.NotNil(t, foundEmail)
		assert.Equal(t, user.Name, foundEmail.Name)

		foundID, err := repo.FindByID(ctx, user.ID)
		assert.NoError(t, err)
		require.NotNil(t, foundID)
		assert.Equal(t, "integration@example.com", foundID.Email)
	})

	t.Run("Duplicate email returns ErrConflict", func(t *testing.T) {
		dupUser := &models.User{
			Name:         "Duplicate User",
			Email:        "integration@example.com",
			PasswordHash: "hashedsecret",
		}

		err := repo.Create(ctx, dupUser)
		assert.ErrorIs(t, err, apperrors.ErrConflict)
	})
}
