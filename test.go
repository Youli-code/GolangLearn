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

}
func main() {

	http.HandleFunc("/tasks", getTaskHandler)

	fmt.Println("Server running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", nil)
	}

}
