package project

type CreateProjectRequest struct {
	Name					string		`json:"name" binding:"required"`
	Description		string		`json:"description"`
	Color					string		`json:"color"`
}