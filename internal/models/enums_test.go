package models_test

import (
	"testing"
	"task-management-api/internal/models"
)

func TestTaskStatus_IsValid(t *testing.T) {
	tests := []struct {
		status   models.TaskStatus
		expected bool
	}{
		{models.Todo, true},
		{models.Doing, true},
		{models.Done, true},
		{models.TaskStatus("invalid"), false},
		{models.TaskStatus(""), false},
	}

	for _, tt := range tests {
		if got := tt.status.IsValid(); got != tt.expected {
			t.Errorf("TaskStatus(%s).IsValid() = %v, want %v", tt.status, got, tt.expected)
		}
	}
}

func TestTaskPriority_IsValid(t *testing.T) {
	tests := []struct {
		priority models.TaskPriority
		expected bool
	}{
		{models.Low, true},
		{models.Medium, true},
		{models.High, true},
		{models.TaskPriority("invalid"), false},
		{models.TaskPriority(""), false},
	}

	for _, tt := range tests {
		if got := tt.priority.IsValid(); got != tt.expected {
			t.Errorf("TaskPriority(%s).IsValid() = %v, want %v", tt.priority, got, tt.expected)
		}
	}
}
