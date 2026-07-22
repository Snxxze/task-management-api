package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"task-management-api/internal/apperrors"
	taskdto "task-management-api/internal/dto/task"
	"task-management-api/internal/middleware"
	"task-management-api/internal/models"
	"task-management-api/internal/repositories/filters"
	"task-management-api/internal/services"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskService services.TaskService
}

func NewTaskHandler(
	taskService services.TaskService,
) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

// Create godoc
// @Summary Create task
// @Description Create a new task in a project.
// @Tags Tasks
// @Accept json
// @Produce json
// @Param request body task.CreateTaskRequest true "Create task request"
// @Success 201 {object} task.TaskResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks [post]
func (h *TaskHandler) Create(c *gin.Context) {
	var req taskdto.CreateTaskRequest
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

	res, err := h.taskService.Create(c.Request.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{
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

// Find godoc
// @Summary Get tasks
// @Description Get tasks with optional filters (project_id, status, priority).
// @Tags Tasks
// @Produce json
// @Param project_id query int false "Project ID"
// @Param status query string false "Status (todo, doing, done)"
// @Param priority query string false "Priority (low, medium, high)"
// @Success 200 {array} task.TaskResponse
// @Failure 500 {object} map[string]string
// @Router /tasks [get]
func (h *TaskHandler) Find(c *gin.Context) {
	var filter filters.TaskFilter

	if projectIDStr := c.Query("project_id"); projectIDStr != "" {
		if pid, err := strconv.ParseUint(projectIDStr, 10, 32); err == nil {
			p := uint(pid)
			filter.ProjectID = &p
		}
	}

	if statusStr := c.Query("status"); statusStr != "" {
		st := models.TaskStatus(statusStr)
		filter.Status = &st
	}

	if priorityStr := c.Query("priority"); priorityStr != "" {
		pr := models.TaskPriority(priorityStr)
		filter.Priority = &pr
	}

	tasks, err := h.taskService.Find(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// FindByID godoc
// @Summary Get task by id
// @Description Get a single task by id.
// @Tags Tasks
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} task.TaskResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id} [get]
func (h *TaskHandler) FindByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task id",
		})
		return
	}

	res, err := h.taskService.FindByID(c.Request.Context(), uint(id))
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
// @Summary Update task
// @Description Update a task by id.
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param request body task.UpdateTaskRequest true "Update task request"
// @Success 200 {object} task.TaskResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id} [patch]
func (h *TaskHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task id",
		})
		return
	}

	var req taskdto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	res, err := h.taskService.Update(c.Request.Context(), uint(id), req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{
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
// @Summary Delete task
// @Description Delete a task by id.
// @Tags Tasks
// @Param id path int true "Task ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id} [delete]
func (h *TaskHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task id",
		})
		return
	}

	err = h.taskService.Delete(c.Request.Context(), uint(id))
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