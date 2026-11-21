package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"pulse-backend/internal/database"
	"pulse-backend/internal/middleware"
	"pulse-backend/internal/models"
)

type LogTimeRequest struct {
	ProjectID   string  `json:"project_id" binding:"required"`
	TaskID      string  `json:"task_id"`
	Hours       float64 `json:"hours" binding:"required,gt=0"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

// LogTime records work hours
func LogTime(c *gin.Context) {
	userID, _ := middleware.GetCurrentUser(c)

	var req LogTimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	dateStr := req.Date
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	entry := models.TimeEntry{
		ID:          "time_" + uuid.New().String()[:8],
		ProjectID:   req.ProjectID,
		TaskID:      req.TaskID,
		UserID:      userID,
		Hours:       req.Hours,
		Description: req.Description,
		Date:        dateStr,
		CreatedAt:   time.Now(),
	}

	if err := database.DB.Create(&entry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log time"})
		return
	}

	database.DB.Preload("User").First(&entry, "id = ?", entry.ID)

	c.JSON(http.StatusCreated, entry)
}

// ListTimeEntries returns logged time entries
func ListTimeEntries(c *gin.Context) {
	projectID := c.Query("project_id")

	query := database.DB.Preload("User").Order("date desc")
	if projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}

	var entries []models.TimeEntry
	if err := query.Find(&entries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch timesheet entries"})
		return
	}

	c.JSON(http.StatusOK, entries)
}

// GetTeamWorkload aggregates workload hours per user
func GetTeamWorkload(c *gin.Context) {
	projectID := c.Query("project_id")

	type WorkloadResult struct {
		UserID     string  `json:"user_id"`
		UserName   string  `json:"user_name"`
		TotalHours float64 `json:"total_hours"`
		TasksCount int64   `json:"tasks_count"`
	}

	var results []WorkloadResult
	query := database.DB.Table("time_entries").
		Select("time_entries.user_id, users.name as user_name, SUM(time_entries.hours) as total_hours, count(distinct time_entries.task_id) as tasks_count").
		Joins("JOIN users ON users.id = time_entries.user_id").
		Group("time_entries.user_id, users.name")

	if projectID != "" {
		query = query.Where("time_entries.project_id = ?", projectID)
	}

	query.Scan(&results)

	c.JSON(http.StatusOK, results)
}
