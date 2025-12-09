package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"ecom-go-backend/internal/database"
)

func AddToCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	var body struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}

	json.NewDecoder(r.Body).Decode(&body)
	if body.Quantity <= 0 {
		body.Quantity = 1
	}

	query := `
		INSERT INTO cart_items (user_id, product_id, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, product_id)
		DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity;
	`

	_, err := database.DB.Exec(context.Background(), query, userID, body.ProductID, body.Quantity)
	if err != nil {
		http.Error(w, "Failed to add to cart: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "added_to_cart"})
}

func RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	var body struct {
		ProductID int `json:"product_id"`
	}

	json.NewDecoder(r.Body).Decode(&body)

	_, err := database.DB.Exec(context.Background(),
		"DELETE FROM cart_items WHERE user_id=$1 AND product_id=$2",
		userID, body.ProductID)
	if err != nil {
		http.Error(w, "Failed to remove from cart", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "removed_from_cart"})
}

func GetCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	query := `
		SELECT p.id, p.name, p.description, p.category_id, p.price, p.stock, p.image_url, c.quantity
		FROM cart_items c
		JOIN products p ON c.product_id = p.id
		WHERE c.user_id = $1
		ORDER BY c.added_at DESC;
	`

	rows, err := database.DB.Query(context.Background(), query, userID)
	if err != nil {
		http.Error(w, "Failed to fetch cart", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var cart []map[string]interface{}

	for rows.Next() {
		var item struct {
			ID          int
			Name        string
			Description string
			CategoryID  int
			Price       float64
			Stock       int
			ImageUrl    string
			Quantity    int
		}

		err := rows.Scan(
			&item.ID, &item.Name, &item.Description, &item.CategoryID,
			&item.Price, &item.Stock, &item.ImageUrl, &item.Quantity,
		)
		if err != nil {
			http.Error(w, "Failed to scan cart item: "+err.Error(), http.StatusInternalServerError)
			return
		}

		cart = append(cart, map[string]interface{}{
			"id":          item.ID,
			"name":        item.Name,
			"description": item.Description,
			"category_id": item.CategoryID,
			"price":       item.Price,
			"stock":       item.Stock,
			"image_url":   item.ImageUrl,
			"quantity":    item.Quantity,
		})
	}

	json.NewEncoder(w).Encode(cart)
}
