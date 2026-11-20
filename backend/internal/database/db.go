package database

import (
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"pulse-backend/internal/models"
)

var DB *gorm.DB

// InitDB initializes SQLite database and auto-migrates all models
func InitDB(dbPath string) (*gorm.DB, error) {
	if dbPath == "" {
		dbPath = "pulse.db"
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}

	// AutoMigrate all schemas including project_members pivot table
	err = db.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.ProjectMember{},
		&models.Milestone{},
		&models.Task{},
		&models.SubTask{},
		&models.Comment{},
	)
	if err != nil {
		return nil, err
	}

	DB = db
	log.Println("[Database] Initialized and auto-migrated schemas successfully")
	return db, nil
}
