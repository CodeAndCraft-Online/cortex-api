package database

import (
	"fmt"
	"log"
	"os"
	"time"

	models "github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {

	// ✅ Read values from environment variables
	dbUser := os.Getenv("POSTGRES_USER")
	dbPass := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPass, dbName, dbPort)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Database connection successful")

	// AutoMigrate ensures tables are created automatically
	err = DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}, &models.Vote{}, &models.SubInvitation{}, &models.Sub{}, &models.SubMembership{}, &models.PasswordResetToken{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	} else {
		log.Println("Database migrated successfully!")
	}

	go DeleteExpiredTokens()
}

// DeleteExpiredTokens removes tokens that are past expiration
func DeleteExpiredTokens() {
	for {
		time.Sleep(1 * time.Hour) // Runs every hour
		result := DB.Where("expires_at < ?", time.Now()).Delete(&models.PasswordResetToken{})
		if result.Error != nil {
			log.Println("Error deleting expired tokens:", result.Error)
		} else if result.RowsAffected > 0 {
			log.Println("Deleted", result.RowsAffected, "expired reset tokens.")
		}
	}
}
