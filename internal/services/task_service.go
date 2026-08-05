package services

import (
	"context"
	"strings"
	"time"

	"task-management-api/internal/apperrors"
	taskdto "task-management-api/internal/dto/task"
	"task-management-api/internal/mappers"
	"task-management-api/internal/models"
	"task-management-api/internal/repositories"
	"task-management-api/internal/repositories/filters"
)

type TaskService interface {
	Create(ctx context.Context, userID uint, req taskdto.CreateTaskRequest) (*taskdto.TaskResponse, error)
	Find(ctx context.Context, userID uint, filter filters.TaskFilter) ([]taskdto.TaskResponse, error)
	FindByID(ctx context.Context, id uint, userID uint) (*taskdto.TaskResponse, error)
	Update(ctx context.Context, id uint, userID uint, req taskdto.UpdateTaskRequest) (*taskdto.TaskResponse, error)
	Delete(ctx context.Context, id uint, userID uint) error
}

type taskService struct {
	taskRepo    repositories.TaskRepository
	projectRepo repositories.ProjectRepository
	now         func() time.Time
}

func NewTaskService(
	taskRepo repositories.TaskRepository,
	projectRepo repositories.ProjectRepository,
) TaskService {
	return &taskService{
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
		now:         time.Now,
	}
}

func (s *taskService) Create(
	ctx context.Context,
	userID uint,
	req taskdto.CreateTaskRequest,
) (*taskdto.TaskResponse, error) {
	if strings.TrimSpace(req.Title) == "" {
		return nil, apperrors.ErrBadRequest
	}

	_, err := s.projectRepo.FindByID(ctx, req.ProjectID, userID)
	if err != nil {
		return nil, err
	}

	priority := models.Medium
	if req.Priority != "" {
		p := models.TaskPriority(strings.ToLower(req.Priority))
		if !p.IsValid() {
			return nil, apperrors.ErrBadRequest
		}
		priority = p
	}

	task := models.Task{
		Title:       strings.TrimSpace(req.Title),
		Description: strings.TrimSpace(req.Description),
		Status:      models.Todo,
		Priority:    priority,
		DueDate:     req.DueDate,
		ProjectID:   req.ProjectID,
	}

	if err := s.taskRepo.Create(ctx, &task); err != nil {
		return nil, err
	}

	res := mappers.ToTaskResponse(task)

	return &res, nil
}

func (s *taskService) Find(
	ctx context.Context,
	userID uint,
	filter filters.TaskFilter,
) ([]taskdto.TaskResponse, error) {
	tasks, err := s.taskRepo.Find(ctx, userID, filter)
	if err != nil {
		return nil, err
	}

	res := make([]taskdto.TaskResponse, 0, len(tasks))
	for _, task := range tasks {
		res = append(res, mappers.ToTaskResponse(task))
	}

	return res, nil
}

func (s *taskService) FindByID(
	ctx context.Context,
	id uint,
	userID uint,
) (*taskdto.TaskResponse, error) {
	task, err := s.taskRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	res := mappers.ToTaskResponse(*task)

	return &res, nil
}

func (s *taskService) Update(
	ctx context.Context,
	id uint,
	userID uint,
	req taskdto.UpdateTaskRequest,
) (*taskdto.TaskResponse, error) {
	task, err := s.taskRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		trimmed := strings.TrimSpace(*req.Title)
		if trimmed == "" {
			return nil, apperrors.ErrBadRequest
		}

		task.Title = trimmed
	}

	if req.Description != nil {
		task.Description = strings.TrimSpace(*req.Description)
	}

	if req.Status != nil {
		statusVal := models.TaskStatus(strings.ToLower(*req.Status))
		if !statusVal.IsValid() {
			return nil, apperrors.ErrBadRequest
		}
		task.Status = statusVal
		if statusVal == models.Done && task.CompletedAt == nil {
			now := s.now()
			task.CompletedAt = &now

		} else if statusVal != models.Done {
			task.CompletedAt = nil

		}
	}

	if req.Priority != nil {
		priorityVal := models.TaskPriority(strings.ToLower(*req.Priority))
		if !priorityVal.IsValid() {
			return nil, apperrors.ErrBadRequest
		}
		task.Priority = priorityVal
	}

	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}

	res := mappers.ToTaskResponse(*task)

	return &res, nil
}

func (s *taskService) Delete(
	ctx context.Context,
	id uint,
	userID uint,
) error {
	return s.taskRepo.Delete(ctx, id, userID)
}