package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"ecom-go-backend/internal/database"
	"ecom-go-backend/internal/models"
)

func AddToWishlist(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	var body struct {
		ProductID int `json:"product_id"`
	}

	json.NewDecoder(r.Body).Decode(&body)

	query := `
        INSERT INTO wishlist (user_id, product_id, added_at)
        VALUES ($1, $2, NOW())
        ON CONFLICT (user_id, product_id)
        DO NOTHING;
    `
	_, err := database.DB.Exec(context.Background(), query, userID, body.ProductID)
	if err != nil {
		http.Error(w, "Failed to add to wishlist", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "added"})
}

func RemoveFromWishlist(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	var body struct {
		ProductID int `json:"product_id"`
	}

	json.NewDecoder(r.Body).Decode(&body)

	query := `
    INSERT INTO wishlist (user_id, product_id)
    VALUES ($1, $2)
    ON CONFLICT DO NOTHING;
	`

	_, err := database.DB.Exec(context.Background(), query, userID, body.ProductID)
	if err != nil {
		http.Error(w, "Failed to remove from wishlist", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "removed"})
}

func GetWishlist(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	query := `
        SELECT p.id, p.name, p.description, p.category_id, p.price, p.stock, p.image_url
        FROM wishlist w
        JOIN products p ON w.product_id = p.id
        WHERE w.user_id = $1
        ORDER BY w.added_at DESC;
    `

	rows, err := database.DB.Query(context.Background(), query, userID)
	if err != nil {
		http.Error(w, "Failed to fetch wishlist", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.CategoryID, &p.Price, &p.Stock, &p.ImageURL)
		if err != nil {
			http.Error(w, "Failed to parse wishlist products", http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
