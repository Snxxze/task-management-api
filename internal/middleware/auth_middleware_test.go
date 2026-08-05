package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"task-management-api/internal/middleware"
	"task-management-api/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtSecret := "test-jwt-secret-key"

	validToken, err := util.GenerateToken(42, jwtSecret, 1*time.Hour)
	assert.NoError(t, err)

	expiredToken, err := util.GenerateToken(42, jwtSecret, -1*time.Hour)
	assert.NoError(t, err)

	invalidSecretToken, err := util.GenerateToken(42, "wrong-secret", 1*time.Hour)
	assert.NoError(t, err)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		checkContext   bool
	}{
		{
			name:           "Happy Path: Valid Bearer Token",
			authHeader:     "Bearer " + validToken,
			expectedStatus: http.StatusOK,
			checkContext:   true,
		},
		{
			name:           "Error Path: Missing Authorization Header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			checkContext:   false,
		},
		{
			name:           "Error Path: Invalid Format (No Bearer Prefix)",
			authHeader:     validToken,
			expectedStatus: http.StatusUnauthorized,
			checkContext:   false,
		},
		{
			name:           "Error Path: Expired Token",
			authHeader:     "Bearer " + expiredToken,
			expectedStatus: http.StatusUnauthorized,
			checkContext:   false,
		},
		{
			name:           "Error Path: Token Signed with Wrong Secret",
			authHeader:     "Bearer " + invalidSecretToken,
			expectedStatus: http.StatusUnauthorized,
			checkContext:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			var capturedUserID uint

			router.Use(middleware.AuthMiddleware(jwtSecret))
			router.GET("/test-protected", func(c *gin.Context) {
				id, err := middleware.GetUserIDFromContext(c)
				if err == nil {
					capturedUserID = id
				}
				c.Status(http.StatusOK)
			})

			req, _ := http.NewRequest(http.MethodGet, "/test-protected", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkContext {
				assert.Equal(t, uint(42), capturedUserID)
			}
		})
	}
}
