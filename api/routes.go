package api

import (
	"ecom-go-backend/internal/handlers"
	"net/http"
)

func SetupRoutes() {
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/categories", handlers.CategoriesHandler)
	http.HandleFunc("/products", handlers.ProductsHandler)
	http.HandleFunc("/most-popular", handlers.ProductsHandler)
	http.HandleFunc("/recently-viewed", handlers.ProductsHandler)
	http.HandleFunc("/recently-viewed/add", handlers.AddRecentlyViewed)
	http.HandleFunc("/cart/add", handlers.AddToCart)
	http.HandleFunc("/cart/remove", handlers.RemoveFromCart)
	http.HandleFunc("/cart/get", handlers.GetCart)
	http.HandleFunc("/checkout", handlers.CheckoutHandler)
	http.HandleFunc("/past-orders", handlers.GetPastOrders)
}
