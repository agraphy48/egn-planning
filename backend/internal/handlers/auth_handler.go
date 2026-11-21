package handlers

import (
	"net/http"
	"os"
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

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token       string       `json:"token"`
	User        *models.User `json:"user"`
	Permissions []string     `json:"permissions"`
}

// Login handles user authentication and JWT token issuance
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate unique session ID for single-device active session enforcement
	sessionID := uuid.New().String()
	token, err := auth.GenerateToken(user.ID, user.Email, user.Role, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Register single active session in memory cache and database
	auth.SetActiveSession(user.ID, sessionID)
	now := time.Now()
	database.DB.Model(&user).Updates(map[string]interface{}{
		"current_session_id": sessionID,
		"last_active_at":     now,
	})

	user.CurrentSessionID = sessionID
	user.LastActiveAt = now

	permissions := auth.GetPermissionsForRole(user.Role)

	c.JSON(http.StatusOK, LoginResponse{
		Token:       token,
		User:        &user,
		Permissions: permissions,
	})
}

// Logout revokes the currently used JWT token server-side and invalidates the session
func Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			tokenString := parts[1]
			claims, err := auth.ValidateToken(tokenString)
			var exp time.Time
			var userID string
			if err == nil && claims != nil {
				userID = claims.UserID
				if claims.ExpiresAt != nil {
					exp = claims.ExpiresAt.Time
				} else {
					exp = time.Now().Add(7 * 24 * time.Hour)
				}
			} else {
				exp = time.Now().Add(7 * 24 * time.Hour)
			}
			_ = auth.RevokeToken(database.DB, tokenString, userID, exp)
			if userID != "" {
				auth.ClearActiveSession(userID)
				database.DB.Model(&models.User{}).Where("id = ?", userID).Update("current_session_id", "")
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
		"status":  "success",
	})
}

// GetMe returns current authenticated user profile + role + flat array of permissions
func GetMe(c *gin.Context) {
	userID, role := middleware.GetCurrentUser(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User profile not found"})
		return
	}

	permissions := auth.GetPermissionsForRole(role)

	c.JSON(http.StatusOK, gin.H{
		"id":          user.ID,
		"name":        user.Name,
		"email":       user.Email,
		"role":        user.Role,
		"avatar_url":  user.AvatarURL,
		"permissions": permissions,
		"created_at":  user.CreatedAt,
	})
}

// GetAuthConfig returns environment configuration regarding dev shortcuts
func GetAuthConfig(c *gin.Context) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		if gin.Mode() == gin.ReleaseMode {
			env = "production"
		} else {
			env = "development"
		}
	}

	allowQuickLogin := env != "production"

	c.JSON(http.StatusOK, gin.H{
		"allow_quick_login": allowQuickLogin,
		"environment":       env,
	})
}

// Heartbeat allows frontend to ping server, update activity timestamp and verify session freshness
func Heartbeat(c *gin.Context) {
	userID, role := middleware.GetCurrentUser(c)
	c.JSON(http.StatusOK, gin.H{
		"status":  "alive",
		"user_id": userID,
		"role":    role,
		"time":    time.Now().Format(time.RFC3339),
	})
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

// ChangePassword allows an authenticated user to securely update their own password
func ChangePassword(c *gin.Context) {
	userID, _ := middleware.GetCurrentUser(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi autentikasi tidak valid atau telah berakhir"})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kata sandi baru minimal harus 6 karakter"})
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Konfirmasi kata sandi baru tidak cocok"})
		return
	}

	if req.CurrentPassword == req.NewPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kata sandi baru tidak boleh sama dengan kata sandi saat ini"})
		return
	}

	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Akun pengguna tidak ditemukan"})
		return
	}

	// Verify current password against stored bcrypt hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kata sandi saat ini yang Anda masukkan salah"})
		return
	}

	// Hash new password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi kata sandi baru"})
		return
	}

	// Update password in database
	now := time.Now()
	if err := database.DB.Model(&user).Updates(map[string]interface{}{
		"password":   string(hashedPassword),
		"updated_at": now,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan kata sandi baru ke basis data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Kata sandi berhasil diperbarui. Silakan gunakan kata sandi baru untuk login berikutnya.",
		"status":  "success",
	})
}
