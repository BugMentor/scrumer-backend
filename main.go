package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/handler"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"scrumer-backend/graph"
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

	// Auto migrate models in dependency order
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to auto migrate User model:", err)
	}

	err = db.AutoMigrate(&models.Project{})
	if err != nil {
		log.Fatal("Failed to auto migrate Project model:", err)
	}

	err = db.AutoMigrate(&models.Sprint{})
	if err != nil {
		log.Fatal("Failed to auto migrate Sprint model:", err)
	}

	err = db.AutoMigrate(&models.Task{})
	if err != nil {
		log.Fatal("Failed to auto migrate Task model:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get sql.DB:", err)
	}
	defer sqlDB.Close()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	schema, err := graph.NewSchema(db)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	// Pass db to the handler
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	r.POST("/graphql", gin.WrapH(h))
	r.GET("/graphql", gin.WrapH(h))

	r.Run() // listen and serve on 0.0.0.0:8080
}
