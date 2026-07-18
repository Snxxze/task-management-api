package bootstrap

import (
	"task-management-api/internal/handlers"
	"task-management-api/internal/repositories"
	"task-management-api/internal/services"

	"gorm.io/gorm"
)

type Application struct {
	ProjectHandler *handlers.ProjectHandler
	UserHandler    *handlers.UserHandler
}

func New(db *gorm.DB) *Application {
	projectRepo := repositories.NewProjectRepository(db)
	projectService := services.NewProjectService(projectRepo)
	projectHandler := handlers.NewProjectHandler(projectService)

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	return &Application{
		ProjectHandler: projectHandler,
		UserHandler:    userHandler,
	}
}
