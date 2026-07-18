package task

import "time"

type CreateTaskRequest struct {
	Title       string
	Description string

	Priority string
	DueDate  *time.Time

	ProjectID		uint
}