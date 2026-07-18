package apperrors

import "errors"

var (
	ErrTaskNotFound = errors.New("task not found")

	ErrInvalidTaskStatus = errors.New("invalid task status")
)