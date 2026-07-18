package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model

	Title			string		`gorm:"type:varchar(255);not null"`
	Description	string	`gorm:"type:text"`

	Status	TaskStatus	`gorm:"type:varchar(20);not null;default:'todo'"`
	Priority TaskPriority `gorm:"type:varchar(20);not null;default:'medium'"`

	DueDate     *time.Time
  CompletedAt *time.Time

	ProjectID		uint	`gorm:"not null"`
	Project		Project	`gorm:"foreignKey:ProjectID"`
}
