package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Task struct {
	ID          int
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
}

var tasks = []Task{
	{
		ID:          1,
		Title:       "Learn Go",
		Description: "Complete the Go programming course",
		Completed:   false,
		CreatedAt:   time.Now(),
	},
	{
		ID:          2,
		Title:       "Build API",
		Description: "Implement task CRUD operation",
		Completed:   false,
		CreatedAt:   time.Now(),
	},
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	err := json.NewEncoder(w).Encode(tasks)
	if err != nil {
		http.Error(w, "Failed to encode tasks", http.StatusInternalServerError)
	}

}

func postTaskHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return

	}

	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, "Invalid Json", http.StatusBadRequest)
		return
	}

	if newTask.Title == "" {
		http.Error(w, "Mising task Title", http.StatusBadRequest)
		return
	}

	newTask.ID = len(tasks) + 1
	newTask.CreatedAt = time.Now()
	newTask.Completed = false

	tasks = append(tasks, newTask)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}
func main() {

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getTaskHandler(w, r)
		case http.MethodPost:
			postTaskHandler(w, r)
		default:
			http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", nil)
	}

}
