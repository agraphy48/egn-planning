package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"pulse-backend/internal/auth"
	"pulse-backend/internal/database"
	"pulse-backend/internal/middleware"
	"pulse-backend/internal/models"
)

// ListTasks returns tasks for a project, enforcing client privacy visibility rules
func ListTasks(c *gin.Context) {
	projectID := c.Query("project_id")
	if projectID == "" {
		projectID = c.GetString("project_id")
	}

	_, role := middleware.GetCurrentUser(c)

	query := database.DB.Preload("Assignee").Preload("Subtasks").Preload("Comments.User")
	if projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}

	// Client role restriction: Clients only see client-visible deliverables
	if role == auth.RoleClient {
		query = query.Where("client_visible = ?", true)
	}

	var tasks []models.Task
	if err := query.Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// GetTask returns a single task
func GetTask(c *gin.Context) {
	taskID := c.Param("id")
	_, role := middleware.GetCurrentUser(c)

	var task models.Task
	if err := database.DB.Preload("Assignee").Preload("Subtasks").Preload("Comments.User").Where("id = ?", taskID).First(&task).Error; err != nil {
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
	ClientVisible *bool  `json:"client_visible"`
}

// CreateTask creates a new task
func CreateTask(c *gin.Context) {
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
		ClientVisible: clientVisible,
	}

	if err := database.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	database.DB.Preload("Assignee").Preload("Subtasks").First(&task, "id = ?", task.ID)

	c.JSON(http.StatusCreated, task)
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// UpdateTaskStatus moves task between Kanban columns
func UpdateTaskStatus(c *gin.Context) {
	taskID := c.Param("id")

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

	task.Status = req.Status
	if err := database.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task status"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// UpdateTask updates general task fields
func UpdateTask(c *gin.Context) {
	taskID := c.Param("id")

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

	c.JSON(http.StatusOK, task)
}

// DeleteTask removes a task
func DeleteTask(c *gin.Context) {
	taskID := c.Param("id")

	if err := database.DB.Where("id = ?", taskID).Delete(&models.Task{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

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
	Content string `json:"content" binding:"required"`
}

// AddComment posts a new comment on a task
func AddComment(c *gin.Context) {
	taskID := c.Param("id")
	userID, _ := middleware.GetCurrentUser(c)

	var req AddCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	comment := models.Comment{
		ID:        "com_" + uuid.New().String()[:8],
		TaskID:    taskID,
		UserID:    userID,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to post comment"})
		return
	}

	database.DB.Preload("User").First(&comment, "id = ?", comment.ID)

	c.JSON(http.StatusCreated, comment)
}
