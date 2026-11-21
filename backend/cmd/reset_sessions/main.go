package main

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"

	"pulse-backend/internal/database"
	"pulse-backend/internal/models"
)

func main() {
	db, err := database.InitDB("pulse.db")
	if err != nil {
		log.Fatalf("Failed to open pulse.db: %v", err)
	}

	// 1. Release all users from all devices (clear current_session_id)
	result := db.Model(&models.User{}).Where("1 = 1").Updates(map[string]interface{}{
		"current_session_id": "",
		"last_active_at":     time.Now(),
	})
	if result.Error != nil {
		log.Fatalf("Failed to clear sessions: %v", result.Error)
	}
	fmt.Printf("[OK] Released all users from devices! Rows updated: %d\n", result.RowsAffected)

	// 2. Clear revoked tokens
	db.Exec("DELETE FROM revoked_tokens")
	fmt.Println("[OK] Cleared revoked_tokens table")

	// 3. Verify all users and their passwords
	var users []models.User
	db.Find(&users)
	fmt.Printf("[OK] Found %d total users in pulse.db:\n", len(users))
	for _, u := range users {
		err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte("Password123!"))
		status := "VALID Password123!"
		if err != nil {
			status = fmt.Sprintf("PASSWORD MISMATCH: %v", err)
		}
		fmt.Printf(" - [%s] %s (%s) -> %s\n", u.Role, u.Name, u.Email, status)
	}
}
