package routes

import (
	"strings"
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

	// Static file serving for uploaded task assets
	r.Static("/uploads", "./uploads")

	// CORS configuration for Vue 3 frontend integration and Cloudflare Tunnel
	corsConfig := cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://localhost:4173",
			"http://localhost:3000",
			"http://localhost:5080",
			"https://project.kedukanaksara.com",
			"http://project.kedukanaksara.com",
		},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasSuffix(origin, "kedukanaksara.com") ||
				strings.Contains(origin, "localhost") ||
				strings.Contains(origin, "127.0.0.1") {
				return true
			}
			return false
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(corsConfig))

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "online",
			"time":   time.Now().Format(time.RFC3339),
			"system": "Pulse Backend RBAC v1.5.0",
		})
	})

	api := r.Group("/api/v1")
	{
		// 1. Public Authentication Routes
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", handlers.Login)
			authGroup.POST("/logout", handlers.Logout)
			authGroup.GET("/config", handlers.GetAuthConfig)
		}

		// 2. Protected Routes (Requires valid JWT Bearer Token)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Current user profile + role + permissions
			protected.GET("/auth/me", handlers.GetMe)
			protected.GET("/auth/heartbeat", handlers.Heartbeat)

			// Real-time Executive Dashboard KPIs & Activity Feed
			dashboardGroup := protected.Group("/dashboard")
			{
				dashboardGroup.GET("/kpis", handlers.GetKPIStats)
				dashboardGroup.GET("/activity", handlers.GetActivityFeed)
			}

			// In-App Notifications
			notifGroup := protected.Group("/notifications")
			{
				notifGroup.GET("", handlers.ListNotifications)
				notifGroup.PATCH("/:id/read", handlers.MarkNotificationRead)
				notifGroup.POST("/read-all", handlers.MarkAllNotificationsRead)
			}

			// Timesheets & Team Workload
			timeGroup := protected.Group("/timesheets")
			{
				timeGroup.POST("", handlers.LogTime)
				timeGroup.GET("", handlers.ListTimeEntries)
				timeGroup.GET("/workload", handlers.GetTeamWorkload)
			}

			// User Management (Admin CRUD & Global role updates)
			usersGroup := protected.Group("/users")
			{
				usersGroup.GET("", middleware.RequirePermission(auth.PermUserManage), handlers.ListUsers)
				usersGroup.GET("/:id", middleware.RequirePermission(auth.PermUserManage), handlers.GetUser)
				usersGroup.POST("", middleware.RequirePermission(auth.PermUserManage), handlers.CreateUser)
				usersGroup.PATCH("/:id", middleware.RequirePermission(auth.PermUserManage), handlers.UpdateUser)
				usersGroup.PUT("/:id", middleware.RequirePermission(auth.PermUserManage), handlers.UpdateUser)
				usersGroup.DELETE("/:id", middleware.RequirePermission(auth.PermUserManage), handlers.DeleteUser)
				usersGroup.PATCH("/:id/role", middleware.RequireRole(auth.RoleSuperAdmin), handlers.UpdateUserRole)
				usersGroup.POST("/:id/reset-session", middleware.RequirePermission(auth.PermUserManage), handlers.ResetUserSession)
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

					// Project Wiki / Knowledge Base
					wikiGroup := projectScoped.Group("/wiki")
					{
						wikiGroup.GET("", handlers.ListWikiPages)
						wikiGroup.POST("", middleware.RequirePermission(auth.PermProjectEdit), handlers.CreateWikiPage)
						wikiGroup.GET("/:page_id", handlers.GetWikiPage)
						wikiGroup.PATCH("/:page_id", middleware.RequirePermission(auth.PermProjectEdit), handlers.UpdateWikiPage)
						wikiGroup.DELETE("/:page_id", middleware.RequirePermission(auth.PermProjectEdit), handlers.DeleteWikiPage)
					}

					// Project Budgets & Financials
					budgetGroup := projectScoped.Group("/budgets")
					{
						budgetGroup.GET("", handlers.GetProjectBudgets)
						budgetGroup.POST("", middleware.RequirePermission(auth.PermProjectEdit), handlers.CreateBudgetItem)
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

				// Attachments Upload & Deletion
				tasksGroup.POST("/:id/attachments", middleware.RequirePermission(auth.PermAttachmentUpload), handlers.UploadAttachment)
				tasksGroup.DELETE("/attachments/:att_id", middleware.RequirePermission(auth.PermTaskEdit), handlers.DeleteAttachment)
			}

			// Milestones Management (Contract Fix: accepts both PATCH and POST for approve/reject)
			milestonesGroup := protected.Group("/milestones")
			{
				milestonesGroup.GET("", handlers.ListMilestones)
				milestonesGroup.PATCH("/:id/approve", middleware.RequirePermission(auth.PermMilestoneApprove), handlers.ApproveMilestone)
				milestonesGroup.POST("/:id/approve", middleware.RequirePermission(auth.PermMilestoneApprove), handlers.ApproveMilestone)
				milestonesGroup.PATCH("/:id/reject", middleware.RequirePermission(auth.PermMilestoneApprove), handlers.RejectMilestone)
				milestonesGroup.POST("/:id/reject", middleware.RequirePermission(auth.PermMilestoneApprove), handlers.RejectMilestone)
			}
		}
	}

	return r
}
