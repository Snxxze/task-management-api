package apperrors

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid email or password")

	ErrInvalidToken       = errors.New("invalid token")
	
	ErrExpiredToken       = errors.New("token expired")
)