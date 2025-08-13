package main

import "errors"

type TaskStore interface {
	GetAllTasks(filterCompleted *bool) ([]Task, error)
	GetTaskByID(id int) (Task, error)
	CreateTask(task *Task) error
	UpdateTask(task *Task) error
	DeleteTask(id int) error
}

var (
	ErrNotFound = errors.New("not found")
	ErrInvalid  = errors.New("invalid input")
)
