package mappers

import (
	projectdto "task-management-api/internal/dto/project"
	"task-management-api/internal/models"
)

func ToProjectResponse(project models.Project) projectdto.ProjectResponse {
	return projectdto.ProjectResponse{
		ID:          	project.ID,
		Name:        	project.Name,
		Description: 	project.Description,
		Color:       	project.Color,
		IsArchived:  	project.IsArchived,
		CreatedAt:   	project.CreatedAt,
		UpdatedAt:   	project.UpdatedAt,
	}
}