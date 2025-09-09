package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type CheckoutRequest struct {
	UserID int `json:"user_id"`
	Items  []struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	} `json:"items"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Only POST method allowed")
		return
	}

	var req CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if len(req.Items) == 0 {
		respondWithError(w, http.StatusBadRequest, "Cart is empty")
		return
	}

	tx, err := Db.Begin(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to start transaction")
		return
	}
	defer tx.Rollback(context.Background())

	total := 0.0

	for _, item := range req.Items {
		var price float64
		var stock int

		err := tx.QueryRow(context.Background(),
			"SELECT price, stock FROM products WHERE id=$1 FOR UPDATE", item.ProductID).Scan(&price, &stock)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Product ID %d not found", item.ProductID))
			return
		}

		if stock < item.Quantity {
			respondWithError(w, http.StatusBadRequest,
				fmt.Sprintf("Insufficient stock for product ID %d. Available: %d, Requested: %d",
					item.ProductID, stock, item.Quantity))
			return
		}

		total += price * float64(item.Quantity)

		// Update stock
		_, err = tx.Exec(context.Background(),
			"UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to update product stock")
			return
		}
	}

	var orderID int
	err = tx.QueryRow(context.Background(),
		"INSERT INTO orders (user_id, total) VALUES ($1, $2) RETURNING id",
		req.UserID, total).Scan(&orderID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create order")
		return
	}

	for _, item := range req.Items {
		_, err = tx.Exec(context.Background(),
			"INSERT INTO order_products (order_id, product_id, quantity) VALUES ($1, $2, $3)",
			orderID, item.ProductID, item.Quantity)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to insert order item")
			return
		}
	}

	if err = tx.Commit(context.Background()); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Order created successfully",
		"order_id": orderID,
		"total":    total,
	})
}
