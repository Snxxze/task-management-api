package routes

import (
	"task-management-api/internal/bootstrap"
	"task-management-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine, app *bootstrap.Application, jwtSecret string) {
	api := router.Group("/api/v1")

	// Public routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", app.AuthHandler.Register)
		auth.POST("/login", app.AuthHandler.Login)
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(jwtSecret))

	projects := protected.Group("/projects")
	{
		projects.POST("", app.ProjectHandler.Create)
		projects.GET("", app.ProjectHandler.FindAll)
		projects.GET("/:id", app.ProjectHandler.FindByID)
		projects.PATCH("/:id", app.ProjectHandler.Update)
		projects.DELETE("/:id", app.ProjectHandler.Delete)
	}

	users := protected.Group("/users")
	{
		users.GET("/:id", app.UserHandler.GetProfile)
		users.PATCH("/:id", app.UserHandler.UpdateProfile)
		users.DELETE("/:id", app.UserHandler.DeleteUser)
	}

	tasks := protected.Group("/tasks")
	{
		tasks.POST("", app.TaskHandler.Create)
		tasks.GET("", app.TaskHandler.Find)
		tasks.GET("/:id", app.TaskHandler.FindByID)
		tasks.PATCH("/:id", app.TaskHandler.Update)
		tasks.DELETE("/:id", app.TaskHandler.Delete)
	}
}