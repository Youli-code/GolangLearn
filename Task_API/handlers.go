package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type errResp struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errResp{Error: msg})
}
func getTaskHandler(store TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var filter *bool
		if val := r.URL.Query().Get("completed"); val != "" {
			if val == "true" {
				t := true
				filter = &t
			} else if val == "false" {
				f := false
				filter = &f
			} else {
				writeErr(w, http.StatusBadRequest, "invalid completed value")
				return
			}
		}

		tasks, err := store.GetAllTasks(filter)
		if err != nil {
			writeErr(w, http.StatusInternalServerError, "internal error")
			return
		}
		writeJSON(w, http.StatusOK, tasks)
	}
}

func postTaskHandler(store TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newTask Task
		if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
			writeErr(w, http.StatusBadRequest, "invalid json")
			return
		}
		if err := store.CreateTask(&newTask); err != nil {
			switch {
			case errors.Is(err, ErrInvalid):
				writeErr(w, http.StatusBadRequest, err.Error())
			default:
				writeErr(w, http.StatusInternalServerError, "internal error")
			}
			return
		}
		writeJSON(w, http.StatusCreated, newTask)
	}
}

func getTaskByIDHandler(store TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("ID")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			writeErr(w, http.StatusBadRequest, "invalid id")
			return
		}
		task, err := store.GetTaskByID(id)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				writeErr(w, http.StatusNotFound, "not found")
			} else {
				writeErr(w, http.StatusInternalServerError, "internal error")
			}
			return
		}
		writeJSON(w, http.StatusOK, task)
	}
}

func deleteTaskByIDHandler(store TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("ID")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			writeErr(w, http.StatusBadRequest, "invalid id")
			return
		}
		if err := store.DeleteTask(id); err != nil {
			if errors.Is(err, ErrNotFound) {
				writeErr(w, http.StatusNotFound, "not found")
			} else {
				writeErr(w, http.StatusInternalServerError, "internal error")
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	}
}

func updateTaskByIDHandler(store TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("ID")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			writeErr(w, http.StatusBadRequest, "invalid id")
			return
		}
		var updated Task
		if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
			writeErr(w, http.StatusBadRequest, "invalid json")
			return
		}
		updated.ID = id
		if err := store.UpdateTask(&updated); err != nil {
			switch {
			case errors.Is(err, ErrInvalid):
				writeErr(w, http.StatusBadRequest, err.Error())
			case errors.Is(err, ErrNotFound):
				writeErr(w, http.StatusNotFound, "not found")
			default:
				writeErr(w, http.StatusInternalServerError, "internal error")
			}
			return
		}
		writeJSON(w, http.StatusOK, updated)
	}
}
