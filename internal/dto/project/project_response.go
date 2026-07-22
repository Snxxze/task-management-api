package project

import "time"

type ProjectResponse struct {
	ID          uint				`json:"id"`
	Name        string			`json:"name"`
	Description string			`json:"description"`
	Color       string			`json:"color"`

	IsArchived  bool				`json:"is_archived"`

	CreatedAt time.Time			`json:"created_at"`
	UpdatedAt	time.Time			`json:"updated_at"`
}