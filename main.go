package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/handler"
	"github.com/joho/godotenv"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"scrumer-backend/graph"
	"scrumer-backend/models"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	_ = godotenv.Load()
	db, err := openDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	if err := db.Exec("SELECT 1").Error; err != nil {
		log.Fatal("Failed to verify database connection:", err)
	}
	log.Println("Database connection verified.")

	if err := migrateDB(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migration complete.")

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get sql.DB:", err)
	}
	defer sqlDB.Close()

	r, err := setupRouter(db)
	if err != nil {
		log.Fatalf("failed to setup router: %v", err)
	}
	r.Run() // listen and serve on 0.0.0.0:8080
}

func openDB() (*gorm.DB, error) {
	driver := getEnv("DB_DRIVER", "sqlite")
	if driver == "postgres" {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			getEnv("DB_HOST", "localhost"),
			getEnv("DB_USER", "user"),
			getEnv("DB_PASSWORD", ""),
			getEnv("DB_NAME", "scrumer"),
			getEnv("DB_PORT", "5432"),
			getEnv("DB_SSLMODE", "disable"))
		return gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}
	path := getEnv("DB_SQLITE_PATH", "scrumer.db")
	return gorm.Open(sqlite.Open(path), &gorm.Config{})
}

func migrateDB(db *gorm.DB) error {
	if getEnv("DB_DRIVER", "sqlite") == "sqlite" {
		return db.AutoMigrate(&models.User{}, &models.Project{}, &models.Sprint{}, &models.Task{})
	}
	if !db.Migrator().HasTable(&models.User{}) {
		if err := db.Migrator().CreateTable(&models.User{}); err != nil {
			return err
		}
	}
	if !db.Migrator().HasTable(&models.Project{}) {
		if err := db.Migrator().CreateTable(&models.Project{}); err != nil {
			return err
		}
	}
	if !db.Migrator().HasTable(&models.Sprint{}) {
		if err := db.Migrator().CreateTable(&models.Sprint{}); err != nil {
			return err
		}
	}
	if !db.Migrator().HasTable(&models.Task{}) {
		if err := db.Migrator().CreateTable(&models.Task{}); err != nil {
			return err
		}
	}
	if !db.Migrator().HasTable("user_projects") {
		return db.Exec(`
			CREATE TABLE user_projects (
				user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				project_id bigint NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
				created_at timestamptz,
				updated_at timestamptz,
				deleted_at timestamptz,
				PRIMARY KEY (user_id, project_id)
			);
		`).Error
	}
	return nil
}

// setupRouter builds the Gin router with /ping and /graphql. Used by main and integration tests.
func setupRouter(db *gorm.DB) (*gin.Engine, error) {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	schema, err := graph.NewSchema(db)
	if err != nil {
		return nil, err
	}
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})
	r.POST("/graphql", gin.WrapH(h))
	r.GET("/graphql", gin.WrapH(h))
	return r, nil
}
