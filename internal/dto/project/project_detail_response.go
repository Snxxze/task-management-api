package project

import "task-management-api/internal/dto/task"

type ProjectDetailResponse struct {
	ID          uint
	Name        string
	Description string
	Color       string
	IsArchived  bool

	Tasks []task.TaskResponse
}