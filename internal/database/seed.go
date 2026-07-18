package database

import (
	"fmt"
	"task-management-api/internal/models"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	user := models.User{
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashed-password",
	}

	err := db.Where("email = ?", user.Email).FirstOrCreate(&user).Error
	if err != nil {
		return fmt.Errorf("seed user: %w", err)
	}

	return nil
}
