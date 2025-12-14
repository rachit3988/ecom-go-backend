package main

import (
	"ecom-go-backend/api"
	"ecom-go-backend/internal/database"
	"ecom-go-backend/internal/handlers"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

func main() {
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
