package models

type TaskStatus string
type TaskPriority string

const (
	Todo 	TaskStatus = "todo"
	Doing	TaskStatus = "doing"
	Done	TaskStatus = "done"
)

const (
	Low		TaskPriority = "low"
	Medium	TaskPriority = "medium"
	High	TaskPriority = "high"
)

func (s TaskStatus) IsValid() bool {
	switch s {
	case Todo, Doing, Done:
		return true
	}
	return false
}

func (p TaskPriority) IsValid() bool {
	switch p {
	case Low, Medium, High:
		return true
	}
	return false
}