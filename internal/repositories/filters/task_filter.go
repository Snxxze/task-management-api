package filters

import "task-management-api/internal/models"

type TaskFilter struct {
	ProjectID *uint
	Status    *models.TaskStatus
	Priority   *models.TaskPriority
}