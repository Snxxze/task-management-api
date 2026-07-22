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
	TaskHandler    *handlers.TaskHandler
	AuthHandler    *handlers.AuthHandler
}

func New(db *gorm.DB, jwtSecret string) *Application {
	projectRepo := repositories.NewProjectRepository(db)
	projectService := services.NewProjectService(projectRepo)
	projectHandler := handlers.NewProjectHandler(projectService)

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	taskRepo := repositories.NewTaskRepository(db)
	taskService := services.NewTaskService(taskRepo, projectRepo)
	taskHandler := handlers.NewTaskHandler(taskService)
	
	authService := services.NewAuthService(userRepo, jwtSecret)
	authHandler := handlers.NewAuthHandler(authService)
	
	return &Application{
		ProjectHandler: projectHandler,
		UserHandler:    userHandler,
		TaskHandler:    taskHandler,
		AuthHandler:    authHandler,
	}
}