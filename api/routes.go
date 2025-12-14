package api

import (
	"ecom-go-backend/internal/handlers"
	"ecom-go-backend/internal/middleware"
	"net/http"
)

func SetupRoutes() {
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/categories", handlers.CategoriesHandler)
	http.HandleFunc("/products", handlers.ProductsHandler)
	http.HandleFunc("/most-popular", handlers.ProductsHandler)
	http.HandleFunc("/recently-viewed", handlers.GetRecentlyViewed)

	// Protected routes: wrap with JWTAuth middleware
	http.Handle("/recently-viewed/add", middleware.JWTAuth(http.HandlerFunc(handlers.AddRecentlyViewed)))
	http.Handle("/cart/add", middleware.JWTAuth(http.HandlerFunc(handlers.AddToCart)))
	http.Handle("/cart/remove", middleware.JWTAuth(http.HandlerFunc(handlers.RemoveFromCart)))
	http.Handle("/cart/get", middleware.JWTAuth(http.HandlerFunc(handlers.GetCart)))
	http.Handle("/checkout", middleware.JWTAuth(http.HandlerFunc(handlers.Checkout)))
	http.Handle("/past-orders", middleware.JWTAuth(http.HandlerFunc(handlers.GetPastOrders)))
}
