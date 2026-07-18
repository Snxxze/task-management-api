package repositories

import (
	"context"
	"errors"
	"fmt"
	"task-management-api/internal/apperrors"
	"task-management-api/internal/models"
	"task-management-api/internal/repositories/filters"

	"gorm.io/gorm"
)

type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	FindByID(ctx context.Context, id uint) (*models.Task, error)
	Find(ctx context.Context, filter filters.TaskFilter) ([]models.Task, error)
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id uint) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(
	ctx context.Context,
	task *models.Task,
) error {
	err := r.db.WithContext(ctx).Create(task).Error

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return apperrors.ErrConflict
		}

		return fmt.Errorf("create task: %w", err)
	}

	return nil
}

func (r *taskRepository) FindByID(
	ctx context.Context,
	id uint,
) (*models.Task, error) {
	var task models.Task
	err := r.db.WithContext(ctx).Take(&task, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}

		return nil, fmt.Errorf("find task by id: %w", err)
	}

	return &task, nil
}

func (r *taskRepository) Find(
	ctx context.Context,
	filter filters.TaskFilter,
) ([]models.Task, error) {
	db := r.db.WithContext(ctx)

	if filter.ProjectID != nil {
		db = db.Where("project_id = ?", *filter.ProjectID)
	}

	if filter.Status != nil {
		db = db.Where("status = ?", *filter.Status)
	}

	if filter.Priority != nil {
		db = db.Where("priority = ?", *filter.Priority)
	}

	db = db.Order("created_at DESC")

	var tasks []models.Task
	err := db.Find(&tasks).Error

	if err != nil {
		return nil, fmt.Errorf("find tasks: %w", err)
	}

	return tasks, nil
}

func (r *taskRepository) Update(
	ctx context.Context,
	task *models.Task,
) error {
	err := r.db.WithContext(ctx).Save(task).Error

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return apperrors.ErrConflict
		}

		return fmt.Errorf("update task: %w", err)
	}

	return nil
}

func (r *taskRepository) Delete(
	ctx context.Context,
	id uint,
) error {
	result := r.db.WithContext(ctx).Delete(&models.Task{}, id)

	if result.Error != nil {
		return fmt.Errorf("delete task: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}