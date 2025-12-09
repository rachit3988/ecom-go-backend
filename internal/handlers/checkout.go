package handlers

import (
	"context"
	"ecom-go-backend/internal/models"
	"encoding/json"
	"net/http"
)

func Checkout(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	// 1. Fetch cart items
	rows, err := Db.Query(context.Background(), `
		SELECT c.product_id, c.quantity, p.price
		FROM cart_items c
		JOIN products p ON p.id = c.product_id
		WHERE c.user_id = $1
	`, userID)
	if err != nil {
		http.Error(w, "Failed to fetch cart: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []models.CheckoutItemInfo
	var total float64 = 0

	for rows.Next() {
		var it models.CheckoutItemInfo
		err := rows.Scan(&it.ProductID, &it.Quantity, &it.Price)
		if err != nil {
			http.Error(w, "Failed to scan cart item: "+err.Error(), http.StatusInternalServerError)
			return
		}
		total += float64(it.Quantity) * it.Price
		items = append(items, it)
	}

	if len(items) == 0 {
		http.Error(w, "Cart is empty", http.StatusBadRequest)
		return
	}

	// 2. Begin transaction
	tx, err := Db.Begin(context.Background())
	if err != nil {
		http.Error(w, "Failed to begin transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(context.Background())

	// 3. Create order
	var orderID int64
	err = tx.QueryRow(context.Background(), `
		INSERT INTO orders (user_id, total_amount)
		VALUES ($1, $2)
		RETURNING id
	`, userID, total).Scan(&orderID)
	if err != nil {
		http.Error(w, "Failed to create order: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Insert order items
	for _, it := range items {
		_, err := tx.Exec(context.Background(), `
			INSERT INTO order_items (order_id, product_id, quantity, price)
			VALUES ($1, $2, $3, $4)
		`, orderID, it.ProductID, it.Quantity, it.Price)
		if err != nil {
			http.Error(w, "Failed to insert order items: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// 5. Clear cart
	_, err = tx.Exec(context.Background(),
		"DELETE FROM cart_items WHERE user_id=$1", userID)
	if err != nil {
		http.Error(w, "Failed to clear cart: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 6. Commit transaction
	if err := tx.Commit(context.Background()); err != nil {
		http.Error(w, "Failed to commit: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 7. Send response
	resp := models.CheckoutResponse{
		OrderID:     orderID,
		TotalAmount: total,
		Items:       items,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
