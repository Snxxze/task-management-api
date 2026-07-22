package mappers

import (
	taskdto "task-management-api/internal/dto/task"
	"task-management-api/internal/models"
)

func ToTaskResponse(task models.Task) taskdto.TaskResponse {
	return taskdto.TaskResponse{
		ID:          	task.ID,
		Title:       	task.Title,
		Description: 	task.Description,
		Status:      	string(task.Status),
		Priority:    	string(task.Priority),
		DueDate:     	task.DueDate,
		CompletedAt: 	task.CompletedAt,
		ProjectID:   	task.ProjectID,
		CreatedAt:   	task.CreatedAt,
		UpdatedAt:   	task.UpdatedAt,
	}
}