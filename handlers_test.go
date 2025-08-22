package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockPotionStore struct {
	GetAllPotionsFunc func() ([]Potion, error)
}

func (m *MockPotionStore) GetAllPotions() ([]Potion, error) {
	return m.GetAllPotionsFunc()
}

func TestGetPotionsHandler_ReturnsList(t *testing.T) {
	mock := &MockPotionStore{
		GetAllPotionsFunc: func() ([]Potion, error) {
			return []Potion{
				{ID: 1, Name: "Healing", Power: "Restore HP"},
			}, nil
		},
	}

	req := httptest.NewRequest("GET", "/potions", nil)
	rec := httptest.NewRecorder()

	handler := getPotionsHandler(mock)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}

	var potions []Potion
	if err := json.NewDecoder(rec.Body).Decode(&potions); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	if len(potions) != 1 || potions[0].Name != "Healing" {
		t.Errorf("unexpected respones: %+v", potions)
	}
}
