package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"pulse-backend/internal/database"
	"pulse-backend/internal/models"
)

// GetKPIStats calculates aggregated real-time workspace KPIs directly from database
func GetKPIStats(c *gin.Context) {
	var activeProjects int64
	database.DB.Model(&models.Project{}).Where("status = ?", "active").Count(&activeProjects)

	var totalTasks, completedTasks, overdueTasks int64
	database.DB.Model(&models.Task{}).Count(&totalTasks)
	database.DB.Model(&models.Task{}).Where("status = ?", "Done").Count(&completedTasks)

	nowStr := time.Now().Format("2006-01-02")
	database.DB.Model(&models.Task{}).Where("due_date < ? AND status != ?", nowStr, "Done").Count(&overdueTasks)

	var tasksDueThisWeek int64
	nextWeekStr := time.Now().AddDate(0, 0, 7).Format("2006-01-02")
	database.DB.Model(&models.Task{}).Where("due_date >= ? AND due_date <= ? AND status != ?", nowStr, nextWeekStr, "Done").Count(&tasksDueThisWeek)

	completionRate := 0
	if totalTasks > 0 {
		completionRate = int((completedTasks * 100) / totalTasks)
	}

	c.JSON(http.StatusOK, gin.H{
		"active_projects":       activeProjects,
		"tasks_due_this_week":   tasksDueThisWeek,
		"overdue_tasks":         overdueTasks,
		"completion_rate":       completionRate,
		"active_projects_delta": fmt.Sprintf("+%d this month", activeProjects),
		"total_tasks":           totalTasks,
		"completed_tasks":       completedTasks,
	})
}

// GetActivityFeed returns recent audit & collaboration logs
func GetActivityFeed(c *gin.Context) {
	projectID := c.Query("project_id")

	var logs []models.ActivityLog
	query := database.DB.Order("created_at desc").Limit(25)
	if projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}

	if err := query.Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activities"})
		return
	}

	// Format for frontend
	type FeedItem struct {
		ID        string `json:"id"`
		User      string `json:"user"`
		AvatarURL string `json:"avatar_url"`
		Action    string `json:"action"`
		Target    string `json:"target"`
		Project   string `json:"project"`
		TimeAgo   string `json:"time_ago"`
	}

	var items []FeedItem
	for _, l := range logs {
		items = append(items, FeedItem{
			ID:        l.ID,
			User:      l.UserName,
			AvatarURL: l.UserAvatar,
			Action:    l.Action,
			Target:    l.Target,
			Project:   l.ProjectID,
			TimeAgo:   l.CreatedAt.Format("15:04 (02 Jan)"),
		})
	}

	c.JSON(http.StatusOK, items)
}
