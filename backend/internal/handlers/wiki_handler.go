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
	"pulse-backend/internal/services"
)

// ListWikiPages returns documents for a project
func ListWikiPages(c *gin.Context) {
	projectID := c.Param("id")
	_, role := middleware.GetCurrentUser(c)

	query := database.DB.Where("project_id = ?", projectID)
	if role == auth.RoleClient {
		query = query.Where("client_visible = ?", true)
	}

	var pages []models.WikiPage
	if err := query.Order("created_at desc").Find(&pages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch wiki pages"})
		return
	}

	c.JSON(http.StatusOK, pages)
}

// GetWikiPage returns a single documentation page
func GetWikiPage(c *gin.Context) {
	pageID := c.Param("page_id")
	_, role := middleware.GetCurrentUser(c)

	var page models.WikiPage
	if err := database.DB.Where("id = ?", pageID).First(&page).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wiki page not found"})
		return
	}

	if role == auth.RoleClient && !page.ClientVisible {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: internal documentation"})
		return
	}

	c.JSON(http.StatusOK, page)
}

type CreateWikiRequest struct {
	Title         string `json:"title" binding:"required"`
	Content       string `json:"content" binding:"required"`
	ClientVisible *bool  `json:"client_visible"`
}

// CreateWikiPage adds a new knowledge base document
func CreateWikiPage(c *gin.Context) {
	projectID := c.Param("id")
	userID, _ := middleware.GetCurrentUser(c)

	var req CreateWikiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	clientVisible := true
	if req.ClientVisible != nil {
		clientVisible = *req.ClientVisible
	}

	var user models.User
	database.DB.Where("id = ?", userID).First(&user)

	page := models.WikiPage{
		ID:            "wiki_" + uuid.New().String()[:8],
		ProjectID:     projectID,
		Title:         req.Title,
		Slug:          req.Title,
		Content:       req.Content,
		AuthorID:      userID,
		AuthorName:    user.Name,
		ClientVisible: clientVisible,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := database.DB.Create(&page).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wiki page"})
		return
	}

	services.LogActivity(projectID, userID, user.Name, user.AvatarURL, "created wiki page", page.Title, "")

	c.JSON(http.StatusCreated, page)
}

// UpdateWikiPage updates a wiki page
func UpdateWikiPage(c *gin.Context) {
	pageID := c.Param("page_id")

	var page models.WikiPage
	if err := database.DB.Where("id = ?", pageID).First(&page).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wiki page not found"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	req["updated_at"] = time.Now()
	database.DB.Model(&page).Updates(req)

	c.JSON(http.StatusOK, page)
}

// DeleteWikiPage deletes a wiki page
func DeleteWikiPage(c *gin.Context) {
	pageID := c.Param("page_id")
	database.DB.Where("id = ?", pageID).Delete(&models.WikiPage{})
	c.JSON(http.StatusOK, gin.H{"message": "Wiki page deleted", "id": pageID})
}
