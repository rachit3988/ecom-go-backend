package models

type User struct {
	ID       int    `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"` // Store hash, not plaintext
}

type Category struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type Product struct {
	ID          int     `db:"id"`
	Name        string  `db:"name"`
	Description string  `db:"description"`
	CategoryID  int     `db:"category_id"`
	Price       float64 `db:"price"`
	Stock       int     `db:"stock"`
}

type Order struct {
	ID         int     `db:"id"`
	UserID     int     `db:"user_id"`
	ProductIDs []int   `db:"product_ids"` // Can later use order_items table for details
	Total      float64 `db:"total"`
}
