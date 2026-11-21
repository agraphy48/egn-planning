package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"pulse-backend/internal/database"
	"pulse-backend/internal/middleware"
	"pulse-backend/internal/models"
)

// ListNotifications returns notifications for the authenticated user
func ListNotifications(c *gin.Context) {
	userID, _ := middleware.GetCurrentUser(c)

	var notifs []models.Notification
	if err := database.DB.Where("user_id = ?", userID).Order("created_at desc").Limit(30).Find(&notifs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	var unreadCount int64
	database.DB.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&unreadCount)

	c.JSON(http.StatusOK, gin.H{
		"notifications": notifs,
		"unread_count":  unreadCount,
	})
}

// MarkNotificationRead marks a single notification as read
func MarkNotificationRead(c *gin.Context) {
	notifID := c.Param("id")
	userID, _ := middleware.GetCurrentUser(c)

	database.DB.Model(&models.Notification{}).Where("id = ? AND user_id = ?", notifID, userID).Update("is_read", true)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// MarkAllNotificationsRead marks all notifications as read for current user
func MarkAllNotificationsRead(c *gin.Context) {
	userID, _ := middleware.GetCurrentUser(c)

	database.DB.Model(&models.Notification{}).Where("user_id = ?", userID).Update("is_read", true)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
