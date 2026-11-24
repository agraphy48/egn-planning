package main

import (
	"log"
	"os"

	"pulse-backend/internal/auth"
	"pulse-backend/internal/database"
	"pulse-backend/internal/mailer"
	"pulse-backend/internal/routes"
)

func main() {
	log.Println("[Pulse] Starting Project Management RBAC Backend Server...")

	// 0. Load Environment Variables (.env)
	mailer.LoadEnv()

	// 1. Initialize SQLite Database & Auto-Migrate Schemas
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "pulse.db"
	}

	_, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatalf("[Pulse] Database initialization failed: %v", err)
	}

	// Initialize token revocation memory cache from database
	auth.InitRevocationStore(database.DB)

	// 2. Seed 4 Core Role Accounts, Projects, and Membership Pivots
	database.SeedDatabase()
	database.SeedAuxiliaryData()

	// 3. Initialize Router & RBAC Middleware
	r := routes.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("[Pulse] Server listening on http://localhost:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("[Pulse] Server failed to start: %v", err)
	}
}
