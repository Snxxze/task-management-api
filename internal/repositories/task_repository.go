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
	FindByID(ctx context.Context, id uint, userID uint) (*models.Task, error)
	Find(ctx context.Context, userID uint, filter filters.TaskFilter) ([]models.Task, error)
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id uint, userID uint) error
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
	userID uint,
) (*models.Task, error) {
	var task models.Task
	err := r.db.WithContext(ctx).
		Joins("JOIN projects ON projects.id = tasks.project_id AND projects.deleted_at IS NULL").
		Where("tasks.id = ? AND projects.user_id = ?", id, userID).
		Take(&task).Error

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
	userID uint,
	filter filters.TaskFilter,
) ([]models.Task, error) {
	query := r.db.WithContext(ctx).
		Joins("JOIN projects ON projects.id = tasks.project_id AND projects.deleted_at IS NULL").
		Where("projects.user_id = ?", userID)

	if filter.ProjectID != nil {
		query = query.Where("tasks.project_id = ?", *filter.ProjectID)
	}

	if filter.Status != nil {
		query = query.Where("tasks.status = ?", *filter.Status)
	}

	if filter.Priority != nil {
		query = query.Where("tasks.priority = ?", *filter.Priority)
	}

	query = query.Order("tasks.created_at DESC")

	var tasks []models.Task
	err := query.Find(&tasks).Error

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
	userID uint,
) error {
	subQuery := r.db.Model(&models.Project{}).Select("id").Where("user_id = ?", userID)

	result := r.db.WithContext(ctx).
		Where("id = ? AND project_id IN (?)", id, subQuery).
		Delete(&models.Task{})

	if result.Error != nil {
		return fmt.Errorf("delete task: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}