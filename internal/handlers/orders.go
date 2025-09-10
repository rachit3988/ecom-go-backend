package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"ecom-go-backend/internal/models"
)

func OrderHistoryHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id query param is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	// Fetch orders for user
	ordersRows, err := Db.Query(context.Background(),
		"SELECT id, total FROM orders WHERE user_id=$1 ORDER BY id DESC", userID)
	if err != nil {
		http.Error(w, "Failed to fetch orders: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer ordersRows.Close()

	var orders []models.OrderDetail

	for ordersRows.Next() {
		var order models.OrderDetail
		err := ordersRows.Scan(&order.OrderID, &order.Total)
		if err != nil {
			http.Error(w, "Failed to scan order: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Fetch products for this order
		productsRows, err := Db.Query(context.Background(), `
            SELECT p.id, p.name, p.description, op.quantity, p.price
            FROM order_products op
            JOIN products p ON p.id = op.product_id
            WHERE op.order_id = $1`, order.OrderID)
		if err != nil {
			http.Error(w, "Failed to fetch order products: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var items []models.OrderProductDetail

		for productsRows.Next() {
			var item models.OrderProductDetail
			err := productsRows.Scan(&item.ProductID, &item.Name, &item.Description, &item.Quantity, &item.Price)
			if err != nil {
				productsRows.Close()
				http.Error(w, "Failed to scan product: "+err.Error(), http.StatusInternalServerError)
				return
			}
			items = append(items, item)
		}
		productsRows.Close()

		order.Items = items
		orders = append(orders, order)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}
