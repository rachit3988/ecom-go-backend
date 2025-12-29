package main

import (
	"log"
	"net/http"
	"os"

	"ecom-go-backend/api"
	"ecom-go-backend/internal/database"
	"ecom-go-backend/internal/handlers"
	"ecom-go-backend/internal/middleware"

	"github.com/joho/godotenv"
)

func main() {

	// Load env FIRST
	_ = godotenv.Load()

	// Init JWT AFTER env
	middleware.InitJWT()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := database.ConnectDB(); err != nil {
		log.Fatal("DB connection failed:", err)
	}
	defer database.DB.Close()

	handlers.Db = database.DB

	api.SetupRoutes()

	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
