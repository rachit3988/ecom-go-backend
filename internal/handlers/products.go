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
	// Optional filters: product ID via ?id=1 or category via ?category=1
	categoryIDStr := r.URL.Query().Get("category")
	productIDStr := r.URL.Query().Get("id")

	var rows pgx.Rows
	var err error

	switch {
	case productIDStr != "":
		productID, convErr := strconv.Atoi(productIDStr)
		if convErr != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}
		rows, err = Db.Query(context.Background(),
			"SELECT id, name, description, category_id, price, stock, image_urls FROM products WHERE id=$1", productID)
	case categoryIDStr != "":
		categoryID, convErr := strconv.Atoi(categoryIDStr)
		if convErr != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}
		rows, err = Db.Query(context.Background(),
			"SELECT id, name, description, category_id, price, stock, image_urls FROM products WHERE category_id=$1", categoryID)
	default:
		rows, err = Db.Query(context.Background(),
			"SELECT id, name, description, category_id, price, stock, image_urls FROM products")
	}

	if err != nil {
		http.Error(w, "Failed to fetch products: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.CategoryID, &p.Price, &p.Stock, &p.ImageURLs)
		if err != nil {
			http.Error(w, "Failed to scan product: "+err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
