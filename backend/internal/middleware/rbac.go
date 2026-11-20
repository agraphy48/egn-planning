package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"pulse-backend/internal/auth"
	"pulse-backend/internal/database"
	"pulse-backend/internal/models"
)

// RequireRole checks whether the user's role matches any of the allowedRoles
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, role := GetCurrentUser(c)
		if role == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":  "Unauthorized: role not found in context",
				"status": http.StatusUnauthorized,
			})
			return
		}

		// super_admin always bypasses global role checks
		if role == auth.RoleSuperAdmin {
			c.Next()
			return
		}

		for _, allowed := range allowedRoles {
			if role == allowed {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error":         "Forbidden: insufficient role privileges",
			"current_role":  role,
			"allowed_roles": allowedRoles,
			"status":        http.StatusForbidden,
		})
	}
}

// RequirePermission checks whether the user's role has the requested granular permission
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, role := GetCurrentUser(c)
		if role == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":  "Unauthorized: role not found in context",
				"status": http.StatusUnauthorized,
			})
			return
		}

		if !auth.HasPermission(role, permission) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":               "Forbidden: missing required permission",
				"required_permission": permission,
				"current_role":        role,
				"status":              http.StatusForbidden,
			})
			return
		}

		c.Next()
	}
}

// RequireProjectAccess ensures user has membership in the requested project (bypass for super_admin)
func RequireProjectAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, role := GetCurrentUser(c)
		if userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":  "Unauthorized: user context missing",
				"status": http.StatusUnauthorized,
			})
			return
		}

		// super_admin has global access to all projects
		if role == auth.RoleSuperAdmin {
			c.Next()
			return
		}

		// 1. Resolve project_id from URL params or query
		projectID := c.Param("project_id")
		if projectID == "" {
			// Check if :id is the project ID (e.g. /projects/:id or /projects/:id/members)
			projectID = c.Param("id")
		}
		if projectID == "" {
			projectID = c.Query("project_id")
		}

		// If still empty and path contains /tasks/:id or /milestones/:id, lookup parent project
		if projectID == "" {
			taskID := c.Param("task_id")
			if taskID == "" && len(c.Param("id")) > 0 {
				taskID = c.Param("id")
			}
			if taskID != "" {
				var task models.Task
				if err := database.DB.Select("project_id").Where("id = ?", taskID).First(&task).Error; err == nil {
					projectID = task.ProjectID
				}
			}
		}

		if projectID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":  "Bad Request: project_id could not be resolved from request",
				"status": http.StatusBadRequest,
			})
			return
		}

		// 2. Query project_members pivot table
		var member models.ProjectMember
		err := database.DB.Where("project_id = ? AND user_id = ?", projectID, userID).First(&member).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":      "Forbidden: you are not an assigned member of this project",
				"project_id": projectID,
				"user_id":    userID,
				"status":     http.StatusForbidden,
			})
			return
		}

		// Store project membership context for downstream handlers
		c.Set("project_id", projectID)
		c.Set("project_role", member.ProjectRole)

		c.Next()
	}
}
