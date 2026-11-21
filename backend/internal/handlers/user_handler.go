package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"pulse-backend/internal/auth"
	"pulse-backend/internal/database"
	"pulse-backend/internal/middleware"
	"pulse-backend/internal/models"
)

type CreateUserRequest struct {
	Name      string `json:"name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	Role      string `json:"role" binding:"required"`
	AvatarURL string `json:"avatar_url"`
}

type UpdateUserRequest struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"` // optional: if not empty, hash & update
	Role      string `json:"role"`
	AvatarURL string `json:"avatar_url"`
}

type UpdateRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

// ListUsers returns list of system users
func ListUsers(c *gin.Context) {
	var users []models.User
	if err := database.DB.Order("created_at asc").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUser returns a single user profile by ID
func GetUser(c *gin.Context) {
	targetUserID := c.Param("id")
	if targetUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user ID parameter"})
		return
	}

	var user models.User
	if err := database.DB.Where("id = ?", targetUserID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser registers a new system user with hashed password and role
func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	cleanEmail := strings.TrimSpace(strings.ToLower(req.Email))
	cleanName := strings.TrimSpace(req.Name)

	if !auth.IsValidRole(req.Role) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid role specified. Allowed: super_admin, pm, member, client",
		})
		return
	}

	// Check email uniqueness
	var existing models.User
	if err := database.DB.Where("LOWER(email) = ?", cleanEmail).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email is already registered by another user"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
		return
	}

	// Default avatar if none provided
	avatar := strings.TrimSpace(req.AvatarURL)
	if avatar == "" {
		avatar = fmt.Sprintf("https://api.dicebear.com/7.x/avataaars/svg?seed=%s", url.QueryEscape(cleanName))
	}

	newID := fmt.Sprintf("usr_%s", strings.ReplaceAll(uuid.New().String()[:12], "-", ""))

	now := time.Now()
	user := models.User{
		ID:               newID,
		Name:             cleanName,
		Email:            cleanEmail,
		Password:         string(hashedPassword),
		Role:             req.Role,
		AvatarURL:        avatar,
		CurrentSessionID: "",
		LastActiveAt:     now,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user in database: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// UpdateUser modifies user details (name, email, role, avatar, and optional password)
func UpdateUser(c *gin.Context) {
	targetUserID := c.Param("id")
	if targetUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user ID parameter"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("id = ?", targetUserID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Safeguard: Protect primary super admin account
	if (user.ID == "usr_super_admin" || user.Email == "admin@pulse.workspace") && req.Role != "" && req.Role != auth.RoleSuperAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Peran akun Root Super Admin tidak dapat diubah demi keamanan sistem"})
		return
	}

	// Email update with uniqueness check
	if req.Email != "" {
		cleanEmail := strings.TrimSpace(strings.ToLower(req.Email))
		if cleanEmail != strings.ToLower(user.Email) {
			var check models.User
			if err := database.DB.Where("LOWER(email) = ? AND id != ?", cleanEmail, user.ID).First(&check).Error; err == nil {
				c.JSON(http.StatusConflict, gin.H{"error": "Email tersebut sudah digunakan oleh akun lain"})
				return
			}
			user.Email = cleanEmail
		}
	}

	if req.Name != "" {
		user.Name = strings.TrimSpace(req.Name)
	}

	if req.Role != "" {
		if !auth.IsValidRole(req.Role) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Peran tidak valid"})
			return
		}
		user.Role = req.Role
	}

	if req.AvatarURL != "" {
		user.AvatarURL = strings.TrimSpace(req.AvatarURL)
	}

	// Optional password update
	if req.Password != "" {
		if len(req.Password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Kata sandi minimal 6 karakter"})
			return
		}
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi kata sandi"})
			return
		}
		user.Password = string(hashed)
		// Revoke current session so user must log in with new password
		user.CurrentSessionID = ""
		auth.ClearActiveSession(user.ID)
	}

	user.UpdatedAt = time.Now()
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data pengguna: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser deletes a user from the system with safety checks
func DeleteUser(c *gin.Context) {
	targetUserID := c.Param("id")
	if targetUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user ID parameter"})
		return
	}

	currentUserID, _ := middleware.GetCurrentUser(c)

	// Safeguard 1: Prevent user from deleting own account
	if currentUserID == targetUserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Anda tidak dapat menghapus akun Anda sendiri saat sedang masuk"})
		return
	}

	// Safeguard 2: Prevent deletion of Root Super Admin
	if targetUserID == "usr_super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akun Root Super Admin dilindungi dan tidak dapat dihapus"})
		return
	}

	var user models.User
	if err := database.DB.Where("id = ?", targetUserID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	if user.Email == "admin@pulse.workspace" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akun Root Super Admin dilindungi dan tidak dapat dihapus"})
		return
	}

	// Revoke active session
	auth.ClearActiveSession(targetUserID)

	// Remove from project members pivot table
	database.DB.Where("user_id = ?", targetUserID).Delete(&models.ProjectMember{})

	// Delete user record
	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus pengguna dari database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pengguna berhasil dihapus dari sistem",
		"id":      targetUserID,
	})
}

// ResetUserSession forces a user session termination / device disconnect
func ResetUserSession(c *gin.Context) {
	targetUserID := c.Param("id")
	if targetUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user ID parameter"})
		return
	}

	var user models.User
	if err := database.DB.Where("id = ?", targetUserID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	// Clear session in memory & database
	auth.ClearActiveSession(targetUserID)
	database.DB.Model(&user).Update("current_session_id", "")

	c.JSON(http.StatusOK, gin.H{
		"message": "Sesi pengguna berhasil dilepas. Semua perangkat terkait telah logout paksa.",
		"user_id": targetUserID,
	})
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

	if (user.ID == "usr_super_admin" || user.Email == "admin@pulse.workspace") && req.Role != auth.RoleSuperAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Peran akun Root Super Admin tidak dapat diubah"})
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
