package task

import "time"

type TaskResponse struct {
	ID uint

	Title       string
	Description string

	Status   string
	Priority string

	DueDate *time.Time
	CompletedAt		*time.Time

	ProjectID		uint
}