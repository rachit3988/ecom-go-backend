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
	http.HandleFunc("/checkout", handlers.CheckoutHandler)
}
