package main

import (
	"database/sql"
	"errors"
	"time"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(db *sql.DB) *SQLiteStore {
	return &SQLiteStore{db: db}
}

func (s *SQLiteStore) GetAllTasks(filterCompleted *bool) ([]Task, error) {
	query := `SELECT id, title, description, completed, created_at, updated_at FROM tasks`
	var args []any

	if filterCompleted != nil {
		query += ` WHERE completed = ?`
		args = append(args, *filterCompleted)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Completed, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func (s *SQLiteStore) GetTaskByID(id int) (Task, error) {
	var t Task
	err := s.db.QueryRow(
		`SELECT id, title, description, completed, created_at, updated_at FROM tasks WHERE id = ?`,
		id,
	).Scan(&t.ID, &t.Title, &t.Description, &t.Completed, &t.CreatedAt, &t.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Task{}, ErrNotFound
		}
		return Task{}, err
	}
	return t, nil
}

func (s *SQLiteStore) CreateTask(task *Task) error {
	if task.Title == "" {
		return ErrInvalid
	}
	now := time.Now()
	res, err := s.db.Exec(
		`INSERT INTO tasks (title, description, completed, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		task.Title, task.Description, false, now, now,
	)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	task.ID = int(id)
	task.Completed = false
	task.CreatedAt = now
	task.UpdatedAt = now
	return nil
}

func (s *SQLiteStore) UpdateTask(task *Task) error {
	if task.ID <= 0 || task.Title == "" {
		return ErrInvalid
	}
	now := time.Now()
	res, err := s.db.Exec(
		`UPDATE tasks SET title = ?, description = ?, completed = ?, updated_at = ? WHERE id = ?`,
		task.Title, task.Description, task.Completed, now, task.ID,
	)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	task.UpdatedAt = now
	return nil
}

func (s *SQLiteStore) DeleteTask(id int) error {
	res, err := s.db.Exec(`DELETE FROM tasks WHERE id = ?`, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}
