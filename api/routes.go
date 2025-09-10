package api

import (
	"ecom-go-backend/internal/handlers"
	"ecom-go-backend/internal/middleware"
	"net/http"
)

func SetupRoutes() {
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)

	// Protected routes: wrap with JWTAuth middleware
	http.Handle("/categories", middleware.JWTAuth(http.HandlerFunc(handlers.CategoriesHandler)))
	http.Handle("/products", middleware.JWTAuth(http.HandlerFunc(handlers.ProductsHandler)))
	http.Handle("/checkout", middleware.JWTAuth(http.HandlerFunc(handlers.CheckoutHandler)))
	http.Handle("/order-history", middleware.JWTAuth(http.HandlerFunc(handlers.OrderHistoryHandler)))
}
