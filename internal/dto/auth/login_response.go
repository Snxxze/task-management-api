package auth

type LoginResponse struct {
	ID		uint		`json:"id"`
	Name	string	`json:"name"`
	Email		string	`json:"email"`
}