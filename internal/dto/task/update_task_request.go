package task

import "time"

type UpdateTaskRequest struct {
	Title       *string
	Description *string

	Status   *string
	Priority *string

	DueDate *time.Time
}