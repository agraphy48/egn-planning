package routes

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"pulse-backend/internal/auth"
	"pulse-backend/internal/handlers"
	"pulse-backend/internal/middleware"
)

// SetupRouter initializes Gin engine with CORS, RBAC middleware, and protected routes
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS configuration for Vue 3 frontend integration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:4173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "online",
			"time":   time.Now().Format(time.RFC3339),
			"system": "Pulse Backend RBAC v1.4.0",
		})
	})

	api := r.Group("/api/v1")
	{
		// 1. Public Authentication Routes
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", handlers.Login)
		}

		// 2. Protected Routes (Requires valid JWT Bearer Token)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Current user profile + role + permissions
			protected.GET("/auth/me", handlers.GetMe)

			// User Management (Admin & Global role updates)
			usersGroup := protected.Group("/users")
			{
				usersGroup.GET("", middleware.RequirePermission(auth.PermUserManage), handlers.ListUsers)
				usersGroup.PATCH("/:id/role", middleware.RequireRole(auth.RoleSuperAdmin), handlers.UpdateUserRole)
			}

			// Projects Management
			projectsGroup := protected.Group("/projects")
			{
				projectsGroup.GET("", middleware.RequirePermission(auth.PermProjectRead), handlers.ListProjects)
				projectsGroup.POST("", middleware.RequirePermission(auth.PermProjectCreate), handlers.CreateProject)

				// Single project operations requiring project membership
				projectScoped := projectsGroup.Group("/:id")
				projectScoped.Use(middleware.RequireProjectAccess())
				{
					projectScoped.GET("", handlers.GetProject)
					projectScoped.PATCH("", middleware.RequirePermission(auth.PermProjectEdit), handlers.UpdateProject)
					projectScoped.DELETE("", middleware.RequirePermission(auth.PermProjectDelete), handlers.DeleteProject)

					// Project Members Sub-resource (Super Admin & PM)
					membersGroup := projectScoped.Group("/members")
					{
						membersGroup.GET("", handlers.ListProjectMembers)
						membersGroup.POST("", middleware.RequirePermission(auth.PermProjectMemberManage), handlers.AddProjectMember)
						membersGroup.DELETE("/:user_id", middleware.RequirePermission(auth.PermProjectMemberManage), handlers.RemoveProjectMember)
					}
				}
			}

			// Tasks Management
			tasksGroup := protected.Group("/tasks")
			{
				tasksGroup.GET("", middleware.RequireProjectAccess(), handlers.ListTasks)
				tasksGroup.POST("", middleware.RequirePermission(auth.PermTaskCreate), middleware.RequireProjectAccess(), handlers.CreateTask)
				tasksGroup.GET("/:id", handlers.GetTask)
				tasksGroup.PATCH("/:id/status", middleware.RequirePermission(auth.PermTaskMove), handlers.UpdateTaskStatus)
				tasksGroup.PATCH("/:id", middleware.RequirePermission(auth.PermTaskEdit), handlers.UpdateTask)
				tasksGroup.DELETE("/:id", middleware.RequirePermission(auth.PermTaskDelete), handlers.DeleteTask)

				// Subtasks & Comments
				tasksGroup.POST("/:id/subtasks", middleware.RequirePermission(auth.PermSubtaskManage), handlers.AddSubtask)
				tasksGroup.PATCH("/:id/subtasks/:sub_id", middleware.RequirePermission(auth.PermSubtaskManage), handlers.ToggleSubtask)
				tasksGroup.POST("/:id/comments", middleware.RequirePermission(auth.PermCommentCreate), handlers.AddComment)
			}

			// Milestones Management
			milestonesGroup := protected.Group("/milestones")
			{
				milestonesGroup.GET("", handlers.ListMilestones)
				milestonesGroup.PATCH("/:id/approve", middleware.RequirePermission(auth.PermMilestoneApprove), handlers.ApproveMilestone)
			}
		}
	}

	return r
}
