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