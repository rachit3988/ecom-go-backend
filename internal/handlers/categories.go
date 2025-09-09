package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"ecom-go-backend/internal/models"
)

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := Db.Query(context.Background(), "SELECT id, name FROM categories")
	if err != nil {
		http.Error(w, "Failed to fetch categories: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		err := rows.Scan(&c.ID, &c.Name)
		if err != nil {
			http.Error(w, "Failed to scan category: "+err.Error(), http.StatusInternalServerError)
			return
		}
		categories = append(categories, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
