# GolangLearn

This repository is a collection of Go-based projects created for the purpose of learning and experimentation. It includes both backend API work and explorations into concurrency, channels, and goroutines.

## Task_API/

A basic task manager API built in Go. The goal here was to explore how to structure a simple REST backend using Go's `net/http` package, connect it to a SQLite database, and gradually introduce features like authentication and testing.

### Key Features

- Full CRUD functionality for task management
- SQLite integration with an abstracted `store.go` layer
- Clear separation between HTTP handlers and business logic
- Request validation and consistent JSON error responses
- JWT authentication (login, token validation middleware)
- Unit testing using mock interfaces for handler logic

This part of the project is meant to represent what a beginner-to-intermediate Go backend project might look like. It focuses more on simplicity, readability, and clarity than performance or advanced design patterns.

## Workers/

This folder contains a standalone program that demonstrates how goroutines, function-passing, and channels work together in Go.

Each worker:
- Is launched as a goroutine
- Takes a piece of data and a function to apply to it
- Sends the result back to the main function via a channel

This is a practical example of concurrency using Go's `sync.WaitGroup` and `chan` types.

### How to Run

From inside the `Workers/` folder:

```bash
go run main.go
