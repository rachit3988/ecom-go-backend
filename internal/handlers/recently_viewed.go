package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"ecom-go-backend/internal/database"
	"ecom-go-backend/internal/models"
)

func AddRecentlyViewed(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	var body struct {
		ProductID int `json:"product_id"`
	}

	json.NewDecoder(r.Body).Decode(&body)

	query := `
        INSERT INTO recently_viewed (user_id, product_id, viewed_at)
        VALUES ($1, $2, NOW())
        ON CONFLICT (user_id, product_id)
        DO UPDATE SET viewed_at = NOW();
    `
	_, err := database.DB.Exec(context.Background(), query, userID, body.ProductID)
	if err != nil {
		http.Error(w, "Failed to add recently viewed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "added"})
}

func GetRecentlyViewed(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	query := `
        SELECT p.id, p.name, p.description, p.category_id, p.price, p.stock, p.image_url
        FROM recently_viewed r
        JOIN products p ON r.product_id = p.id
        WHERE r.user_id = $1
        ORDER BY r.viewed_at DESC
        LIMIT 20;
    `

	rows, err := database.DB.Query(context.Background(), query, userID)
	if err != nil {
		http.Error(w, "Failed to fetch recently viewed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.CategoryID, &p.Price, &p.Stock, &p.ImageURLs)
		if err != nil {
			http.Error(w, "Failed to parse recently viewed products", http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
