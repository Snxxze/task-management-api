package project

type CreateProjectRequest struct {
	Name					string		`json:"name"`
	Description		string		`json:"description"`
	Color					string		`json:"color"`
}