package repositories

import (
	"context"
	"errors"
	"fmt"
	"task-management-api/internal/apperrors"
	"task-management-api/internal/models"

	"gorm.io/gorm"
)

type ProjectRepository interface {
	Create(ctx context.Context, project *models.Project) error
	FindByID(ctx context.Context, id uint, userID uint) (*models.Project, error)
	FindAllByUserID(ctx context.Context, userID uint) ([]models.Project, error)
	Update(ctx context.Context, project *models.Project) error
	Delete(ctx context.Context, id uint, userID uint) error
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(
	ctx context.Context,
	project *models.Project,
) error {
	err := r.db.WithContext(ctx).Create(project).Error

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return apperrors.ErrConflict
		}

		return fmt.Errorf("create project: %w", err)
	}

	return nil
}

func (r *projectRepository) FindByID(
	ctx context.Context,
	id uint,
	userID uint,
) (*models.Project, error) {
	var project models.Project
	err := r.db.WithContext(ctx).Take(&project, "id = ? AND user_id = ?", id, userID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}

		return nil, fmt.Errorf("find project by id: %w", err)
	}

	return &project, nil
}

func (r *projectRepository) FindAllByUserID(
	ctx context.Context,
	userID uint,
) ([]models.Project, error) {
	var projects []models.Project
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&projects).Error

	if err != nil {
		return nil, fmt.Errorf("find projects by user id: %w", err)
	}

	return projects, nil
}

func (r *projectRepository) Update(
	ctx context.Context, 
	project *models.Project,
) error {
	err := r.db.WithContext(ctx).Save(project).Error

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return apperrors.ErrConflict
		}

		return fmt.Errorf("update project by id: %w", err)
	}

	return nil
}

func (r *projectRepository) Delete(
	ctx context.Context,
	id uint,
	userID uint,
) error {
	result := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&models.Project{})

	if result.Error != nil {
		return fmt.Errorf("delete project: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}