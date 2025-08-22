package main

import (
	"encoding/json"
	"net/http"
)

func getPotionsHandler(store PotionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		potions, err := store.GetAllPotions()
		if err != nil {
			http.Error(w, "Failed to fetch potions", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(potions)
	}
}
