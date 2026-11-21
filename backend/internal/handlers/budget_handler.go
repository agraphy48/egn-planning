package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"pulse-backend/internal/database"
	"pulse-backend/internal/models"
)

// GetProjectBudgets returns budget items and total allocation
func GetProjectBudgets(c *gin.Context) {
	projectID := c.Param("id")

	var budgets []models.Budget
	if err := database.DB.Where("project_id = ?", projectID).Find(&budgets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch budget"})
		return
	}

	var totalAllocated, totalSpent float64
	for _, b := range budgets {
		totalAllocated += b.Allocated
		totalSpent += b.Spent
	}

	c.JSON(http.StatusOK, gin.H{
		"budgets":         budgets,
		"total_allocated": totalAllocated,
		"total_spent":     totalSpent,
		"currency":        "IDR",
	})
}

type CreateBudgetRequest struct {
	Category  string  `json:"category" binding:"required"`
	Allocated float64 `json:"allocated" binding:"required,gt=0"`
	Spent     float64 `json:"spent"`
	Currency  string  `json:"currency"`
}

// CreateBudgetItem adds a budget allocation
func CreateBudgetItem(c *gin.Context) {
	projectID := c.Param("id")

	var req CreateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	currency := req.Currency
	if currency == "" {
		currency = "IDR"
	}

	item := models.Budget{
		ID:        "bud_" + uuid.New().String()[:8],
		ProjectID: projectID,
		Category:  req.Category,
		Allocated: req.Allocated,
		Spent:     req.Spent,
		Currency:  currency,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := database.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create budget item"})
		return
	}

	c.JSON(http.StatusCreated, item)
}
