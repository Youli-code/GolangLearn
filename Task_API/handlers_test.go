package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockStore struct {
	GetAllTasksFunc func(filterCompleted *bool) ([]Task, error)
	GetTaskByIDFunc func(id int) (Task, error)
	CreateTaskFunc  func(task *Task) error
	UpdateTaskFunc  func(task *Task) error
	DeleteTaskFunc  func(id int) error
}

func (m *MockStore) GetAllTasks(filterCompleted *bool) ([]Task, error) {
	return m.GetAllTasksFunc(filterCompleted)
}
func (m *MockStore) GetTaskByID(id int) (Task, error) {
	return m.GetTaskByIDFunc(id)
}
func (m *MockStore) CreateTask(task *Task) error {
	return m.CreateTaskFunc(task)
}
func (m *MockStore) UpdateTask(task *Task) error {
	return m.UpdateTaskFunc(task)
}
func (m *MockStore) DeleteTask(id int) error {
	return m.DeleteTaskFunc(id)
}

func TestGetTaskHandler_ReturnsTasks(t *testing.T) {
	mockStore := &MockStore{
		GetAllTasksFunc: func(_ *bool) ([]Task, error) {
			return []Task{
				{ID: 1, Title: "Test Task", Description: "test desc", Completed: false},
			}, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rec := httptest.NewRecorder()

	handler := getTaskHandler(mockStore)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %s", ct)
	}

	var tasks []Task
	if err := json.Unmarshal(rec.Body.Bytes(), &tasks); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(tasks) != 1 || tasks[0].Title != "Test Task" {
		t.Fatalf("unexpected tasks: %+v", tasks)
	}
}

func TestGetTaskHandler_FilterCompleted(t *testing.T) {
	trueVal := true
	mockStore := &MockStore{
		GetAllTasksFunc: func(filter *bool) ([]Task, error) {
			if filter == nil || *filter != trueVal {
				t.Fatalf("expected filter=true, got %+v", filter)
			}
			return []Task{{ID: 1, Title: "Completed Task", Completed: true}}, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/tasks?completed=true", nil)
	rec := httptest.NewRecorder()

	getTaskHandler(mockStore).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}
}

func TestGetTaskHandler_FilterIncomplete(t *testing.T) {
	falseVal := false
	mockStore := &MockStore{
		GetAllTasksFunc: func(filter *bool) ([]Task, error) {
			if filter == nil || *filter != falseVal {
				t.Fatalf("expected filter=false, got %+v", filter)
			}
			return []Task{{ID: 2, Title: "Incomplete Task", Completed: false}}, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/tasks?completed=false", nil)
	rec := httptest.NewRecorder()

	getTaskHandler(mockStore).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}
}

func TestGetTaskByIDHandler_Found(t *testing.T) {
	mockStore := &MockStore{
		GetTaskByIDFunc: func(id int) (Task, error) {
			return Task{ID: id, Title: "Found Task"}, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	rec := httptest.NewRecorder()

	getTaskByIDHandler(mockStore).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}
}

func TestGetTaskByIDHandler_NotFound(t *testing.T) {
	mockStore := &MockStore{
		GetTaskByIDFunc: func(id int) (Task, error) {
			return Task{}, ErrNotFound
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	rec := httptest.NewRecorder()

	getTaskByIDHandler(mockStore).ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 Not Found, got %d", rec.Code)
	}
}

func TestPostTaskHandler_CreatesTask(t *testing.T) {
	mockStore := &MockStore{
		CreateTaskFunc: func(task *Task) error {
			task.ID = 99
			return nil
		},
	}

	body := bytes.NewBufferString(`{"title":"New Task","description":"desc"}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", body)
	rec := httptest.NewRecorder()

	postTaskHandler(mockStore).ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", rec.Code)
	}
}

func TestPostTaskHandler_InvalidJSON(t *testing.T) {
	mockStore := &MockStore{}
	body := bytes.NewBufferString(`{bad json}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", body)
	rec := httptest.NewRecorder()

	postTaskHandler(mockStore).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", rec.Code)
	}
}

func TestPostTaskHandler_MissingTitle(t *testing.T) {
	mockStore := &MockStore{
		CreateTaskFunc: func(task *Task) error {
			return ErrInvalid
		},
	}
	body := bytes.NewBufferString(`{"description":"no title"}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", body)
	rec := httptest.NewRecorder()

	postTaskHandler(mockStore).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", rec.Code)
	}
}

func TestUpdateTaskHandler_UpdatesTask(t *testing.T) {
	mockStore := &MockStore{
		UpdateTaskFunc: func(task *Task) error {
			return nil
		},
	}

	body := bytes.NewBufferString(`{"title":"Updated Task","completed":true}`)
	req := httptest.NewRequest(http.MethodPut, "/tasks/1", body)
	rec := httptest.NewRecorder()

	updateTaskByIDHandler(mockStore).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}
}

func TestUpdateTaskHandler_InvalidJSON(t *testing.T) {
	mockStore := &MockStore{}
	body := bytes.NewBufferString(`{bad json}`)
	req := httptest.NewRequest(http.MethodPut, "/tasks/1", body)
	rec := httptest.NewRecorder()

	updateTaskByIDHandler(mockStore).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", rec.Code)
	}
}

func TestUpdateTaskHandler_MissingTitle(t *testing.T) {
	mockStore := &MockStore{
		UpdateTaskFunc: func(task *Task) error {
			return ErrInvalid
		},
	}
	body := bytes.NewBufferString(`{"description":"no title"}`)
	req := httptest.NewRequest(http.MethodPut, "/tasks/1", body)
	rec := httptest.NewRecorder()

	updateTaskByIDHandler(mockStore).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", rec.Code)
	}
}

func TestDeleteTaskHandler_DeletesTask(t *testing.T) {
	mockStore := &MockStore{
		DeleteTaskFunc: func(id int) error {
			return nil
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	rec := httptest.NewRecorder()

	deleteTaskByIDHandler(mockStore).ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204 No Content, got %d", rec.Code)
	}
}

// DELETE /tasks/{id} - not found
func TestDeleteTaskHandler_NotFound(t *testing.T) {
	mockStore := &MockStore{
		DeleteTaskFunc: func(id int) error {
			return ErrNotFound
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	rec := httptest.NewRecorder()

	deleteTaskByIDHandler(mockStore).ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 Not Found, got %d", rec.Code)
	}
}
