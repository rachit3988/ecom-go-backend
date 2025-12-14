package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"ecom-go-backend/internal/models"
)

func GetMostPopular(w http.ResponseWriter, r *http.Request) {
	rows, err := Db.Query(context.Background(),
		`SELECT 
            p.id,
            p.name,
            p.description,
            p.category_id,
            p.price,
            p.stock,
            p.image_url,
            COUNT(oi.product_id) AS popularity_score
        FROM order_items oi
        JOIN products p ON p.id = oi.product_id
        GROUP BY p.id
        ORDER BY popularity_score DESC
        LIMIT 20`,
	)
	if err != nil {
		http.Error(w, "Failed to fetch popular products: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		var p models.Product
		var popularity int

		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.CategoryID,
			&p.Price,
			&p.Stock,
			&p.ImageURLs,
			&popularity,
		)
		if err != nil {
			http.Error(w, "Failed to scan product: "+err.Error(), http.StatusInternalServerError)
			return
		}

		products = append(products, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
