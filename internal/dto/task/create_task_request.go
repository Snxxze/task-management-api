package task

import "time"

type CreateTaskRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	
	Priority    string     `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
	
	ProjectID   uint       `json:"project_id" binding:"required"`
}