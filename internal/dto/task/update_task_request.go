package task

import "time"

type UpdateTaskRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	
	Status      *string    `json:"status"`
	Priority    *string    `json:"priority"`
	
	DueDate     *time.Time `json:"due_date"`
}