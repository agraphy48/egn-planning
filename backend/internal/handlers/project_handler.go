package handlers

import (
	"fmt"
	"net/http"
	"strings"
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
	PMUserID   string `json:"pm_user_id"`
}

// CreateProject creates a new project and attaches creator and specified PM
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

	// Determine assigned Project Manager (PM)
	pmID := userID
	if req.PMUserID != "" {
		pmID = req.PMUserID
	}

	pmMember := models.ProjectMember{
		ProjectID:   project.ID,
		UserID:      pmID,
		ProjectRole: auth.RolePM,
		JoinedAt:    time.Now(),
	}
	database.DB.Create(&pmMember)

	// If the creator is different from the assigned PM, also ensure creator retains PM access
	if userID != "" && userID != pmID {
		creatorMember := models.ProjectMember{
			ProjectID:   project.ID,
			UserID:      userID,
			ProjectRole: auth.RolePM,
			JoinedAt:    time.Now(),
		}
		database.DB.Where(models.ProjectMember{ProjectID: project.ID, UserID: userID}).FirstOrCreate(&creatorMember)
	}

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

	// Upsert / Create project member
	var existing models.ProjectMember
	if err := database.DB.Where("project_id = ? AND user_id = ?", projectID, req.UserID).First(&existing).Error; err == nil {
		existing.ProjectRole = req.ProjectRole
		if err := database.DB.Save(&existing).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project member role"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":      "Peran anggota proyek berhasil diperbarui",
			"project_id":   projectID,
			"user_id":      req.UserID,
			"project_role": req.ProjectRole,
		})
		return
	}

	member := models.ProjectMember{
		ProjectID:   projectID,
		UserID:      req.UserID,
		ProjectRole: req.ProjectRole,
		JoinedAt:    time.Now(),
	}

	if err := database.DB.Create(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign project member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Anggota berhasil ditambahkan ke proyek",
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

type CopyProjectRequest struct {
	Name            string `json:"name" binding:"required"`
	ClientName      string `json:"client_name"`
	Deadline        string `json:"deadline"`
	ColorDot        string `json:"color_dot"`
	CopyMembers     bool   `json:"copy_members"`
	ResetTaskStatus bool   `json:"reset_task_status"`
	CopyWiki        bool   `json:"copy_wiki"`
}

// CopyProject duplicates an entire project with all its tasks, subtasks, milestones, and members
func CopyProject(c *gin.Context) {
	sourceID := c.Param("id")
	userID, _ := middleware.GetCurrentUser(c)

	var req CopyProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama proyek baru wajib diisi: " + err.Error()})
		return
	}

	cleanName := strings.TrimSpace(req.Name)
	if cleanName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama proyek baru tidak boleh kosong"})
		return
	}

	// 1. Fetch source project
	var srcProject models.Project
	if err := database.DB.Where("id = ?", sourceID).First(&srcProject).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proyek sumber tidak ditemukan"})
		return
	}

	// 2. Prepare new project fields
	newProjectID := "proj_" + uuid.New().String()[:8]
	slugBase := strings.ToLower(cleanName)
	slug := strings.ReplaceAll(slugBase, " ", "-")
	slug = strings.ReplaceAll(slug, "/", "-")
	if slug == "" {
		slug = newProjectID
	}

	// Ensure unique slug
	var existingSlugCount int64
	database.DB.Model(&models.Project{}).Where("slug = ?", slug).Count(&existingSlugCount)
	if existingSlugCount > 0 {
		slug = fmt.Sprintf("%s-%s", slug, uuid.New().String()[:4])
	}

	clientName := req.ClientName
	if clientName == "" {
		clientName = srcProject.ClientName
	}

	deadline := req.Deadline
	if deadline == "" {
		deadline = srcProject.Deadline
	}

	colorDot := req.ColorDot
	if colorDot == "" {
		colorDot = srcProject.ColorDot
	}

	now := time.Now()
	newProject := models.Project{
		ID:                 newProjectID,
		Name:               cleanName,
		Slug:               slug,
		ClientName:         clientName,
		Status:             "active",
		ProgressPercentage: 0,
		ColorDot:           colorDot,
		Deadline:           deadline,
		OwnerID:            userID,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := database.DB.Create(&newProject).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menduplikasi proyek: " + err.Error()})
		return
	}

	// 3. Project Members
	// Always attach creator as PM
	creatorMember := models.ProjectMember{
		ProjectID:   newProject.ID,
		UserID:      userID,
		ProjectRole: auth.RolePM,
		JoinedAt:    now,
	}
	database.DB.Create(&creatorMember)

	if req.CopyMembers {
		var srcMembers []models.ProjectMember
		database.DB.Where("project_id = ?", sourceID).Find(&srcMembers)
		for _, m := range srcMembers {
			if m.UserID != userID && m.UserID != "" {
				newM := models.ProjectMember{
					ProjectID:   newProject.ID,
					UserID:      m.UserID,
					ProjectRole: m.ProjectRole,
					JoinedAt:    now,
				}
				database.DB.Where(models.ProjectMember{ProjectID: newProject.ID, UserID: m.UserID}).FirstOrCreate(&newM)
			}
		}
	}

	// 4. Duplicate Milestones (Map old milestone IDs to new milestone IDs)
	var srcMilestones []models.Milestone
	database.DB.Where("project_id = ?", sourceID).Find(&srcMilestones)

	for _, ms := range srcMilestones {
		newMsID := "ms_" + uuid.New().String()[:8]

		msStatus := ms.Status
		if req.ResetTaskStatus {
			msStatus = "pending"
		}

		newMs := models.Milestone{
			ID:            newMsID,
			ProjectID:     newProject.ID,
			Title:         ms.Title,
			Description:   ms.Description,
			DueDate:       ms.DueDate,
			Status:        msStatus,
			FeedbackNotes: ms.FeedbackNotes,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		if !req.ResetTaskStatus {
			newMs.ApprovedBy = ms.ApprovedBy
			newMs.ApprovedAt = ms.ApprovedAt
		}
		database.DB.Create(&newMs)
	}

	// 5. Duplicate Tasks and Subtasks
	var srcTasks []models.Task
	database.DB.Preload("Subtasks").Where("project_id = ?", sourceID).Find(&srcTasks)

	taskCount := 0
	for _, t := range srcTasks {
		newTaskID := "task_" + uuid.New().String()[:8]
		tStatus := t.Status
		if req.ResetTaskStatus {
			tStatus = "Backlog"
		}

		newTask := models.Task{
			ID:            newTaskID,
			ProjectID:     newProject.ID,
			Key:           t.Key,
			Title:         t.Title,
			Description:   t.Description,
			Category:      t.Category,
			Priority:      t.Priority,
			Status:        tStatus,
			AssigneeID:    t.AssigneeID,
			DueDate:       t.DueDate,
			StartDate:     t.StartDate,
			TimeEstimate:  t.TimeEstimate,
			Fase:          t.Fase,
			EpicModul:     t.EpicModul,
			Deliverable:   t.Deliverable,
			ClientVisible: t.ClientVisible,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		if err := database.DB.Create(&newTask).Error; err == nil {
			taskCount++
			// Duplicate Subtasks
			for _, st := range t.Subtasks {
				newSubID := "sub_" + uuid.New().String()[:8]
				isDone := st.IsDone
				if req.ResetTaskStatus {
					isDone = false
				}
				newSub := models.SubTask{
					ID:        newSubID,
					TaskID:    newTask.ID,
					Title:     st.Title,
					IsDone:    isDone,
					CreatedAt: now,
				}
				database.DB.Create(&newSub)
			}
		}
	}

	// 6. Duplicate Wiki Pages (if requested)
	if req.CopyWiki {
		var srcWikis []models.WikiPage
		database.DB.Where("project_id = ?", sourceID).Find(&srcWikis)
		for _, w := range srcWikis {
			newWikiID := "wiki_" + uuid.New().String()[:8]
			newWiki := models.WikiPage{
				ID:            newWikiID,
				ProjectID:     newProject.ID,
				Title:         w.Title,
				Slug:          w.Slug,
				Content:       w.Content,
				AuthorID:      userID,
				AuthorName:    w.AuthorName,
				ClientVisible: w.ClientVisible,
				CreatedAt:     now,
				UpdatedAt:     now,
			}
			database.DB.Create(&newWiki)
		}
	}

	// 7. Audit Log
	var user models.User
	database.DB.Select("name, avatar_url").Where("id = ?", userID).First(&user)
	logMsg := fmt.Sprintf("Menduplikasi %d deliverable/tugas dari proyek '%s' ke '%s'", taskCount, srcProject.Name, newProject.Name)
	auditLog := models.ActivityLog{
		ID:         "act_" + uuid.New().String()[:8],
		ProjectID:  newProject.ID,
		UserID:     userID,
		UserName:   user.Name,
		UserAvatar: user.AvatarURL,
		Action:     "DUPLICATE_PROJECT",
		Target:     newProject.Name,
		Details:    logMsg,
		CreatedAt:  now,
	}
	database.DB.Create(&auditLog)

	newProject.TotalTasks = taskCount
	newProject.CompletedTasks = 0

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Proyek berhasil disalin beserta seluruh deliverable dan tugasnya",
		"project":      newProject,
		"tasks_copied": taskCount,
	})
}
