package routes

import (
	"task-management-api/internal/bootstrap"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine, app *bootstrap.Application) {
	api := router.Group("/api/v1")

	projects := api.Group("/projects")
	{
		projects.POST("", app.ProjectHandler.Create)
		projects.GET("", app.ProjectHandler.FindAll)
		projects.PATCH("/:id", app.ProjectHandler.Update)
		projects.DELETE("/:id", app.ProjectHandler.Delete)
	}

	users := api.Group("/users")
	{
		users.GET("/:id", app.UserHandler.GetProfile)
		users.PATCH("/:id", app.UserHandler.UpdateProfile)
		users.DELETE("/:id", app.UserHandler.DeleteUser)
	}
}
