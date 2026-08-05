package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"task-management-api/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequestIDMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Happy Path: Generates new Request ID if header missing", func(t *testing.T) {
		router := gin.New()
		router.Use(middleware.RequestIDMiddleware())
		router.GET("/test-id", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest(http.MethodGet, "/test-id", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		respHeader := w.Header().Get(middleware.RequestIDHeader)
		assert.NotEmpty(t, respHeader)
	})

	t.Run("Happy Path: Preserves incoming Request ID if header provided", func(t *testing.T) {
		router := gin.New()
		router.Use(middleware.RequestIDMiddleware())
		router.GET("/test-id", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		customReqID := "my-custom-trace-id-12345"
		req, _ := http.NewRequest(http.MethodGet, "/test-id", nil)
		req.Header.Set(middleware.RequestIDHeader, customReqID)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, customReqID, w.Header().Get(middleware.RequestIDHeader))
	})
}
