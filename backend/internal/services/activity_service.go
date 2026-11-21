package services

import (
	"time"

	"github.com/google/uuid"

	"pulse-backend/internal/database"
	"pulse-backend/internal/models"
)

// LogActivity records an immutable audit log entry in database
func LogActivity(projectID, userID, userName, userAvatar, action, target, details string) {
	if database.DB == nil {
		return
	}

	entry := models.ActivityLog{
		ID:         "act_" + uuid.New().String()[:8],
		ProjectID:  projectID,
		UserID:     userID,
		UserName:   userName,
		UserAvatar: userAvatar,
		Action:     action,
		Target:     target,
		Details:    details,
		CreatedAt:  time.Now(),
	}

	database.DB.Create(&entry)
}

// CreateNotification sends an in-app notification to a user
func CreateNotification(userID, title, message, notifType, link string) {
	if database.DB == nil || userID == "" {
		return
	}

	if notifType == "" {
		notifType = "info"
	}

	notif := models.Notification{
		ID:        "notif_" + uuid.New().String()[:8],
		UserID:    userID,
		Title:     title,
		Message:   message,
		Type:      notifType,
		Link:      link,
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	database.DB.Create(&notif)
}
