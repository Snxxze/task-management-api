package task

import "time"

type TaskResponse struct {
	ID          uint       	`json:"id"`
	
	Title       string     	`json:"title"`
	Description string     	`json:"description"`
	
	Status      string     	`json:"status"`
	Priority    string     	`json:"priority"`
	
	DueDate     *time.Time 	`json:"due_date"`
	CompletedAt *time.Time 	`json:"completed_at"`
	
	ProjectID   uint       	`json:"project_id"`
	
	CreatedAt   time.Time  	`json:"created_at"`
	UpdatedAt   time.Time  	`json:"updated_at"`
}