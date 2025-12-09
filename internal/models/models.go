package models

import "time"

type User struct {
	ID       int    `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"` // Store hash, not plaintext
}

type Category struct {
	ID       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	ImageURL string `json:"image_url" db:"image_url"`
}

type Product struct {
	ID          int     `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	Description string  `json:"description" db:"description"`
	CategoryID  int     `json:"category_id" db:"category_id"`
	Price       float64 `json:"price" db:"price"`
	Stock       int     `json:"stock" db:"stock"`
	ImageURL    string  `json:"image_url" db:"image_url"`
}

type OrderProductDetail struct {
	ProductID   int     `json:"product_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

type OrderDetail struct {
	OrderID int                  `json:"order_id"`
	Total   float64              `json:"total"`
	Items   []OrderProductDetail `json:"items"`
}

type RecentlyViewed struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	ProductID int       `json:"product_id" db:"product_id"`
	ViewedAt  time.Time `json:"viewed_at" db:"viewed_at"`
}

type WishlistItem struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	ProductID int    `json:"product_id"`
	AddedAt   string `json:"added_at"`
}

type CheckoutResponse struct {
	OrderID     int64              `json:"order_id"`
	TotalAmount float64            `json:"total_amount"`
	Items       []CheckoutItemInfo `json:"items"`
}

type CheckoutItemInfo struct {
	ProductID int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type PastOrderItem struct {
	ProductID int     `json:"product_id"`
	Name      string  `json:"name"`
	ImageURL  string  `json:"image_url"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type PastOrder struct {
	OrderID     int64           `json:"order_id"`
	TotalAmount float64         `json:"total_amount"`
	CreatedAt   string          `json:"created_at"`
	Items       []PastOrderItem `json:"items"`
}
