package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"task-management-api/internal/apperrors"
	projectdto "task-management-api/internal/dto/project"
	"task-management-api/internal/middleware"
	"task-management-api/internal/services"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	projectService services.ProjectService
}

func NewProjectHandler(
	projectService services.ProjectService,
) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

// Create godoc
// @Summary Create project
// @Description Create a new project for the current authenticated user.
// @Tags Projects
// @Accept json
// @Produce json
// @Param request body project.CreateProjectRequest true "Create project request"
// @Success 201 {object} project.ProjectResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /projects [post]
func (h *ProjectHandler) Create(
	c *gin.Context,
) {
	var req projectdto.CreateProjectRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	res, err := h.projectService.Create(
		c.Request.Context(),
		userID,
		req,
	)

	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, apperrors.ErrConflict):
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, apperrors.ErrBadRequest):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}

		return
	}

	c.JSON(http.StatusCreated, res)
}

// FindAll godoc
// @Summary Get projects
// @Description Get all projects for the current authenticated user.
// @Tags Projects
// @Produce json
// @Success 200 {array} project.ProjectResponse
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /projects [get]
func (h *ProjectHandler) FindAll(
	c *gin.Context,
) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	projects, err := h.projectService.FindAll(
		c.Request.Context(),
		userID,
	)

	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}

		return
	}

	c.JSON(http.StatusOK, projects)
}

// FindByID godoc
// @Summary Get project by id
// @Description Get a single project by id for the current authenticated user.
// @Tags Projects
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} project.ProjectResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /projects/{id} [get]
func (h *ProjectHandler) FindByID(
	c *gin.Context,
) {
	idParam := c.Param("id")

	id, err := strconv.ParseUint(
		idParam,
		10,
		32,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project id",
		})
		return
	}

	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	res, err := h.projectService.FindByID(
		c.Request.Context(),
		uint(id),
		userID,
	)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}

		return
	}

	c.JSON(http.StatusOK, res)
}

// Update godoc
// @Summary Update project
// @Description Update a project by id.
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param request body project.UpdateProjectRequest true "Update project request"
// @Success 200 {object} project.ProjectResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /projects/{id} [patch]
func (h *ProjectHandler) Update(
	c *gin.Context,
) {
	idParam := c.Param("id")

	id, err := strconv.ParseUint(
		idParam,
		10,
		32,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project id",
		})
		return
	}

	var req projectdto.UpdateProjectRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	res, err := h.projectService.Update(
		c.Request.Context(),
		uint(id),
		userID,
		req,
	)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, apperrors.ErrConflict):
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, apperrors.ErrBadRequest):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}

		return
	}

	c.JSON(http.StatusOK, res)
}

// Delete godoc
// @Summary Delete project
// @Description Delete a project by id.
// @Tags Projects
// @Param id path int true "Project ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /projects/{id} [delete]
func (h *ProjectHandler) Delete(
	c *gin.Context,
) {
	idParam := c.Param("id")

	id, err := strconv.ParseUint(
		idParam,
		10,
		32,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project id",
		})
		return
	}

	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	err = h.projectService.Delete(
		c.Request.Context(),
		uint(id),
		userID,
	)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}

		return
	}

	c.Status(http.StatusNoContent)
}
