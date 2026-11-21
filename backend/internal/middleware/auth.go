package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"pulse-backend/internal/auth"
	"pulse-backend/internal/database"
	"pulse-backend/internal/models"
)

const (
	ContextUserID    = "user_id"
	ContextEmail     = "email"
	ContextRole      = "role"
	ContextSessionID = "session_id"
)

// AuthMiddleware validates JWT Bearer tokens, enforces single active session & idle timeout
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

		// 1. Single Active Session Enforcement (Single Device / Concurrent Login Prevention)
		if claims.SessionID != "" {
			activeSessionID, ok := auth.GetActiveSession(claims.UserID)
			if !ok {
				// Not cached in memory yet, load from DB
				var dbUser models.User
				if err := database.DB.Select("current_session_id", "last_active_at").Where("id = ?", claims.UserID).First(&dbUser).Error; err == nil {
					activeSessionID = dbUser.CurrentSessionID
					auth.SetActiveSession(claims.UserID, dbUser.CurrentSessionID)
				}
			}

			if activeSessionID != "" && claims.SessionID != activeSessionID {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error":   "Akun Anda telah dibuka di perangkat atau jendela lain. Sesi ini telah dihentikan demi keamanan.",
					"code":    "SESSION_CONCURRENT_KICKED",
					"status":  http.StatusUnauthorized,
				})
				return
			}
		}

		// 2. Server-side Inactivity Idle Timeout (15 minutes of zero activity)
		lastActive := auth.GetLastActive(claims.UserID)
		if !lastActive.IsZero() && time.Since(lastActive) > 15*time.Minute {
			_ = auth.RevokeToken(database.DB, parts[1], claims.UserID, time.Now().Add(24*time.Hour))
			auth.ClearActiveSession(claims.UserID)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Sesi Anda telah berakhir otomatis karena tidak ada aktivitas selama 15 menit.",
				"code":    "SESSION_IDLE_TIMEOUT",
				"status":  http.StatusUnauthorized,
			})
			return
		}

		// Update activity timestamp in memory
		auth.UpdateLastActive(claims.UserID)

		// Inject into Gin context
		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextEmail, claims.Email)
		c.Set(ContextRole, claims.Role)
		c.Set(ContextSessionID, claims.SessionID)

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
