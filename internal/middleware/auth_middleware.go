package middleware

import (
	"errors"
	"net/http"
	"strings"

	"task-management-api/internal/util"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "userID"

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header is required",
			})
			c.Abort()
			
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header format must be Bearer {token}",
			})
			c.Abort()

			return
		}

		tokenString := parts[1]
		userID, err := util.ParseToken(tokenString, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			c.Abort()

			return
		}

		c.Set(UserIDKey, userID)
		c.Next()
	}
}

func GetUserIDFromContext(c *gin.Context) (uint, error) {
	val, exists := c.Get(UserIDKey)
	if !exists {
		return 0, errors.New("user ID not found in context")
	}

	userID, ok := val.(uint)
	if !ok {
		return 0, errors.New("invalid user ID type in context")
	}

	return userID, nil
}
