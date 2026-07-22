package services

import (
	"context"
	"strings"
	"task-management-api/internal/apperrors"
	projectdto "task-management-api/internal/dto/project"
	"task-management-api/internal/mappers"
	"task-management-api/internal/models"
	"task-management-api/internal/repositories"
)

type ProjectService interface {
	Create(ctx context.Context, userID uint, req projectdto.CreateProjectRequest) (*projectdto.ProjectResponse, error)
	FindAll(ctx context.Context, userID uint) ([]projectdto.ProjectResponse, error)
	FindByID(ctx context.Context, id uint, userID uint) (*projectdto.ProjectResponse, error)
	Update(ctx context.Context, id uint, userID uint, req projectdto.UpdateProjectRequest) (*projectdto.ProjectResponse, error)
	Delete(ctx context.Context, id uint, userID uint) error
}

type projectService struct {
	projectRepo repositories.ProjectRepository
}

func NewProjectService(
	repo repositories.ProjectRepository,
) ProjectService {
	return &projectService{
		projectRepo: repo,
	}
}

func (s *projectService) Create(
	ctx context.Context, 
	userID uint, 
	req projectdto.CreateProjectRequest,
) (*projectdto.ProjectResponse, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, apperrors.ErrBadRequest
	}

	project := models.Project{
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		Color:       req.Color,
		UserID:      userID,
	}

	err := s.projectRepo.Create(ctx, &project)
	if err != nil {
		return nil, err
	}

	res := mappers.ToProjectResponse(project)

	return &res, nil
}

func (s *projectService) FindAll(
	ctx context.Context, 
	userID uint,
) ([]projectdto.ProjectResponse, error) {
	projects, err := s.projectRepo.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	res := make([]projectdto.ProjectResponse, 0, len(projects))

	for _, project := range projects {
		res = append(
			res, 
			mappers.ToProjectResponse(project),
		)
	}

	return res, nil
}

func (s *projectService) FindByID(
	ctx context.Context,
	id uint,
	userID uint,
) (*projectdto.ProjectResponse, error) {
	project, err := s.projectRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	res := mappers.ToProjectResponse(*project)
	return &res, nil
}

func (s *projectService) Update(
	ctx context.Context, 
	id uint, 
	userID uint,
	req projectdto.UpdateProjectRequest,
) (*projectdto.ProjectResponse, error) {
	project, err := s.projectRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		project.Name = strings.TrimSpace(*req.Name)
	}

	if req.Description != nil {
		project.Description = strings.TrimSpace(*req.Description)
	}

	if req.Color != nil {
		project.Color = *req.Color
	}

	if req.IsArchived != nil {
		project.IsArchived = *req.IsArchived
	}

	if strings.TrimSpace(project.Name) == "" {
		return nil, apperrors.ErrBadRequest
	}

	err = s.projectRepo.Update(ctx, project)
	if err != nil {
		return nil, err
	}

	res := mappers.ToProjectResponse(*project)

	return &res, nil
}

func (s *projectService) Delete(
	ctx context.Context, 
	id uint,
	userID uint,
) error {
	err := s.projectRepo.Delete(ctx, id, userID)
	if err != nil {
		return err
	}

	return nil
}
