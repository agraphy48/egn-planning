package handlers

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"pulse-backend/internal/database"
	"pulse-backend/internal/models"
)

type StatusWithUsage struct {
	models.StatusConfig
	UsageCount int64 `json:"usage_count"`
}

// ListStatuses returns all configured statuses with usage counters
func ListStatuses(c *gin.Context) {
	category := c.Query("category")

	query := database.DB.Model(&models.StatusConfig{})
	if category != "" {
		query = query.Where("category = ?", category)
	}

	var statuses []models.StatusConfig
	if err := query.Order("category ASC, position ASC, name ASC").Find(&statuses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch statuses: " + err.Error()})
		return
	}

	// Calculate live usage counts
	result := make([]StatusWithUsage, len(statuses))
	for i, s := range statuses {
		var count int64
		if s.Category == "task" {
			database.DB.Model(&models.Task{}).Where("status = ?", s.Name).Count(&count)
		} else if s.Category == "project" {
			database.DB.Model(&models.Project{}).Where("status = ?", s.Slug).Count(&count)
		} else if s.Category == "milestone" {
			database.DB.Model(&models.Milestone{}).Where("status = ?", s.Slug).Count(&count)
		}
		result[i] = StatusWithUsage{
			StatusConfig: s,
			UsageCount:   count,
		}
	}

	c.JSON(http.StatusOK, result)
}

// GetStatus returns a single status by ID
func GetStatus(c *gin.Context) {
	id := c.Param("id")

	var status models.StatusConfig
	if err := database.DB.Where("id = ? OR slug = ?", id, id).First(&status).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Status tidak ditemukan"})
		return
	}

	var count int64
	if status.Category == "task" {
		database.DB.Model(&models.Task{}).Where("status = ?", status.Name).Count(&count)
	} else if status.Category == "project" {
		database.DB.Model(&models.Project{}).Where("status = ?", status.Slug).Count(&count)
	}

	c.JSON(http.StatusOK, StatusWithUsage{
		StatusConfig: status,
		UsageCount:   count,
	})
}

type CreateStatusRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug"`
	Category    string `json:"category"`
	ColorDot    string `json:"color_dot"`
	BadgeBg     string `json:"badge_bg"`
	BadgeText   string `json:"badge_text"`
	Position    int    `json:"position"`
	IsCompleted bool   `json:"is_completed"`
	Description string `json:"description"`
}

// CreateStatus creates a new status config
func CreateStatus(c *gin.Context) {
	var req CreateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama status wajib diisi: " + err.Error()})
		return
	}

	cleanName := strings.TrimSpace(req.Name)
	if cleanName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama status tidak boleh kosong"})
		return
	}

	category := strings.ToLower(strings.TrimSpace(req.Category))
	if category == "" {
		category = "task"
	}
	if category != "task" && category != "project" && category != "milestone" {
		category = "task"
	}

	slug := strings.TrimSpace(strings.ToLower(req.Slug))
	if slug == "" {
		slugRegex := regexp.MustCompile(`[^a-z0-9\-]+`)
		slug = strings.ToLower(cleanName)
		slug = strings.ReplaceAll(slug, " ", "-")
		slug = slugRegex.ReplaceAllString(slug, "")
	}

	// Check if already exists in same category
	var existing models.StatusConfig
	if err := database.DB.Where("category = ? AND (name = ? OR slug = ?)", category, cleanName, slug).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Status dengan nama '%s' atau slug '%s' sudah terdaftar dalam kategori ini", cleanName, slug)})
		return
	}

	color := strings.TrimSpace(req.ColorDot)
	if color == "" {
		color = "#6366F1"
	}

	badgeBg := strings.TrimSpace(req.BadgeBg)
	if badgeBg == "" {
		badgeBg = "bg-indigo-50"
	}
	badgeText := strings.TrimSpace(req.BadgeText)
	if badgeText == "" {
		badgeText = "text-indigo-700"
	}

	// Calculate max position if not given
	pos := req.Position
	if pos <= 0 {
		var maxPos int
		row := database.DB.Model(&models.StatusConfig{}).Where("category = ?", category).Select("COALESCE(MAX(position), 0)").Row()
		_ = row.Scan(&maxPos)
		pos = maxPos + 1
	}

	id := fmt.Sprintf("st_%s_%s", category, strings.ReplaceAll(uuid.New().String()[:8], "-", ""))

	status := models.StatusConfig{
		ID:          id,
		Name:        cleanName,
		Slug:        slug,
		Category:    category,
		ColorDot:    color,
		BadgeBg:     badgeBg,
		BadgeText:   badgeText,
		Position:    pos,
		IsCompleted: req.IsCompleted,
		IsSystem:    false,
		Description: strings.TrimSpace(req.Description),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := database.DB.Create(&status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan status: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, StatusWithUsage{
		StatusConfig: status,
		UsageCount:   0,
	})
}

