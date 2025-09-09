package handlers

import (
	"context"
	"ecom-go-backend/internal/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
)

func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	// Optional: Filter by category ID via query param ?category=1
	categoryIDStr := r.URL.Query().Get("category")
	var rows pgx.Rows
	var err error

	if categoryIDStr != "" {
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}
		rows, err = Db.Query(context.Background(),
			"SELECT id, name, description, category_id, price, stock FROM products WHERE category_id=$1", categoryID)
	} else {
		rows, err = Db.Query(context.Background(),
			"SELECT id, name, description, category_id, price, stock FROM products")
	}

	if err != nil {
		http.Error(w, "Failed to fetch products: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.CategoryID, &p.Price, &p.Stock)
		if err != nil {
			http.Error(w, "Failed to scan product: "+err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
