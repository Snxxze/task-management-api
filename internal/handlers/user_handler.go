package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"task-management-api/internal/apperrors"
	userdto "task-management-api/internal/dto/user"
	"task-management-api/internal/middleware"
	"task-management-api/internal/services"
	"task-management-api/internal/util"

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

// GetProfile godoc
// @Summary Get user profile
// @Description Get user profile by id.
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Failure 401 {object} util.Response
// @Failure 403 {object} util.Response
// @Failure 404 {object} util.Response
// @Failure 500 {object} util.Response
// @Security BearerAuth
// @Router /users/{id} [get]
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
		util.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	authUserID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		util.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	if uint(id) != authUserID {
		util.Error(c, http.StatusForbidden, "forbidden")
		return
	}

	res, err := h.userService.GetProfile(
		c.Request.Context(),
		uint(id),
	)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrNotFound):
			util.Error(c, http.StatusNotFound, err.Error())

		default:
			util.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	util.Success(c, http.StatusOK, res)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update user profile by id.
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body user.UpdateUserRequest true "Update user request"
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Failure 401 {object} util.Response
// @Failure 403 {object} util.Response
// @Failure 404 {object} util.Response
// @Failure 409 {object} util.Response
// @Failure 500 {object} util.Response
// @Security BearerAuth
// @Router /users/{id} [patch]
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
		util.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	authUserID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		util.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	if uint(id) != authUserID {
		util.Error(c, http.StatusForbidden, "forbidden")
		return
	}

	var req userdto.UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, http.StatusBadRequest, err.Error())
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
			util.Error(c, http.StatusNotFound, err.Error())

		case errors.Is(err, apperrors.ErrConflict):
			util.Error(c, http.StatusConflict, err.Error())

		case errors.Is(err, apperrors.ErrBadRequest):
			util.Error(c, http.StatusBadRequest, err.Error())

		default:
			util.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	util.Success(c, http.StatusOK, res)
}

// DeleteUser godoc
// @Summary Delete user account
// @Description Delete user account by id.
// @Tags Users
// @Param id path int true "User ID"
// @Success 204
// @Failure 400 {object} util.Response
// @Failure 401 {object} util.Response
// @Failure 403 {object} util.Response
// @Failure 404 {object} util.Response
// @Failure 500 {object} util.Response
// @Security BearerAuth
// @Router /users/{id} [delete]
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
		util.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	authUserID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		util.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	if uint(id) != authUserID {
		util.Error(c, http.StatusForbidden, "forbidden")
		return
	}

	err = h.userService.DeleteUser(
		c.Request.Context(),
		uint(id),
	)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrNotFound):
			util.Error(c, http.StatusNotFound, err.Error())

		default:
			util.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	c.Status(http.StatusNoContent)
}
