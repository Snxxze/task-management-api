package services

import (
	"context"
	taskdto "task-management-api/internal/dto/task"
	"task-management-api/internal/repositories/filters"
)

type TaskService interface {
	Create(ctx context.Context, userID uint, req taskdto.CreateTaskRequest) (*taskdto.TaskResponse, error)
	Find(ctx context.Context, filter filters.TaskFilter) ([]taskdto.TaskResponse, error)
	Update(ctx context.Context, id uint, req taskdto.UpdateTaskRequest) (*taskdto.TaskResponse, error)
	Delete(ctx context.Context, id uint) error
}