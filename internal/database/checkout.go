package database

import (
	"context"
	"ecom-go-backend/internal/models"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func Checkout(ctx context.Context, db *pgx.Conn, userID int) (*models.CheckoutResponse, error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Get cart items with price
	rows, err := tx.Query(ctx, `
		SELECT c.product_id, c.quantity, p.price
		FROM cart_items c
		JOIN products p ON p.id = c.product_id
		WHERE c.user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}

	var items []models.CheckoutItemInfo
	var total float64 = 0

	for rows.Next() {
		var item models.CheckoutItemInfo
		err := rows.Scan(&item.ProductID, &item.Quantity, &item.Price)
		if err != nil {
			return nil, err
		}
		total += float64(item.Quantity) * item.Price
		items = append(items, item)
	}
	rows.Close()

	if len(items) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	// Create order
	var orderID int64
	err = tx.QueryRow(ctx, `
		INSERT INTO orders (user_id, total_amount)
		VALUES ($1, $2)
		RETURNING id
	`, userID, total).Scan(&orderID)
	if err != nil {
		return nil, err
	}

	// Insert order items
	for _, it := range items {
		_, err := tx.Exec(ctx, `
			INSERT INTO order_items (order_id, product_id, quantity, price)
			VALUES ($1, $2, $3, $4)
		`, orderID, it.ProductID, it.Quantity, it.Price)

		if err != nil {
			return nil, err
		}
	}

	// Clear cart
	_, err = tx.Exec(ctx, `DELETE FROM cart_items WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &models.CheckoutResponse{
		OrderID:     orderID,
		TotalAmount: total,
		Items:       items,
	}, nil
}
