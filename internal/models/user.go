package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Name		string	`gorm:"type:varchar(100);not null"`
	Email		string	`gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash	string	`gorm:"type:text;not null"`

	Projects []Project		`gorm:"foreignKey:UserID"`
}