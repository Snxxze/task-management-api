package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"task-management-api/internal/apperrors"
	projectdto "task-management-api/internal/dto/project"
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
// @Description Create a new project for the current test user.
// @Tags Projects
// @Accept json
// @Produce json
// @Param request body project.CreateProjectRequest true "Create project request"
// @Success 201 {object} project.ProjectResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
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

	// TODO: Replace with authenticated user ID from JWT middleware.
	userID := uint(1)

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
// @Description Get all projects for the current test user.
// @Tags Projects
// @Produce json
// @Success 200 {array} project.ProjectResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /projects [get]
func (h *ProjectHandler) FindAll(
	c *gin.Context,
) {
	userID := uint(1)

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
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
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

	res, err := h.projectService.Update(
		c.Request.Context(),
		uint(id),
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
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
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

	err = h.projectService.Delete(
		c.Request.Context(),
		uint(id),
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
