package project

type UpdateProjectRequest struct {
	Name					*string		`json:"name"`
	Description		*string		`json:"description"`
	Color					*string		`json:"color"`
	IsArchived		*bool			`json:"is_archived"`
}