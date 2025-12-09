package handlers

import (
	"context"
	"ecom-go-backend/internal/models"
	"encoding/json"
	"net/http"
)

func GetPastOrders(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	// Fetch user's orders
	rows, err := Db.Query(context.Background(), `
		SELECT id, total_amount, created_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		http.Error(w, "Failed to fetch orders: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var orders []models.PastOrder

	for rows.Next() {
		var o models.PastOrder
		err := rows.Scan(&o.OrderID, &o.TotalAmount, &o.CreatedAt)
		if err != nil {
			http.Error(w, "Failed to scan order: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Fetch order items for each order
		itemRows, err := Db.Query(context.Background(), `
			SELECT 
				oi.product_id,
				p.name,
				p.image_url,
				oi.quantity,
				oi.price
			FROM order_items oi
			JOIN products p ON p.id = oi.product_id
			WHERE oi.order_id = $1
		`, o.OrderID)
		if err != nil {
			http.Error(w, "Failed to fetch order items: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var items []models.PastOrderItem
		for itemRows.Next() {
			var it models.PastOrderItem
			err := itemRows.Scan(&it.ProductID, &it.Name, &it.ImageURL, &it.Quantity, &it.Price)
			if err != nil {
				http.Error(w, "Failed to scan order item: "+err.Error(), http.StatusInternalServerError)
				return
			}
			items = append(items, it)
		}
		itemRows.Close()

		o.Items = items
		orders = append(orders, o)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}
