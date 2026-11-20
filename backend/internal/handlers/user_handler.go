package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"pulse-backend/internal/auth"
	"pulse-backend/internal/database"
	"pulse-backend/internal/models"
)

type UpdateRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

// UpdateUserRole updates a user's global role (restricted to super_admin)
func UpdateUserRole(c *gin.Context) {
	targetUserID := c.Param("id")
	if targetUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user ID parameter"})
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	if !auth.IsValidRole(req.Role) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid role specified. Allowed: super_admin, pm, member, client",
		})
		return
	}

	var user models.User
	if err := database.DB.Where("id = ?", targetUserID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Role = req.Role
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "User role updated successfully",
		"user_id":     user.ID,
		"new_role":    user.Role,
		"permissions": auth.GetPermissionsForRole(user.Role),
	})
}

// ListUsers returns list of system users
func ListUsers(c *gin.Context) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}
