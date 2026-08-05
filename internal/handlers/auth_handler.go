package handlers

import (
	"errors"
	"net/http"

	"task-management-api/internal/apperrors"
	authdto "task-management-api/internal/dto/auth"
	"task-management-api/internal/services"
	"task-management-api/internal/util"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(
	authService services.AuthService,
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary Register user
// @Description Register a new user account.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body auth.RegisterRequest true "Register request"
// @Success 201 {object} util.Response
// @Failure 400 {object} util.Response
// @Failure 409 {object} util.Response
// @Failure 500 {object} util.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req authdto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrConflict):
			util.Error(c, http.StatusConflict, "email already registered")
		case errors.Is(err, apperrors.ErrBadRequest):
			util.Error(c, http.StatusBadRequest, err.Error())
		default:
			util.Error(c, http.StatusInternalServerError, "internal server error")
		}
		
		return
	}

	util.Success(c, http.StatusCreated, res)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user with email and password.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body auth.LoginRequest true "Login request"
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Failure 401 {object} util.Response
// @Failure 500 {object} util.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req authdto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInvalidCredentials):
			util.Error(c, http.StatusUnauthorized, apperrors.ErrInvalidCredentials.Error())

		case errors.Is(err, apperrors.ErrBadRequest):
			util.Error(c, http.StatusBadRequest, err.Error())

		default:
			util.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	util.Success(c, http.StatusOK, res)
}