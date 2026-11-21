package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"pulse-backend/internal/database"
	"pulse-backend/internal/middleware"
	"pulse-backend/internal/models"
	"pulse-backend/internal/services"
)

// ListMilestones returns milestones for a project
func ListMilestones(c *gin.Context) {
	projectID := c.Query("project_id")

	var milestones []models.Milestone
	query := database.DB
	if projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}

	if err := query.Find(&milestones).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch milestones"})
		return
	}

	c.JSON(http.StatusOK, milestones)
}

type ApproveMilestoneRequest struct {
	ClientName string `json:"client_name"`
}

// ApproveMilestone signs off a milestone (permitted for client, pm, super_admin)
func ApproveMilestone(c *gin.Context) {
	milestoneID := c.Param("id")
	userID, _ := middleware.GetCurrentUser(c)

	var req ApproveMilestoneRequest
	_ = c.ShouldBindJSON(&req)

	var milestone models.Milestone
	if err := database.DB.Where("id = ?", milestoneID).First(&milestone).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Milestone not found"})
		return
	}

	clientName := req.ClientName
	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err == nil && clientName == "" {
		clientName = user.Name
	} else if clientName == "" {
		clientName = "Authorized Client Stakeholder"
	}

	now := time.Now()
	milestone.Status = "approved"
	milestone.ApprovedBy = clientName
	milestone.ApprovedAt = &now

	if err := database.DB.Save(&milestone).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve milestone"})
		return
	}

	services.LogActivity(milestone.ProjectID, userID, clientName, user.AvatarURL, "approved milestone", milestone.Title, "Status: Approved")

	c.JSON(http.StatusOK, gin.H{
		"message":     "Milestone approved successfully",
		"milestone":   milestone,
		"approved_by": clientName,
		"approved_at": now,
	})
}

type RejectMilestoneRequest struct {
	ClientName string `json:"client_name"`
	Reason     string `json:"reason" binding:"required"`
}

// RejectMilestone requests changes / flags revisions on a milestone
func RejectMilestone(c *gin.Context) {
	milestoneID := c.Param("id")
	userID, _ := middleware.GetCurrentUser(c)

	var req RejectMilestoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Revision reason is required"})
		return
	}

	var milestone models.Milestone
	if err := database.DB.Where("id = ?", milestoneID).First(&milestone).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Milestone not found"})
		return
	}

	var user models.User
	clientName := req.ClientName
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err == nil && clientName == "" {
		clientName = user.Name
	}

	milestone.Status = "revision"
	milestone.FeedbackNotes = req.Reason

	if err := database.DB.Save(&milestone).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update milestone"})
		return
	}

	services.LogActivity(milestone.ProjectID, userID, clientName, user.AvatarURL, "requested milestone revision", milestone.Title, "Feedback: "+req.Reason)

	c.JSON(http.StatusOK, gin.H{
		"message":   "Milestone revision requested",
		"milestone": milestone,
	})
}
