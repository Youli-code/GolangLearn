package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"createdAt"`
}

var db *sql.DB

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query("SELECT id, title, description, completed, created_at FROM tasks")
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var allTasks []Task

	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Completed, &t.CreatedAt)
		if err != nil {
			http.Error(w, "Failed to scan task", http.StatusInternalServerError)
			return
		}
		allTasks = append(allTasks, t)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error reading rows", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(allTasks)
}

func postTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}

	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if newTask.Title == "" {
		http.Error(w, "Missing task Title", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO tasks (title, description, completed, created_at) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(query, newTask.Title, newTask.Description, false, time.Now())
	if err != nil {
		http.Error(w, "Failed to insert task", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	newTask.ID = int(id)
	newTask.Completed = false
	newTask.CreatedAt = time.Now()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func getTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("ID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var t Task
	query := "SELECT id, title, description, completed, created_at FROM tasks WHERE id = ?"
	err = db.QueryRow(query, id).Scan(&t.ID, &t.Title, &t.Description, &t.Completed, &t.CreatedAt)
	if err == sql.ErrNoRows {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func deleteTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("ID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func updateTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("ID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var updated Task
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if updated.Title == "" {
		http.Error(w, "Missing task Title", http.StatusBadRequest)
		return
	}

	query := `UPDATE tasks SET title = ?, description = ?, completed = ? WHERE id = ?`
	result, err := db.Exec(query, updated.Title, updated.Description, updated.Completed, id)
	if err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	updated.ID = id
	updated.CreatedAt = time.Now()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		panic(err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		completed BOOLEAN,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(createTable); err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", getTaskHandler)
	mux.HandleFunc("POST /tasks", postTaskHandler)
	mux.HandleFunc("GET /tasks/{ID}", getTaskByIDHandler)
	mux.HandleFunc("DELETE /tasks/{ID}", deleteTaskByIDHandler)
	mux.HandleFunc("PUT /tasks/{ID}", updateTaskByIDHandler)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
