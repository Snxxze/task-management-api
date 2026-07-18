package models

import (
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model

	Name		string		`gorm:"type:varchar(100);not null"`
	Description string	`gorm:"type:text"`

	Color		string	`gorm:"type:varchar(20)"`
	IsArchived	bool

	UserID uint `gorm:"not null"`
	User	User	`gorm:"foreignKey:UserID"`

	Tasks		[]Task	`gorm:"foreignKey:ProjectID"`
}