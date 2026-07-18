package apperrors

import "errors"

var (
	ErrInternal = errors.New("internal server error")

	ErrNotFound = errors.New("resource not found")

	ErrConflict = errors.New("resource conflict")

	ErrUnauthorized = errors.New("unauthorized")

	ErrForbidden = errors.New("forbidden")

	ErrBadRequest = errors.New("bad request")
)