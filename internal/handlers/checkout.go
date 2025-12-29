package handlers

import (
	"context"
	"encoding/json"
	"net/http"
)

type CheckoutItemRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type CheckoutRequest struct {
	Items []CheckoutItemRequest `json:"items"`
}

type CheckoutItemResponse struct {
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type CheckoutResponse struct {
	OrderID     int64                  `json:"order_id"`
	TotalAmount float64                `json:"total_amount"`
	Items       []CheckoutItemResponse `json:"items"`
}

func CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	userIDVal := r.Context().Value("user_id")
	if userIDVal == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := userIDVal.(int)
	if !ok {
		http.Error(w, "Invalid user context", http.StatusUnauthorized)
		return
	}

	var req CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Items) == 0 {
		http.Error(w, "No items provided", http.StatusBadRequest)
		return
	}

	tx, err := Db.Begin(context.Background())
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(context.Background())

	var total float64
	var orderItems []CheckoutItemResponse

	for _, item := range req.Items {
		if item.Quantity <= 0 {
			http.Error(w, "Invalid quantity", http.StatusBadRequest)
			return
		}

		var price float64
		err := tx.QueryRow(context.Background(),
			"SELECT price FROM products WHERE id=$1",
			item.ProductID,
		).Scan(&price)

		if err != nil {
			http.Error(w, "Invalid product_id", http.StatusBadRequest)
			return
		}

		total += price * float64(item.Quantity)

		orderItems = append(orderItems, CheckoutItemResponse{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     price,
		})
	}

	var orderID int64
	err = tx.QueryRow(context.Background(),
		`INSERT INTO orders (user_id, total_amount)
		 VALUES ($1, $2)
		 RETURNING id`,
		userID, total,
	).Scan(&orderID)

	if err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	for _, it := range orderItems {
		_, err := tx.Exec(context.Background(),
			`INSERT INTO order_items (order_id, product_id, quantity, price)
			 VALUES ($1, $2, $3, $4)`,
			orderID, it.ProductID, it.Quantity, it.Price,
		)
		if err != nil {
			http.Error(w, "Failed to insert order items", http.StatusInternalServerError)
			return
		}
	}

	// optional: clear cart
	_, _ = tx.Exec(context.Background(),
		"DELETE FROM cart_items WHERE user_id=$1", userID)

	if err := tx.Commit(context.Background()); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	resp := CheckoutResponse{
		OrderID:     orderID,
		TotalAmount: total,
		Items:       orderItems,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
