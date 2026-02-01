package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"scrumer-backend/models"
)

func main() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Purging database data...")

	// Drop tables in reverse order of dependency, including the join table
	err = db.Migrator().DropTable(&models.Task{}, &models.Sprint{}, "user_projects", &models.Project{}, &models.User{})
	if err != nil {
		log.Fatal("Failed to drop tables:", err)
	}

	log.Println("Database purging completed.")
}