type UpdateStatusRequest struct {
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Category     string `json:"category"`
	ColorDot     string `json:"color_dot"`
	BadgeBg      string `json:"badge_bg"`
	BadgeText    string `json:"badge_text"`
	Position     *int   `json:"position"`
	IsCompleted  *bool  `json:"is_completed"`
	Description  string `json:"description"`
	MigrateTasks bool   `json:"migrate_tasks"`
}

// UpdateStatus modifies an existing status
func UpdateStatus(c *gin.Context) {
	id := c.Param("id")

	var status models.StatusConfig
	if err := database.DB.Where("id = ?", id).First(&status).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Status tidak ditemukan"})
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format payload tidak valid: " + err.Error()})
		return
	}

	oldName := status.Name
	if req.Name != "" {
		status.Name = strings.TrimSpace(req.Name)
	}
	if req.Slug != "" {
		status.Slug = strings.TrimSpace(strings.ToLower(req.Slug))
	}
	if req.ColorDot != "" {
		status.ColorDot = strings.TrimSpace(req.ColorDot)
	}
	if req.BadgeBg != "" {
		status.BadgeBg = strings.TrimSpace(req.BadgeBg)
	}
	if req.BadgeText != "" {
		status.BadgeText = strings.TrimSpace(req.BadgeText)
	}
	if req.Position != nil {
		status.Position = *req.Position
	}
	if req.IsCompleted != nil {
		status.IsCompleted = *req.IsCompleted
	}
	if req.Description != "" {
		status.Description = strings.TrimSpace(req.Description)
	}
	status.UpdatedAt = time.Now()

	if err := database.DB.Save(&status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui status: " + err.Error()})
		return
	}

	// If name changed and category is task, migrate task statuses
	var migratedCount int64
	if status.Category == "task" && oldName != status.Name && req.MigrateTasks {
		res := database.DB.Model(&models.Task{}).Where("status = ?", oldName).Update("status", status.Name)
		migratedCount = res.RowsAffected
	}

	var usageCount int64
	if status.Category == "task" {
		database.DB.Model(&models.Task{}).Where("status = ?", status.Name).Count(&usageCount)
	} else if status.Category == "project" {
		database.DB.Model(&models.Project{}).Where("status = ?", status.Slug).Count(&usageCount)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":         status,
		"usage_count":    usageCount,
		"migrated_tasks": migratedCount,
		"message":        "Status berhasil diperbarui",
	})
}

// DeleteStatus deletes a status config safely
func DeleteStatus(c *gin.Context) {
	id := c.Param("id")
	reassignTo := c.Query("reassign_to")

	var status models.StatusConfig
	if err := database.DB.Where("id = ?", id).First(&status).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Status tidak ditemukan"})
		return
	}

	// Check usage
	var usageCount int64
	if status.Category == "task" {
		database.DB.Model(&models.Task{}).Where("status = ?", status.Name).Count(&usageCount)
	} else if status.Category == "project" {
		database.DB.Model(&models.Project{}).Where("status = ?", status.Slug).Count(&usageCount)
	}

	if usageCount > 0 {
		if reassignTo == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":       fmt.Sprintf("Status ini masih digunakan oleh %d item. Tentukan parameter 'reassign_to' untuk memindahkan data sebelum menghapus.", usageCount),
				"usage_count": usageCount,
				"requires_reassign": true,
			})
			return
		}

		// Execute reassignment
		if status.Category == "task" {
			database.DB.Model(&models.Task{}).Where("status = ?", status.Name).Update("status", reassignTo)
		} else if status.Category == "project" {
			database.DB.Model(&models.Project{}).Where("status = ?", status.Slug).Update("status", reassignTo)
		}
	}

	// Hard delete from database so config is cleanly removed
	if err := database.DB.Unscoped().Delete(&status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus status: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Status '%s' berhasil dihapus.", status.Name),
		"reassigned_count": usageCount,
	})
}

type ReorderItem struct {
	ID       string `json:"id" binding:"required"`
	Position int    `json:"position"`
}

type ReorderRequest struct {
	Items []ReorderItem `json:"items" binding:"required"`
}

// ReorderStatuses updates positions of multiple statuses
func ReorderStatuses(c *gin.Context) {
	var req ReorderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload tidak valid: " + err.Error()})
		return
	}

	tx := database.DB.Begin()
	for _, item := range req.Items {
		if err := tx.Model(&models.StatusConfig{}).Where("id = ?", item.ID).Update("position", item.Position).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan urutan status: " + err.Error()})
			return
		}
	}
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "Urutan status berhasil diperbarui",
		"count":   len(req.Items),
	})
}
