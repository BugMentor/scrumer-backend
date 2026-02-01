package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-faker/faker/v4"
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

	log.Println("Seeding database with synthetic data...")

	// Create users
	var users []models.User
	for i := 0; i < 5; i++ {
		user := models.User{
			Username: faker.Username(),
			Email:    faker.Email(),
			Password: "password", // Static password for simplicity
		}
		if err := db.Create(&user).Error; err != nil {
			log.Printf("Failed to create user: %v", err)
			continue
		}
		users = append(users, user)
	}
	log.Printf("Created %d users", len(users))

	// Create projects
	var projects []models.Project
	for i := 0; i < 3; i++ {
		project := models.Project{
			Name:        faker.Word(),
			Description: faker.Sentence(),
		}
		if err := db.Create(&project).Error; err != nil {
			log.Printf("Failed to create project: %v", err)
			continue
		}
		// Associate users with projects
		for _, user := range users {
			db.Model(&project).Association("Users").Append(&user)
		}
		projects = append(projects, project)
	}
	log.Printf("Created %d projects", len(projects))

	// Create sprints and tasks
	for _, project := range projects {
		for i := 0; i < 2; i++ { // 2 sprints per project
			dateString := faker.Date()
			parsedDate, err := time.Parse("2006-01-02", dateString) // Assuming "YYYY-MM-DD" format
			if err != nil {
				log.Printf("Failed to parse date string %s: %v", dateString, err)
				continue // Skip sprint creation if date parsing fails
			}
			startDate := parsedDate.Add(time.Duration(i*7*24) * time.Hour) // Start date spaced 7 days apart
			endDate := startDate.Add(14 * 24 * time.Hour)                   // 14-day sprint

			sprint := models.Sprint{
				Name:      fmt.Sprintf("%s Sprint %d", project.Name, i+1),
				StartDate: startDate,
				EndDate:   endDate,
				ProjectID: project.ID,
			}
			if err := db.Create(&sprint).Error; err != nil {
				log.Printf("Failed to create sprint: %v", err)
				continue
			}

			for j := 0; j < 5; j++ { // 5 tasks per sprint
				task := models.Task{
					Title:       faker.Sentence(),
					Description: faker.Paragraph(),
					Status:      faker.Word(), // Will be a random word, could make this more structured
					Priority:    faker.Word(), // Will be a random word, could make this more structured
					SprintID:    sprint.ID,
					AssigneeID:  users[j%len(users)].ID, // Assign tasks to users in a round-robin fashion
				}
				if err := db.Create(&task).Error; err != nil {
					log.Printf("Failed to create task: %v", err)
					continue
				}
			}
			log.Printf("Created sprint '%s' with tasks", sprint.Name)
		}
	}

	log.Println("Database seeding completed.")
}
