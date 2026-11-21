package handlers

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"pulse-backend/internal/auth"
	"pulse-backend/internal/database"
	"pulse-backend/internal/middleware"
	"pulse-backend/internal/models"
	"pulse-backend/internal/services"
)

// ListTasks returns tasks for a project, strictly enforcing client privacy visibility rules
func ListTasks(c *gin.Context) {
	projectID := c.Query("project_id")
	if projectID == "" {
		projectID = c.GetString("project_id")
	}

	_, role := middleware.GetCurrentUser(c)

	query := database.DB.Preload("Assignee").Preload("Subtasks").Preload("Attachments")
	if projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}

	// CRITICAL SECURITY FIX: Filter both tasks and comments for Client role
	if role == auth.RoleClient {
		query = query.Where("client_visible = ?", true).
			Preload("Comments", "client_visible = ?", true).
			Preload("Comments.User")
	} else {
		query = query.Preload("Comments.User")
	}

	var tasks []models.Task
	if err := query.Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// GetTask returns a single task with client data isolation
func GetTask(c *gin.Context) {
	taskID := c.Param("id")
	_, role := middleware.GetCurrentUser(c)

	query := database.DB.Preload("Assignee").Preload("Subtasks").Preload("Attachments")
	if role == auth.RoleClient {
		query = query.Preload("Comments", "client_visible = ?", true).Preload("Comments.User")
	} else {
		query = query.Preload("Comments.User")
	}

	var task models.Task
	if err := query.Where("id = ?", taskID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Prevent client from viewing internal-only tasks
	if role == auth.RoleClient && !task.ClientVisible {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: task is marked internal only"})
		return
	}

	c.JSON(http.StatusOK, task)
}

type CreateTaskRequest struct {
	ProjectID     string `json:"project_id" binding:"required"`
	Title         string `json:"title" binding:"required"`
	Description   string `json:"description"`
	Category      string `json:"category"`
	Priority      string `json:"priority"`
	Status        string `json:"status"`
	AssigneeID    string `json:"assignee_id"`
	DueDate       string `json:"due_date"`
	StartDate     string `json:"start_date"`
	TimeEstimate  string `json:"time_estimate"`
	Fase          string `json:"fase"`
	EpicModul     string `json:"epic_modul"`
	Deliverable   string `json:"deliverable"`
	ClientVisible *bool  `json:"client_visible"`
}

// CreateTask creates a new task
func CreateTask(c *gin.Context) {
	userID, _ := middleware.GetCurrentUser(c)
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	taskID := "task_" + uuid.New().String()[:8]
	key := fmt.Sprintf("#EGN-%d", 100+rand.Intn(900))

	category := req.Category
	if category == "" {
		category = "General"
	}
	priority := req.Priority
	if priority == "" {
		priority = "Medium"
	}
	status := req.Status
	if status == "" {
		status = "Backlog"
	}

	clientVisible := true
	if req.ClientVisible != nil {
		clientVisible = *req.ClientVisible
	}

	task := models.Task{
		ID:            taskID,
		ProjectID:     req.ProjectID,
		Key:           key,
		Title:         req.Title,
		Description:   req.Description,
		Category:      category,
		Priority:      priority,
		Status:        status,
		AssigneeID:    req.AssigneeID,
		DueDate:       req.DueDate,
		StartDate:     req.StartDate,
		TimeEstimate:  req.TimeEstimate,
		Fase:          req.Fase,
		EpicModul:     req.EpicModul,
		Deliverable:   req.Deliverable,
		ClientVisible: clientVisible,
	}

	if err := database.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	database.DB.Preload("Assignee").Preload("Subtasks").Preload("Attachments").First(&task, "id = ?", task.ID)

	// Audit Log
	var author models.User
	database.DB.Where("id = ?", userID).First(&author)
	services.LogActivity(req.ProjectID, userID, author.Name, author.AvatarURL, "created task", task.Key+" - "+task.Title, "Status: "+status)

	c.JSON(http.StatusCreated, task)
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// UpdateTaskStatus moves task between Kanban columns
func UpdateTaskStatus(c *gin.Context) {
	taskID := c.Param("id")
	userID, _ := middleware.GetCurrentUser(c)

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	var task models.Task
	if err := database.DB.Where("id = ?", taskID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	oldStatus := task.Status
	task.Status = req.Status
	if err := database.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task status"})
		return
	}

	// Audit Log
	var author models.User
	database.DB.Where("id = ?", userID).First(&author)
	services.LogActivity(task.ProjectID, userID, author.Name, author.AvatarURL, "moved task", task.Key, fmt.Sprintf("%s -> %s", oldStatus, req.Status))

	c.JSON(http.StatusOK, task)
}

// UpdateTask updates general task fields
func UpdateTask(c *gin.Context) {
	taskID := c.Param("id")
	userID, _ := middleware.GetCurrentUser(c)

	var task models.Task
	if err := database.DB.Where("id = ?", taskID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := database.DB.Model(&task).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	database.DB.Preload("Assignee").Preload("Subtasks").Preload("Comments.User").Preload("Attachments").First(&task, "id = ?", taskID)

	var author models.User
	database.DB.Where("id = ?", userID).First(&author)
	services.LogActivity(task.ProjectID, userID, author.Name, author.AvatarURL, "updated task", task.Key, "Task attributes updated")

	c.JSON(http.StatusOK, task)
}

// DeleteTask removes a task
func DeleteTask(c *gin.Context) {
	taskID := c.Param("id")
	userID, _ := middleware.GetCurrentUser(c)

	var task models.Task
	if err := database.DB.Where("id = ?", taskID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if err := database.DB.Where("id = ?", taskID).Delete(&models.Task{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	var author models.User
	database.DB.Where("id = ?", userID).First(&author)
	services.LogActivity(task.ProjectID, userID, author.Name, author.AvatarURL, "deleted task", task.Key, task.Title)

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully", "task_id": taskID})
}

type AddSubtaskRequest struct {
	Title string `json:"title" binding:"required"`
}

// AddSubtask creates a subtask item
func AddSubtask(c *gin.Context) {
	taskID := c.Param("id")

	var req AddSubtaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	subtask := models.SubTask{
		ID:        "sub_" + uuid.New().String()[:8],
		TaskID:    taskID,
		Title:     req.Title,
		IsDone:    false,
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&subtask).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subtask"})
		return
	}

	c.JSON(http.StatusCreated, subtask)
}

type ToggleSubtaskRequest struct {
	IsDone bool `json:"is_done"`
}

// ToggleSubtask updates subtask checkmark
func ToggleSubtask(c *gin.Context) {
	subtaskID := c.Param("sub_id")

	var req ToggleSubtaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var subtask models.SubTask
	if err := database.DB.Where("id = ?", subtaskID).First(&subtask).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subtask not found"})
		return
	}

	subtask.IsDone = req.IsDone
	database.DB.Save(&subtask)

	c.JSON(http.StatusOK, subtask)
}

type AddCommentRequest struct {
	Content       string `json:"content" binding:"required"`
	ClientVisible *bool  `json:"client_visible"`
}

// AddComment posts a new comment on a task with client visibility isolation
func AddComment(c *gin.Context) {
	taskID := c.Param("id")
	userID, role := middleware.GetCurrentUser(c)

	var req AddCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	clientVisible := true
	if req.ClientVisible != nil {
		clientVisible = *req.ClientVisible
	}
	// If a client posts, it's always client visible
	if role == auth.RoleClient {
		clientVisible = true
	}

	comment := models.Comment{
		ID:            "com_" + uuid.New().String()[:8],
		TaskID:        taskID,
		UserID:        userID,
		Content:       req.Content,
		ClientVisible: clientVisible,
		CreatedAt:     time.Now(),
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to post comment"})
		return
	}

	database.DB.Preload("User").First(&comment, "id = ?", comment.ID)

	c.JSON(http.StatusCreated, comment)
}

// UploadAttachment uploads a file attachment for a task
func UploadAttachment(c *gin.Context) {
	taskID := c.Param("id")
	userID, _ := middleware.GetCurrentUser(c)

	var task models.Task
	if err := database.DB.Where("id = ?", taskID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing file upload: " + err.Error()})
		return
	}
	defer file.Close()

	// 1. Strict File Size Validation (Max 10 MB)
	const maxFileSize = 10 * 1024 * 1024 // 10 MB
	if header.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran file melebihi batas maksimum 10 MB"})
		return
	}

	// 2. Strict Extension Whitelist (OWASP ASVS 12.1.1 - Insecure File Upload Prevention)
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExtensions := map[string]bool{
		".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".webp": true,
		".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true,
		".csv": true, ".txt": true, ".zip": true,
	}
	if !allowedExtensions[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipe file tidak diizinkan. Format yang didukung: Dokumen (PDF, Word, Excel, CSV, TXT), Gambar (PNG, JPG, WEBP), atau Arsip ZIP."})
		return
	}

	cleanFileName := filepath.Base(filepath.Clean(header.Filename))

	uploadDir := filepath.Join(".", "uploads")
	_ = os.MkdirAll(uploadDir, os.ModePerm)

	safeFileName := fmt.Sprintf("%s_%d%s", uuid.New().String()[:8], time.Now().Unix(), ext)
	dstPath := filepath.Join(uploadDir, safeFileName)

	out, err := os.Create(dstPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	defer out.Close()

	if _, err = io.Copy(out, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write file"})
		return
	}

	sizeMB := fmt.Sprintf("%.1f MB", float64(header.Size)/(1024*1024))
	if header.Size < 1024*1024 {
		sizeMB = fmt.Sprintf("%.0f KB", float64(header.Size)/1024)
	}

	att := models.Attachment{
		ID:         "att_" + uuid.New().String()[:8],
		TaskID:     taskID,
		FileName:   cleanFileName,
		FileSize:   sizeMB,
		FileType:   ext,
		FileURL:    "/uploads/" + safeFileName,
		UploadedBy: userID,
		CreatedAt:  time.Now(),
	}

	if err := database.DB.Create(&att).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to persist attachment"})
		return
	}

	var author models.User
	database.DB.Where("id = ?", userID).First(&author)
	services.LogActivity(task.ProjectID, userID, author.Name, author.AvatarURL, "uploaded attachment", att.FileName, task.Key)

	c.JSON(http.StatusCreated, att)
}

// DeleteAttachment removes a file attachment
func DeleteAttachment(c *gin.Context) {
	attID := c.Param("att_id")

	var att models.Attachment
	if err := database.DB.Where("id = ?", attID).First(&att).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attachment not found"})
		return
	}

	database.DB.Delete(&att)
	c.JSON(http.StatusOK, gin.H{"message": "Attachment deleted", "id": attID})
}
