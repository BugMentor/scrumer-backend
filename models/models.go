package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
	// Projects []Project `gorm:"many2many:user_projects;"`
}

type Project struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex;not null"`
	Description string
	// Users       []User `gorm:"many2many:user_projects;"`
	Sprints     []Sprint
}

type Sprint struct {
	gorm.Model
	Name      string `gorm:"not null"`
	StartDate time.Time
	EndDate   time.Time
	ProjectID uint
	Tasks     []Task
}

type Task struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string
	Status      string `gorm:"default:'Todo'"`
	Priority    string `gorm:"default:'Medium'"`
	SprintID    uint
	AssigneeID  uint
	Assignee    User `gorm:"foreignkey:AssigneeID"`
}
