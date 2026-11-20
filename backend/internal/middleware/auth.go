package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"pulse-backend/internal/auth"
)

const (
	ContextUserID = "user_id"
	ContextEmail  = "email"
	ContextRole   = "role"
)

// AuthMiddleware validates JWT Bearer tokens and injects user context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized: missing Authorization header",
				"code":    "AUTH_HEADER_MISSING",
				"status":  http.StatusUnauthorized,
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized: invalid Authorization header format, expected 'Bearer <token>'",
				"code":    "AUTH_HEADER_MALFORMED",
				"status":  http.StatusUnauthorized,
			})
			return
		}

		claims, err := auth.ValidateToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized: " + err.Error(),
				"code":    "AUTH_TOKEN_INVALID",
				"status":  http.StatusUnauthorized,
			})
			return
		}

		// Inject into Gin context
		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextEmail, claims.Email)
		c.Set(ContextRole, claims.Role)

		c.Next()
	}
}

// GetCurrentUser returns the user_id and role from Gin context
func GetCurrentUser(c *gin.Context) (userID string, role string) {
	uVal, _ := c.Get(ContextUserID)
	rVal, _ := c.Get(ContextRole)

	userID, _ = uVal.(string)
	role, _ = rVal.(string)
	return
}
