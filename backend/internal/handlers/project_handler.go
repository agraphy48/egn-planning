package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"pulse-backend/internal/auth"
	"pulse-backend/internal/database"
	"pulse-backend/internal/middleware"
	"pulse-backend/internal/models"
)

// ListProjects returns projects accessible to the user
func ListProjects(c *gin.Context) {
	userID, role := middleware.GetCurrentUser(c)

	var projects []models.Project
	if role == auth.RoleSuperAdmin {
		if err := database.DB.Preload("Members.User").Find(&projects).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects"})
			return
		}
	} else {
		// Filter by project_members pivot
		subQuery := database.DB.Table("project_members").Select("project_id").Where("user_id = ?", userID)
		if err := database.DB.Where("id IN (?)", subQuery).Preload("Members.User").Find(&projects).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects"})
			return
		}
	}

	// Calculate task counts
	for i := range projects {
		var total, completed int64
		database.DB.Model(&models.Task{}).Where("project_id = ?", projects[i].ID).Count(&total)
		database.DB.Model(&models.Task{}).Where("project_id = ? AND status = ?", projects[i].ID, "Done").Count(&completed)
		projects[i].TotalTasks = int(total)
		projects[i].CompletedTasks = int(completed)
		if total > 0 {
			projects[i].ProgressPercentage = int((completed * 100) / total)
		}
	}

	c.JSON(http.StatusOK, projects)
}

// GetProject returns a single project
func GetProject(c *gin.Context) {
	projectID := c.Param("id")

	var project models.Project
	if err := database.DB.Preload("Members.User").Preload("Milestones").Where("id = ?", projectID).First(&project).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	var total, completed int64
	database.DB.Model(&models.Task{}).Where("project_id = ?", project.ID).Count(&total)
	database.DB.Model(&models.Task{}).Where("project_id = ? AND status = ?", project.ID, "Done").Count(&completed)
	project.TotalTasks = int(total)
	project.CompletedTasks = int(completed)
	if total > 0 {
		project.ProgressPercentage = int((completed * 100) / total)
	}

	c.JSON(http.StatusOK, project)
}

type CreateProjectRequest struct {
	Name       string `json:"name" binding:"required"`
	Slug       string `json:"slug"`
	ClientName string `json:"client_name"`
	ColorDot   string `json:"color_dot"`
	Deadline   string `json:"deadline"`
}

// CreateProject creates a new project and attaches creator as PM member
func CreateProject(c *gin.Context) {
	userID, _ := middleware.GetCurrentUser(c)

	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	projectID := "proj_" + uuid.New().String()[:8]
	slug := req.Slug
	if slug == "" {
		slug = projectID
	}

	colorDot := req.ColorDot
	if colorDot == "" {
		colorDot = "#4F46E5"
	}

	project := models.Project{
		ID:         projectID,
		Name:       req.Name,
		Slug:       slug,
		ClientName: req.ClientName,
		Status:     "active",
		ColorDot:   colorDot,
		Deadline:   req.Deadline,
		OwnerID:    userID,
	}

	if err := database.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	// Add creator as PM member in project_members
	member := models.ProjectMember{
		ProjectID:   project.ID,
		UserID:      userID,
		ProjectRole: auth.RolePM,
		JoinedAt:    time.Now(),
	}
	database.DB.Create(&member)

	c.JSON(http.StatusCreated, project)
}

// UpdateProject updates project metadata
func UpdateProject(c *gin.Context) {
	projectID := c.Param("id")

	var project models.Project
	if err := database.DB.Where("id = ?", projectID).First(&project).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := database.DB.Model(&project).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}

	c.JSON(http.StatusOK, project)
}

// DeleteProject removes a project (restricted to super_admin and PM)
func DeleteProject(c *gin.Context) {
	projectID := c.Param("id")

	if err := database.DB.Where("id = ?", projectID).Delete(&models.Project{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	// Cascade delete project_members
	database.DB.Where("project_id = ?", projectID).Delete(&models.ProjectMember{})

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully", "project_id": projectID})
}

// ListProjectMembers returns all members assigned to the project
func ListProjectMembers(c *gin.Context) {
	projectID := c.Param("id")

	var members []models.ProjectMember
	if err := database.DB.Preload("User").Where("project_id = ?", projectID).Find(&members).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch project members"})
		return
	}

	c.JSON(http.StatusOK, members)
}

type AddMemberRequest struct {
	UserID      string `json:"user_id" binding:"required"`
	ProjectRole string `json:"project_role" binding:"required"` // pm, member, client
}

// AddProjectMember adds or invites a user to the project
func AddProjectMember(c *gin.Context) {
	projectID := c.Param("id")

	var req AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	if !auth.IsValidRole(req.ProjectRole) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project role"})
		return
	}

	// Verify target user exists
	var user models.User
	if err := database.DB.Where("id = ?", req.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User to add not found"})
		return
	}

	member := models.ProjectMember{
		ProjectID:   projectID,
		UserID:      req.UserID,
		ProjectRole: req.ProjectRole,
		JoinedAt:    time.Now(),
	}

	// Upsert / Create
	if err := database.DB.Where("project_id = ? AND user_id = ?", projectID, req.UserID).Assign(member).FirstOrCreate(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign project member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Member added to project successfully",
		"project_id":   projectID,
		"user_id":      req.UserID,
		"project_role": req.ProjectRole,
	})
}

// RemoveProjectMember removes a user from the project
func RemoveProjectMember(c *gin.Context) {
	projectID := c.Param("id")
	targetUserID := c.Param("user_id")

	result := database.DB.Where("project_id = ? AND user_id = ?", projectID, targetUserID).Delete(&models.ProjectMember{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove member"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found in this project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Member removed from project",
		"project_id": projectID,
		"user_id":    targetUserID,
	})
}
