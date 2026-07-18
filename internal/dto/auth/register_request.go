package auth

type RegisterRequest struct {
	Name		string		`json:"name"`
	Email		string		`json:"email"`
	PasswordHash		string		`json:"password_hash"`
}