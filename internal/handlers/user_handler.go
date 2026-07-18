package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"task-management-api/internal/apperrors"
	userdto "task-management-api/internal/dto/user"
	"task-management-api/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(
	userService services.UserService,
) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetProfile(
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
			"error": "invalid user id",
		})
		return
	}

	res, err := h.userService.GetProfile(
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

	c.JSON(http.StatusOK, res)
}

func (h *UserHandler) UpdateProfile(
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
			"error": "invalid user id",
		})
		return
	}

	var req userdto.UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	res, err := h.userService.UpdateProfile(
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

func (h *UserHandler) DeleteUser(
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
			"error": "invalid user id",
		})
		return
	}

	err = h.userService.DeleteUser(
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
